package kithttp

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
	"sort"
	"strings"
	"unicode"

	"github.com/dave/jennifer/jen"
	"github.com/fitan/genx/common"
	"github.com/samber/lo"
	"golang.org/x/tools/go/packages"
)

const (
	RequestParamTagName string = "param"
	QueryTag            string = "query"
	RawQueryTag         string = "rawQuery"
	HeaderTag           string = "header"
	PathTag             string = "path"
	BodyTag             string = "body"
	FileTag             string = "file"
	CtxTag              string = "ctx"
	EndpointCtxTag      string = "endpointCtx"
	FormTag             string = "form"

	DocKitHttpParamMark     string = "@kit-http-param"
	DocKitEndpointParamMark string = "@kit-endpoint-param"
)

type KitRequest struct {
	pkg *packages.Package

	ServiceName   string
	RequestTypeOf *types.Struct
	RequestName   string
	RequestIsBody bool
	RequestIsNil  bool

	NamedMap    map[string]string
	RawQuery    OrderRequestParamMap
	Query       OrderRequestParamMap
	Path        OrderRequestParamMap
	Body        OrderRequestParamMap
	File        OrderRequestParamMap
	Header      OrderRequestParamMap
	Ctx         OrderRequestParamMap
	EndpointCtx OrderRequestParamMap
	Form        OrderRequestParamMap
	Empty       OrderRequestParamMap
}

type OrderRequestParamMap map[string]RequestParam

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

func NewKitRequest(pkg *packages.Package, serviceName, requestName string, requestIsBody bool) *KitRequest {
	return &KitRequest{
		pkg:           pkg,
		ServiceName:   serviceName,
		RequestName:   requestName,
		RequestIsBody: requestIsBody,
		NamedMap:      make(map[string]string),
		RawQuery:      make(OrderRequestParamMap),
		Query:         make(OrderRequestParamMap),
		Path:          make(OrderRequestParamMap),
		Body:          make(OrderRequestParamMap),
		File:          make(OrderRequestParamMap),
		Header:        make(OrderRequestParamMap),
		Ctx:           make(OrderRequestParamMap),
		EndpointCtx:   make(OrderRequestParamMap),
		Form:          make(OrderRequestParamMap),
		Empty:         make(OrderRequestParamMap),
	}
}

type RequestParam struct {
	ParamDoc *ast.CommentGroup

	ParamPath string

	FieldName string

	ParamName string
	// path, query, header, body, empty
	ParamSource string
	// basic, map, slice,ptr
	ParamType string
	// time [time.Time]
	ParamTypeName string
	// []string map[string]string
	RawParamType string
	// int,string,bool,float
	BasicType string

	HasPtr bool

	HasNamed bool

	XType *common.Type

	XTypeID string
}

func (r RequestParam) FormDataSwagType() string {
	switch r.XTypeID {
	case "*multipart.FileHeader", "interface{io.Reader; io.ReaderAt; io.Seeker; io.Closer}", "[]*multipart.FileHeader":
		return "file"
	default:
		return "string"
	}
}

func (r RequestParam) ParamNameAlias() string {
	return "_" + strings.ReplaceAll(r.ParamName, "-", "_")
}

func (r RequestParam) Comment() string {
	var item []string
	if r.ParamDoc == nil {
		return ""
	}
	for _, v := range r.ParamDoc.List {
		item = append(item, v.Text)
	}
	return strings.Join(item, "\n")
}

func (r RequestParam) Annotations() string {
	if r.ParamDoc == nil {
		return `" "`
	}
	for _, v := range r.ParamDoc.List {
		docFormat := DocFormat(v.Text)
		if strings.HasPrefix(docFormat, "// "+r.FieldName) {
			return fmt.Sprintf(`"%s"`, strings.TrimPrefix(docFormat, "// "+r.FieldName))
		}
	}
	return fmt.Sprintf(`"%s"`, strings.TrimPrefix(r.ParamDoc.List[0].Text, "// "))
}

func (r RequestParam) ToVal() jen.Code {
	return jen.Var().Id(r.ParamNameAlias()).Id(r.ParamTypeName)
	// return jen.Var().Id(r.ParamNameAlias()).Id(r.ParamTypeName)
	//switch r.ParamType {
	//case "basic":
	//	return jen.Var().Id(r.ParamName).Id(r.BasicType)
	//case "map":
	//	return jen.Var().Id(r.ParamName).Map(jen.String()).Id(r.BasicType).Values()
	//case "slice":
	//	return jen.Var().Id(r.ParamName).Index().Id(r.BasicType)
	//case "struct":
	//	return jen.Var().Id(r.ParamName).Id(r.BasicType)
	//
	//}
	//return nil
}

func (k *KitRequest) ParamPath(paramName string) (res string) {
	defer func() {
		if res != "" {
			res = "." + res
		}
	}()

	if strings.ToUpper(paramName) == strings.ToUpper(k.RequestName) {
		return ""
	}
	param, ok := k.RawQuery[paramName]
	if ok {
		return param.ParamPath
	}
	param, ok = k.Query[paramName]
	if ok {
		return param.ParamPath
	}

	param, ok = k.Path[paramName]
	if ok {
		return param.ParamPath
	}

	param, ok = k.File[paramName]
	if ok {
		return param.ParamPath
	}

	param, ok = k.Header[paramName]
	if ok {
		return param.ParamPath
	}

	param, ok = k.Body[paramName]
	if ok {
		return param.ParamPath
	}

	param, ok = k.Ctx[paramName]
	if ok {
		return param.ParamPath
	}

	param, ok = k.EndpointCtx[paramName]
	if ok {
		return param.ParamPath
	}

	param, ok = k.Form[paramName]
	if ok {
		return param.ParamPath
	}

	param, ok = k.Empty[paramName]
	if ok {
		return param.ParamPath
	}

	if paramName == "req" {
		return ""
	}

	panic("param not found: " + paramName)
}

