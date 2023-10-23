package gen

import "github.com/fitan/genx/common"

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
