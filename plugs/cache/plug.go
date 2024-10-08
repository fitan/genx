package cache

import (
	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"golang.org/x/exp/slog"
)

type Plug struct {
}

func (p *Plug) Name() string {
	return "@cache"
}

func (p *Plug) Gen(option gen.Option, implGoTypeMetes []gen.InterfaceGoTypeMeta) error {

	parseImpl := common.NewInterfaceSerialize(option.Pkg)
	for _, v := range implGoTypeMetes {
		meta, err := parseImpl.Parse(v.Obj, &v.Doc)
		if err != nil {
			slog.Error("parseImpl.Parse", err, slog.String("name", v.Obj.String()))
			return err
		}

		Gen(option.Pkg, v.Doc, meta.Methods)
	}
	return nil
}
