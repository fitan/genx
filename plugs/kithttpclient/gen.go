package kithttpclient

import (
	"fmt"
	"go/types"
	"log/slog"
	"regexp"
	"sort"
	"strings"

	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
	"github.com/samber/lo"
)

type StructMeta struct {
}

type StructFieldMeta struct {
}

type KitHttpClient struct {
	methods         []*MethodMeta
	option          gen.Option
	implGoTypeMetes []gen.InterfaceGoTypeMeta
}

type MethodMeta struct {
	Name             string
	InterfaceFunc    *types.Func
	BasePath         string
	HttpUrl          string
	HttpMethod       string
	RequestTypeName  string
	ResponseTypeName string
	RequestBody      bool
	HasResponseErr   bool

	PathParams   OrderRequestParamMap
	BodyParams   OrderRequestParamMap
	QueryParams  OrderRequestParamMap
	HeaderParams OrderRequestParamMap

	StructFieldMetaData []common.StructFieldMetaDataV2
}

func (m *MethodMeta) Parse() {
	m.structFieldMetaData2RequestParamMap(*m.StructFieldMetaData[0].NextFields, []int{})
}

func (m *MethodMeta) structFieldMetaData2RequestParamMap(fields []common.StructFieldMetaDataV2, index []int) {
	for i, v := range fields {
		tag, ok := v.Tag.Lookup("param")
		if !ok {
			continue
		}

		tagValues := strings.Split(tag, ",")

		if len(tagValues) != 2 {
			err := fmt.Errorf("param tag error: %s", tag)
			panic(err)
		}

		paramType := tagValues[0]
		paramName := tagValues[1]
		pathList := []string{}
		indexes := append(index, i)
		current := fields
		lo.ForEach(indexes, func(item int, index int) {
			pathList = append(pathList, current[item].Name)
			current = *current[item].NextFields
		})

		paramPathStr := strings.Join(pathList, ".")

		switch paramType {
		case "path":
			m.PathParams[paramName] = RequestParam{
				StructFieldMetaDataIndex: indexes,
				ParamPath:                paramPathStr,
				TypeX:                    v.TypeX,
			}
		case "query":
			m.QueryParams[paramName] = RequestParam{
				ParamPath:                paramPathStr,
				TypeX:                    v.TypeX,
				StructFieldMetaDataIndex: indexes,
			}
		case "header":
			m.HeaderParams[paramName] = RequestParam{
				ParamPath:                paramPathStr,
				TypeX:                    v.TypeX,
				StructFieldMetaDataIndex: indexes,
			}
		case "body":
			m.BodyParams[paramName] = RequestParam{
				ParamPath:                paramPathStr,
				TypeX:                    v.TypeX,
				StructFieldMetaDataIndex: indexes,
			}
		}

		if len(*v.NextFields) > 0 {
			m.structFieldMetaData2RequestParamMap(*v.NextFields, append(index, i))
		}
	}
}

func (o OrderRequestParamMap) OrderSlice() (res []RequestParam) {
	keys := make([]string, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		res = append(res, o[k])
	}

	return res
}

func (m MethodMeta) ParamPath(v string) (res string) {
	defer func() {
		if res != "" {
			res = "." + res
		}
	}()

	param, ok := m.PathParams[v]
	if ok {
		return param.ParamPath
	}

	param, ok = m.BodyParams[v]
	if ok {
		return param.ParamPath
	}

	param, ok = m.QueryParams[v]
	if ok {
		return param.ParamPath
	}

	param, ok = m.HeaderParams[v]
	if ok {
		return param.ParamPath
	}

	if v == "req" {
		return ""
	}

	panic("param not found: " + v)
}

type OrderRequestParamMap map[string]RequestParam

type RequestParam struct {
	ParamPath                string
	TypeX                    *common.Type
	StructFieldMetaDataIndex []int
}

