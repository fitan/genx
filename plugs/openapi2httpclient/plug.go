package openapi2httpclient

import (
	"io/ioutil"

	"github.com/davecgh/go-spew/spew"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/renderer"
)

func Load() {
	// create a new JSON mock generator

	// create a new JSON mock generator
	mg := renderer.NewMockGenerator(renderer.)

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

	paths := v3Model.Model.Paths
	schemas := v3Model.Model.Components.Schemas

	item, _ := paths.PathItems.OrderedMap.Get("/v1/assets/list")
	
	
	mg.GenerateMock()

}
