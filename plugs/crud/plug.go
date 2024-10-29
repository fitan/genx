package crud

import (
	"fmt"
	"path/filepath"

	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
)

type Plug struct {
}

func (s Plug) Name() string {
	return "@crud"
}

func (s Plug) http(option gen.Option, structGoTypeMetas []gen.StructGoTypeMeta) (res []gen.GenResult, err error) {

	crudHttpJ := jen.NewFile(option.Pkg.Name)
	crudHttpTypeJ := jen.NewFile(option.Pkg.Name)

	for _, v := range option.Config.Imports {
		crudHttpJ.AddImport(v.Path, v.Alias)
		crudHttpTypeJ.AddImport(v.Path, v.Alias)
	}

	crudHttpJ.AddImport("github.com/pkg/errors", "")
	crudHttpJ.AddImport("github.com/samber/lo", "")

	crud := &HttpCrud{
		CrudHttpServiceJ:     crudHttpJ,
		CrudHttpServiceTypeJ: crudHttpTypeJ,
		Option:               option,
		StructGoTypeMetas:    structGoTypeMetas,
	}

	err = crud.Gen()
	if err != nil {
		return
	}

	res = append(res, gen.GenResult{
		FileName: filepath.Join(option.Dir, "crud_http_types.go"),
		FileStr:  crudHttpTypeJ.GoString(),
		Cover:    false,
	})

	res = append(res, gen.GenResult{
		FileName: filepath.Join(option.Dir, "crud_http_service.go"),
		FileStr:  crudHttpJ.GoString(),
		Cover:    false,
	})

	return

}

func (s Plug) gorm(option gen.Option, structGoTypeMetas []gen.StructGoTypeMeta) (res []gen.GenResult, err error) {
	crudGormJ := jen.NewFile(option.Pkg.Name)
	crudGormTypeJ := jen.NewFile(option.Pkg.Name)
	crudGormScopeJ := jen.NewFile(option.Pkg.Name)

	crudGormJ.AddImport("gorm.io/gorm", "")
	crudGormScopeJ.AddImport("gorm.io/gorm", "")

	for _, v := range option.Config.Imports {
		crudGormJ.AddImport(v.Path, v.Alias)
		crudGormTypeJ.AddImport(v.Path, v.Alias)
	}
	crud := &GormCrud{
		CrudGormJ:         crudGormJ,
		CrudGormTypeJ:     crudGormTypeJ,
		Option:            option,
		StructGoTypeMetas: structGoTypeMetas,
	}

	err = crud.Gen()
	if err != nil {
		return
	}

	res = append(res, gen.GenResult{
		FileName: filepath.Join(option.Dir, "crud_gorm_types.go"),
		FileStr:  crudGormTypeJ.GoString(),
		Cover:    false,
	})

	res = append(res, gen.GenResult{
		FileName: filepath.Join(option.Dir, "crud_gorm_service.go"),
		FileStr:  crudGormJ.GoString(),
		Cover:    false,
	})

	return
}

func (s Plug) Gen(option gen.Option, structGoTypeMetas []gen.StructGoTypeMeta) (res []gen.GenResult, err error) {
	for _, v := range structGoTypeMetas {

		typeName, has := v.Doc.ByFuncNameAndArgName("@crud", "type")
		if !has {
			err = fmt.Errorf("@crud type must be set")
			return
		}

		switch typeName {
		case "http":
			return s.http(option, structGoTypeMetas)
		case "gorm":
			return s.gorm(option, structGoTypeMetas)
		default:
			err = fmt.Errorf("@crud type %s not support", typeName)
			return
		}
	}

	return
}
