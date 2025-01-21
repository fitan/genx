package alert

import (
	"github.com/fitan/genx/common"
	"github.com/fitan/jennifer/jen"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/go/packages"
)

func Gen(pkg *packages.Package, methods []common.InterfaceMethod) string {
	slog.Info("gen alert", slog.Any("pkg", pkg), slog.Any("methods", methods))
	j := jen.NewFile(pkg.Name)
	common.JenAddImports(pkg, j)
	j.AddImport("github.com/spf13/cast", "")
	j.Add(genAlertStructAndNew())
	j.Add()

	funcList := make([]jen.Code, 0)

	for _, method := range methods {
		funcList = append(funcList, genAlertFunc(pkg.Name, method))
	}

	j.Add(funcList...)
	return j.GoString()

}

func genAlertFunc(pkgName string, method common.InterfaceMethod) jen.Code {

	code := jen.Null()

	methodParamCode := make([]jen.Code, 0)
	methodResultCode := make([]jen.Code, 0)
	nextMethodParamCode := make([]jen.Code, 0)

	for _, param := range method.Params {
		if param.ID == "context.Context" {
			methodParamCode = append(methodParamCode, jen.Id(param.Name).Qual("context", "Context"))
		} else {
			methodParamCode = append(methodParamCode, jen.Id(param.Name).Id(param.ID))
		}

		nextMethodParamCode = append(nextMethodParamCode, jen.Id(param.Name))
	}
	for _, param := range method.Results {
		methodResultCode = append(methodResultCode, jen.Id(param.Name).Id(param.ID))
	}

	code.Func().Params(
		jen.Id("s").Op("*").Id("alert")).Id(method.Name).Params(
		methodParamCode...,
	).Params(
		methodResultCode...,
	).BlockFunc(func(fg *jen.Group) {
		if method.Doc.ByFuncName("@alert-enable") != nil {

			fg.Defer().Func().Params().BlockFunc(func(g *jen.Group) {
				g.Id("if err == nil {return}")

				g.Id(`_traceId,_ := ctx.Value("traceId").(string)`)
				g.Id("_title").Op(":=").Lit(pkgName + "." + method.Name)

				var level string
				var metrics string
				method.Doc.ByFuncNameAndArgs("@alert-level", &level)
				method.Doc.ByFuncNameAndArgs("@alert-metrics", &metrics)

				if level == "" {
					g.Id("_level").Op(":=").Id("s.level")
				} else {
					switch level {
					case "info":
						g.Id("_level").Op(":=").Id("alarm.LevelInfo")
					case "warn":
						g.Id("_level").Op(":=").Id("alarm.LevelWarning")
					case "error":
						g.Id("_level").Op(":=").Id("alarm.LevelError")
					}
				}

				g.Id("_err").Op(":=").Id("s").Dot("api").Dot("Alarm").Call().Dot("Push").Call(
					jen.Id("ctx"),
					jen.Id("_title"),
					jen.Id("_traceId").Op("+").Id("err").Dot("Error").Call(),
					jen.Id("strings.Join").Call(jen.Id("[]string").ValuesFunc(func(g *jen.Group) {
						g.Lit(pkgName + "." + method.Name)
						if lo.IsNotEmpty(metrics) {
							g.Id(metrics)
						}
					}), jen.Lit(".")),
					jen.Id("_level"),
					jen.Id("s").Dot("silencePeriod"),
				)

				g.If(jen.Id("_err").Op("!=").Nil()).Block(
					jen.Id(`level.Error(s.logger).Log("alarm.push.error",_err.Error())`),
				)
			}).Call()

		}

		if method.HasResults() {
			fg.Return().Id("s").Dot("next").Dot(method.Name).Call(
				nextMethodParamCode...,
			)
		} else {
			fg.Id("s").Dot("next").Dot(method.Name).Call(
				nextMethodParamCode...,
			)
		}
	})

	return code.Line()

}

func genAlertStructAndNew() jen.Code {
	return jen.Null().Id(`
type alert struct {
	level         alarm.Level
	silencePeriod int
	api           api.Service
	next          Service
	logger        log.Logger
}

func NewAlert(level alarm.Level, silencePeriod int, api api.Service, log log.Logger) Middleware {
	return func(next Service) Service {
		return &alert{
			level:         level,
			silencePeriod: silencePeriod,
			api:           api,
			next:          next,
			logger:        log,
		}
	}
}
`)
}
