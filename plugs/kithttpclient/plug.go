package kithttpclient

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
)

type Plug struct {
}

func (s Plug) Name() string {
	return "@kit-http-client"
}

func (s Plug) Gen(option gen.Option, implGoTypeMetes []gen.InterfaceGoTypeMeta) error {
	j := jen.NewFile(option.Pkg.Name)

	kitHttpClient := &KitHttpClient{implGoTypeMetes: implGoTypeMetes, option: option}

	kitHttpClient.Parse()

	slog.Info("kit http client", "methods", kitHttpClient.methods)
	b, _ := json.Marshal(kitHttpClient.methods)
	fmt.Println(string(b))

	kitHttpClient.Gen(j)

	common.WriteGO("kit_http_client.go", j.GoString())
	return nil

}