func (k *KitRequest) SetParam(param RequestParam) {
	switch param.ParamSource {
	case RawQueryTag:
		k.RawQuery[param.ParamName] = param
	case QueryTag:
		k.Query[param.ParamName] = param
	case PathTag:
		k.Path[param.ParamName] = param
	case HeaderTag:
		k.Header[param.ParamName] = param
	case BodyTag:
		k.Body[param.ParamName] = param
	case CtxTag:
		k.Ctx[param.ParamName] = param
	case EndpointCtxTag:
		k.EndpointCtx[param.ParamName] = param
	case FormTag:
		k.Form[param.ParamName] = param
	case FileTag:
		k.File[param.ParamName] = param
	case "":
		k.Empty[param.ParamName] = param

	default:
		panic("param source error: " + param.ParamSource + "," + param.ParamName)
	}
}

func (k *KitRequest) ParseParamTag(fieldName, tag string) (paramSource string, paramName string) {

	split := strings.Split(tag, ",")
	if len(split) == 1 {
		return split[0], downFirst(fieldName)
	}

	if len(split) == 2 {
		return split[0], split[1]
	}

	return "", ""

}

func (k *KitRequest) DeepCopyRequest() *jen.Statement {
	if k.RequestIsBody {
		return jen.Type().Id(k.RequestName + "HttpBody").Id(common.TypeOf(k.RequestTypeOf).TypeAsJen().GoString())
	}

	if len(k.Body) > 0 {
		return jen.Type().Id(k.RequestName + "HttpBody").Id(k.Body.OrderSlice()[0].XType.TypeAsJen().GoString())
	}

	return jen.Null()
}

func (k *KitRequest) DecodeRequest() (s string) {
	listCode := make([]jen.Code, 0, 0)
	// req := Request{}
	if k.RequestName != "nil" {
		listCode = append(listCode, jen.Id("req").Op(":=").Id(k.RequestName).Block())
		listCode = append(listCode, k.DefineVal()...)
		listCode = append(listCode, k.BindBodyParam()...)
		listCode = append(listCode, k.BindPathParam()...)
		listCode = append(listCode, k.BindRawQueryPram()...)
		listCode = append(listCode, k.BindQueryParam()...)
		listCode = append(listCode, k.BindHeaderParam()...)
		listCode = append(listCode, k.BindFileParam()...)
		listCode = append(listCode, k.BindFormParam()...)
		listCode = append(listCode, k.BindCtxParam()...)
		listCode = append(listCode, k.BindRequest()...)
		listCode = append(listCode, jen.Return(jen.Id("req"), jen.Id("err")))
	} else {
		listCode = append(listCode, jen.Return(jen.Id("nil"), jen.Id("nil")))
	}
	var LineListCode []jen.Code
	for _, v := range listCode {
		LineListCode = append(LineListCode, jen.Line(), v)
	}

	fn := jen.Func().Id("decode"+upFirst(k.ServiceName)+"Request").Params(
		jen.Id("ctx").Id("context.Context"),
		jen.Id("r").Id("*http").Dot("Request"),
	).Call(
		jen.Id("res").Interface(),
		jen.Id("err").Id("error"),
	).Block(
		LineListCode...,
	)
	return fn.GoString()
}

func (k *KitRequest) DefineVal() []jen.Code {
	listCode := make([]jen.Code, 0, 0)
	for _, v := range k.RawQuery.OrderSlice() {
		listCode = append(listCode, v.ToVal())
	}
	for _, v := range k.Query.OrderSlice() {
		listCode = append(listCode, v.ToVal())
	}
	for _, v := range k.Path.OrderSlice() {
		listCode = append(listCode, v.ToVal())
	}
	for _, v := range k.Header.OrderSlice() {
		listCode = append(listCode, v.ToVal())
	}

	for _, v := range k.Form.OrderSlice() {
		listCode = append(listCode, v.ToVal())
	}

	for _, v := range k.Ctx.OrderSlice() {
		listCode = append(listCode, v.ToVal())
	}
	return listCode
}

func (k *KitRequest) Validate() []jen.Code {
	list := make([]jen.Code, 0, 0)
	list = append(
		list,
		jen.List(jen.Id("_"), jen.Id("err")).Op("=").Id("valid").Dot("ValidateStruct").Call(jen.Id("req")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Err().Op("=").Id("errors.Wrap").Call(jen.Id("err"), jen.Lit("valid.ValidateStruct")),
			jen.Return(),
		),
		//jen.If(jen.Id("!validRes")).Block(
		//	jen.Err().Op("=").Id("fmt.Errorf").Call(jen.Lit("valid false")),
		//	jen.Return(),
		//),
	)
	return list
}

