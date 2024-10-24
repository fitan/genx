package kithttp

import (
	"bytes"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
	"github.com/pkg/errors"
	"github.com/sourcegraph/conc"
)

func ObserverGen(opt gen.Option, imd common.InterfaceMetaDate) error {

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
			panic(err)
		}

		methodMap[v.Name] = *m
	}

	tInput := &TemplateInputInterface{
		Name:    "Service",
		Methods: methodMap,
		Doc:     *imd.Doc,
		Opt:     opt,
	}

	logStr, err := common.GenFileByTemplate(opt.Static, "ce_log", tInput)
	if err != nil {
		err = errors.Wrapf(err, "gen file %s", logFile)
		return err
	}
	traceStr, err := common.GenFileByTemplate(opt.Static, "ce_trace", tInput)
	if err != nil {
		err = errors.Wrapf(err, "gen file %s", traceFile)
		return err
	}

	lJen := jen.NewFile(opt.Pkg.Name)
	tJen := jen.NewFile(opt.Pkg.Name)

	importMap := map[string]string{
		"encoding/json":                             "json",
		"net/http":                                  "http",
		"strings":                                   "strings",
		"github.com/asaskevich/govalidator":         "valid",
		"github.com/go-kit/kit/endpoint":            "endpoint",
		"github.com/go-kit/kit/transport/http":      "kithttp",
		"github.com/gorilla/mux":                    "mux",
		"github.com/pkg/errors":                     "errors",
		"github.com/spf13/cast":                     "cast",
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

	var wg conc.WaitGroup

	wg.Go(func() {
		common.WriteGO(filepath.Join(opt.Dir, "logging.go"), lJen.GoString())
	})

	wg.Go(func() {
		common.WriteGO(filepath.Join(opt.Dir, "tracing.go"), tJen.GoString())
	})

	wg.Wait()

	return nil
}

func Gen(opt gen.Option, imd common.InterfaceMetaDate) {
	endpointFile, err := fs.ReadFile(opt.Static, "static/template/kit_endpoint.tmpl")
	if err != nil {
		panic(err)
	}

	httpFile, err := fs.ReadFile(opt.Static, "static/template/kit_http.tmpl")
	if err != nil {
		panic(err)
	}

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
		kit, err := NewKit(m.Doc)
		if err != nil {
			panic(err)
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

		methodMap[v.Name] = *m
	}

	tInput := &TemplateInputInterface{
		Name:    "Service",
		Methods: methodMap,
		Doc:     *imd.Doc,
		Opt:     opt,
	}

	et, err := template.New("kit_endpoint").Funcs(helperFuncs).Parse(string(endpointFile))
	if err != nil {
		panic(err)
	}

	ht, err := template.New("kit_http").Funcs(helperFuncs).Parse(string(httpFile))
	if err != nil {
		panic(err)
	}

	lt, err := template.New("ce_log").Funcs(helperFuncs).Parse(string(logFile))
	if err != nil {
		panic(err)
	}

	tt, err := template.New("ce_trace").Funcs(helperFuncs).Parse(string(traceFile))
	if err != nil {
		panic(err)
	}

	var eBuffer bytes.Buffer
	var hBuffer bytes.Buffer
	var lBuffer bytes.Buffer
	var tBuffer bytes.Buffer
	err = et.Execute(&eBuffer, tInput)
	if err != nil {
		panic(err)
	}

	err = ht.Execute(&hBuffer, tInput)
	if err != nil {
		panic(err)
	}

	err = lt.Execute(&lBuffer, tInput)

	if err != nil {
		panic(err)
	}

	err = tt.Execute(&tBuffer, tInput)
	if err != nil {
		panic(err)
	}

	eJen := jen.NewFile(opt.Pkg.Name)
	hJen := jen.NewFile(opt.Pkg.Name)
	lJen := jen.NewFile(opt.Pkg.Name)
	tJen := jen.NewFile(opt.Pkg.Name)

	importMap := map[string]string{
		"encoding/json":                             "json",
		"net/http":                                  "http",
		"strings":                                   "strings",
		"github.com/asaskevich/govalidator":         "valid",
		"github.com/go-kit/kit/endpoint":            "endpoint",
		"github.com/go-kit/kit/transport/http":      "kithttp",
		"github.com/gorilla/mux":                    "mux",
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
		lJen.AddImport(k, v)
		tJen.AddImport(k, v)
	}

	for _, v := range opt.Config.Imports {
		eJen.AddImport(v.Path, v.Alias)
		hJen.AddImport(v.Path, v.Alias)
		lJen.AddImport(v.Path, v.Alias)
		tJen.AddImport(v.Path, v.Alias)
	}

	// for k, _ := range opt.Pkg.Imports {
	// eJen.AddImport(k, "")
	// hJen.AddImport(k, "")
	// }

	eJen.Id(eBuffer.String())
	hJen.Id(hBuffer.String())
	lJen.Id(lBuffer.String())
	tJen.Id(tBuffer.String())

	var wg conc.WaitGroup

	wg.Go(func() {
		common.WriteGO(filepath.Join(opt.Dir, "endpoint.go"), eJen.GoString())
	})

	wg.Go(func() {
		common.WriteGO(filepath.Join(opt.Dir, "http.go"), hJen.GoString())
	})

	wg.Go(func() {
		common.WriteGO(filepath.Join(opt.Dir, "logging.go"), lJen.GoString())
	})

	wg.Go(func() {
		common.WriteGO(filepath.Join(opt.Dir, "tracing.go"), tJen.GoString())
	})

	wg.Wait()
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
