package gen

import (
	"github.com/fitan/genx/common"
	"go/types"
)

type InterfacePlugImpl interface {
	Name() string
	Gen(option Option, implGoTypeMetes []InterfaceGoTypeMeta) error
}

type ImplMeta struct {
	NameGoTypeMap map[string][]InterfaceGoTypeMeta
}

type InterfaceGoTypeMeta struct {
	Doc    *common.Doc
	Params []string
	Obj    *types.Interface
}
