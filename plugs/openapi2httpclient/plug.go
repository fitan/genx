package openapi2httpclient

import (
	"fmt"
	"io/ioutil"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/renderer"
)

func Load() {
	// create a new JSON mock generator

	// create a new JSON mock generator
	mg := renderer.NewMockGenerator(renderer.JSON)

	// tell the mock generator to pretty print the output
	mg.SetPretty()
	f, _ := ioutil.ReadFile("paas-assets.openapi.json")
	document, err := libopenapi.NewDocument(f)

	if err != nil {
		panic(err)
	}

	// create a new document from specification and build a v3 model.
	v3Model, errs := document.BuildV3Model()
	if errs != nil {
		panic(errs)
	}

	// create a mock of the Fries model
	friesModel, ok := v3Model.Model.Components.Schemas.Get("src_pkg_collect.TaskListResponse")
	if !ok {
		panic("Fries model not found")
	}

	// generate a mock of the fries schema
	mock, err := mg.GenerateMock(friesModel, "")

	if err != nil {
		panic(err)
	}

	// print the mock to stdout
	fmt.Println(string(mock))
}
