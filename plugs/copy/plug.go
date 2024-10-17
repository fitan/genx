package mapstruct

import (
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
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
	current := map[string]struct{}{}
	j := jen.NewFile(option.Pkg.Name)
	for _, v := range callGoTypeMetes {
		if _, ok := current[v.Name]; ok {
			continue
		} else {
			current[v.Name] = struct{}{}
		}
		if len(v.Params) != 2 {
			spew.Dump(v.Params)
			panic("copy: params must be two")
			return errors.New("copy: params must be two")
		}

		destType := v.Params[0]
		srcType := v.Params[1]

		if !destType.Pointer {
			panic("copy: params must pointer")
			return errors.New("copy: params must pointer")
		}

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

		j.Func().Id(v.Name).Params(jen.Id("dest").Add(destType.TypeAsJenComparePkgName(pkg)), jen.Id("src").Add(srcType.TypeAsJenComparePkgName(pkg))).Block(
			jen.Id(objName).Block().Dot("Copy").Call(jen.Id("dest"), jen.Id("src")),
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
			Head:           true,
			DefaultFn: jen.Func().Params(jen.Id("d").Id(objName)).
				Id("Copy").Params(jen.Id("dest").Add(destType.TypeAsJenComparePkgName(pkg)), jen.Id("src").Add(srcType.TypeAsJenComparePkgName(pkg)))}
		cp.Gen()

	}

	common.WriteGO(filepath.Join(option.Dir, "copy.go"), j.GoString())

	return nil
}
