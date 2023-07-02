package gormq

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/fitan/genx/common"
	"golang.org/x/tools/go/packages"
	"gorm.io/gorm"
)

type GormQ struct {
	Model string
	Fields []common.Field
}

type MetaData struct {
	FieldPath string
	Op string
	Value common.Type
}

func (g GormQ) ToQuery() (res string ,err error) {
	for _,f := range g.Fields {
		metaData := MetaData{}
		has := f.Doc.ByFuncNameAndArgs("gormq", &metaData.FieldPath, &metaData.Op)
		if !has {
			continue
		}

		if metaData.FieldPath == "" {
			err = fmt.Errorf("gormq: %s is empty", f.Path)
			return
		}

		if metaData.FieldPath == "" {
			metaData.FieldPath = "="
		}

		




	}
}



// gorm scope

func (g GormQ) ToQuery() string {
	op := "="
	if g.Option != "" {
		op = g.Option
	}

	if len(fields) == 1 {
		fields = append(fields, op, "?")
		return strings.Join(fields, " ")
	}
}



func Gen(pkg *packages.Package, doc common.Doc,methods []common.StructMetaDate) {


	for _,v := range methods {
		for _,f := range v.Fields {
			if f.Doc != nil {
				q := GormQ{}
				f.Doc.ByFuncNameAndArgs("gormq", &q.Fields, &q.Option)

				if q.Fields != ""  {

				}
			}
		}
	}

}