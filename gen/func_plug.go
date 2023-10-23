package gen

import (
	"github.com/fitan/genx/common"
)

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
