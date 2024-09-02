package mapstruct

import (
	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
	"github.com/pkg/errors"
)

type Plug struct {
}

func (s Plug) Name() string {
	return "@copy"
}

func (s Plug) Gen(option gen.Option, callGoTypeMetes []gen.CallGoTypeMeta) error {
	j := jen.NewFile(option.Pkg.Name)
	for _, v := range callGoTypeMetes {
		if len(v.Params) != 1 || len(v.Results) != 1 {
			return errors.New("copy: params must be one and results must be one")
		}

		srcType := v.Params[0]
		destType := v.Results[0]
		pkg := option.Pkg
		objName := v.Name + "Copy"

		// refT := common.TypesType2ReflectType(destType.T)
		// tSchema, err := schema.Parse(reflect.New(refT).Elem().Interface(), &sync.Map{}, schema.NamingStrategy{})
		// if err != nil {
		// 	return errors.Wrap(err, "schema.Parse")
		// }
		// spew.Dump(tSchema.Fields)
		// jsonb, err := json.Marshal(reflect.New(refT).Elem().Interface())
		// if err != nil {
		// 	return err
		// }
		// slog.Info("destType 2 json: ", "json", string(jsonb))

		j.Func().Id(v.Name).Params(jen.Id("src").Add(srcType.TypeAsJenComparePkgName(pkg))).Params(jen.Id("dest").Add(destType.TypeAsJenComparePkgName(pkg))).Block(
			jen.Id("dest").Op("=").Id(objName).Block().Dot("Copy").Call(jen.Id("src")),
			jen.Return(),
		)

		j.Type().Id(objName).Struct()

		cp := Copy{
			Pkg:            option.Pkg,
			StructName:     objName,
			JenF:           j,
			Recorder:       NewRecorder(),
			SrcParentPath:  []string{},
			SrcPath:        []string{},
			Src:            NewDataFieldMap(option.Pkg, []string{}, []string{}, "src", srcType),
			DestParentPath: []string{},
			DestPath:       []string{},
			Dest:           NewDataFieldMap(option.Pkg, []string{}, []string{}, "dest", destType),
			DefaultFn: jen.Func().Params(jen.Id("d").Id(objName)).
				Id("Copy").Params(jen.Id("src").Add(srcType.TypeAsJenComparePkgName(pkg))).Params(jen.Id("dest").Add(destType.TypeAsJenComparePkgName(pkg)))}
		cp.Gen()

	}

	common.WriteGO("copy.go", j.GoString())

	return nil
}
