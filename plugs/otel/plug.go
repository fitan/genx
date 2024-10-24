package otel

import (
	"path/filepath"

	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"golang.org/x/exp/slog"
)

type Plug struct {
}

func (p *Plug) Name() string {
	return "@otel"
}

func (p *Plug) Gen(option gen.Option, implGoTypeMetes []gen.InterfaceGoTypeMeta) error {

	parseImpl := common.NewInterfaceSerialize(option.Pkg)
	for _, v := range implGoTypeMetes {
		meta, err := parseImpl.Parse(v.Obj, &v.Doc)
		if err != nil {
			slog.Error("parseImpl.Parse", err, slog.String("name", v.Obj.String()))
			return err
		}

		f := Gen(option.Pkg, meta.Methods)

		common.WriteGO(filepath.Join(option.Dir, "otel.go"), f)
	}
	return nil
}
