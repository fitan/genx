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
	j.AddImport("github.com/samber/do/v2", "do")
	common.JenAddImports(pkg, j)
	j.Add(genLoggingStruct())
	j.Add()

	funcList := make([]jen.Code, 0)

	for _, method := range methods {
		funcList = append(funcList, genLoggingFunc(method))
	}

	j.Add(funcList...)
	j.Add(genNewLogging(pkg.PkgPath))
	common.WriteGO("logging.go", j.GoString())

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
				jen.Id(param.Name+"Json").Op(":=").Id("string").Call(jen.Id(param.Name+"Byte")),
			)

			logParamCode = append(logParamCode, jen.Lit(param.Name), jen.Id(param.Name+"Json"))
		} else {
			logParamCode = append(logParamCode, jen.Lit(param.Name), jen.Id(param.Name))
		}
	}

	// logParamStatement := jen.List(logParamCode...)
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
			BlockFunc(func(g *jen.Group) {
				g.Var().Id("level").Op("=").Id("slog.LevelInfo")

				if method.ReturnsError {
					g.If(jen.Id("err != nil")).Block(
						jen.Id("level").Op("=").Id("slog.LevelError"),
					)
				}

				g.Id("s").Dot("logger").Dot("Log").CallFunc(func(g *jen.Group) {
					g.Id("ctx")
					g.Id("level")
					g.Lit("")
					g.Lit("method")
					g.Lit(method.Name)
					for _, v := range logParamCode {
						g.Add(v)
					}
					g.Lit("took")
					g.Qual("time", "Since").Call(jen.Id("begin")).Dot("String").Call()
					if method.ReturnsError {
						g.Lit("err").Op(",").Id("err")
					} else {
						g.Add(jen.Return())
					}
				})
			}).Call(jen.Qual("time", "Now").Call()),
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
	return jen.Null().Type().Id("logging").Struct(
		jen.Id("logger").Op("*").Qual("log/slog", "Logger"),
		jen.Id("next").Id("Service"),
	)
}

func genNewLogging(logPrefix string) jen.Code {
	return jen.Func().Id("NewLogging").Params(jen.Id("i").Id("do").Dot("Injector")).Params(jen.Id("Middleware")).Block(
		jen.Return().Func().Params(jen.Id("next").Id("Service")).Params(jen.Id("Service")).Block(
			jen.Return().Op("&").Id("logging").Values(
				jen.Id("logger").Op(":").Id(`do.MustInvoke[*slog.Logger](i)`),
				jen.Id("next").Op(":").Id("next")),
		),
	)
}
