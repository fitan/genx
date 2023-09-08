package gormq

import (
	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
	"golang.org/x/exp/slog"
)

//go:generate gowrap gen -g -p ./
type Plug struct {
}

func (p *Plug) Name() string {
	return "@gormModel"
}

func (p *Plug) Gen(option gen.Option, implGoTypeMetes []gen.StructGoTypeMeta) (err error) {
	slog.Info("gormModelGen", slog.Any("option", option), slog.Any("implGoTypeMetes", implGoTypeMetes))
	j := jen.NewFile(option.Pkg.Name)
	parseStruct := common.NewStructSerialize(option.Pkg)
	slog.Info("implGoTypeMets", implGoTypeMetes)
	for _, v := range implGoTypeMetes {
		var meta common.StructMetaData
		meta, err = parseStruct.Parse(v.Obj)
		if err != nil {
			slog.Error("parseImpl.Parse", err, slog.String("name", v.Obj.String()))
			return err
		}

		var modelS string

		v.Doc.ByFuncNameAndArgs(p.Name(), &modelS)

		err = Gen(j, option.Pkg, v.Name, modelS, v.Obj, meta.Fields)
		if err != nil {
			return
		}
	}
	slog.Info("scope.go", j.GoString())
	common.WriteGO("gorm_scope.go", j.GoString())
	return
}
