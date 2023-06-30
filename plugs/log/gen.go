package log

import (
	"github.com/fitan/genx/common"
	"github.com/fitan/jennifer/jen"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/go/packages"
)

func Gen(pkg *packages.Package, methods []common.InterfaceMethod) {
	slog.Info("gen logging", slog.Any("pkg", pkg), slog.Any("methods", methods))
	j := jen.NewFile(pkg.Name)
	common.JenAddImports(pkg, j)
	j.Add(genLoggingStruct())
	j.Add()

	funcList := make([]jen.Code, 0)

	for _, method := range methods {
		funcList = append(funcList, genLoggingFunc(method))
	}

	j.Add(funcList...)
	j.Add(genNewLogging(pkg.PkgPath))
	common.WriteGO("log.go", j.GoString())

}

func genLoggingFunc(method common.InterfaceMethod) jen.Code {

	methodParamCode := make([]jen.Code, 0)
	methodResultCode := make([]jen.Code, 0)

	logParamCode := make([]jen.Code, 0)
	jsonParamCode := make([]jen.Code, 0)
	nextMethodParamCode := make([]jen.Code, 0)

	for _, param := range method.Params {
		if param.ID == "context.Context" {
			methodParamCode = append(methodParamCode, jen.Id(param.Name).Qual("context", "Context"))
		} else {
			methodParamCode = append(methodParamCode, jen.Id(param.Name).Id(param.ID))
		}

		nextMethodParamCode = append(nextMethodParamCode, jen.Id(param.Name))
	}

	for _, param := range method.ParamsExcludeCtx() {
		if !param.Basic() {
			jsonParamCode = append(jsonParamCode,
				jen.List(jen.Id(param.Name+"Byte"), jen.Id("_")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Id(param.Name)).Line(),
				jen.Id(param.Name+"Json").Op(":=").Id("string").Call(jen.Id(param.Name+"Byte")).Line(),
			)

			logParamCode = append(logParamCode, jen.Lit(param.Name), jen.Id(param.Name+"Json"))
		} else {
			logParamCode = append(logParamCode, jen.Lit(param.Name), jen.Id(param.Name))
		}
	}

	logParamStatement := jen.List(logParamCode...)
	jsonStatement := jen.Statement(jsonParamCode)

	for _, param := range method.Results {
		methodResultCode = append(methodResultCode, jen.Id(param.Name).Id(param.ID))
	}

	return jen.Func().Params(
		jen.Id("s").Op("*").Id("logging")).Id(method.Name).Params(
		methodParamCode...,
	).Params(
		methodResultCode...,
	).Block(
		&jsonStatement,
		jen.Defer().Func().Params(
			jen.Id("begin").Qual("time", "Time")).
			Block(
				jen.Id("_").Op("=").Id("s").Dot("logger").Dot("Log").Call(
					jen.Id("s").Dot("traceId"),
					jen.Id("ctx").Dot("Value").Call(jen.Id("s").Dot("traceId")),
					jen.Lit("method"), jen.Lit(method.Name),
					logParamStatement,
					jen.Lit("took"), jen.Qual("time", "Since").Call(jen.Id("begin")),
					func() jen.Code {
						if method.ReturnsError {
							return jen.Lit("err").Op(",").Id("err")
						}
						return nil
					}(),
				)).Call(jen.Qual("time", "Now").Call()),
		func() jen.Code {
			if method.HasResults() {
				return jen.Return().Id("s").Dot("next").Dot(method.Name).Call(
					nextMethodParamCode...,
				)
			} else {
				return jen.Id("s").Dot("next").Dot(method.Name).Call(
					nextMethodParamCode...,
				)
			}
		}()).Line()
}

func genLoggingStruct() jen.Code {
	return jen.Null().Type().Id("logging").Struct(jen.Id("logger").Qual("github.com/go-kit/kit/log", "Logger"), jen.Id("next").Id("Service"), jen.Id("traceId").Id("string"))
}

func genNewLogging(logPrefix string) jen.Code {
	return jen.Func().Id("NewLogging").Params(jen.Id("logger").Id("log").Dot("Logger"), jen.Id("traceId").Id("string")).Params(jen.Id("Middleware")).Block(
		jen.Id("logger").Op("=").Id("log").Dot("With").Call(
			jen.Id("logger"),
			jen.Lit(logPrefix),
			jen.Lit("logging"),
		), jen.Return().Func().Params(jen.Id("next").Id("Service")).Params(jen.Id("Service")).Block(
			jen.Return().Op("&").Id("logging").Values(jen.Id("logger").Op(":").Qual("github.com/go-kit/kit/log/level", "Info").Call(
				jen.Id("logger")),
				jen.Id("next").Op(":").Id("next"),
				jen.Id("traceId").Op(":").Id("traceId")),
		),
	)
}
