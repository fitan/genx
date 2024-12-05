package gen

import (
	"github.com/fitan/genx/common"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc"
)

func NewGlobalX(xs []*X, tui *Model) (res *GlobalX) {
	return &GlobalX{
		Xs: xs,
		Plugs: GlobalPlugs{
			Impl:     []GlobalInterfacePlugImpl{},
			Type:     []GlobalTypePlugImpl{},
			TypeSpec: []GlobalTypeSpecPlugImpl{},
			Struct:   []GlobalStructPlugImpl{},
			Func:     []GlobalFuncPlugImpl{},
			Call:     []GlobalCallPlugImpl{},
		},
		TUI: tui,
		WG:  conc.NewWaitGroup(),
	}

}

type GlobalX struct {
	Xs    []*X
	Plugs GlobalPlugs
	TUI   *Model
	WG    *conc.WaitGroup
}

type GlobalPlugs struct {
	Impl     []GlobalInterfacePlugImpl
	Type     []GlobalTypePlugImpl
	TypeSpec []GlobalTypeSpecPlugImpl
	Struct   []GlobalStructPlugImpl
	Func     []GlobalFuncPlugImpl
	Call     []GlobalCallPlugImpl
}

func (g *GlobalX) RegCall(plug GlobalCallPlugImpl) {
	g.Plugs.Call = append(g.Plugs.Call, plug)
}

func (g *GlobalX) RegTypeSpec(plug GlobalTypeSpecPlugImpl) {
	g.Plugs.TypeSpec = append(g.Plugs.TypeSpec, plug)
}

func (g *GlobalX) RegType(plug GlobalTypePlugImpl) {
	g.Plugs.Type = append(g.Plugs.Type, plug)
}

func (g *GlobalX) RegImpl(plug GlobalInterfacePlugImpl) {
	g.Plugs.Impl = append(g.Plugs.Impl, plug)
}

func (g *GlobalX) RegStruct(plug GlobalStructPlugImpl) {
	g.Plugs.Struct = append(g.Plugs.Struct, plug)
}

func (g *GlobalX) RegFunc(plug GlobalFuncPlugImpl) {
	g.Plugs.Func = append(g.Plugs.Func, plug)
}

func (g *GlobalX) implByName(name string) (res []GlobalInterfaceGoTypeMeta, has bool) {
	lo.ForEach(g.Xs, func(item *X, index int) {
		meta, ok := item.implByName(name)
		if ok {
			has = true
			res = append(res, GlobalInterfaceGoTypeMeta{Metas: meta, Option: item.Option})
		}
	})

	return
}

func (g *GlobalX) typeSpecByName(name string) (res []GlobalTypeSpecGoTypeMeta, has bool) {
	lo.ForEach(g.Xs, func(item *X, index int) {
		meta, ok := item.typeSpecByName(name)
		if ok {
			has = true
			res = append(res, GlobalTypeSpecGoTypeMeta{Metas: meta, Option: item.Option})
		}
	})
	return
}

func (g *GlobalX) typeByName(name string) (res []GlobalTypeGoTypeMeta, has bool) {
	lo.ForEach(g.Xs, func(item *X, index int) {
		meta, ok := item.typeByName(name)
		if ok {
			has = true
			res = append(res, GlobalTypeGoTypeMeta{Metas: meta, Option: item.Option})
		}
	})
	return
}

func (g *GlobalX) structByName(name string) (res []GlobalStructGoTypeMeta, has bool) {
	lo.ForEach(g.Xs, func(item *X, index int) {
		meta, ok := item.structByName(name)
		if ok {
			has = true
			res = append(res, GlobalStructGoTypeMeta{Metas: meta, Option: item.Option})
		}
	})
	return
}

func (g *GlobalX) funcByName(name string) (res []GlobalFuncGoTypeMeta, has bool) {
	lo.ForEach(g.Xs, func(item *X, index int) {
		meta, ok := item.funcByName(name)
		if ok {
			has = true
			res = append(res, GlobalFuncGoTypeMeta{Metas: meta, Option: item.Option})
		}
	})
	return
}

