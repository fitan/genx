package kithttpclient

import (
	"path/filepath"

	"github.com/fitan/genx/common"
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

	// 使用安全的解析方法，避免 panic
	if parseErr := kitHttpClient.Parse(); parseErr != nil {
		// 如果是 GenxError，直接返回
		if genxErr, ok := parseErr.(*common.GenxError); ok {
			return nil, genxErr
		}
		// 否则包装为 GenxError
		return nil, common.PluginError("kit-http-client plugin failed").
			WithCause(parseErr).
			WithPlugin("@kit-http-client").
			WithDetails("failed to parse interface definitions").
			Build()
	}

	kitHttpClient.Gen(j)

	res = append(res, gen.GenResult{
		FileName: filepath.Join(option.Dir, "kit_http_client.go"),
		FileStr:  j.GoString(),
		Cover:    true,
	})

	return

}