func (k *KitRequest) V10() []jen.Code {
	list := make([]jen.Code, 0, 0)
	list = append(
		list,
		jen.List(jen.Id("err")).Op("=").Id("validate").Dot("Struct").Call(jen.Id("req")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Err().Op("=").Id("encode.InvalidParams.Wrap").Call(jen.Id("err")),
			jen.Return(),
		),
	)
	return list
}

func (k *KitRequest) BindBodyParam() []jen.Code {
	listCode := make([]jen.Code, 0, 0)
	if k.RequestIsNil {
		return listCode
	}
	returnCode := jen.If(jen.Err().Op("!=").Nil()).Block(
		jen.Err().Op("=").Id("errors.Wrap").Call(jen.Id("err"), jen.Lit("decode body")),
		jen.Return(),
	)
	if k.RequestIsBody {
		// err = json.NewDecoder(r.Body).Decode(&req)
		decode := jen.Id("err").Op("=").Id("json.NewDecoder").Call(jen.Id("r.Body")).Dot("Decode").Parens(jen.Id("&req"))
		listCode = append(listCode, decode, returnCode)

		return listCode
	}

	if len(k.Body) == 0 {
		return listCode
	}

	if len(k.Body) != 1 {
		panic("body param count error " + fmt.Sprint(len(k.Body)))
	}

	for _, v := range k.Body.OrderSlice() {
		if v.ParamTypeName == "[]byte" {
			decode := jen.List(jen.Id("req."+v.ParamPath), jen.Id("err")).Op("=").Qual("io/ioutil", "ReadAll").Call(jen.Id("r.Body"))
			listCode = append(listCode, decode, returnCode)
		} else {
			decode := jen.Id("err").Op("=").Id("json.NewDecoder").Call(jen.Id("r.Body")).Dot("Decode").Parens(jen.Id("&req." + v.ParamPath))
			listCode = append(listCode, decode, returnCode)
		}
	}

	return listCode
}

func (k *KitRequest) BindFileParam() []jen.Code {
	list := make([]jen.Code, 0, 0)
	for _, v := range k.File.OrderSlice() {
		list = append(list, jen.List(jen.Id("_"), jen.Id("req."+v.ParamPath), jen.Id("err")).Op("=").Id("r.FormFile").Call(jen.Lit(v.ParamName)).Line().
			If(jen.Err().Op("!=").Nil()).Block(
			jen.Err().Op("=").Id("errors.Wrap").Call(jen.Id("err"), jen.Lit("r.FormFile")),
		))
	}
	return list
}

func (k *KitRequest) BindHeaderParam() []jen.Code {
	list := make([]jen.Code, 0, 0)

	for _, v := range k.Header.OrderSlice() {
		//r.Header.Get("project")
		varBind := jen.Id("r.Header.Get").Call(jen.Lit(v.ParamName))
		if v.BasicType != "string" {
			// cast.ToInt(vars["id"])
			varBind = jen.Id("cast").Dot("To" + upFirst(v.BasicType) + "E").Call(varBind)
			// id, err := cast.ToIntE(vars["id"])
			varBind = jen.List(jen.Id(v.ParamNameAlias()), jen.Err()).Op("=").Add(varBind)
			// if err != nil {
			// 	return err
			// }
			returnCode := jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(),
			)
			list = append(list, varBind, returnCode)
			continue
		}
		// id = vars["id"]
		val := jen.Id(v.ParamNameAlias()).Op("=").Add(varBind)
		list = append(list, val)
	}

	return list

}

func (k *KitRequest) BindRawQueryPram() []jen.Code {
	list := make([]jen.Code, 0, 0)
	if len(k.RawQuery) == 0 {
		return list
	}
	for _, v := range k.RawQuery.OrderSlice() {
		urlValues, err := UrlValues(v.ParamName, v.ParamType, v.ParamTypeName)
		if err != nil {
			panic(err)
		}

		list = append(list, urlValues...)
		continue
	}

	return list
}

func (k *KitRequest) BindQueryParam() []jen.Code {
	list := make([]jen.Code, 0, 0)

	if len(k.Query) == 0 {
		return list
	}

	for _, v := range k.Query.OrderSlice() {
		//r.URL.Query().Get("project")
		varBind := jen.Id("r.URL.Query().Get").Call(jen.Lit(v.ParamName))

		if !(v.ParamType == "basic" && v.BasicType == "string") {
			castCode, err := CastMap(v.XType, v.ParamNameAlias(), v.ParamType, v.ParamTypeName, varBind)
			if err != nil {
				panic(err)
			}
			list = append(list, castCode...)
			continue
		}
		// id = vars["id"]
		val := jen.Id(v.ParamNameAlias()).Op("=").Add(varBind)
		list = append(list, val)
	}

	return list
}

func (k *KitRequest) BindPathParam() []jen.Code {
	list := make([]jen.Code, 0, 0)

	if len(k.Path) == 0 {
		return list
	}

	// vars := mux.Vars(r)
	vars := jen.Id("vars").Op(":=").Qual("github.com/gorilla/mux", "Vars").Call(jen.Id("r"))
	list = append(list, vars)
	for _, v := range k.Path.OrderSlice() {
		// vars["id"]
		varBind := jen.Id("vars").Index(jen.Lit(v.ParamName))
		if v.BasicType != "string" {
			// cast.ToInt(vars["id"])
			varBind = jen.Id("cast").Dot("To" + upFirst(v.BasicType) + "E").Call(varBind)
			// id, err := cast.ToIntE(vars["id"])
			varBind = jen.List(jen.Id(v.ParamNameAlias()), jen.Err()).Op("=").Add(varBind)
			// if err != nil {
			// 	return err
			// }
			returnCode := jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(),
			)
			list = append(list, varBind, returnCode)
			continue
		}
		// id = vars["id"]
		val := jen.Id(v.ParamNameAlias()).Op("=").Add(varBind)
		list = append(list, val)
	}

	return list
}

