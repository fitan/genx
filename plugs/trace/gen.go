package trace

import (
	"github.com/fitan/genx/common"
	"github.com/fitan/jennifer/jen"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/go/packages"
)

func Gen(pkg *packages.Package, methods []common.InterfaceMethod) {
	slog.Info("gen tracing", slog.Any("pkg", pkg), slog.Any("methods", methods))
	j := jen.NewFile(pkg.Name)
	common.JenAddImports(pkg, j)
	j.Add(genTracingStruct())
	j.Add()

	funcList := make([]jen.Code, 0)

	for _, method := range methods {
		funcList = append(funcList, genTracingFunc(pkg.PkgPath, method))
	}

	j.Add(funcList...)
	j.Add(genNewTracing())

	common.WriteGO("tracing.go", j.GoString())
}

func genTracingFunc(tracingPrefix string, method common.InterfaceMethod) jen.Code {
	methodParamCode := make([]jen.Code, 0)
	methodParamCode = append(methodParamCode, jen.Id("ctx").Qual("context", "Context"))
	methodResultCode := make([]jen.Code, 0)

	tracingParamCode := make([]jen.Code, 0)
	jsonParamCode := make([]jen.Code, 0)
	nextMethodParamCode := make([]jen.Code, 0)
	for _, param := range method.Params {
		nextMethodParamCode = append(nextMethodParamCode, jen.Id(param.Name))
	}

	for _, param := range method.ParamsExcludeCtx() {
		methodParamCode = append(methodParamCode, jen.Id(param.Name).Id(param.ID))
		if !param.Basic() {
			jsonParamCode = append(jsonParamCode,
				jen.List(jen.Id(param.Name+"Byte"), jen.Id("_")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Id(param.Name)).Line(),
				jen.Id(param.Name+"Json").Op(":=").Id("string").Call(jen.Id(param.Name+"Byte")).Line(),
			)
			tracingParamCode = append(tracingParamCode, jen.Lit(param.Name), jen.Id(param.Name+"Json"))

		} else {
			tracingParamCode = append(tracingParamCode, jen.Lit(param.Name), jen.Id(param.Name))
		}
	}

	// jsonStatement := jen.Statement(jsonParamCode)
	// tracingParamStatement := jen.List(tracingParamCode...)

	for _, param := range method.Results {
		methodResultCode = append(methodResultCode, jen.Id(param.Name).Id(param.ID))
	}

	return jen.Func().Params(jen.Id("s").Op("*").Id("tracing")).Id(method.Name).Params(
		methodParamCode...,
	).Params(
		methodResultCode...,
	).BlockFunc(func(g *jen.Group) {
		g.List(jen.Id("ctx"), jen.Id("span")).Op(":=").Id("s.tracer.Tracer").Call(jen.Lit(method.Name)).Dot("Start").Call(jen.Id("ctx"), jen.Lit(method.Name))
		g.Defer().Func().Params().BlockFunc(func(g *jen.Group) {

			g.Id("_params").Op(":=").Map(jen.String()).Interface().ValuesFunc(func(g *jen.Group) {
				for _, v := range method.ParamsExcludeCtx() {
					g.Lit(v.Name).Op(":").Add(jen.Id(v.Name))
				}
			})
			g.List(jen.Id("_paramsB"), jen.Id("_")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Id("_params"))

			g.Id("span").Dot("SetAttributes").CallFunc(func(g *jen.Group) {
				g.Id("attribute.String").Call(jen.Lit("params"), jen.Id("string(_paramsB)"))
			})
			if method.ReturnsError {
				g.If(jen.Err().Op("!=").Nil()).Block(
					jen.Id("span.RecordError").Call(jen.Id("err")),
					jen.Id("span").Dot("SetStatus").Call(jen.Qual("go.opentelemetry.io/otel/codes", "Error"), jen.Id("err").Dot("Error").Call()),
				)
			}

			g.Id("span.End()")
		}).Call()
		if method.HasResults() {
			g.Return().Id("s").Dot("next").Dot(method.Name).Call(
				nextMethodParamCode...,
			)
		} else {
			g.Id("s").Dot("next").Dot(method.Name).Call(
				nextMethodParamCode...,
			)
		}
	}).Line()
}

func genTracingStruct() jen.Code {
	return jen.Null().Type().Id("tracing").Struct(jen.Id("next").Id("Service"), jen.Id("tracer").Id("*sdktrace.TracerProvider"))
}

func genNewTracing() jen.Code {
	return jen.Func().Id("NewTracing").Params(jen.Id("i").Id("do.Injector")).Params(jen.Id("Middleware")).Block(jen.Return().Func().Params(jen.Id("next").Id("Service")).Params(jen.Id("Service")).Block(jen.Return().Op("&").Id("tracing").Values(jen.Id("next").Op(":").Id("next"), jen.Id("tracer").Op(":").Id("do.MustInvoke[*sdktrace.TracerProvider](i)"))))
}
