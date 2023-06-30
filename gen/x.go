package gen

import (
	"github.com/fitan/genx/common"
	"go/ast"
	"go/types"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/go/packages"
)

type TypeName int

type X struct {
	Option Option
	Metas  Metas
	Plugs  Plugs
}

type Metas struct {
	Impl   ImplMeta
	Type   TypeMeta
	Struct StructMeta
}

type Plugs struct {
	Impl   []InterfacePlugImpl
	Type   []TypePlugImpl
	Struct []StructPlugImpl
}

func (x *X) RegImpl(plug InterfacePlugImpl) {
	x.Plugs.Impl = append(x.Plugs.Impl, plug)
}

func (x *X) RegType(plug TypePlugImpl) {
	x.Plugs.Type = append(x.Plugs.Type, plug)
}

func (x *X) RegStruct(plug StructPlugImpl) {
	x.Plugs.Struct = append(x.Plugs.Struct, plug)
}

func (x *X) Gen() {
	x.parse()
	slog.Info("parse success", slog.Any("type", x.Metas.Type.NameGoTypeMap), slog.Any("impl", x.Metas.Impl.NameGoTypeMap))
	x.typeGen()
	x.implGen()
	x.structGen()
	return
}

func (x *X) parse() {
	for _, v := range x.Option.Pkg.Syntax {
		ast.Inspect(v, func(node ast.Node) bool {
			switch t := node.(type) {
			case *ast.ImportSpec:
				x.Option.Imports = append(x.Option.Imports, t)
			case *ast.GenDecl:
				if t.Doc.Text() == "" {
					return true
				}
				doc, err := common.ParseDoc(t.Doc.Text())
				if err != nil {
					slog.Error("parse doc error", err, slog.String("doc", v.Doc.Text()))
					panic(err)
				}

				for _, typeSpec := range t.Specs {
					switch st := typeSpec.(type) {
					case *ast.TypeSpec:
						for _, line := range doc.Funcs {
							if line.Func != nil {
								x.Metas.Type.NameGoTypeMap[line.Func.FuncName] = append(x.Metas.Type.NameGoTypeMap[line.Func.FuncName], TypeGoTypeMeta{
									Doc: doc,
									Obj: x.Option.Pkg.TypesInfo.TypeOf(st.Type),
								})
							}
						}

						switch st.Type.(type) {
						case *ast.InterfaceType:
							for _, line := range doc.Funcs {
								if line.Func != nil {
									x.Metas.Impl.NameGoTypeMap[line.Func.FuncName] = append(x.Metas.Impl.NameGoTypeMap[line.Func.FuncName], InterfaceGoTypeMeta{
										Doc: doc,
										Obj: x.Option.Pkg.TypesInfo.TypeOf(st.Type).(*types.Interface),
									})
								}
							}
						case *ast.StructType:
							for _, line := range doc.Funcs {
								if line.Func != nil {
									x.Metas.Struct.NameGoTypeMap[line.Func.FuncName] = append(x.Metas.Struct.NameGoTypeMap[line.Func.FuncName], StructGoTypeMeta{
										Doc: doc,
										Obj: x.Option.Pkg.TypesInfo.TypeOf(st.Type).(*types.Struct),
									})
								}
							}

						}
					}
				}
			default:
				return true
			}
			return true
		})
	}
}

func (x *X) implGen() {
	for _, v := range x.Plugs.Impl {
		metas, ok := x.Metas.Impl.NameGoTypeMap[v.Name()]
		if ok {
			err := v.Gen(x.Option, metas)
			if err != nil {
				slog.Error("impl gen error", err, slog.String("name", v.Name()))
			}
		}
	}
}

func (x *X) typeGen() {
	for _, v := range x.Plugs.Type {
		metas, ok := x.Metas.Type.NameGoTypeMap[v.Name()]
		if ok {
			err := v.Gen(x.Option, metas)
			if err != nil {
				slog.Error("type gen error", err, slog.String("name", v.Name()))
			}
		}
	}
}

func (x *X) structGen() {
	for _, v := range x.Plugs.Struct {
		metas, ok := x.Metas.Struct.NameGoTypeMap[v.Name()]
		if ok {
			err := v.Gen(x.Option, metas)
			if err != nil {
				slog.Error("struct gen error", err, slog.String("name", v.Name()))
			}
		}
	}
}

func NewX(dir string) (*X, error) {
	p, err := common.LoadPkg(dir)
	if err != nil {
		return nil, err
	}
	return &X{
		Option: Option{
			Pkg:             p,
			Imports:         make([]*ast.ImportSpec, 0),
			MainExtraImport: make([][]string, 0),
		},
		Metas: Metas{
			Impl: ImplMeta{
				NameGoTypeMap: make(map[string][]InterfaceGoTypeMeta),
			},
			Type: TypeMeta{
				NameGoTypeMap: make(map[string][]TypeGoTypeMeta),
			},
		},
		Plugs: Plugs{
			Impl: make([]InterfacePlugImpl, 0),
			Type: make([]TypePlugImpl, 0),
		},
	}, nil
}

type Option struct {
	Pkg             *packages.Package
	Imports         []*ast.ImportSpec
	MainExtraImport [][]string
}
