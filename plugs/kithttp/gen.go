package kithttp

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
	"github.com/pkg/errors"
)

func ObserverGen(opt gen.Option, imd common.InterfaceMetaDate) (res []gen.GenResult, err error) {

	logFile, err := fs.ReadFile(opt.Static, "static/template/ce_log.tmpl")
	if err != nil {
		panic(err)
	}

	traceFile, err := fs.ReadFile(opt.Static, "static/template/ce_trace.tmpl")
	if err != nil {
		panic(err)
	}

	methodMap := make(map[string]Method, 0)
	f := func(m common.InterfaceMethod) (*Method, error) {

		method := Method{Name: m.Name}
		method.IMethod = m
		method.Doc = m.Doc
		method.RawDoc = m.RawDoc
		method.AcceptsContext = m.AcceptsContext
		method.ReturnsError = m.ReturnsError
		return &method, nil
	}
	for _, v := range imd.Methods {
		m, err := f(v)
		if err != nil {
			return nil, err
		}

		methodMap[v.Name] = *m
	}

	tInput := &TemplateInputInterface{
		Name:    "Service",
		Methods: methodMap,
		Doc:     *imd.Doc,
		RawDoc:  *imd.RawDoc,
		Opt:     opt,
	}

	logStr, err := common.GenFileByTemplate(opt.Static, "ce_log", tInput)
	if err != nil {
		err = errors.Wrapf(err, "gen file %s", logFile)
		return
	}
	traceStr, err := common.GenFileByTemplate(opt.Static, "ce_trace", tInput)
	if err != nil {
		err = errors.Wrapf(err, "gen file %s", traceFile)
		return
	}

	lJen := jen.NewFile(opt.Pkg.Name)
	tJen := jen.NewFile(opt.Pkg.Name)

	importMap := map[string]string{
		"encoding/json":                             "",
		"net/http":                                  "",
		"strings":                                   "",
		"github.com/asaskevich/govalidator":         "valid",
		"github.com/go-kit/kit/endpoint":            "",
		"github.com/go-kit/kit/transport/http":      "",
		"github.com/gorilla/mux":                    "",
		"github.com/pkg/errors":                     "",
		"github.com/spf13/cast":                     "",
		"github.com/go-playground/validator/v10":    "validator",
		"context":                                   "",
		"github.com/opentracing/opentracing-go/ext": "",
		"github.com/opentracing/opentracing-go":     "",
		"github.com/samber/lo":                      "",
	}

	for k, v := range importMap {
		lJen.AddImport(k, v)
		tJen.AddImport(k, v)
	}

	for _, v := range opt.Config.Imports {
		lJen.AddImport(v.Path, v.Alias)
		tJen.AddImport(v.Path, v.Alias)
	}

	// for k, _ := range opt.Pkg.Imports {
	// eJen.AddImport(k, "")
	// hJen.AddImport(k, "")
	// }

	lJen.Id(logStr)
	tJen.Id(traceStr)

	res = append(res, gen.GenResult{
		FileName: filepath.Join(opt.Dir, "logging.go"),
		FileStr:  lJen.GoString(),
		Cover:    true,
	})

	res = append(res, gen.GenResult{
		FileName: filepath.Join(opt.Dir, "tracing.go"),
		FileStr:  tJen.GoString(),
		Cover:    true,
	})

	return
}

