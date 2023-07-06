package gen

import (
	"github.com/fitan/genx/common"
	"go/types"
)

type StructPlugImpl interface {
	Name() string
	Gen(option Option, structGoTypeMetes []StructGoTypeMeta) error
}

type StructMeta struct {
	NameGoTypeMap map[string][]StructGoTypeMeta
}

type StructGoTypeMeta struct {
	Name   string
	Doc    *common.Doc
	Params []string
	Obj    *types.Struct
}
