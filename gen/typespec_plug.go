package gen

import (
	"github.com/fitan/genx/common"
	"go/ast"
)

type TypeSpecPlugImpl interface {
	Name() string
	Gen(option Option, typeSpecMetas []TypeSpecGoTypeMeta) error
}

type TypeSpecMeta struct {
	NameGoTypeMap map[string][]TypeSpecGoTypeMeta
}

type TypeSpecGoTypeMeta struct {
	Doc    *common.Doc
	Params []string
	Obj    *ast.TypeSpec
}
