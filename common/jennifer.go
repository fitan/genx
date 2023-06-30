package common

import (
	"github.com/fitan/jennifer/jen"
	"golang.org/x/tools/go/packages"
	"strings"
)

func JenAddImports(p *packages.Package, f *jen.File) {
	for _, s := range p.Syntax {
		for _, v := range s.Imports {
			var path, pathName string
			if v.Path != nil {
				path = strings.Trim(v.Path.Value, `"`)
			}
			if v.Name != nil {
				pathName = strings.Trim(v.Name.Name, `"`)
			}
			f.AddImport(path, pathName)
		}
	}
}
