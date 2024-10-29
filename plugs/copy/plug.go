package mapstruct

import (
	"path/filepath"

	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
	"github.com/pkg/errors"
)

type Plug struct {
}

func (s Plug) Name() string {
	return "@copy"
}

func (s Plug) Gen(option gen.Option, callGoTypeMetes []gen.CallGoTypeMeta) (res []gen.GenResult, err error) {
	current := map[string]struct{}{}
	j := jen.NewFile(option.Pkg.Name)
	for _, v := range callGoTypeMetes {
		if _, ok := current[v.Name]; ok {
			continue
		} else {
			current[v.Name] = struct{}{}
		}
		if len(v.Params) != 2 {
			return nil, errors.New("copy: params must be two")
		}

		destType := v.Params[0]
		srcType := v.Params[1]

		if !destType.Pointer {
			return nil, errors.New("copy: params must pointer")
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

	res = append(res, gen.GenResult{
		FileName: filepath.Join(option.Dir, "copy.go"),
		FileStr:  j.GoString(),
		Cover:    true,
	})

	return
}
