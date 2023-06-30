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

	common.WriteGO("trace.go", j.GoString())
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

	jsonStatement := jen.Statement(jsonParamCode)
	tracingParamStatement := jen.List(tracingParamCode...)

	for _, param := range method.Results {
		methodResultCode = append(methodResultCode, jen.Id(param.Name).Id(param.ID))
	}

	return jen.Func().Params(jen.Id("s").Op("*").Id("tracing")).Id(method.Name).Params(
		methodParamCode...,
	).Params(
		methodResultCode...,
	).Block(
		jen.List(jen.Id("span"), jen.Id("ctx")).Op(":=").Id("opentracing").Dot("StartSpanFromContextWithTracer").Call(
			jen.Id("ctx"),
			jen.Id("s").Dot("tracer"),
			jen.Lit(method.Name),
			jen.Id("opentracing").Dot("Tag").Values(jen.Id("Key").Op(":").Id("string").Call(
				jen.Qual("github.com/opentracing/opentracing-go/ext", "Component")),
				jen.Id("Value").Op(":").Lit(tracingPrefix),
			),
		),
		jen.Defer().Func().Params().Block(
			&jsonStatement,
			jen.Id("span").Dot("LogKV").Call(
				tracingParamStatement,
				func() jen.Code {
					if method.ReturnsError {
						return jen.Lit("err").Op(",").Id("err")
					}
					return nil
				}(),
			),
			func() jen.Code {
				if method.ReturnsError {
					return jen.Id("span").Dot("SetTag").Call(
						jen.Id("string").Call(
							jen.Id("ext").Dot("Error"),
						),
						jen.Id("err != nil"),
					)
				}
				return nil
			}(),
			jen.Id("span").Dot("Finish").Call(),
		).Call(),
		func() jen.Code {
			if method.HasResults() {
				return jen.Return().Id("s").Dot("next").Dot(method.Name).Call(
					nextMethodParamCode...,
				)
			}
			return jen.Id("s").Dot("next").Dot(method.Name).Call(
				nextMethodParamCode...,
			)
		}(),
	).Line()
}

func genTracingStruct() jen.Code {
	return jen.Null().Type().Id("tracing").Struct(jen.Id("next").Id("Service"), jen.Id("tracer").Qual("github.com/opentracing/opentracing-go", "Tracer"))
}

func genNewTracing() jen.Code {
	return jen.Func().Id("NewTracing").Params(jen.Id("otTracer").Id("opentracing").Dot("Tracer")).Params(jen.Id("Middleware")).Block(jen.Return().Func().Params(jen.Id("next").Id("Service")).Params(jen.Id("Service")).Block(jen.Return().Op("&").Id("tracing").Values(jen.Id("next").Op(":").Id("next"), jen.Id("tracer").Op(":").Id("otTracer"))))
}
