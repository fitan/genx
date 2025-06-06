package fx

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
	"strings"

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
	return "@fx"
}

// FxStructPlug 处理 fx.In 和 fx.Out 结构体
type FxStructPlug struct {
	names map[string]int
}

func NewFxStructPlug() *FxStructPlug {
	return &FxStructPlug{
		names: make(map[string]int, 0),
	}
}

func (p *FxStructPlug) Name() string {
	return "@fx-struct"
}

func (p *FxStructPlug) getPkgName(name string) string {
	s, ok := p.names[name]
	if ok {
		p.names[name] = s + 1
		return fmt.Sprintf("%s%d", name, p.names[name])
	} else {
		p.names[name] = 0
		return name
	}
}

func (p *FxStructPlug) Gen(option gen.Option, structGoTypeMetas []gen.StructGoTypeMeta) (res []gen.GenResult, err error) {
	f := jen.NewFile(option.Pkg.Name)
	f.AddImport("go.uber.org/fx", "")

	hasGenerated := false

	for _, v := range structGoTypeMetas {
		// 检查是否有 @fx-struct 注解
		fxDoc := v.Doc.ByFuncName(p.Name())
		if fxDoc == nil {
			continue
		}

		// 解析参数
		var structType string
		if len(fxDoc.Args) > 0 {
			structType = strings.Trim(fxDoc.Args[0].Value, `"`)
		}

		slog.Info("processing fx struct",
			"struct", v.Name,
			"type", structType)

		switch structType {
		case "in":
			p.generateFxInStruct(f, v)
			hasGenerated = true
		case "out":
			p.generateFxOutStruct(f, v)
			hasGenerated = true
		default:
			slog.Error("unsupported fx struct type",
				"struct", v.Name,
				"type", structType,
				"supported_types", "in|out")
		}
	}

	if hasGenerated {
		res = append(res, gen.GenResult{
			FileName: filepath.Join(option.Dir, "fx_structs.go"),
			FileStr:  f.GoString(),
			Cover:    true,
		})
	}

	return
}

func (p *FxStructPlug) generateFxInStruct(f *jen.File, meta gen.StructGoTypeMeta) {
	// 生成 fx.In 结构体
	f.Type().Id(meta.Name + "In").StructFunc(func(g *jen.Group) {
		g.Id("fx.In")

		// 遍历原结构体的字段，生成对应的依赖注入字段
		for i := 0; i < meta.Obj.NumFields(); i++ {
			field := meta.Obj.Field(i)
			if !field.IsField() {
				continue
			}

			// 获取字段类型
			fieldType := field.Type().String()

			// 检查字段是否有特殊标签
			tag := meta.Obj.Tag(i)
			var fxTag string

			if strings.Contains(tag, `group:`) {
				// 提取group标签
				fxTag = tag
			} else if strings.Contains(tag, `name:`) {
				// 提取name标签
				fxTag = tag
			}

			if fxTag != "" {
				g.Id(field.Name()).Id(fieldType).Tag(map[string]string{"fx": fxTag})
			} else {
				g.Id(field.Name()).Id(fieldType)
			}
		}
	}).Line()
}

func (p *FxStructPlug) generateFxOutStruct(f *jen.File, meta gen.StructGoTypeMeta) {
	// 生成 fx.Out 结构体
	f.Type().Id(meta.Name + "Out").StructFunc(func(g *jen.Group) {
		g.Id("fx.Out")

		// 遍历原结构体的字段，生成对应的输出字段
		for i := 0; i < meta.Obj.NumFields(); i++ {
			field := meta.Obj.Field(i)
			if !field.IsField() {
				continue
			}

			// 获取字段类型
			fieldType := field.Type().String()

			// 检查字段是否有特殊标签
			tag := meta.Obj.Tag(i)
			var fxTag string

			if strings.Contains(tag, `group:`) {
				// 提取group标签
				fxTag = tag
			} else if strings.Contains(tag, `name:`) {
				// 提取name标签
				fxTag = tag
			}

			if fxTag != "" {
				g.Id(field.Name()).Id(fieldType).Tag(map[string]string{"fx": fxTag})
			} else {
				g.Id(field.Name()).Id(fieldType)
			}
		}
	}).Line()
}