func (g *GlobalX) callByName(name string) (res []GlobalCallGoTypeMeta, has bool) {
	lo.ForEach(g.Xs, func(item *X, index int) {
		meta, ok := item.callByName(name)
		if ok {
			has = true
			res = append(res, GlobalCallGoTypeMeta{Metas: meta, Option: item.Option})
		}
	})
	return
}

func (g *GlobalX) implGen() {
	for _, v := range g.Plugs.Impl {
		metas, ok := g.implByName(v.Name())
		if ok {
			modelName := v.Name()
			g.WG.Go(func() {
				g.UpdateTUI(modelName, func() (gens []GenResult, err error) { return v.Gen(metas) })
			})
		}
	}
}

func (g *GlobalX) typeGen() {
	for _, v := range g.Plugs.Type {
		metas, ok := g.typeByName(v.Name())
		if ok {
			modelName := v.Name()
			g.WG.Go(func() {
				g.UpdateTUI(modelName, func() (gens []GenResult, err error) { return v.Gen(metas) })
			})
		}
	}
}

func (g *GlobalX) typeSpecGen() {
	for _, v := range g.Plugs.TypeSpec {
		metas, ok := g.typeSpecByName(v.Name())
		if ok {
			modelName := v.Name()
			g.WG.Go(func() {
				g.UpdateTUI(modelName, func() (gens []GenResult, err error) { return v.Gen(metas) })
			})
		}
	}
}

func (g *GlobalX) structGen() {
	for _, v := range g.Plugs.Struct {
		metas, ok := g.structByName(v.Name())
		if ok {
			modelName := v.Name()
			g.WG.Go(func() {
				g.UpdateTUI(modelName, func() (gens []GenResult, err error) { return v.Gen(metas) })
			})
		}
	}
}

func (g *GlobalX) funcGen() {
	for _, v := range g.Plugs.Func {
		metas, ok := g.funcByName(v.Name())
		if ok {
			modelName := v.Name()
			g.WG.Go(func() {
				g.UpdateTUI(modelName, func() (gens []GenResult, err error) { return v.Gen(metas) })
			})
		}
	}
}

func (g *GlobalX) callGen() {
	for _, v := range g.Plugs.Call {
		metas, ok := g.callByName(v.Name())
		if ok {
			modelName := v.Name()
			g.WG.Go(func() {
				g.UpdateTUI(modelName, func() (gens []GenResult, err error) { return v.Gen(metas) })
			})
		}
	}
}

func (g *GlobalX) Gen() {
	g.typeGen()
	g.implGen()
	g.structGen()
	g.typeSpecGen()
	g.funcGen()
	g.callGen()
	g.WG.Wait()

	g.TUI.PkgEnd(UpdateTreeReq{
		PkgName: "global",
	})
}

func (g *GlobalX) UpdateTUI(plugName string, f func() (gens []GenResult, err error)) {

	g.TUI.PlugStart(UpdateTreeReq{
		PkgName:  "global",
		PlugName: plugName,
	})

	gens, err := f()
	if err != nil {
		g.TUI.PlugEnd(UpdateTreeReq{
			PkgName:  "global",
			PlugName: plugName,
			FileName: "",
			Status:   2,
			Err:      err.Error(),
		})
		return
	}

	gw := conc.NewWaitGroup()

	for _, gen := range gens {
		gw.Go(func() {
			g.TUI.FileStart(UpdateTreeReq{
				PkgName:  "global",
				PlugName: plugName,
				FileName: gen.FileName,
				Status:   0,
				Err:      "",
			})

			cover, err := common.WriteGoWithOpt(gen.FileName, gen.FileStr, common.WriteOpt{
				Cover: gen.Cover,
			})

			g.TUI.FileEnd(UpdateTreeReq{
				PkgName:  "global",
				PlugName: plugName,
				FileName: gen.FileName,
				Status:   lo.Ternary(cover, 1, 3),
				Err: lo.TernaryF(err != nil, func() string {
					return err.Error()
				}, func() string {
					return ""
				}),
			})

		})
	}

	gw.Wait()

	g.TUI.PlugEnd(UpdateTreeReq{
		PkgName:  "global",
		PlugName: plugName,
		Status:   1,
	})

	return
}
