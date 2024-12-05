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

func (p *Plug) Gen(option gen.Option, implGoTypeMetes []gen.InterfaceGoTypeMeta) (res []gen.GenResult, err error) {

	parseImpl := common.NewInterfaceSerialize(option.Pkg)
	for _, v := range implGoTypeMetes {
		meta, err := parseImpl.Parse(v.Obj, v.RawDoc, &v.Doc)
		if err != nil {
			slog.Error("parseImpl.Parse", err, slog.String("name", v.Obj.String()))
			return nil, err
		}

		f := Gen(option.Pkg, meta.Methods)

		res = append(res, gen.GenResult{
			FileName: filepath.Join(option.Dir, "otel.go"),
			FileStr:  f,
			Cover:    true,
		})

	}
	return
}