func (p *Plug) Gen(req []gen.GlobalFuncGoTypeMeta) (res []gen.GenResult, err error) {
	slog.Info("fx plugin starting", "total_packages", len(req))

	appMetas := []gen.GlobalFuncGoTypeMeta{}
	provideMetas := []provideMeta{}
	invokeMetas := []invokeMeta{}
	decorateMetas := []decorateMeta{}

	// 扫描所有函数，按类型分类
	lo.ForEach(req, func(item gen.GlobalFuncGoTypeMeta, index int) {
		slog.Info("processing package", "package", item.Option.Pkg.PkgPath, "functions", len(item.Metas))
		appMeta := gen.GlobalFuncGoTypeMeta{
			Option: item.Option,
			Metas:  []gen.FuncGoTypeMeta{},
		}
		provideMeta := provideMeta{
			Option: item.Option,
			Metas:  []provideMetaItem{},
		}
		invokeMeta := invokeMeta{
			Option: item.Option,
			Metas:  []invokeMetaItem{},
		}
		decorateMeta := decorateMeta{
			Option: item.Option,
			Metas:  []decorateMetaItem{},
		}

		var appMetaOk, provideMetaOk, invokeMetaOk, decorateMetaOk bool

		lo.ForEach(item.Metas, func(meta gen.FuncGoTypeMeta, index int) {
			// 检查是否有 @fx 注解
			fxDoc := meta.Doc.ByFuncName(p.Name())
			if fxDoc == nil {
				return // 没有 @fx 注解，跳过
			}

			// 解析键值对参数
			var fxType, fxName, fxGroup string

			// 从参数中提取 type
			if typeValue, ok := meta.Doc.ByFuncNameAndArgName(p.Name(), "type"); ok {
				fxType = typeValue
			}

			// 从参数中提取 name (接口名)
			if nameValue, ok := meta.Doc.ByFuncNameAndArgName(p.Name(), "name"); ok {
				fxName = nameValue
			}

			// 从参数中提取 group
			if groupValue, ok := meta.Doc.ByFuncNameAndArgName(p.Name(), "group"); ok {
				fxGroup = groupValue
			}

			slog.Info("processing fx function",
				"function", meta.Name,
				"type", fxType,
				"name", fxName,
				"group", fxGroup)

			switch fxType {
			case "app":
				appMeta.Metas = append(appMeta.Metas, meta)
				appMetaOk = true
			case "provide":
				provideMeta.Metas = append(provideMeta.Metas, provideMetaItem{
					FxName:  fxName,
					FxGroup: fxGroup,
					Meta:    meta,
				})
				provideMetaOk = true
			case "invoke":
				invokeMeta.Metas = append(invokeMeta.Metas, invokeMetaItem{
					FxName: fxName,
					Meta:   meta,
				})
				invokeMetaOk = true
			case "decorate":
				decorateMeta.Metas = append(decorateMeta.Metas, decorateMetaItem{
					FxName: fxName,
					Meta:   meta,
				})
				decorateMetaOk = true
			default:
				if fxType != "" {
					slog.Error("unsupported fx function type",
						"function", meta.Name,
						"fx_type", fxType,
						"supported_types", "app|provide|invoke|decorate")
				}
			}
		})

		if appMetaOk {
			appMetas = append(appMetas, appMeta)
		}
		if provideMetaOk {
			provideMetas = append(provideMetas, provideMeta)
		}
		if invokeMetaOk {
			invokeMetas = append(invokeMetas, invokeMeta)
		}
		if decorateMetaOk {
			decorateMetas = append(decorateMetas, decorateMeta)
		}
	})

	// 按包路径排序，确保生成代码的一致性
	sort.Slice(provideMetas, func(i, j int) bool {
		return provideMetas[i].Option.Pkg.PkgPath > provideMetas[j].Option.Pkg.PkgPath
	})
	sort.Slice(invokeMetas, func(i, j int) bool {
		return invokeMetas[i].Option.Pkg.PkgPath > invokeMetas[j].Option.Pkg.PkgPath
	})
	sort.Slice(decorateMetas, func(i, j int) bool {
		return decorateMetas[i].Option.Pkg.PkgPath > decorateMetas[j].Option.Pkg.PkgPath
	})

	// 为每个标记为 app 的包生成 fx_app.go
	lo.ForEach(appMetas, func(item gen.GlobalFuncGoTypeMeta, index int) {
		f := jen.NewFile(item.Option.Pkg.Name)
		f.AddImport("go.uber.org/fx", "")

		// 生成 FxOptions 函数
		f.Func().Id("FxOptions").Params().Index().Id("fx.Option").BlockFunc(func(g *jen.Group) {
			g.Return(jen.Index().Id("fx.Option").ValuesFunc(func(optGroup *jen.Group) {
				// 添加所有 Provide 选项
				p.addProvideOptions(f, optGroup, provideMetas)
				// 添加所有 Invoke 选项
				p.addInvokeOptions(f, optGroup, invokeMetas)
				// 添加所有 Decorate 选项
				p.addDecorateOptions(f, optGroup, decorateMetas)
			}))
		}).Line()

		// 生成便捷的 NewApp 函数
		f.Func().Id("NewApp").Params(jen.Id("opts").Op("...").Id("fx.Option")).Op("*").Id("fx.App").BlockFunc(func(g *jen.Group) {
			g.Id("allOpts").Op(":=").Append(jen.Id("FxOptions").Call(), jen.Id("opts").Op("..."))
			g.Return(jen.Id("fx.New").Call(jen.Id("allOpts").Op("...")))
		}).Line()

		res = append(res, gen.GenResult{
			FileName: filepath.Join(item.Option.Dir, "fx_app.go"),
			FileStr:  f.GoString(),
			Cover:    true,
		})
	})

	return
}

