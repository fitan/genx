package kithttp

import (
	"path"
	"strings"

	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/genx/parser"
	"github.com/samber/lo"
)

type TemplateInputInterface struct {
	Name    string
	Methods map[string]Method
	Doc     common.Doc
	Opt     gen.Option
}

func (t *TemplateInputInterface) Instance() string {
	return common.Last2DirName(t.Opt.Dir)
}

func (t *TemplateInputInterface) KitServerOption() (res string) {
	line := t.Doc.ByFuncName("@kit-server-option")

	if line == nil {
		return ""
	}

	return strings.Join(lo.Map(line.Args, func(item parser.FuncArg, index int) string {
		return item.Value
	}), ",")
}

func (t *TemplateInputInterface) Tags() string {
	var tag string
	t.Doc.ByFuncNameAndArgs("@tags", &tag)

	return tag
}

func (t *TemplateInputInterface) ValidVersion() string {
	var validVersion string
	t.Doc.ByFuncNameAndArgs("@validVersion", &validVersion)

	return validVersion
}

func (t *TemplateInputInterface) BasePath() string {
	var basePath string
	t.Doc.ByFuncNameAndArgs("@basePath", &basePath)

	return basePath
}

func (t *TemplateInputInterface) EnableSwag(name string) bool {
	var swag string
	t.Doc.ByFuncNameAndArgs("@swag", &swag)
	if swag == "false" {
		return false
	}

	return t.Methods[name].EnableSwag()
}

func (t *TemplateInputInterface) HasMethodPath(name string) bool {
	return t.Methods[name].RawKit.Conf.Url != ""
}

func (t *TemplateInputInterface) MethodPath(name string) string {
	return strings.TrimSuffix(path.Join(t.BasePath(), t.Methods[name].RawKit.Conf.Url), "/")
}