func (k *KitRequest) BindFormParam() []jen.Code {
	list := make([]jen.Code, 0, 0)
	if len(k.Form) != 0 {
		parse := jen.Id("err").Op("=").Id("r.").Id("ParseMultipartForm(32 << 20)").Line().
			If(jen.Err().Op("!=").Nil()).Block(
			jen.Err().Op("=").Id("errors.Wrap").Call(jen.Id("err"), jen.Lit("r.ParseMultipartForm")),
			jen.Return(),
		)
		list = append(list, parse)
	}
	for _, v := range k.Form.OrderSlice() {
		tID := v.XType.TypeAsJenComparePkgNameString(k.pkg)
		if tID == "[]*multipart.FileHeader" {
			code := jen.Id(v.ParamNameAlias()).Op("=").Id("r.MultipartForm.File").Index(jen.Lit(v.ParamName)).Line()
			list = append(list, code)
			continue
		}

		if tID == "interface{io.Reader; io.ReaderAt; io.Seeker; io.Closer}" {
			code := jen.Id(v.ParamNameAlias() + ", _, err").Op("=").Id("r.FormFile").Call(jen.Lit(v.ParamName)).Line()
			code = code.If(jen.Err().Op("!=").Nil().Op("&&").Id("!errors.Is").Call(jen.Id("err"), jen.Qual("net/http", "ErrMissingFile"))).Block(
				jen.Id("err").Op("=").Id("errors.Wrap").Call(jen.Id("err"), jen.Lit("FormFile")),
				jen.Return(),
			)
			list = append(list, code)
			continue
		}

		if tID == "*multipart.FileHeader" {
			code := jen.Id("_ ," + v.ParamNameAlias() + ", err").Op("=").Id("r.FormFile").Call(jen.Lit(v.ParamName)).Line()
			code = code.If(jen.Err().Op("!=").Nil()).Block(
				jen.Id("err").Op("=").Id("errors.Wrap").Call(jen.Id("err"), jen.Lit("FormFile")),
				jen.Return(),
			)
			list = append(list, code)
			continue
		}

		if v.XType.Basic {
			if v.XType.BasicType.Kind() == types.String {
				code := jen.Id(v.ParamNameAlias()).Op("=").Id("r.FormValue").Call(jen.Lit(v.ParamName)).Line()
				list = append(list, code)
				continue
			} else {
				code, err := CastMap(v.XType, v.ParamNameAlias(), v.ParamType, v.ParamTypeName, jen.Id("r.FormValue").Call(jen.Lit(v.ParamName)))
				if err != nil {
					panic("not support form param type " + v.ParamNameAlias() + "tID: " + tID)
				}
				list = append(list, code...)
				continue
			}
		}

		if v.XType.Struct || v.XType.List || v.XType.Map {
			code := jen.Id("err").Op("=").Id("json.Unmarshal").Call(jen.Id("[]byte").Call(jen.Id("r.FormValue").Call(jen.Lit(v.ParamName))), jen.Op("&").Id(v.ParamNameAlias())).Line()
			code = code.If(jen.Err().Op("!=").Nil()).Block(
				jen.Id("err").Op("=").Id("errors.Wrap").Call(jen.Id("err"), jen.Lit("FormValue json.Unmarshal")),
				jen.Return(),
			)
			list = append(list, code)
			continue
		}

		panic("not support form param type " + v.ParamNameAlias() + "tID: " + tID)
	}

	return list
}

func (k *KitRequest) BindCtxParam() []jen.Code {
	list := make([]jen.Code, 0, 0)
	for _, v := range k.Ctx.OrderSlice() {
		var ctxKey string

		if v.ParamDoc == nil {
			panic("ctx param doc is nil")
		}
		for _, d := range v.ParamDoc.List {
			fields := strings.Fields(d.Text)
			if fields[1] == DocKitHttpParamMark {
				if len(fields) < 3 {
					panic("ctx param doc error: " + d.Text)
				}

				if fields[2] == "ctx" {
					ctxKey = fields[3]
				}

			}
		}
		if ctxKey == "" {
			panic("not find ctx param doc error: " + v.ParamDoc.Text())
		}
		ctxVal := jen.Var().Id(v.ParamName + "OK").Bool()
		varBind := jen.List(jen.Id(v.ParamNameAlias()), jen.Id(v.ParamName+"OK")).Op("=").Id("ctx.Value").Call(jen.Id(ctxKey)).Assert(jen.Id(v.RawParamType))
		ifBind := jen.If(jen.Id(v.ParamName+"OK")).Op("==").False().Block(
			jen.Err().Op("=").Id("errors.New").Call(jen.Lit("ctx param "+v.ParamName+" is not found")),
			jen.Return(),
		)
		list = append(list, ctxVal, varBind, ifBind)
	}
	return list
}

