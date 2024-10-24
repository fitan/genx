package crud

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
	"github.com/sourcegraph/conc"
)

type Plug struct {
}

func (s Plug) Name() string {
	return "@crud"
}

func (s Plug) http(option gen.Option, structGoTypeMetas []gen.StructGoTypeMeta) error {

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

	err := crud.Gen()
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	c := conc.NewWaitGroup()
	c.Go(func() {
		common.WriteGoByOpt(filepath.Join(option.Dir, "crud_http_types.go"), crudHttpTypeJ.GoString(), common.WriteOpt{
			Cover: false,
		})
		// common.WriteGO(filepath.Join(option.Dir, "crud_http_types.go"), crudHttpTypeJ.GoString())
	})

	c.Go(func() {
		common.WriteGoByOpt(filepath.Join(option.Dir, "crud_http_service.go"), crudHttpJ.GoString(), common.WriteOpt{
			Cover: false,
		})
	})

	c.Wait()

	return nil

}

func (s Plug) gorm(option gen.Option, structGoTypeMetas []gen.StructGoTypeMeta) error {
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

	err := crud.Gen()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	c := conc.NewWaitGroup()

	c.Go(func() {
		common.WriteGO(filepath.Join(option.Dir, "crud_gorm_types.go"), crudGormTypeJ.GoString())
	})

	c.Go(func() {
		common.WriteGO(filepath.Join(option.Dir, "crud_gorm_service.go"), crudGormJ.GoString())
	})

	c.Wait()

	return nil
}

func (s Plug) Gen(option gen.Option, structGoTypeMetas []gen.StructGoTypeMeta) error {
	for _, v := range structGoTypeMetas {

		typeName, has := v.Doc.ByFuncNameAndArgName("@crud", "type")
		if !has {
			slog.Error("@crud type must be set")
			return fmt.Errorf("@crud type must be set")
		}

		switch typeName {
		case "http":
			err := s.http(option, structGoTypeMetas)
			if err != nil {
				return err
			}
		case "gorm":
			err := s.gorm(option, structGoTypeMetas)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("@crud tag value must be http or gorm")
		}
	}

	return nil
}
