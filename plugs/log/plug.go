package log

import (
	"path/filepath"

	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"golang.org/x/exp/slog"
)

type Plug struct {
}

func (p *Plug) Name() string {
	return "@log"
}

func (p *Plug) Gen(option gen.Option, implGoTypeMetes []gen.InterfaceGoTypeMeta) (res []gen.GenResult, err error) {

	parseImpl := common.NewInterfaceSerialize(option.Pkg)

	for _, v := range implGoTypeMetes {
		meta, err := parseImpl.Parse(v.Obj, v.RawDoc, &v.Doc)
		if err != nil {
			slog.Error("parseImpl.Parse", err, slog.String("name", v.Obj.String()))
			return nil, common.ParseError("failed to parse interface for log plugin").
				WithCause(err).
				WithPlugin("@log").
				WithInterface(v.Name).
				WithDetails("unable to parse interface metadata for logging instrumentation").
				Build()
		}

		f := Gen(option.Pkg, meta.Methods)

		res = append(res, gen.GenResult{
			FileName: filepath.Join(option.Dir, "logging.go"),
			FileStr:  f,
			Cover:    true,
		})
	}
	return
}
