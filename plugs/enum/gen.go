package enum

import (
	"strconv"
	"strings"

	"github.com/fitan/genx/common"
	"github.com/fitan/genx/parser"
	"github.com/fitan/jennifer/jen"
)

type Enum struct {
	Param    []parser.FuncArg
	Args     []Arg
	TypeName string
}

type Arg struct {
	Key   string
	Value string
}

//func (e *Enum) ToJenFString() (s string) {
//	state := jen.Statement{}
//	for _, f := range e.Pkg.Syntax {
//		tsList, gdDoc := typeSpecs(f)
//
//		for tsIndex, ts := range tsList {
//			if i, ok := ts.Type.(*ast.Ident); ok {
//				if i.String() == "int" {
//					text := gdDoc[tsIndex].Text()
//					if strings.Contains(text, "@enum") {
//						parseEnum, err := parse.Parse(text)
//						if err != nil {
//							panic(err)
//						}
//						enum := &Enum{
//							TypeName: ts.Name.String(),
//							Doc:      parseEnum,
//						}
//						codes, err := enum.Gen()
//						if err != nil {
//							panic(err)
//						}
//						jenF.Add(codes...)
//						state.Add(codes...)
//					}
//				}
//			}
//		}
//	}
//
//	return state.GoString()
//}

func (e *Enum) Init() (err error) {

	args := make([]Arg, 0, 0)
	for _, arg := range e.Param {

		argSplit := strings.Split(arg.Value, ":")
		var key, value string
		if len(argSplit) == 1 {
			key = argSplit[1]
			value = ""
		} else if len(argSplit) >= 2 {
			key = argSplit[0]
			value = argSplit[1]
		} else {
			return common.ValidationError("invalid enum argument format").
				WithPlugin("@enum").
				WithExtra("argument", arg.Value).
				WithDetails("enum argument must be in format 'key:value' or 'key'. Example: @enum active:1 inactive:0").
				Build()
		}
		args = append(args, Arg{
			Key:   key,
			Value: value,
		})
	}
	e.Args = args
	return nil
}

func (e *Enum) Gen() (codes []jen.Code, err error) {
	codes = make([]jen.Code, 0, 100)
	err = e.Init()
	if err != nil {
		return
	}
	if len(e.Args) == 0 {
		return
	}
	codes = append(codes, e.Const()...)
	codes = append(codes, e.Value()...)
	//codes = append(codes, e.Json()...)
	//codes = append(codes, e.GormSerialize()...)
	codes = append(codes, e.Remark()...)
	codes = append(codes, e.String()...)
	return
}

func (e *Enum) Const() (codes []jen.Code) {
	codes = append(codes, jen.Id("_").Op("=").Id("iota"))
	for _, arg := range e.Args {
		codes = append(codes, jen.Id(strings.ToTitle(e.TypeName+"_"+arg.Key)))
	}

	aliasCode := jen.Const().DefsFunc(func(g *jen.Group) {
		for _, arg := range e.Args {
			g.Id(strings.ToUpper(e.TypeName + "_" + arg.Key + "alias")).Op("=").Lit(arg.Key)
		}
	}).Line().Line()

	remarkCode := jen.Const().DefsFunc(func(g *jen.Group) {
		for _, arg := range e.Args {
			g.Id(strings.ToUpper(e.TypeName + "_" + arg.Key + "remark")).Op("=").Lit(arg.Value)
		}
	}).Line().Line()

	return []jen.Code{jen.Const().Defs(codes...).Line().Line(), aliasCode, remarkCode}
}

func (e *Enum) Value() (codes []jen.Code) {
	varCode := jen.Var().Id("_" + e.TypeName + "Value").Op("=").Map(jen.Int()).Id(e.TypeName).Values(jen.DictFunc(func(d jen.Dict) {
		for index, arg := range e.Args {
			d[jen.Id(strconv.Itoa(index+1))] = jen.Id(strings.ToTitle(e.TypeName + "_" + arg.Key))
		}
	})).Line().Line()

	parseCode := jen.Func().Id("Parse"+e.TypeName).Params(jen.Id("id").Int()).Params(jen.Id(e.TypeName), jen.Error()).Block(
		jen.If(jen.Id("x").Op(",").Id("ok").Op(":=").Id("_"+e.TypeName+"Value").Index(jen.Id("id")), jen.Id("ok")).Block(
			jen.Return(jen.Id("x"), jen.Nil()),
		),
		jen.Return(jen.Lit(0), jen.Qual("fmt", "Errorf").Call(jen.Lit("unknown enum value: %s"), jen.Id("id"))),
	).Line().Line()

	codes = append(codes, varCode, parseCode)
	return
}

