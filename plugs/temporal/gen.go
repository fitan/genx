package temporal

import (
	"log/slog"

	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
)

func Gen(j *jen.File, option gen.Option, implGoTypeMetes []gen.InterfaceGoTypeMeta) {
	j.AddImport("github.com/samber/do/v2", "")
	j.Type().Id("Temporal").Struct(
		jen.Id("next").Id("Service"),
		jen.Id("w").Qual("go.temporal.io/sdk/worker", "Worker"),
	)

	/* j.Line().Func().Id("NewTemporal").Params(jen.Id("i").Id("do.Injector")).Block(
		jen.Id("do").Dot("Provide").Call(jen.Id("i"), jen.Id("initTemporal")),
	) */

	parseImpl := common.NewInterfaceSerialize(option.Pkg)

	j.Line().Func().Id("initTemporal").Params(jen.Id("i").Id("do.Injector")).Params(
		jen.Id("t").Op("*").Id("Temporal"), jen.Id("err").Id("error"),
	).BlockFunc(func(g *jen.Group) {
		g.Id(`w := do.MustInvoke[worker.Worker](i)`).Line()
		g.Id(`next := do.MustInvoke[Service](i)`).Line()

		g.Id("t = &Temporal{next: next, w: w}").Line()

		for _, v := range implGoTypeMetes {
			meta, err := parseImpl.Parse(v.Obj, &v.Doc)
			if err != nil {
				slog.Error("parseImpl.Parse", err, slog.String("name", v.Obj.String()))
				return
			}

			for _, m := range meta.Methods {
				if m.Doc.ByFuncName("@temporal-activity") != nil {
					g.Id(`w.RegisterActivity(next.` + m.Name + ")").Line()
				}
			}
		}

		g.Return(jen.Id("t"), jen.Id("nil"))
	})

	for _, v := range implGoTypeMetes {
		meta, err := parseImpl.Parse(v.Obj, &v.Doc)
		if err != nil {
			slog.Error("parseImpl.Parse", err, slog.String("name", v.Obj.String()))
			return
		}

		for _, m := range meta.Methods {
			if m.Doc.ByFuncName("@temporal-activity") == nil {
				continue
			}
			j.Line().Func().Params(
				jen.Id("t").Op("*").Id("Temporal"),
			).Id(m.Name).ParamsFunc(func(g *jen.Group) {
				for _, p := range m.Params {
					if p.ID == "context.Context" {
						g.Add(jen.Id(p.Name).Qual("go.temporal.io/sdk/workflow", "Context"))
					} else {
						g.Add(jen.Id(p.Name).Id(p.ID))
					}
				}
			}).ParamsFunc(func(g *jen.Group) {
				for _, r := range m.Results {
					g.Add(jen.Id(r.Name).Id(r.ID))
				}
			}).Block(
				jen.Id("err").Op("=").Id("workflow.ExecuteActivity").CallFunc(func(g *jen.Group) {
					g.Id("ctx")
					g.Id("t.next." + m.Name)
					for i, p := range m.Params {
						if i != 0 {
							g.Id(p.Name)
						}
					}
				}).Dot("Get").Call(jen.Id("ctx"), jen.Id("&"+m.Results[0].Name)),
				jen.Return(),
			).Line()
		}
	}
}