func Gen(opt gen.Option, imd common.InterfaceMetaDate) (res []gen.GenResult, err error) {
	methodMap := make(map[string]Method, 0)
	f := func(m common.InterfaceMethod) (*Method, error) {
		kit, err := NewKit(m.Doc)
		if err != nil {
			panic(err)
		}
		// 如果没有不生成
		if kit.Conf.Url == "" {
			return nil, nil
		}
		kitRequest := NewKitRequest(opt.Pkg, m.Name, kit.Conf.HttpRequestName, kit.Conf.HttpRequestBody)
		kitRequest.ParseRequest()

		method := Method{Name: m.Name}
		method.IMethod = m
		method.RawKit = kit
		method.Doc = m.Doc
		method.RawDoc = m.RawDoc
		method.KitRequest = kitRequest
		method.AcceptsContext = m.AcceptsContext
		method.ReturnsError = m.ReturnsError
		method.KitRequestDecode = kitRequest.DecodeRequest()
		return &method, nil
	}
	for _, v := range imd.Methods {
		m, err := f(v)
		if err != nil {
			panic(err)
		}

		if m != nil {
			methodMap[v.Name] = *m
		}
	}

	tInput := &TemplateInputInterface{
		Name:    "Service",
		Methods: methodMap,
		Doc:     *imd.Doc,
		RawDoc:  *imd.RawDoc,
		Opt:     opt,
	}

	eFile, err := common.GenFileByTemplate(opt.Static, "kit_endpoint", tInput)
	if err != nil {
		return
	}

	hFile, err := common.GenFileByTemplate(opt.Static, "kit_http", tInput)
	if err != nil {
		return
	}

	eJen := jen.NewFile(opt.Pkg.Name)
	hJen := jen.NewFile(opt.Pkg.Name)

	importMap := map[string]string{
		"encoding/json":                             "",
		"net/http":                                  "",
		"strings":                                   "",
		"github.com/asaskevich/govalidator":         "valid",
		"github.com/go-kit/kit/endpoint":            "",
		"github.com/go-kit/kit/transport/http":      "kithttp",
		"github.com/gorilla/mux":                    "",
		"github.com/pkg/errors":                     "errors",
		"github.com/spf13/cast":                     "cast",
		"github.com/go-playground/validator/v10":    "validator",
		"context":                                   "",
		"github.com/opentracing/opentracing-go/ext": "",
		"github.com/opentracing/opentracing-go":     "",
		"github.com/samber/lo":                      "",
	}

	for k, v := range importMap {
		eJen.AddImport(k, v)
		hJen.AddImport(k, v)
	}

	for _, v := range opt.Config.Imports {
		eJen.AddImport(v.Path, v.Alias)
		hJen.AddImport(v.Path, v.Alias)
	}

	// for k, _ := range opt.Pkg.Imports {
	// eJen.AddImport(k, "")
	// hJen.AddImport(k, "")
	// }

	eJen.Id(eFile)
	hJen.Id(hFile)

	res = append(res, gen.GenResult{
		FileName: filepath.Join(opt.Dir, "endpoint.go"),
		FileStr:  eJen.GoString(),
		Cover:    true,
	})

	res = append(res, gen.GenResult{
		FileName: filepath.Join(opt.Dir, "http.go"),
		FileStr:  hJen.GoString(),
		Cover:    true,
	})

	return
}

func genEndpointConst(methodNameList []string) jen.Code {
	j := jen.Null()

	for _, methodName := range methodNameList {
		j.Const().Id(methodName + "MethodName").Op("=").Lit(methodName).Line()
	}

	j.Var().Id("MethodNameList").Op("=").Index().String().ValuesFunc(func(g *jen.Group) {
		for _, methodName := range methodNameList {
			g.Id(methodName + "MethodName")
		}
	})

	return j
}

func genEndpoints(methodNameList []string) jen.Code {
	listCode := make([]jen.Code, 0, len(methodNameList))
	for _, methodName := range methodNameList {
		listCode = append(listCode, jen.Id(methodName+"Endpoint").Qual("github.com/go-kit/kit/endpoint", "Endpoint"))
	}
	return jen.Null().Type().Id("Endpoints").Struct(
		listCode...,
	)
}

func genNewEndpoint(methodNameList []string) jen.Code {
	endpointVarList := make([]jen.Code, 0, len(methodNameList))
	endpointForList := make([]jen.Code, 0, len(methodNameList))

	for _, methodName := range methodNameList {
		endpointVarList = append(endpointVarList, jen.Id(methodName+"Endpoint").Op(":").Id("make"+methodName+"Endpoint").Call(jen.Id("s")))

		endpointForList = append(endpointForList, jen.For(jen.List(jen.Id("_"), jen.Id("m")).Op(":=").Range().Id("dmw").Index(jen.Id(methodName+"MethodName"))).Block(jen.Id("eps").Dot(methodName+"Endpoint").Op("=").Id("m").Call(jen.Id("eps").Dot(methodName+"Endpoint"))).Line())
	}

	endpointForListStatement := jen.Statement(endpointForList)

	return jen.Func().Id("NewEndpoint").Params(jen.Id("s").Id("Service"), jen.Id("dmw").Map(jen.Id("string")).Index().Qual("github.com/go-kit/kit/endpoint", "Middleware")).Params(jen.Id("Endpoints")).Block(
		jen.Id("eps").Op(":=").Id("Endpoints").Values(
			endpointVarList...,
		),
		&endpointForListStatement,
		jen.Return().Id("eps"),
	)
}

var helperFuncs = template.FuncMap{
	"up":        strings.ToUpper,
	"down":      strings.ToLower,
	"upFirst":   upFirst,
	"downFirst": downFirst,
	"replace":   strings.ReplaceAll,
	"snake":     toSnakeCase,
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	result := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	result = matchAllCap.ReplaceAllString(result, "${1}_${2}")
	return strings.ToLower(result)
}
