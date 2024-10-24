package crud

import (
	"fmt"
	"strings"

	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
	"github.com/samber/lo"
)

type GormCrud struct {
	CrudGormJ         *jen.File
	CrudGormTypeJ     *jen.File
	CrudScopeJ        *jen.File
	Option            gen.Option
	StructGoTypeMetas []gen.StructGoTypeMeta
}

func (s *GormCrud) Gen() error {

	for _, i := range s.StructGoTypeMetas {
		err := s.genImpl(i)
		if err != nil {
			return err
		}
	}

	return nil

}

func (s *GormCrud) genImpl(i gen.StructGoTypeMeta) (err error) {
	modelId, modelHas := i.Doc.ByFuncNameAndArgName("@crud", "model")
	idName, idHas := i.Doc.ByFuncNameAndArgName("@crud", "idName")
	idType, idTypeHas := i.Doc.ByFuncNameAndArgName("@crud", "idType")
	if !idHas || !modelHas || !idTypeHas {
		err = fmt.Errorf("@crud model,id,idType must be set")
		return
	}

	modelName := strings.ReplaceAll(modelId, ".", "")

	input := CrudTemplateInput{
		InterfaceName: modelName + "GormCrudBaseImpl",
		StructName:    i.Name,
		ModelName:     modelId,
		IdName:        idName,
		IdType:        idType,
	}

	crudGormFile, err := common.GenFileByTemplate(s.Option.Static, "crud_gorm", input)
	if err != nil {
		return
	}

	crudGormTypeFile, err := common.GenFileByTemplate(s.Option.Static, "crud_gorm_types", input)
	if err != nil {
		return
	}

	s.CrudGormJ.Id(crudGormFile)

	s.CrudGormTypeJ.Id(crudGormTypeFile)

	return
}

type CrudTemplateInput struct {
	InterfaceName string
	StructName    string
	ModelName     string
	IdName        string
	IdType        string
}

func (c CrudTemplateInput) ModelPkgName() string {
	return strings.Split(c.ModelName, ".")[0]
}

func (c CrudTemplateInput) GetModelFnName() string {
	s, _ := lo.Last(strings.Split(c.ModelName, "."))
	return s
}

type HttpCrud struct {
	CrudHttpServiceJ     *jen.File
	CrudHttpServiceTypeJ *jen.File
	Option               gen.Option
	StructGoTypeMetas    []gen.StructGoTypeMeta
}

func (s *HttpCrud) Gen() error {

	for _, i := range s.StructGoTypeMetas {
		err := s.genImpl(i)
		if err != nil {
			return err
		}
	}

	return nil

}

func (s *HttpCrud) genImpl(i gen.StructGoTypeMeta) (err error) {
	modelId, modelHas := i.Doc.ByFuncNameAndArgName("@crud", "model")
	idName, idHas := i.Doc.ByFuncNameAndArgName("@crud", "idName")
	idType, idTypeHas := i.Doc.ByFuncNameAndArgName("@crud", "idType")
	if !idHas || !modelHas || !idTypeHas {
		err = fmt.Errorf("@crud model,id,idType must be set")
		return
	}

	modelName := strings.ReplaceAll(modelId, ".", "")

	input := CrudTemplateInput{
		InterfaceName: modelName + "HttpCrudBaseImpl",
		StructName:    i.Name,
		ModelName:     modelId,
		IdName:        idName,
		IdType:        idType,
	}

	crudHttpFile, err := common.GenFileByTemplate(s.Option.Static, "crud_http", input)
	if err != nil {
		return
	}

	crudHttpTypesFile, err := common.GenFileByTemplate(s.Option.Static, "crud_http_types", input)
	if err != nil {
		return
	}

	/* crudGormTypeFile, err := common.GenFileByTemplate(s.Option.Static, "crud_gorm_type", input)
	if err != nil {
		return
	}

	crudGormScopeFile, err := common.GenFileByTemplate(s.Option.Static, "crud_gorm_scope", input) */

	s.CrudHttpServiceJ.Id(crudHttpFile)

	s.CrudHttpServiceTypeJ.Id(crudHttpTypesFile)

	return
}