func (k *KitHttpClient) Parse() error {
	parseImpl := common.NewInterfaceSerialize(k.option.Pkg)

	for _, v := range k.implGoTypeMetes {
		interfaceMetaDate, err := parseImpl.Parse(v.Obj, v.RawDoc, &v.Doc)
		if err != nil {
			return common.ParseError("failed to parse interface").
				WithCause(err).
				WithPlugin("@kit-http-client").
				WithInterface(v.Obj.String()).
				WithDetails("unable to parse interface metadata").
				Build()
		}

		var basePath string
		interfaceMetaDate.Doc.ByFuncNameAndArgs("@basePath", &basePath)
		slog.Info("parsing interface", "interface", v.Obj.String(), "basePath", basePath)

		for _, method := range interfaceMetaDate.Methods {
			var httpUrl, httpMethod string
			var requestTypeName string
			var requestBody string

			has := method.Doc.ByFuncNameAndArgs("@kit-http", &httpUrl, &httpMethod)
			if !has {
				return common.ValidationError("missing required annotation").
					WithPlugin("@kit-http-client").
					WithInterface(v.Obj.String()).
					WithMethod(method.Name).
					WithAnnotation("@kit-http").
					WithDetails("@kit-http annotation is required for HTTP client generation. Format: @kit-http <url> <method>").
					Build()
			}

			has = method.Doc.ByFuncNameAndArgs("@kit-http-request", &requestTypeName, &requestBody)
			if !has {
				return common.ValidationError("missing required annotation").
					WithPlugin("@kit-http-client").
					WithInterface(v.Obj.String()).
					WithMethod(method.Name).
					WithAnnotation("@kit-http-request").
					WithDetails("@kit-http-request annotation is required. Format: @kit-http-request <RequestType> [body_flag]").
					Build()
			}

			if len(method.Results) != 2 {
				return common.ValidationError("invalid method signature").
					WithPlugin("@kit-http-client").
					WithInterface(v.Obj.String()).
					WithMethod(method.Name).
					WithDetails(fmt.Sprintf("method must return exactly 2 values (response, error), got %d", len(method.Results))).
					Build()
			}

			v2 := common.NewStructSerializeV2(k.option.Pkg, method.Params[1].Type)

			mm := &MethodMeta{
				Name:                method.Name,
				InterfaceFunc:       &types.Func{},
				BasePath:            basePath,
				HttpUrl:             httpUrl,
				HttpMethod:          httpMethod,
				RequestTypeName:     requestTypeName,
				ResponseTypeName:    method.Results[0].ID,
				RequestBody:         requestBody != "" && requestBody != "false",
				HasResponseErr:      method.ReturnsError,
				PathParams:          map[string]RequestParam{},
				BodyParams:          map[string]RequestParam{},
				QueryParams:         map[string]RequestParam{},
				HeaderParams:        map[string]RequestParam{},
				StructFieldMetaData: *v2.Fields,
			}

			// 使用 recovery 机制来捕获 Parse 方法中可能的 panic
			if err := common.WithRecovery(func() error {
				mm.Parse()
				return nil
			}); err != nil {
				return common.ParseError("failed to parse method metadata").
					WithCause(err).
					WithPlugin("@kit-http-client").
					WithInterface(v.Obj.String()).
					WithMethod(method.Name).
					WithDetails("error occurred while parsing method metadata").
					Build()
			}

			k.methods = append(k.methods, mm)
		}
	}

	return nil
}