func (k *KitRequest) BindEndpointCtxParam() string {
	list := make([]jen.Code, 0, 0)
	list = append(list, k.EndpointCtxDefineVal()...)
	for _, v := range k.EndpointCtx.OrderSlice() {
		var ctxKey string

		if v.ParamDoc == nil {
			panic("ctx param doc is nil")
		}
		for _, d := range v.ParamDoc.List {
			fields := strings.Fields(d.Text)
			if fields[1] == DocKitEndpointParamMark {
				if len(fields) < 3 {
					panic("ctx param doc error: " + d.Text)
				}

				if fields[2] == "ctx" {
					ctxKey = fields[3]
				}

			}
		}
		if ctxKey == "" {
			panic("not find ctx param doc error: " + v.ParamDoc.Text())
		}
		ctxVal := jen.Var().Id(v.ParamName + "OK").Bool()
		varBind := jen.List(jen.Id(v.ParamNameAlias()), jen.Id(v.ParamName+"OK")).Op("=").Id("ctx.Value").Call(jen.Id(ctxKey)).Assert(jen.Id(v.RawParamType))
		ifBind := jen.If(jen.Id(v.ParamName+"OK")).Op("==").False().Block(
			jen.Err().Op("=").Id("errors.New").Call(jen.Lit("ctx param "+v.ParamName+" is not found")),
			jen.Return(),
		)
		list = append(list, ctxVal, varBind, ifBind)
	}
	list = append(list, k.BindEndpointCtxRequest()...)
	j := jen.Null()
	for _, v := range list {
		j.Add(v, jen.Line())
	}
	return j.GoString()
}

func (k *KitRequest) EndpointCtxDefineVal() []jen.Code {
	listCode := make([]jen.Code, 0, 0)

	for _, v := range k.EndpointCtx.OrderSlice() {
		listCode = append(listCode, v.ToVal())
	}
	return listCode
}

func (k *KitRequest) BindEndpointCtxRequest() []jen.Code {
	list := make([]jen.Code, 0, 0)
	for _, v := range k.EndpointCtx.OrderSlice() {
		code := lo.Ternary(v.HasNamed,
			jen.Call(lo.Ternary(v.HasPtr, jen.Id("*"+v.ParamTypeName), jen.Id(v.ParamTypeName))).Call(jen.Id(lo.Ternary(v.HasPtr, "&", "")+v.ParamNameAlias())),
			jen.Id(lo.Ternary(v.HasPtr, "&", "")+v.ParamNameAlias()),
		)
		reqBindVal := jen.Id("req").Dot(v.ParamPath).Op("=").Add(code)
		list = append(list, reqBindVal)
	}

	return list
}

func (k *KitRequest) BindRequest() []jen.Code {
	list := make([]jen.Code, 0, 0)
	for _, v := range k.Form.OrderSlice() {
		code := lo.Ternary(v.HasNamed,
			jen.Call(lo.Ternary(v.HasPtr, jen.Id("*"+v.ParamTypeName), jen.Id(v.ParamTypeName))).Call(jen.Id(lo.Ternary(v.HasPtr, "&", "")+v.ParamNameAlias())),
			jen.Id(lo.Ternary(v.HasPtr, "&", "")+v.ParamNameAlias()),
		)
		reqBindVal := jen.Id("req").Dot(v.ParamPath).Op("=").Add(code)
		list = append(list, reqBindVal)
	}

	for _, v := range k.RawQuery.OrderSlice() {
		code := lo.Ternary(v.HasNamed,
			jen.Call(lo.Ternary(v.HasPtr, jen.Id("*"+v.ParamTypeName), jen.Id(v.ParamTypeName))).Call(jen.Id(lo.Ternary(v.HasPtr, "&", "")+v.ParamNameAlias())),
			jen.Id(lo.Ternary(v.HasPtr, "&", "")+v.ParamNameAlias()),
		)
		reqBindVal := jen.Id("req").Dot(v.ParamPath).Op("=").Add(code)
		list = append(list, reqBindVal)
	}
	for _, v := range k.Path.OrderSlice() {
		code := lo.Ternary(v.HasNamed,
			jen.Call(lo.Ternary(v.HasPtr, jen.Id("*"+v.ParamTypeName), jen.Id(v.ParamTypeName))).Call(jen.Id(lo.Ternary(v.HasPtr, "&", "")+v.ParamNameAlias())),
			jen.Id(lo.Ternary(v.HasPtr, "&", "")+v.ParamNameAlias()),
		)

		reqBindVal := jen.Id("req").Dot(v.ParamPath).Op("=").Add(code)
		list = append(list, reqBindVal)
	}

	for _, v := range k.Query.OrderSlice() {
		code := lo.Ternary(v.HasNamed,
			jen.Call(lo.Ternary(v.HasPtr, jen.Id("*"+v.ParamTypeName), jen.Id(v.ParamTypeName))).Call(jen.Id(lo.Ternary(v.HasPtr, "&", "")+v.ParamNameAlias())),
			jen.Id(lo.Ternary(v.HasPtr, "&", "")+v.ParamNameAlias()),
		)
		reqBindVal := jen.Id("req").Dot(v.ParamPath).Op("=").Add(code)

		if v.HasPtr {
			reqBindVal = jen.If(jen.Id(v.ParamNameAlias() + "Str" + `!= ""`)).Block(reqBindVal)
		}
		list = append(list, reqBindVal)
	}

	for _, v := range k.Header.OrderSlice() {
		code := lo.Ternary(v.HasNamed,
			jen.Call(lo.Ternary(v.HasPtr, jen.Id("*"+v.ParamTypeName), jen.Id(v.ParamTypeName))).Call(jen.Id(lo.Ternary(v.HasPtr, "&", "")+v.ParamNameAlias())),
			jen.Id(lo.Ternary(v.HasPtr, "&", "")+v.ParamNameAlias()),
		)
		reqBindVal := jen.Id("req").Dot(v.ParamPath).Op("=").Add(code)
		list = append(list, reqBindVal)
	}

	for _, v := range k.Ctx.OrderSlice() {
		code := lo.Ternary(v.HasNamed,
			jen.Call(lo.Ternary(v.HasPtr, jen.Id("*"+v.ParamTypeName), jen.Id(v.ParamTypeName))).Call(jen.Id(lo.Ternary(v.HasPtr, "&", "")+v.ParamNameAlias())),
			jen.Id(lo.Ternary(v.HasPtr, "&", "")+v.ParamNameAlias()),
		)
		reqBindVal := jen.Id("req").Dot(v.ParamPath).Op("=").Add(code)
		list = append(list, reqBindVal)
	}

	// var c string
	// for k, v := range k.NamedMap {
	// 	c = c + fmt.Sprintf("k: %s, v: %s", k, v)
	// }

	// list = append(list, jen.Comment(c))

	return list
}

