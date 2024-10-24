package gen

import (
	"go/ast"
	"go/types"

	"github.com/fitan/genx/common"
)

type CallPlugImpl interface {
	Name() string
	Gen(option Option, callGoTypeMetes []CallGoTypeMeta) error
}

type CallMeta struct {
	NameGoTypeMap map[string][]CallGoTypeMeta
}

type CallGoTypeMeta struct {
	Name    string
	Doc     common.Doc
	Params  []*common.Type
	Results []*common.Type
}

type FuncPlugImpl interface {
	Name() string
	Gen(option Option, funcGoTypeMetes []FuncGoTypeMeta) error
}

type FuncMeta struct {
	NameGoTypeMap map[string][]FuncGoTypeMeta
}

type FuncGoTypeMeta struct {
	Name    string
	Doc     common.Doc
	Params  []*common.Type
	Results []*common.Type
}

type InterfacePlugImpl interface {
	Name() string
	Gen(option Option, implGoTypeMetes []InterfaceGoTypeMeta) error
}

type ImplMeta struct {
	NameGoTypeMap map[string][]InterfaceGoTypeMeta
}

type InterfaceGoTypeMeta struct {
	Name   string
	Doc    common.Doc
	Params []string
	Obj    *types.Interface
}

type StructPlugImpl interface {
	Name() string
	Gen(option Option, structGoTypeMetes []StructGoTypeMeta) error
}

type StructMeta struct {
	NameGoTypeMap map[string][]StructGoTypeMeta
}

type StructGoTypeMeta struct {
	Name   string
	Doc    common.Doc
	Params []string
	Obj    *types.Struct
}

type TypePlugImpl interface {
	Name() string
	Gen(option Option, typeGoTypeMetes []TypeGoTypeMeta) error
}

type TypeMeta struct {
	NameGoTypeMap map[string][]TypeGoTypeMeta
}

type TypeGoTypeMeta struct {
	Doc    common.Doc
	Params []string
	Obj    types.Type
}

type TypeSpecPlugImpl interface {
	Name() string
	Gen(option Option, typeSpecMetas []TypeSpecGoTypeMeta) error
}

type TypeSpecMeta struct {
	NameGoTypeMap map[string][]TypeSpecGoTypeMeta
}

type TypeSpecGoTypeMeta struct {
	Doc    common.Doc
	Params []string
	Obj    *ast.TypeSpec
}

type GenResult struct {
	FileName string
	FileStr  string
	Cover    bool
}
