package gormq

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/fitan/genx/common"
	"github.com/fitan/jennifer/jen"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/go/packages"
)

/*
@enum(
"Eq:=",
"Neq:<>",
"In:IN",
"NotIn:NOT IN",
"Lt:<",
"Lte:<=",
"Gt:>",
"Gte:>=",
"OR:OR",
"AND:AND"
)
*/
type SqlWhereOp int

// @enum("customize:自定义","column:数据库字段","associate:关联表")
type WhereFieldType int

type SelectGen struct {
}

type IncludeGen struct {
}

type WhereGen struct {
	FieldType WhereFieldType
	Column    string
	Op        SqlWhereOp
	Path      []string
	Next      []WhereGen
}

func (w WhereGen) Gen() jen.Statement {
	return jen.Id("s").Dot(strings.Join(w.Path, "."))
}

func (w WhereGen) ColumnGen() jen.Statement {
	codes := make([]jen.Code, 0)
	codes = append(codes, jen.Lit(fmt.Sprintf("%s %s ?", w.Column, w.Op.String())))
	codes = append(codes, jen.Id("s").Dot(strings.Join(w.Path, ".")))
	return codes
}

func (w WhereGen)

type OrderByGen struct {
}

// FindMany
type PMMany struct {
	Where    struct{}
	OrderBy  struct{}
	Select   struct{}
	Include  struct{}
	Distinct struct{}
	Skip     int
	Cursor   struct{}
	Take     int
}

// 查询物理机 findFirst
type PMFirst struct {
	Select struct {
	}
	Include struct {
	}
	Where struct {
		IP struct {
			Eq         string
			Not        string
			In         []string
			NotIn      string
			Lt         string
			Lte        string
			Gt         string
			Gte        string
			Contains   string
			StartsWith string
			EndsWith   string

			OR struct {
			}
			AND struct {
			}
		}
	}
	OrderBy  struct{}
	Cursor   struct{}
	Take     struct{}
	Skip     struct{}
	Distinct struct{}
}

type GormQ struct {
	J           *jen.File
	Pkg         *packages.Package
	ObjName     string
	Obj         *types.Struct
	Model       string
	GenMetaData []MetaData
	MetaDatas   []common.StructFieldMetaData
}

type MetaData struct {
	GormModelFieldPath string
	Op                 string
	ObjFieldPath       string
}

func (g GormQ) gen() (err error) {
	slog.Info("gen", slog.Any("g", g.MetaDatas))
	for _, f := range g.MetaDatas {
		metaData := MetaData{}
		has := f.Doc.ByFuncNameAndArgs("@gormq", &metaData.GormModelFieldPath, &metaData.Op)
		if !has {
			continue
		}

		if metaData.GormModelFieldPath == "" {
			err = fmt.Errorf("@gormq: %s is empty", f.Path)
			return
		}

		if metaData.Op == "" {
			metaData.Op = "="
		}

		metaData.ObjFieldPath = strings.Join(f.Path, ".")

		g.GenMetaData = append(g.GenMetaData, metaData)
	}
	slog.Info("genmetadata", slog.Any("GenMetaData", g.GenMetaData))

	g.J.Func().Params(
		jen.Id("g").Id(g.ObjName),
	).Id("GormQScopes").Params(jen.Id("res").Id("[]func(*gorm.DB) *gorm.DB"), jen.Id("err").Error()).Block(
		jen.Id("req").Op(":=").Make(jen.Index().Qual("github.com/fitan/mykit/mygorm", "GenxScopesReq"), jen.Lit(0)),
		func() jen.Code {
			code := jen.Line()
			for _, v := range g.GenMetaData {
				code.Id("req").Op("=").Append(
					jen.Id("req"), jen.Qual("github.com/fitan/mykit/mygorm", "GenxScopesReq").Values(jen.Dict{
						jen.Id("Field"): jen.Lit(v.GormModelFieldPath),
						jen.Id("Op"):    jen.Lit(v.Op),
						jen.Id("Value"): jen.Id("g").Dot(v.ObjFieldPath),
					})).Line()
			}
			return code
		}(),
		jen.Return(jen.Qual("github.com/fitan/mykit/mygorm", "GenxScopes").Call(jen.Id(g.Model), jen.Id("req"))),
	)

	return nil

}

func Gen(j *jen.File, pkg *packages.Package, objName string, modelS string, obj *types.Struct, data []common.StructFieldMetaData) error {
	q := GormQ{
		J:           j,
		Pkg:         pkg,
		ObjName:     objName,
		Obj:         obj,
		Model:       modelS,
		GenMetaData: make([]MetaData, 0),
		MetaDatas:   data,
	}
	return q.gen()

}
