package gen

import (
	"go/ast"
	"go/types"

	"github.com/fitan/genx/common"
)

type CallPlugImpl interface {
	Name() string
	Gen(option Option, callGoTypeMetes []CallGoTypeMeta) ([]GenResult, error)
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

type GlobalCallPlugImpl interface {
	Name() string
	Gen(req []GlobalCallGoTypeMeta) ([]GenResult, error)
}

type GlobalCallGoTypeMeta struct {
	Option Option
	Metas  []CallGoTypeMeta
}

type FuncPlugImpl interface {
	Name() string
	Gen(option Option, funcGoTypeMetes []FuncGoTypeMeta) ([]GenResult, error)
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

type GlobalFuncPlugImpl interface {
	Name() string
	Gen(req []GlobalFuncGoTypeMeta) ([]GenResult, error)
}

type GlobalFuncGoTypeMeta struct {
	Option Option
	Metas  []FuncGoTypeMeta
}

type InterfacePlugImpl interface {
	Name() string
	Gen(option Option, implGoTypeMetes []InterfaceGoTypeMeta) ([]GenResult, error)
}

type ImplMeta struct {
	NameGoTypeMap map[string][]InterfaceGoTypeMeta
}

type InterfaceGoTypeMeta struct {
	Name   string
	Doc    common.Doc
	RawDoc *ast.CommentGroup
	Params []string
	Obj    *types.Interface
}

type GlobalInterfacePlugImpl interface {
	Name() string
	Gen(req []GlobalInterfaceGoTypeMeta) ([]GenResult, error)
}

type GlobalInterfaceGoTypeMeta struct {
	Option Option
	Metas  []InterfaceGoTypeMeta
}

type StructPlugImpl interface {
	Name() string
	Gen(option Option, structGoTypeMetes []StructGoTypeMeta) ([]GenResult, error)
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

type GlobalStructPlugImpl interface {
	Name() string
	Gen(req []GlobalStructGoTypeMeta) ([]GenResult, error)
}

type GlobalStructGoTypeMeta struct {
	Option Option
	Metas  []StructGoTypeMeta
}

type TypePlugImpl interface {
	Name() string
	Gen(option Option, typeGoTypeMetes []TypeGoTypeMeta) ([]GenResult, error)
}

type TypeMeta struct {
	NameGoTypeMap map[string][]TypeGoTypeMeta
}

type TypeGoTypeMeta struct {
	Doc    common.Doc
	Params []string
	Obj    types.Type
}

type GlobalTypePlugImpl interface {
	Name() string
	Gen(req []GlobalTypeGoTypeMeta) ([]GenResult, error)
}

type GlobalTypeGoTypeMeta struct {
	Option Option
	Metas  []TypeGoTypeMeta
}

type TypeSpecPlugImpl interface {
	Name() string
	Gen(option Option, typeSpecMetas []TypeSpecGoTypeMeta) ([]GenResult, error)
}

type TypeSpecMeta struct {
	NameGoTypeMap map[string][]TypeSpecGoTypeMeta
}

type TypeSpecGoTypeMeta struct {
	Doc    common.Doc
	Params []string
	Obj    *ast.TypeSpec
}

type GlobalTypeSpecPlugImpl interface {
	Name() string
	Gen(req []GlobalTypeSpecGoTypeMeta) ([]GenResult, error)
}

type GlobalTypeSpecGoTypeMeta struct {
	Option Option
	Metas  []TypeSpecGoTypeMeta
}

type GenResult struct {
	PkgPath  string
	FileName string
	FileStr  string
	Cover    bool
}
