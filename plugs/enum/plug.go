package enum

import (
	"path/filepath"

	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
)

type Plug struct {
}

func (p Plug) Name() string {
	return "@enum"
}

func (p Plug) Gen(option gen.Option, typeSpecMetas []gen.TypeSpecGoTypeMeta) error {
	j := jen.NewFile(option.Pkg.Name)
	j.HeaderComment("Code generated . DO NOT EDIT.")
	common.JenAddImports(option.Pkg, j)
	for _, v := range typeSpecMetas {
		f := v.Doc.ByFuncName("@enum")
		if f == nil {
			continue
		}

		enum := Enum{
			Param:    f.Args,
			TypeName: v.Obj.Name.Name,
		}
		codes, err := enum.Gen()
		if err != nil {
			return err
		}

		j.Add(codes...)
	}

	common.WriteGO(filepath.Join(option.Dir, "enum.go"), j.GoString())
	return nil
}
