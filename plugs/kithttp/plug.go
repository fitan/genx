package kithttp

import (
	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"golang.org/x/exp/slog"
)

type CEPermissionSqlPlug struct {
}

func (p *CEPermissionSqlPlug) Name() string {
	return "@cep"
}

func (p *CEPermissionSqlPlug) Gen(option gen.Option, implGoTypeMetes []gen.InterfaceGoTypeMeta) (res []gen.GenResult, err error) {
	parseImpl := common.NewInterfaceSerialize(option.Pkg)
	for _, v := range implGoTypeMetes {
		meta, err := parseImpl.Parse(v.Obj, v.RawDoc, &v.Doc)
		if err != nil {
			slog.Error("parseImpl.Parse", err, slog.String("name", v.Obj.String()))
			return nil, err
		}

		gens, err := CEPGen(option, meta)
		if err != nil {
			return nil, err
		}

		res = append(res, gens...)
	}
	return
}

type Plug struct {
}

func (p *Plug) Name() string {
	return "@kit"
}

func (p *Plug) Gen(option gen.Option, implGoTypeMetes []gen.InterfaceGoTypeMeta) (res []gen.GenResult, err error) {

	parseImpl := common.NewInterfaceSerialize(option.Pkg)
	for _, v := range implGoTypeMetes {
		meta, err := parseImpl.Parse(v.Obj, v.RawDoc, &v.Doc)
		if err != nil {
			slog.Error("parseImpl.Parse", err, slog.String("name", v.Obj.String()))
			return nil, err
		}

		gens, err := Gen(option, meta)
		if err != nil {
			return nil, err
		}

		res = append(res, gens...)
	}
	return
}

type ObserverPlug struct {
}

func (p *ObserverPlug) Name() string {
	return "@observer"
}

func (p *ObserverPlug) Gen(option gen.Option, implGoTypeMetes []gen.InterfaceGoTypeMeta) (res []gen.GenResult, err error) {
	parseImpl := common.NewInterfaceSerialize(option.Pkg)
	for _, v := range implGoTypeMetes {
		meta, err := parseImpl.Parse(v.Obj, v.RawDoc, &v.Doc)
		if err != nil {
			slog.Error("parseImpl.Parse", err, slog.String("name", v.Obj.String()))
			return nil, err
		}

		gens, err := ObserverGen(option, meta)

		if err != nil {
			return nil, err
		}

		res = append(res, gens...)
	}
	return
}