func (k *KitRequest) ParseRequest() {
	var hasFindRequest bool
	for _, s := range k.pkg.Syntax {
		ast.Inspect(s, func(node ast.Node) bool {
			switch nodeT := node.(type) {
			case *ast.GenDecl:
				for _, spec := range nodeT.Specs {
					if specT, ok := spec.(*ast.TypeSpec); ok {
						if specT.Name.Name == k.RequestName {
							hasFindRequest = true
							doc := nodeT.Doc
							k.Doc(doc)
							t := k.pkg.TypesInfo.TypeOf(specT.Type).(*types.Struct)
							k.RequestTypeOf = t
							k.RequestType([]string{}, k.RequestName, t, "", doc, false)
							k.CheckRequestIsNil()

							return false
						}
					}
				}

			}
			return true
		})
	}
	if !hasFindRequest {
		panic("not find request" + k.RequestName)
	}
}

func (k *KitRequest) Doc(doc *ast.CommentGroup) {
}

func (k *KitRequest) ParseFieldComment(pos token.Pos) (s *ast.CommentGroup) {
	fieldFileName := k.pkg.Fset.Position(pos).Filename
	fieldLine := k.pkg.Fset.Position(pos).Line
	var fieldComment *ast.CommentGroup
	for _, syntax := range k.pkg.Syntax {
		fileName := k.pkg.Fset.Position(syntax.Pos()).Filename
		if fieldFileName == fileName {
			for _, c := range syntax.Comments {
				if k.pkg.Fset.Position(c.End()).Line+1 == fieldLine {
					fieldComment = c
				}
			}
			break
		}
	}
	return fieldComment

	//if fieldComment == nil {
	//	return ""
	//}
	//
	//for _, c := range fieldComment {
	//	commentField := strings.Fields(c.Text)
	//	if len(commentField) < 3 {
	//		panic("comment error: " + c.Text)
	//	}
	//	fmt.Println("commentField", commentField)
	//	if commentField[0] == "@kit-request" && commentField[1] == "ctx" {
	//		return commentField[2]
	//	}
	//}
	//return ""
}

func (k *KitRequest) CheckRequestIsNil() {
	if k.RequestIsBody {
		if k.RequestTypeOf.NumFields() == 0 {
			k.RequestIsNil = true
		}
	}
}