func (e *Enum) Remark() (code []jen.Code) {
	stringCode := jen.Func().Params(jen.Id("e").Id(e.TypeName)).Id("Remark").Params().String().Block(
		jen.Switch(jen.Id("e")).BlockFunc(func(g *jen.Group) {
			for argIndex, arg := range e.Args {
				g.Case(jen.Lit(argIndex + 1)).Block(
					jen.Return(jen.Id(strings.ToUpper(e.TypeName + "_" + arg.Key + "remark"))),
				)
			}
		}),
		jen.Return(jen.Qual("fmt", "Sprintf").Call(jen.Lit("unknown %d"), jen.Id("e"))),
	).Line().Line()
	code = append(code, stringCode)
	return
}

func (e *Enum) String() (code []jen.Code) {
	stringCode := jen.Func().Params(jen.Id("e").Id(e.TypeName)).Id("String").Params().String().Block(
		jen.Switch(jen.Id("e")).BlockFunc(func(g *jen.Group) {
			for argIndex, arg := range e.Args {
				g.Case(jen.Lit(argIndex + 1)).Block(
					jen.Return(jen.Id(strings.ToUpper(e.TypeName + "_" + arg.Key + "alias"))),
				)
			}
		}),
		jen.Return(jen.Qual("fmt", "Sprintf").Call(jen.Lit("unknown %d"), jen.Id("e"))),
	).Line().Line()
	code = append(code, stringCode)
	return
}

func (e *Enum) Json() (codes []jen.Code) {
	marshalCode := jen.Func().Params(jen.Id("e").Id("*"+e.TypeName)).Id("MarshalJSON").Params().Params(jen.Id("[]byte"), jen.Id("error")).Block(
		jen.Switch(jen.Id("*e")).BlockFunc(func(g *jen.Group) {
			for index, arg := range e.Args {
				g.Case(jen.Id(strings.ToTitle(e.TypeName + "_" + arg.Key))).Block(
					jen.Return(jen.Index().Byte().Call(jen.Id("`"+strconv.Itoa(index+1)+"`")), jen.Nil()),
				)
			}
		}),
		jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("unknown enum value: %v"), jen.Id("e"))),
	).Line().Line()

	unmarshalCode := jen.Func().Params(jen.Id("e").Id("*"+e.TypeName)).Id("UnmarshalJSON").Params(jen.Id("data").Index().Byte()).Params(jen.Id("error")).Block(
		jen.List(jen.Id("dataInt"), jen.Id("err")).Op(":=").Id("strconv.Atoi").Call(jen.Id("string").Call(jen.Id("data"))),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Qual("github.com/pkg/errors", "Wrap").Call(jen.Id("err"), jen.Lit("strconv.Atoi"))),
		),
		jen.List(jen.Id("v"), jen.Id("err")).Op(":=").Id("Parse"+e.TypeName).Call(jen.Id("dataInt")),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Id("err")),
		),
		jen.Id("*e").Op("=").Id("v"),
		jen.Return(jen.Nil()),
	).Line().Line()

	codes = append(codes, marshalCode, unmarshalCode)
	return
}

func (e *Enum) GormSerialize() (codes []jen.Code) {
	scanCode := jen.Func().Params(jen.Id("e").Id("*"+e.TypeName)).Id("Scan").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("field").Op("*").Qual("gorm.io/gorm/schema", "Field"),
			jen.Id("dst").Qual("reflect", "Value"),
			jen.Id("dbValue").Interface(),
		).Params(jen.Id("error")).Block(
		jen.Switch(jen.Id("value").Op(":=").Id("dbValue").Assert(jen.Id("type")).Block(
			jen.Case(jen.String()).Block(
				jen.Id("*e").Op("=").Id("_"+e.TypeName+"Value").Index(jen.Id("value")),
			),
			jen.Default().Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("unknown enum value: %v"), jen.Id("value"))),
			),
		),

			jen.Return(jen.Nil()),
		)).Line().Line()
	valueCode := jen.Func().Params(jen.Id("e").Id(e.TypeName)).Id("Value").Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("field").Op("*").Qual("gorm.io/gorm/schema", "Field"),
		jen.Id("dst").Qual("reflect", "Value"),
		jen.Id("fieldValue").Interface(),
	).Params(jen.Id("driver.Value"), jen.Id("error")).Block(
		jen.Return(jen.Id("e").Dot("String").Call(), jen.Nil()),
	).Line().Line()

	codes = append(codes, scanCode, valueCode)
	return
}