func (k KitHttpClient) Gen(j *jen.File) {
	j.AddImport("encoding/json", "")
	j.AddImport("strings", "")
	j.AddImport("github.com/asaskevich/govalidator", "valid")
	j.AddImport("github.com/go-kit/kit/endpoint", "endpoint")
	j.AddImport("github.com/go-kit/kit/transport/http", "kithttp")
	j.AddImport("github.com/gorilla/mux", "")
	j.AddImport("github.com/pkg/errors", "")
	j.AddImport("github.com/spf13/cast", "")
	j.AddImport("github.com/go-kit/log", "")
	j.AddImport("io", "")
	j.AddImport("net/http", "")
	j.AddImport("net/url", "")
	j.AddImport("time", "")
	j.AddImport("fmt", "")
	j.AddImport("context", "")
	j.AddImport("github.com/go-kit/kit/sd", "")
	j.AddImport("github.com/go-kit/kit/sd/lb", "")

	j.Id(`
type HttpClientService struct {
	Option
}

type Option struct {
	PrePath      string
	Logger       log.Logger
	Instancer    sd.Instancer
	RetryMax     int
	RetryTimeout time.Duration
	EndpointOpts []sd.EndpointerOption
	ClientOpts   []kithttp.ClientOption
	Encode       kithttp.EncodeRequestFunc
	Decode       func(i interface{}) func(ctx context.Context, res *http.Response) (response interface{}, err error)
}

func (s HttpClientService) mergeOpt(option *Option) Option {
	if option == nil {
		return s.Option
	}

	if option.RetryMax == 0 {
		option.RetryMax = s.RetryMax
	}
	if option.RetryTimeout == 0 {
		option.RetryTimeout = s.RetryTimeout
	}
	if option.EndpointOpts != nil {
		option.EndpointOpts = append(s.EndpointOpts, option.EndpointOpts...)
	} else {
		option.EndpointOpts = s.EndpointOpts
	}
	if option.ClientOpts != nil {
		option.ClientOpts = append(s.ClientOpts, option.ClientOpts...)
	} else {
		option.ClientOpts = s.ClientOpts
	}

	if option.Encode == nil {
		option.Encode = s.Encode
	}

	if option.Decode == nil {
		option.Decode = s.Decode
	}
	return *option
}

func NewHttpClientService(opt Option) *HttpClientService {
	return &HttpClientService{
		Option: opt,
	}
}
	`)

	j.Type().Id("HttpClientImpl").InterfaceFunc(func(g *jen.Group) {
		for _, v := range k.methods {
			g.Id(v.Name).Params(jen.Id("ctx context.Context"), jen.Id("req "+v.RequestTypeName), jen.Id("option *Option")).Params(jen.Id("res "+v.ResponseTypeName), jen.Id("err error")).Line()
		}

	})

	for _, v := range k.methods {
		code := make([]jen.Code, 0)
		queryCode := make([]jen.Code, 0)

		if len(v.QueryParams) != 0 {
			code = append(code,
				jen.Id("q").Op(":=").Qual("net/url", "Values").Values(),
			)

			for k, v := range v.QueryParams {
				if !v.TypeX.Basic {
					queryCode = append(queryCode,
						jen.Id("q.Add").Call(jen.Lit(k), jen.Qual("github.com/spf13/cast", "ToString").Call(jen.Id("req."+v.ParamPath))),
					)
					continue
				}

				if v.TypeX.BasicType.Kind() != types.String {
					queryCode = append(queryCode,
						jen.Id("q.Add").Call(jen.Lit(k), jen.Qual("github.com/spf13/cast", "ToString").Call(jen.Id("req."+v.ParamPath))),
					)
					continue
				}

				queryCode = append(queryCode,
					jen.If(jen.Id("req."+v.ParamPath).Op("!=").Lit("")).Block(
						jen.Id("q.Add").Call(jen.Lit(k), jen.Id("req."+v.ParamPath)),
					),
				)
			}
			code = append(code, queryCode...)
		}

		re, _ := regexp.Compile(`\{(.*?)\}`)
		matches := re.FindAllStringSubmatch(v.HttpUrl, -1)
		params := make([]string, 0)
		for _, match := range matches {
			if len(match) > 1 {
				params = append(params, match[1])
			}
		}

		urlCode := make([]jen.Code, 0)
		fmtSprintParam := make([]jen.Code, 0)
		fmtSprintParam = append(fmtSprintParam, jen.Lit(re.ReplaceAllString(v.HttpUrl, "%s")))
		for _, param := range params {
			fmtSprintParam = append(fmtSprintParam, jen.Id("req"+v.ParamPath(param)))
		}
		if len(params) != 0 {
			urlCode = append(
				urlCode,
				jen.List(jen.Id("urlStr"), jen.Err()).Op(":=").Qual("net/url", "JoinPath").Call(jen.Id("s.PrePath"), jen.Lit(v.BasePath), jen.Id("fmt.Sprintf").Call(
					fmtSprintParam...,
				)),
				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Err().Op("=").Qual("fmt", "Errorf").Call(jen.Lit("parse prePath %s basePath %s error: %v"), jen.Id("s.PrePath"), jen.Lit(v.BasePath), jen.Err()),
					jen.Return(),
				),
			)
		} else {
			urlCode = append(
				urlCode,
				jen.List(jen.Id("urlStr"), jen.Err()).Op(":=").Qual("net/url", "JoinPath").Call(jen.Id("s.PrePath"), jen.Lit(v.BasePath), jen.Lit(re.ReplaceAllString(v.HttpUrl, "%s"))),
				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Err().Op("=").Qual("fmt", "Errorf").Call(jen.Lit("parse prePath %s basePath %s error: %v"), jen.Id("s.PrePath"), jen.Lit(v.BasePath), jen.Err()),
					jen.Return(),
				),
			)
		}

		urlCode = append(urlCode,
			jen.List(jen.Id("u"), jen.Err()).Op(":=").Qual("net/url", "Parse").CallFunc(func(group *jen.Group) {
				if len(queryCode) != 0 {
					group.Qual("fmt", "Sprintf").Call(jen.Lit("http://%s%s?%s"), jen.Id("instance"), jen.Qual("net/url", "PathEscape").Call(jen.Qual("strings", "TrimRight").Call(jen.Id("urlStr"), jen.Lit("/"))), jen.Id("q.Encode").Call())
				} else {
					group.Qual("fmt", "Sprintf").Call(jen.Lit("http://%s%s"), jen.Id("instance"), jen.Qual("strings", "TrimRight").Call(jen.Id("urlStr"), jen.Lit("/")))
				}
			}),
		)

		urlCode = append(urlCode,
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Err().Op("=").Qual("fmt", "Errorf").Call(jen.Lit("parse url %s error: %v"), jen.Id("urlStr"), jen.Err()),
				jen.Return(),
			),
		)

		var setHeaderCode []jen.Code
		for k, v := range v.HeaderParams {
			setHeaderCode = append(setHeaderCode,
				jen.Id("r.Header.Set").Call(jen.Lit(k), jen.Id("req."+v.ParamPath)),
			)
		}

		if len(setHeaderCode) != 0 {
			code = append(code,
				jen.Id(`headerOpt := kithttp.ClientBefore(func(ctx context.Context, r *http.Request) context.Context`).BlockFunc(func(g *jen.Group) {
					g.Add(setHeaderCode...)
					g.Return(jen.Id("ctx"))
				}),
				jen.Id("opt.ClientOpts").Op("=").Append(jen.Id("opt.ClientOpts"), jen.Id("headerOpt")),
			)
		}

		encode := "kithttp.EncodeJSONRequest"

		// if len(formCode) != 0 {
		// 	encode = "EncodeFormRequest(form.FormDataContentType())"
		// }

		code = append(code,
			jen.Id("factory").Op(":=").Id(`func (instance string) (ee endpoint.Endpoint,ic io.Closer,err error)`).BlockFunc(func(group *jen.Group) {
				for _, v := range urlCode {
					group.Add(v)
				}
				group.Return(jen.Id("kithttp.NewClient").Call(jen.Lit(v.HttpMethod), jen.Id("u"), jen.Id(encode), func() jen.Code {
					if v.HasResponseErr {
						return jen.Id("opt.Decode(&res)")
					}
					return jen.Id("opt.Decode(nil)")
				}(),
					jen.Id("opt.ClientOpts..."),
				).Dot("Endpoint").Call(), jen.Nil(), jen.Nil())
			}),
			jen.Id(`e := sd.NewEndpointer(s.Instancer, factory, s.Logger,s.EndpointOpts...)`),
			jen.Id(`balancer := lb.NewRoundRobin(e)`),
			jen.Id(`retry := lb.Retry(opt.RetryMax, opt.RetryTimeout, balancer)`),
		)

		code = append(code,
			jen.List(jen.Id("_"), jen.Err()).Op("=").Id("retry").Call(jen.Id("ctx"), func() jen.Code {
				if v.RequestBody {
					return jen.Id("req")
				}

				if len(v.BodyParams) != 0 {
					return jen.Id("req." + v.BodyParams.OrderSlice()[0].ParamPath)
				}

				/* if len(m.KitRequest.Form) != 0 {
					return jen.Id("&formB")
				} */

				return jen.Nil()
			}()),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Id("err").Op("=").Qual("fmt", "Errorf").Call(jen.Lit("endpoint error: %v"), jen.Err()),
				jen.Return(),
			),

			jen.Return(),
		)

		j.Func().Params(jen.Id("s").Op("*").Id("HttpClientService")).Id(v.Name).Params(jen.Id("ctx context.Context"), jen.Id("req "+v.RequestTypeName), jen.Id("option *Option")).Params(jen.Id("res "+v.ResponseTypeName), jen.Id("err error")).BlockFunc(func(g *jen.Group) {
			g.Id(`
			_, err = valid.ValidateStruct(req)
			if err != nil {
				err = fmt.Errorf("validate request error: %v", err)
				return
			}
			opt := s.mergeOpt(option)
			`)
			lo.ForEach(code, func(item jen.Code, index int) {
				g.Add(item).Line()
			})
		})

	}
}
