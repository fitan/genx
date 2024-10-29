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

func (p Plug) Gen(option gen.Option, typeSpecMetas []gen.TypeSpecGoTypeMeta) (res []gen.GenResult, err error) {
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
			return nil, err
		}

		j.Add(codes...)
	}

	res = append(res, gen.GenResult{
		FileName: filepath.Join(option.Dir, "enum.go"),
		FileStr:  j.GoString(),
		Cover:    true,
	})

	return
}
