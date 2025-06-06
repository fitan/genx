package gormq

import (
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
)

const FuncName = "@gq"

type Value string

type Op string

func (o Op) ConvertValue(path string) (res []jen.Code) {
	switch o {
	case "><", "!><":
		res = append(res, jen.Id(path).Index(jen.Id("0")))
		res = append(res, jen.Id(path).Index(jen.Id("1")))
		return
	case "like":
		res = append(res, jen.Lit("%").Op("+").Id(path).Op("+").Lit("%"))
		return

	default:
		res = append(res, jen.Id(path))
		return
	}
}

func (o Op) Convert() string {
	switch o {
	case "=":
		return "="
	case "!=":
		return "!="
	case ">":
		return ">"
	case ">=":
		return ">="
	case "<":
		return "<"
	case "<=":
		return "<="
	case "><":
		return "better"
	case "!><":
		return "not better"
	case "like":
		return "like"
	case "in":
		return "in"
	case "!in":
		return "not in"
	//case "between":
	//	return "between"
	//case "!between":
	//	return "not between"
	case "null":
		return "is null"
	case "!null":
		return "is not null"
	default:
		// 记录不支持的操作符，但不 panic
		slog.Error("not support op", "op", string(o))
		return "unknown"
	}
}

//go:generate gowrap gen -g -p ./
type Plug struct {
}

func (p *Plug) Name() string {
	return FuncName
}

func (p *Plug) Gen(option gen.Option, implGoTypeMetes []gen.StructGoTypeMeta) (res []gen.GenResult, err error) {
	j := jen.NewFile(option.Pkg.Name)
	j.AddImport("gorm.io/gorm", "")
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
		FileName: filepath.Join(option.Dir, "gorm_scope.go"),
		FileStr:  j.GoString(),
		Cover:    true,
	})
	return
}
