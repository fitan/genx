package gormq

import (
	"fmt"
	"github.com/fitan/genx/common"
	"github.com/fitan/jennifer/jen"
	"go/types"
	"golang.org/x/tools/go/packages"
	"strings"
)

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
	for _, f := range g.MetaDatas {
		metaData := MetaData{}
		has := f.Doc.ByFuncNameAndArgs("gormq", &metaData.GormModelFieldPath, &metaData.Op)
		if !has {
			continue
		}

		if metaData.GormModelFieldPath == "" {
			err = fmt.Errorf("gormq: %s is empty", f.Path)
			return
		}

		metaData.ObjFieldPath = strings.Join(f.Path, ".")

		g.GenMetaData = append(g.GenMetaData, metaData)
	}

	g.J.Func().Params(
		jen.Id("g").Id(g.ObjName),
	).Id("GormQScopes").Params(jen.Id("res").Id("[]func(*gorm.DB) *gorm.DB"), jen.Id("err").Error()).Block(
		jen.Id("req").Op(":=").Make(jen.Index().Qual("github.com/fitan/mykit/mygorm", "GenxScopesReq")), jen.Lit(0),
		func() jen.Code {
			code := jen.Line()
			for _, v := range g.GenMetaData {
				code.Id("req").Op("=").Append(
					jen.Id("req"), jen.Qual("github.com/fitan/mykit/mygorm", "GenxScopesReq").Values(jen.Dict{
						jen.Id("Field"): jen.Lit(v.GormModelFieldPath),
						jen.Id("Op"):    jen.Lit(v.Op),
						jen.Id("Value"): jen.Id("g").Dot(v.ObjFieldPath),
					}))
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
