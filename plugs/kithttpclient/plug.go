package kithttpclient

import (
	"path/filepath"

	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
)

type Plug struct {
}

func (s Plug) Name() string {
	return "@kit-http-client"
}

func (s Plug) Gen(option gen.Option, implGoTypeMetes []gen.InterfaceGoTypeMeta) (res []gen.GenResult, err error) {
	j := jen.NewFile(option.Pkg.Name)

	kitHttpClient := &KitHttpClient{implGoTypeMetes: implGoTypeMetes, option: option}

	kitHttpClient.Parse()

	kitHttpClient.Gen(j)

	res = append(res, gen.GenResult{
		FileName: filepath.Join(option.Dir, "kit_http_client.go"),
		FileStr:  j.GoString(),
		Cover:    true,
	})

	return

}
