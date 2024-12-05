package otel

import (
	"github.com/fitan/genx/common"
	"github.com/fitan/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

func Gen(pkg *packages.Package, methods []common.InterfaceMethod) string {
	j := jen.NewFile(pkg.Name)
	common.JenAddImports(pkg, j)
	j.AddImport("context", "")
	j.AddImport("encoding/json", "")
	j.AddImport("log/slog", "")
	j.AddImport("time", "")
	j.AddImport("github.com/samber/do/v2", "")
	j.AddImport("go.opentelemetry.io/otel/attribute", "")
	j.AddImport("go.opentelemetry.io/otel/codes", "")
	j.AddImport("go.opentelemetry.io/otel/metric", "")
	j.AddImport("go.opentelemetry.io/otel/trace", "")
	j.Add(genOtelStruct())
	j.Add()

	funcList := make([]jen.Code, 0)

	for _, method := range methods {
		funcList = append(funcList, genOtelFunc(method))
	}

	j.Add(funcList...)
	j.Add(genNewOtel(pkg.Name))

	return j.GoString()
}

func genOtelFunc(method common.InterfaceMethod) jen.Code {
	methodParamCode := make([]jen.Code, 0)
	methodParamCode = append(methodParamCode, jen.Id("ctx").Qual("context", "Context"))
	methodResultCode := make([]jen.Code, 0)

	nextMethodParamCode := make([]jen.Code, 0)
	for _, param := range method.Params {
		nextMethodParamCode = append(nextMethodParamCode, jen.Id(param.Name))
	}

	for _, param := range method.ParamsExcludeCtx() {
		methodParamCode = append(methodParamCode, jen.Id(param.Name).Id(param.ID))
	}

	for _, param := range method.Results {
		methodResultCode = append(methodResultCode, jen.Id(param.Name).Id(param.ID))
	}

	return jen.Func().Params(jen.Id("s").Op("*").Id("otel")).Id(method.Name).Params(
		methodParamCode...,
	).Params(
		methodResultCode...,
	).BlockFunc(func(g *jen.Group) {
		g.Id("_method").Op(":=").Lit(method.Name)
		g.List(jen.Id("ctx"), jen.Id("span")).Op(":=").Id("s.tracer.Start").Call(jen.Id("ctx"), jen.Id(`s.pkgName + "." + _method`))

		g.Defer().Func().Params(jen.Id("begin time.Time")).BlockFunc(func(g *jen.Group) {
			g.Id("_endTime").Op(":=").Id(`time.Since(begin)`)
			g.Var().Id("_level").Op("=").Id("slog.LevelInfo")

			g.If(jen.Err().Op("!=").Nil()).Block(
				jen.Id("_level").Op("=").Id("slog.LevelError"),
				jen.Id(`s.meterServiceStatus.Add(ctx, 1,
				metric.WithAttributes(
					attribute.String("service", s.pkgName),
					attribute.String("method", _method),
					attribute.String("status", "fail"),
				))`),
			).Id("else").Block(
				jen.Id(`s.meterServiceStatus.Add(ctx, 1,
				metric.WithAttributes(
					attribute.String("service", s.pkgName),
					attribute.String("method", _method),
					attribute.String("status", "success"),
				))`),
			)

			g.Id(`s.meterServiceDuration.Record(ctx, float64(_endTime.Milliseconds()),
			metric.WithAttributes(
				attribute.String("service", s.pkgName),
				attribute.String("method", _method),
			))`)

			g.Id(`s.meterServiceCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("service", s.pkgName),
				attribute.String("method", _method),
			))`)

			g.Id("_params").Op(":=").Map(jen.String()).Interface().ValuesFunc(func(g *jen.Group) {
				for _, v := range method.ParamsExcludeCtx() {
					g.Lit(v.Name).Op(":").Add(jen.Id(v.Name))
				}
			})
			g.List(jen.Id("_paramsB"), jen.Id("_")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Id("_params"))

			g.Id("s").Dot("logger").Dot("Log").CallFunc(func(g *jen.Group) {
				g.Id("ctx")
				g.Id("_level")
				g.Lit("")
				g.Lit("pkg")
				g.Id("s.pkgName")
				g.Lit("method")
				g.Id("_method")
				g.Lit("params")
				g.Id("string(_paramsB)")
				g.Lit("took")
				g.Id("_endTime.String()")
				if method.ReturnsError {
					g.Lit("err").Op(",").Id("err")
				}
			})

			g.Id("span").Dot("SetAttributes").CallFunc(func(g *jen.Group) {
				g.Id("attribute.String").Call(jen.Lit("params"), jen.Id("string(_paramsB)"))
				g.Id("attribute.String").Call(jen.Lit("took"), jen.Id("_endTime.String()"))
			})
			if method.ReturnsError {
				g.If(jen.Err().Op("!=").Nil()).Block(
					jen.Id("span.RecordError").Call(jen.Id("err")),
					jen.Id("span").Dot("SetStatus").Call(jen.Qual("go.opentelemetry.io/otel/codes", "Error"), jen.Id("err").Dot("Error").Call()),
				)
			}

			g.Id("span.End()")

		}).Call(jen.Id("time.Now()"))
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

func genOtelStruct() jen.Code {
	return jen.Null().Type().Id("otel").Struct(
		jen.Id("pkgName").String(),
		jen.Id("next").Id("Service"),
		jen.Id("tracer").Id("trace.Tracer"),
		jen.Id("logger").Id("*slog.Logger"),
		jen.Id("meter").Id("metric.Meter"),
		jen.Id("meterServiceStatus").Id("metric.Int64Counter"),
		jen.Id("meterServiceCounter").Id("metric.Int64Counter"),
		jen.Id("meterServiceDuration").Id("metric.Float64Histogram"),
	)
}

func genNewOtel(pkgName string) jen.Code {
	return jen.Func().Id("NewOtel").Params(jen.Id("i").Id("do.Injector")).Params(jen.Id("Middleware")).BlockFunc(
		func(g *jen.Group) {
			g.Return().Func().Params(jen.Id("next").Id("Service")).Params(jen.Id("Service")).Block(
				jen.Return().Op("&").Id("otel").Values(
					jen.Id("pkgName").Op(":").Lit(pkgName),
					jen.Id("next").Op(":").Id("next"),
					jen.Id("tracer").Op(":").Id("do.MustInvoke[trace.Tracer](i)"),
					jen.Id("logger").Op(":").Id(`do.MustInvoke[*slog.Logger](i)`),
					jen.Id("meterServiceStatus").Op(":").Id(`do.MustInvokeNamed[metric.Int64Counter](i, "serviceStatus")`),
					jen.Id("meterServiceCounter").Op(":").Id(`do.MustInvokeNamed[metric.Int64Counter](i, "serviceCounter")`),
					jen.Id("meterServiceDuration").Op(":").Id(`do.MustInvokeNamed[metric.Float64Histogram](i, "serviceDuration")`),
				),
			)
		},
	)
}
