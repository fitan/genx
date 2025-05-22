package kithttp

import (
	"strings"

	"github.com/fitan/genx/common"
	"github.com/fitan/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
	"golang.org/x/tools/go/packages"
)

func NewType2Ast(pkg *packages.Package) *Type2ast {
	return &Type2ast{
		current: make(map[string]struct{}),
		pkg:     pkg,
	}
}

type Type2ast struct {
	pkg     *packages.Package
	current map[string]struct{}
}

func (t *Type2ast) xtypeParse(codes *[]*jen.Statement, names []string, xt *common.Type) *jen.Statement {
	if xt.Basic {
		return jen.Id(xt.BasicType.String())
	}

	if xt.Map {
		return jen.Map(t.xtypeParse(codes, names, xt.MapKey)).Add(t.xtypeParse(codes, names, xt.MapValue))
	}

	if xt.List {
		return jen.Index().Add(t.xtypeParse(codes, names, xt.ListInner))
	}

	if xt.Interface {
		return jen.Interface()
	}

	if xt.Pointer {
		return jen.Op("*").Add(t.xtypeParse(codes, names, xt.PointerInner))
	}

	var typeName string

	if xt.Named {
		// 说明是go的基本类型
		if !strings.Contains(xt.T.String(), "/") {
			return jen.Id(xt.T.String())
		} else {
			xtSplit := strings.Split(xt.T.String(), "/")
			lastName, _ := lo.Last(xtSplit)
			typeName = strcase.ToCamel(lastName)
			// typeName = strings.ReplaceAll(lastName, ".", "")
			// typeName = strings.ReplaceAll(typeName, "_", "")
			// typeName = strings.ReplaceAll(strings.TrimPrefix(xt.T.String(), t.pkg.PkgPath), ",", "")
			// fmt.Println("named:", xt.T.String())
			// fmt.Println("typeName:", typeName)
			// if !strings.Contains(xt.T.String(), "/") {
			// 	// spew.Dump(xt.NamedType.Obj())
			// 	return jen.Id(xt.NamedType.Obj().Pkg().Name()).Dot(xt.NamedType.Obj().Name())
			// }
		}
	}

	if xt.Struct {
		// structNewName := strings.Join(names, "")
		if _, ok := t.current[typeName]; !ok {
			t.current[typeName] = struct{}{}
			*codes = append(*codes, jen.Type().Id(typeName).StructFunc(func(g *jen.Group) {
				for i := 0; i < xt.StructType.NumFields(); i++ {
					var newNames []string
					field := xt.StructType.Field(i)
					xtField := common.TypeOf(field.Type())
					if xtField.Struct || xtField.Pointer || xtField.List || xtField.Map {
						newNames = append(names, field.Name())
					} else {
						newNames = names
					}

					tagMap := ParseTagIntoMap(xt.StructType.Tag(i))

					if !CheckTag(tagMap) {
						continue
					}

					if t, ok := CheckTagSwaggertype(tagMap); ok {
						g.Id(field.Name()).Add(jen.Id(t).Tag(tagMap))
						continue
					}

					g.Id(field.Name()).Add(t.xtypeParse(codes, newNames, xtField)).Tag(tagMap)
				}
			}))
		}

		return jen.Id(typeName)
	}

	return nil
}

func ParseTagIntoMap(tag string) map[string]string {
	result := make(map[string]string)
	tags := strings.Split(tag, " ")
	for _, t := range tags {
		kv := strings.Split(t, ":")
		if len(kv) == 2 {
			// 去除引号
			value := strings.Trim(kv[1], "`\"")
			result[kv[0]] = value
		}
	}
	return result
}

func CheckTag(tags map[string]string) bool {
	if _, ok := tags["param"]; !ok {
		return true
	}

	if strings.HasPrefix(tags["param"], `ctx,`) {
		return false
	}

	return true
}

func CheckTagSwaggertype(tags map[string]string) (string, bool) {
	t, ok := tags["swaggertype"]
	tItem := strings.Split(t, ",")
	if ok {
		switch len(tItem) {
		case 1:
			switch t {
			case "string":
				return "string", true
			case "integer":
				return "int", true
			}
		case 2:
			switch tItem[0] {
			case "array":
				switch tItem[1] {
				case "string":
					return "[]string", true
				case "integer":
					return "[]int", true
				}
			}
		}
	}
	return t, ok
}

func (t *Type2ast) Parse(xt *common.Type, name string) string {
	if _, ok := t.current[name]; ok {
		return ""
	}

	t.current[xt.T.String()] = struct{}{}

	codes := make([]*jen.Statement, 0)
	code := jen.Type().Id(name).StructFunc(func(g *jen.Group) {
		if xt.Struct {
			for i := 0; i < xt.StructType.NumFields(); i++ {
				field := xt.StructType.Field(i)
				xtField := common.TypeOf(field.Type())
				tag := xt.StructType.Tag(i)
				tagMap := ParseTagIntoMap(tag)
				if !CheckTag(tagMap) {
					continue
				}

				if t, ok := CheckTagSwaggertype(tagMap); ok {
					g.Id(field.Name()).Add(jen.Id(t).Tag(tagMap))
					continue
				}

				names := []string{name, field.Name()}
				g.Id(field.Name()).Add(t.xtypeParse(&codes, names, xtField)).Tag(tagMap)
			}
		}
	})
	// code = code.Type().Id(name).Add(k.xtypeParse(code, []string{name}, xt))
	for _, c := range codes {
		code.Line().Add(c)
	}
	return code.GoString()

}