func (k *KitRequest) RequestType(prefix []string, requestName string, requestType types.Type, requestParamTagTypeTag string, doc *ast.CommentGroup, ptr bool) {
	xt := common.TypeOf(requestType)
	rawParamType := requestType.String()
	paramSource, paramName := k.ParseParamTag(requestName, requestParamTagTypeTag)
	// fmt.Println("requestName:", requestName, "rawParamType: ", rawParamType)
	// fmt.Println(xt.Basic, xt.Named, xt.Pointer)

	switch rt := requestType.(type) {
	case *types.Named:
		//fmt.Println("paramName", paramName)
		//fmt.Println("obj.name",rt.Obj().Name())
		//fmt.Println("obj.pkg",rt.Obj().Pkg().Path())
		//fmt.Println("obj.pos",rt.Obj().Pos())
		//fmt.Println("obj.type",rt.Obj().Type())
		//fmt.Println("obj.type.string",rt.Obj().Type().String())
		//fmt.Println("obj.id",rt.Obj().Id())
		//fmt.Println("local.pkg.pkgPath", k.pkg.PkgPath)
		split := strings.Split(strings.TrimPrefix(rt.Obj().Type().String(), k.pkg.PkgPath+"."), "/")
		named := split[len(split)-1]

		k.NamedMap[paramName] = named
		//k.SetParam(RequestParam{
		//	ParamDoc:     doc,
		//	ParamPath:    strings.Join(prefix, "."),
		//	FieldName:    requestName,
		//	ParamName:    paramName,
		//	ParamSource:  paramSource,
		//	ParamType:    "named",
		//	RawParamType: rawParamType,
		//	BasicType:    rt.Underlying().String(),
		//	HasPtr:       false,
		//})
		k.RequestType(prefix, requestName, rt.Underlying(), requestParamTagTypeTag, doc, ptr)
	case *types.Struct:
		var paramTypeName string
		var hasNamed bool

		if paramTypeName, hasNamed = k.NamedMap[paramName]; !hasNamed {
			paramTypeName = common.TypeOf(rt).TypeAsJenComparePkgNameString(k.pkg)
		}
		k.SetParam(RequestParam{
			FieldName:     requestName,
			ParamDoc:      doc,
			ParamPath:     strings.Join(prefix, "."),
			ParamName:     paramName,
			ParamSource:   paramSource,
			ParamType:     "struct",
			ParamTypeName: paramTypeName,
			RawParamType:  rawParamType,
			BasicType:     k.NamedMap[paramName],
			HasPtr:        ptr,
			HasNamed:      hasNamed,
			XType:         xt,
			XTypeID:       xt.TypeAsJenComparePkgNameString(k.pkg),
		})
		for i := 0; i < rt.NumFields(); i++ {
			field := rt.Field(i)
			fieldName := field.Name()
			fieldType := field.Type()
			// fmt.Println("fieldName:", fieldName, "fieldType:", fieldType)
			tag, _ := reflect.StructTag(rt.Tag(i)).Lookup(RequestParamTagName)
			k.RequestType(append(prefix, fieldName), fieldName, fieldType, tag, k.ParseFieldComment(field.Pos()), false)
		}
	case *types.Pointer:
		// fmt.Println("fieldName: ", requestName, rt.Elem().String(), rt.Elem().Underlying().String())
		k.RequestType(prefix, requestName, rt.Elem(), requestParamTagTypeTag, doc, true)
	case *types.Interface:
		var paramTypeName string
		var hasNamed bool

		if paramTypeName, hasNamed = k.NamedMap[paramName]; !hasNamed {
			paramTypeName = common.TypeOf(rt).TypeAsJenComparePkgNameString(k.pkg)
		}
		k.SetParam(RequestParam{
			FieldName:     requestName,
			ParamDoc:      doc,
			ParamPath:     strings.Join(prefix, "."),
			ParamName:     paramName,
			ParamSource:   paramSource,
			ParamType:     "interface",
			ParamTypeName: paramTypeName,
			RawParamType:  rawParamType,
			BasicType:     rt.String(),
			HasPtr:        ptr,
			HasNamed:      hasNamed,
			XType:         xt,
			XTypeID:       xt.TypeAsJenComparePkgNameString(k.pkg),
		})
	case *types.Slice:
		var paramTypeName string
		var hasNamed bool
		if paramTypeName, hasNamed = k.NamedMap[paramName]; !hasNamed {
			//split := strings.Split(strings.TrimPrefix(rt.Elem().String(), k.pkg.PkgPath+"."), "/")
			//paramTypeName = "[]" + split[len(split)-1]
			paramTypeName = common.TypeOf(rt).TypeAsJenComparePkgNameString(k.pkg)
		}
		k.SetParam(RequestParam{
			FieldName:     requestName,
			ParamDoc:      doc,
			ParamPath:     strings.Join(prefix, "."),
			ParamName:     paramName,
			ParamSource:   paramSource,
			ParamType:     "slice",
			ParamTypeName: paramTypeName,
			RawParamType:  rawParamType,
			BasicType:     rt.Elem().Underlying().String(),
			HasPtr:        ptr,
			HasNamed:      hasNamed,
			XType:         xt,
			XTypeID:       xt.TypeAsJenComparePkgNameString(k.pkg),
		})
	case *types.Map:
		var paramTypeName string
		var hasNamed bool
		if paramTypeName, hasNamed = k.NamedMap[paramName]; !hasNamed {
			split := strings.Split(strings.TrimPrefix(rt.Elem().String(), k.pkg.PkgPath+"."), "/")
			paramTypeName = split[len(split)-1]
		}
		k.SetParam(RequestParam{
			FieldName:     requestName,
			ParamDoc:      doc,
			ParamPath:     strings.Join(prefix, "."),
			ParamName:     paramName,
			ParamSource:   paramSource,
			ParamType:     "map",
			ParamTypeName: paramTypeName,
			RawParamType:  rawParamType,
			BasicType:     rt.Elem().Underlying().String(),
			HasPtr:        ptr,
			HasNamed:      hasNamed,
			XType:         xt,
			XTypeID:       xt.TypeAsJenComparePkgNameString(k.pkg),
		})
	case *types.Basic:
		var paramTypeName string
		var hasNamed bool
		if paramTypeName, hasNamed = k.NamedMap[paramName]; !hasNamed {
			paramTypeName = rt.Name()
		}
		k.SetParam(RequestParam{
			FieldName:     requestName,
			ParamDoc:      doc,
			ParamPath:     strings.Join(prefix, "."),
			ParamName:     paramName,
			ParamSource:   paramSource,
			ParamType:     "basic",
			ParamTypeName: paramTypeName,
			RawParamType:  rawParamType,
			BasicType:     rt.Name(),
			HasPtr:        ptr,
			HasNamed:      hasNamed,
			XType:         xt,
			XTypeID:       xt.TypeAsJenComparePkgNameString(k.pkg),
		})
	default:
		return
	}

	return
}

func downFirst(s string) string {
	for _, v := range s {
		return string(unicode.ToLower(v)) + s[len(string(v)):]
	}
	return ""
}

