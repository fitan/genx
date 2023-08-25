package mapstruct

import (
	"encoding/json"
	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
	"reflect"
)

type Plug struct {
}

func (s Plug) Name() string {
	return "@MapStruct"
}

func (s Plug) Gen(option gen.Option, callGoTypeMetes []gen.CallGoTypeMeta) error {
	j := jen.NewFile(option.Pkg.Name)
	for _, v := range callGoTypeMetes {
		if len(v.Params) != 1 || len(v.Results) != 1 {
			return errors.New("MapStruct: params must be one and results must be one")
		}

		srcType := v.Params[0]
		destType := v.Results[0]
		pkg := option.Pkg
		objName := v.Name + "MapStruct"

		refT := common.TypesType2ReflectType(destType.T)
		slog.Info("spew.Dump")
		slog.Info("reflect. type .", reflect.TypeOf(refT).String())
		slog.Info("ref valueof", reflect.ValueOf(refT).Interface())
		jsonb, err := json.Marshal(reflect.ValueOf(refT).Interface())
		if err != nil {
			return err
		}
		slog.Info("destType 2 json: ", "json", string(jsonb))

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
			Src:            NewDataFieldMap(option.Pkg, []string{}, "src", srcType),
			DestParentPath: []string{},
			DestPath:       []string{},
			Dest:           NewDataFieldMap(option.Pkg, []string{}, "dest", destType),
			DefaultFn: jen.Func().Params(jen.Id("d").Id(objName)).
				Id("Copy").Params(jen.Id("src").Add(srcType.TypeAsJenComparePkgName(pkg))).Params(jen.Id("dest").Add(destType.TypeAsJenComparePkgName(pkg)))}
		cp.Gen()
	}

	common.WriteGO("map_struct.go", j.GoString())

	return nil
}
