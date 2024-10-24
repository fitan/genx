package kithttp

import (
	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"golang.org/x/exp/slog"
)

type Plug struct {
}

func (p *Plug) Name() string {
	return "@kit"
}

func (p *Plug) Gen(option gen.Option, implGoTypeMetes []gen.InterfaceGoTypeMeta) error {

	parseImpl := common.NewInterfaceSerialize(option.Pkg)
	for _, v := range implGoTypeMetes {
		meta, err := parseImpl.Parse(v.Obj, &v.Doc)
		if err != nil {
			slog.Error("parseImpl.Parse", err, slog.String("name", v.Obj.String()))
			return err
		}

		Gen(option, meta)
	}
	return nil
}

type ObserverPlug struct {
}

func (p *ObserverPlug) Name() string {
	return "@observer"
}

func (p *ObserverPlug) Gen(option gen.Option, implGoTypeMetes []gen.InterfaceGoTypeMeta) error {
	parseImpl := common.NewInterfaceSerialize(option.Pkg)
	for _, v := range implGoTypeMetes {
		meta, err := parseImpl.Parse(v.Obj, &v.Doc)
		if err != nil {
			slog.Error("parseImpl.Parse", err, slog.String("name", v.Obj.String()))
			return err
		}

		err = ObserverGen(option, meta)
		if err != nil {
			return err
		}
	}
	return nil
}
