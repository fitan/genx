package gen

import (
	"github.com/fitan/genx/common"
	"go/types"
)

type TypePlugImpl interface {
	Name() string
	Gen(option Option, typeGoTypeMetes []TypeGoTypeMeta) error
}

type TypeMeta struct {
	NameGoTypeMap map[string][]TypeGoTypeMeta
}

type TypeGoTypeMeta struct {
	Doc    *common.Doc
	Params []string
	Obj    types.Type
}
