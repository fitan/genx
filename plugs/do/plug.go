package do

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"

	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
	"github.com/samber/lo"
)

func New() *Plug {
	return &Plug{
		names: make(map[string]int, 0),
	}
}

type Plug struct {
	names map[string]int
}

func (p *Plug) getPkgName(name string) string {
	s, ok := p.names[name]
	if ok {
		p.names[name] = s + 1
		return fmt.Sprintf("%s%d", name, p.names[name])
	} else {
		p.names[name] = 0
		return name
	}

}

func (p *Plug) Name() string {
	return "@do"
}

func (p *Plug) Gen(req []gen.GlobalFuncGoTypeMeta) (res []gen.GenResult, err error) {
	initMetas := []gen.GlobalFuncGoTypeMeta{}
	provideMetas := []provideMeta{}
	lo.ForEach(req, func(item gen.GlobalFuncGoTypeMeta, index int) {
		initMeta := gen.GlobalFuncGoTypeMeta{
			Option: item.Option,
			Metas:  []gen.FuncGoTypeMeta{},
		}
		provideMeta := provideMeta{
			Option: item.Option,
			Metas:  []provideMetaItem{},
		}
		var initMetaOk, provideMetaOk bool
		lo.ForEach(item.Metas, func(meta gen.FuncGoTypeMeta, index int) {
			var doType, doName string
			meta.Doc.ByFuncNameAndArgs(p.Name(), &doType, &doName)

			switch doType {
			case "init":
				initMeta.Metas = append(initMeta.Metas, meta)
				initMetaOk = true
			case "provide":
				provideMeta.Metas = append(provideMeta.Metas, provideMetaItem{
					DoName: doName,
					Meta:   meta,
				})

				provideMetaOk = true
			default:
				// 记录不支持的 do 函数类型，但不中断处理
				slog.Error("unsupported do function type",
					"function", meta.Name,
					"do_type", doType,
					"supported_types", "init|provide")
				// 跳过不支持的类型，继续处理其他函数
			}
		})
		if initMetaOk {
			initMetas = append(initMetas, initMeta)
		}
		if provideMetaOk {
			provideMetas = append(provideMetas, provideMeta)
		}

	})

	sort.Slice(provideMetas, func(i, j int) bool {
		return provideMetas[i].Option.Pkg.PkgPath > provideMetas[j].Option.Pkg.PkgPath
	})

	lo.ForEach(initMetas, func(item gen.GlobalFuncGoTypeMeta, index int) {
		f := jen.NewFile(item.Option.Pkg.Name)
		f.AddImport("github.com/samber/do/v2", "")

		f.Func().Id("doInit").Params(jen.Id("i").Id("do.Injector")).BlockFunc(func(g *jen.Group) {
			lo.ForEach(provideMetas, func(v provideMeta, index int) {
				pkgPath := v.Option.Pkg.PkgPath
				pkgName := v.Option.Pkg.Name
				pkgName = p.getPkgName(pkgName)
				f.AddImport(pkgPath, pkgName)
				sort.Slice(v.Metas, func(i, j int) bool {
					return v.Metas[i].Meta.Name > v.Metas[j].Meta.Name
				})
				lo.ForEach(v.Metas, func(vv provideMetaItem, index int) {
					if vv.DoName != "" {
						g.Id("do.ProvideNamed").Call(jen.Id("i"), jen.Lit(vv.DoName), jen.Id(pkgName+"."+vv.Meta.Name))
					} else {
						g.Id("do.Provide").Call(jen.Id("i"), jen.Id(pkgName+"."+vv.Meta.Name))
					}
				})
			})
		}).Line()

		res = append(res, gen.GenResult{
			FileName: filepath.Join(item.Option.Dir, "do_init.go"),
			FileStr:  f.GoString(),
			Cover:    true,
		})
	})

	return
}

type provideMeta struct {
	Option gen.Option
	Metas  []provideMetaItem
}

type provideMetaItem struct {
	DoName string
	Meta   gen.FuncGoTypeMeta
}