func upFirst(s string) string {
	for _, v := range s {
		return string(unicode.ToUpper(v)) + s[len(string(v)):]
	}
	return ""
}

func UrlValues(paramName, t, paramTypeName string) (res []jen.Code, err error) {
	key := t + "." + paramTypeName
	switch key {
	case "map.string":
		res = append(res, jen.Id(paramName+"Map").Op(":=").Make(jen.Map(jen.String()).String(), jen.Id("0")))
		res = append(res, jen.For(jen.Id("k,v := range").Id("r.URL.Query()")).Block(
			jen.Id(paramName+"Map").Index(jen.Id("k")).Op("=").Id(`strings.Join(v,",")`)),
		)
		res = append(res, jen.Id(paramName).Op("=").Id(paramName+"Map"))
	default:
		err = fmt.Errorf("not support type %s", paramTypeName)
	}

	return
}

func CastMap(p *common.Type, paramName, t, paramTypeName string, code jen.Code) (res []jen.Code, err error) {
	if t == "slice" && paramTypeName == "[]string" {
		res = append(res, jen.Id(paramName+"Str").Op(":=").Add(code))
		res = append(res, jen.If(jen.Id(paramName+`Str != ""`)).Block(
			jen.Id(paramName).Op("=").Id("strings.Split").Call(jen.Id(paramName+"Str"), jen.Lit(",")),
		))
		return
	}

	if t == "slice" && p.ListInner.Basic {
		res = append(res, jen.Id(paramName+"Str").Op(":=").Add(code))
		res = append(res, jen.If(jen.Id(paramName+`Str != ""`)).Block(
			jen.Id(fmt.Sprintf(`lo.ForEachWhile(strings.Split(%s, ","), func(item string, index int) bool {
			var %sValue %s
			%sValue, err = cast.To%sE(item)
			if err != nil {
				return false
			}

			%s = append(%s, %sValue)
			return true
		})

		if err != nil {
			return
		}`, paramName+"Str", paramName, p.ListInner.BasicType.String(), paramName, upFirst(p.ListInner.BasicType.String()), paramName, paramName, paramName)),
		))

		return
	}
	var m = map[string]string{
		"basic.int":     "cast.ToIntE",
		"basic.int8":    "cast.ToInt8E",
		"basic.int16":   "cast.ToInt16E",
		"basic.int32":   "cast.ToInt32E",
		"basic.int64":   "cast.ToInt64E",
		"basic.uint":    "cast.ToUintE",
		"basic.uint8":   "cast.ToUint8E",
		"basic.uint16":  "cast.ToUint16E",
		"basic.uint32":  "cast.ToUint32E",
		"basic.uint64":  "cast.ToUint64E",
		"basic.float32": "cast.ToFloat32E",
		"basic.float64": "cast.ToFloat64E",
		"basic.string":  "cast.ToStringE",
		"basic.bool":    "cast.ToBoolE",

		"map.int":   "cast.ToStringMapIntE",
		"map.int64": "cast.ToStringMapInt64E",
		"map.bool":  "cast.ToStringMapBoolE",

		"struct.time.Time":    "cast.ToTimeE",
		"basic.time.Duration": "cast.ToDurationE",
	}
	var ok bool
	var hasType bool
	mKey := t + "." + paramTypeName
	fnStr, ok := m[mKey]
	if !ok {
		hasType = true
		//return
		fnStr, ok = m[t+"."+p.BasicType.String()]
		if !ok {
			err = fmt.Errorf("CastMap not found %s %s", t, paramTypeName)
			return
		}
	}

	paramStr := paramName + "Str"
	varParamStr := jen.Id(paramStr).Op(":=").Add(code)
	paramStrCode := jen.Id(paramStr)
	if t == "slice" {
		paramStrCode = jen.Id("strings.Split").Call(paramStrCode, jen.Lit(","))
	}

	switch mKey {
	case "struct.time.Time":
		paramStrCode = jen.Id("cast.ToInt64").Call(paramStrCode)
	case "basic.time.Duration":
		paramStrCode = jen.Id("cast.ToInt64").Call(paramStrCode)
	}

	ifParamStr := jen.If(jen.Id(paramStr).Op("!=").Lit("")).BlockFunc(func(group *jen.Group) {
		if hasType {
			group.List(jen.Id(paramName+"Asser"), jen.Id(paramName+"Err")).Op(":=").Id(fnStr).Call(paramStrCode).Line()
			group.If(jen.Id(paramName+"Err").Op("!=").Nil()).Block(
				jen.Err().Op("=").Id(paramName+"Err"),
				jen.Return(),
			)
			group.Id(paramName).Op("=").Id(paramTypeName).Call(jen.Id(paramName + "Asser"))
		} else {
			group.List(jen.Id(paramName), jen.Err()).Op("=").Id(fnStr).Call(paramStrCode).Line()
			group.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(),
			)
		}
	})
	//(
	//	jen.BlockFunc(func(group *jen.Group) {
	//		group
	//	}),
	//	jen.List(jen.Id(paramName), jen.Err()).Op("=").Id(fnStr).Call(paramStrCode),
	//	// if err != nil {
	//	// 	return err
	//	// }
	//	jen.If(jen.Err().Op("!=").Nil()).Block(
	//		jen.Return(),
	//	),
	//)

	res = append(res, varParamStr, ifParamStr)

	return
}
