package temporal

import (
	"path/filepath"
	"strings"

	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
)

const FuncName = "@temporal"

type Plug struct {
}

func (p *Plug) Name() string {
	return FuncName
}

func (p *Plug) Gen(option gen.Option, implGoTypeMetes []gen.InterfaceGoTypeMeta) (res []gen.GenResult, err error) {
	j := jen.NewFile(option.Pkg.Name)
	for _, v := range option.Imports {
		if v.Name != nil {
			j.AddImport(strings.Trim(v.Path.Value, `"`), strings.Trim(v.Name.String(), `"`))
		} else {
			j.AddImport(strings.Trim(v.Path.Value, `"`), "")
		}
	}
	//parseStruct := common.NewStructSerialize(option.Pkg)
	//slog.Info("implGoTypeMets", implGoTypeMetes)
	//for _, v := range implGoTypeMetes {
	//	var meta common.StructMetaData
	//	meta, err = parseStruct.Parse(v.Obj)
	//	if err != nil {
	//		slog.Error("parseImpl.Parse", err, slog.String("name", v.Obj.String()))
	//		return err
	//	}
	//
	//	var modelS string
	//
	//	v.Doc.ByFuncNameAndArgs(p.Name(), &modelS)
	//
	//	err = Gen(j, option.Pkg, v.Name, modelS, v.Obj, meta.Fields)
	//	if err != nil {
	//		return
	//	}
	//}

	Gen(j, option, implGoTypeMetes)

	res = append(res, gen.GenResult{
		FileName: filepath.Join(option.Dir, "temporal.go"),
		FileStr:  j.GoString(),
		Cover:    true,
	})
	return
}
