package crud

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"

	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
	"github.com/samber/lo"
	"golang.org/x/tools/go/packages"
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
	InterfaceName       string
	StructName          string
	ModelName           string
	IdName              string
	IdType              string
	GetResponseStruct   string
	CreateRequestStruct string
	UpdateBodyStruct    string
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

func (s *HttpCrud) depthType(t types.Type) types.Type {
	switch tt := t.(type) {
	case *types.Named:
		return s.depthType(tt.Underlying())
	case *types.Pointer:
		return s.depthType(tt.Elem())
	}

	return t
}

func (s *HttpCrud) findTypeStruct(modelId string) (*types.Struct, *packages.Package) {
	modelIdSplit := strings.Split(modelId, ".")
	pkgName, typeName := modelIdSplit[0], modelIdSplit[1]

	for _, v := range s.Option.Pkg.Imports {
		if v.Name == pkgName && v.TypesInfo != nil {
			for e, t := range v.TypesInfo.Uses {
				// if t.Type != nil &&  lo.LastOrEmpty(strings.Split(t., "/")) == modelId {
				if e != nil && e.Name == typeName {
					if structT, ok := t.Type().Underlying().(*types.Struct); ok {
						return structT, v
					}
					// fmt.Println("name", e.Name, "Id", t.Id(), "Name", t.Name(), "TypeString", t.Type().String(), "underlying", t.Type().Underlying())
				}
			}
		}
	}

	return nil, nil
}

func (s *HttpCrud) genGormStruct(structName, modelId string, preload []string) *jen.Statement {
	t, p := s.findTypeStruct(modelId)
	if t == nil {
		return jen.Null()
	}

	return s.copyGormStruct(structName, p, t, preload)
}

var depth = 0

func (s *HttpCrud) copyGormStruct(structName string, pkg *packages.Package, structT *types.Struct, preload []string) *jen.Statement {
	depth = depth + 1
	j := jen.Type().Id(structName)

	nexts := make([]jen.Code, 0)

	j.StructFunc(func(g *jen.Group) {
		for i := 0; i < structT.NumFields(); i++ {
			fieldI := structT.Field(i)
			xType := common.TypeOf(fieldI.Type())
			tags := reflect.StructTag(structT.Tag(i))
			cg := common.GetCommentByTokenPos(pkg, fieldI.Pos()).List
			var cgs []string
			for _, v := range cg {
				cgs = append(cgs, v.Text)
			}
			comment := strings.Join(cgs, "\n")
			if !fieldI.Exported() {
				continue
			}

			if fieldI.Embedded() {
				g.Comment(comment).Line().Id(xType.TypeAsJen().GoString())
				continue
			}
			_, isSerializer := tags.Lookup("serializer")
			depthType := s.depthType(fieldI.Type())
			if isSerializer || common.TypeOf(depthType).Basic {
				g.Comment(comment).Line().Id(fieldI.Name()).Id(xType.TypeAsJen().GoString()).Tag(map[string]string{
					"json": lo.Ternary(lo.IsNotEmpty(tags.Get("json")), tags.Get("json"), common.DownFirst(fieldI.Name())),
				})
				continue
			}

			var isPreload bool
			var nextPreload []string
			var nextStructName = structName + common.UpFirst(fieldI.Name())

			lo.ForEach(preload, func(item string, index int) {
				if lo.FirstOrEmpty(strings.Split(item, ".")) == fieldI.Name() {
					isPreload = true
					nextPreload = strings.Split(item, ".")[1:]
				}
			})

			if !isPreload {
				continue
			}

			if xType.Struct {
				nexts = append(nexts, s.copyGormStruct(nextStructName, pkg, xType.StructType, nextPreload).Line())
				g.Comment(comment).Line().Id(fieldI.Name()).Id(nextStructName).Tag(map[string]string{
					"json": lo.Ternary(lo.IsNotEmpty(tags.Get("json")), tags.Get("json"), common.DownFirst(fieldI.Name())),
				})
				continue
			}

			if xType.List && xType.ListInner.Struct {
				nexts = append(nexts, s.copyGormStruct(nextStructName, pkg, xType.ListInner.StructType, nextPreload).Line())
				g.Comment(comment).Line().Id(fieldI.Name()).Index().Id(nextStructName).Tag(map[string]string{
					"json": lo.Ternary(lo.IsNotEmpty(tags.Get("json")), tags.Get("json"), common.DownFirst(fieldI.Name())),
				})
				continue
			}

		}
	})

	return j.Line().Add(nexts...)

}

func (s *HttpCrud) genImpl(i gen.StructGoTypeMeta) (err error) {
	modelId, modelHas := i.Doc.ByFuncNameAndArgName("@crud", "model")
	idName, idHas := i.Doc.ByFuncNameAndArgName("@crud", "idName")
	idType, idTypeHas := i.Doc.ByFuncNameAndArgName("@crud", "idType")
	preload, _ := i.Doc.ByFuncNameAndArgName("@crud", "preload")
	if !idHas || !modelHas || !idTypeHas {
		err = fmt.Errorf("@crud model,id,idType must be set")
		return
	}

	modelName := strings.ReplaceAll(modelId, ".", "")

	input := CrudTemplateInput{
		InterfaceName:       modelName + "HttpCrudBaseImpl",
		StructName:          i.Name,
		ModelName:           modelId,
		IdName:              idName,
		IdType:              idType,
		GetResponseStruct:   s.genGormStruct("GetResponse", modelId, strings.Split(preload, ",")).GoString(),
		CreateRequestStruct: s.genGormStruct("CreateRequest", modelId, strings.Split("", ",")).GoString(),
		UpdateBodyStruct:    s.genGormStruct("UpdateBody", modelId, strings.Split("", ",")).GoString(),
	}

	/* 	for _, v := range s.Option.Pkg.TypesInfo.Types {
		if lo.LastOrEmpty(strings.Split(v.Type.String(), "/")) == modelId {
			fmt.Println(v.Type.String())
		}
	} */

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