// 数据结构定义
type provideMeta struct {
	Option gen.Option
	Metas  []provideMetaItem
}

type provideMetaItem struct {
	FxName  string
	FxGroup string
	Meta    gen.FuncGoTypeMeta
}

type invokeMeta struct {
	Option gen.Option
	Metas  []invokeMetaItem
}

type invokeMetaItem struct {
	FxName string
	Meta   gen.FuncGoTypeMeta
}

type decorateMeta struct {
	Option gen.Option
	Metas  []decorateMetaItem
}

type decorateMetaItem struct {
	FxName string
	Meta   gen.FuncGoTypeMeta
}

// 添加 Provide 选项的辅助方法
func (p *Plug) addProvideOptions(f *jen.File, optGroup *jen.Group, provideMetas []provideMeta) {
	lo.ForEach(provideMetas, func(v provideMeta, index int) {
		pkgPath := v.Option.Pkg.PkgPath
		pkgName := v.Option.Pkg.Name
		pkgName = p.getPkgName(pkgName)
		f.AddImport(pkgPath, pkgName)

		sort.Slice(v.Metas, func(i, j int) bool {
			return v.Metas[i].Meta.Name > v.Metas[j].Meta.Name
		})

		lo.ForEach(v.Metas, func(vv provideMetaItem, index int) {
			if vv.FxGroup != "" {
				// fx.Provide(fx.Annotate(constructor, fx.ResultTags(`group:"groupname"`)))
				optGroup.Id("fx.Provide").Call(
					jen.Id("fx.Annotate").Call(
						jen.Id(pkgName+"."+vv.Meta.Name),
						jen.Id("fx.ResultTags").Call(jen.Lit(fmt.Sprintf(`group:"%s"`, vv.FxGroup))),
					),
				)
			} else if vv.FxName != "" {
				// fx.Provide(fx.Annotate(constructor, fx.As(new(Interface))))
				optGroup.Id("fx.Provide").Call(
					jen.Id("fx.Annotate").Call(
						jen.Id(pkgName+"."+vv.Meta.Name),
						jen.Id("fx.As").Call(jen.New(jen.Id(vv.FxName))),
					),
				)
			} else {
				// fx.Provide(constructor)
				optGroup.Id("fx.Provide").Call(jen.Id(pkgName + "." + vv.Meta.Name))
			}
		})
	})
}

// 添加 Invoke 选项的辅助方法
func (p *Plug) addInvokeOptions(f *jen.File, optGroup *jen.Group, invokeMetas []invokeMeta) {
	lo.ForEach(invokeMetas, func(v invokeMeta, index int) {
		pkgPath := v.Option.Pkg.PkgPath
		pkgName := v.Option.Pkg.Name
		pkgName = p.getPkgName(pkgName)
		f.AddImport(pkgPath, pkgName)

		sort.Slice(v.Metas, func(i, j int) bool {
			return v.Metas[i].Meta.Name > v.Metas[j].Meta.Name
		})

		lo.ForEach(v.Metas, func(vv invokeMetaItem, index int) {
			optGroup.Id("fx.Invoke").Call(jen.Id(pkgName + "." + vv.Meta.Name))
		})
	})
}

// 添加 Decorate 选项的辅助方法
func (p *Plug) addDecorateOptions(f *jen.File, optGroup *jen.Group, decorateMetas []decorateMeta) {
	lo.ForEach(decorateMetas, func(v decorateMeta, index int) {
		pkgPath := v.Option.Pkg.PkgPath
		pkgName := v.Option.Pkg.Name
		pkgName = p.getPkgName(pkgName)
		f.AddImport(pkgPath, pkgName)

		sort.Slice(v.Metas, func(i, j int) bool {
			return v.Metas[i].Meta.Name > v.Metas[j].Meta.Name
		})

		lo.ForEach(v.Metas, func(vv decorateMetaItem, index int) {
			if vv.FxName != "" {
				// fx.Decorate(fx.Annotate(decorator, fx.As(new(Interface))))
				optGroup.Id("fx.Decorate").Call(
					jen.Id("fx.Annotate").Call(
						jen.Id(pkgName+"."+vv.Meta.Name),
						jen.Id("fx.As").Call(jen.New(jen.Id(vv.FxName))),
					),
				)
			} else {
				// fx.Decorate(decorator)
				optGroup.Id("fx.Decorate").Call(jen.Id(pkgName + "." + vv.Meta.Name))
			}
		})
	})
}
