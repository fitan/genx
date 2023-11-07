package gen

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/fitan/genx/common"
	"go/ast"
	"go/types"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"strings"
	"time"
)

type TypeName int

type X struct {
	Option Option
	Metas  Metas
	Plugs  Plugs
}

type Metas struct {
	Impl     ImplMeta
	Type     TypeMeta
	TypeSpec TypeSpecMeta
	Struct   StructMeta
	Func     FuncMeta
	Call     CallMeta
}

type Plugs struct {
	Impl     []InterfacePlugImpl
	Type     []TypePlugImpl
	TypeSpec []TypeSpecPlugImpl
	Struct   []StructPlugImpl
	Func     []FuncPlugImpl
	Call     []CallPlugImpl
}

func (x *X) RegCall(plug CallPlugImpl) {
	x.Plugs.Call = append(x.Plugs.Call, plug)
}

func (x *X) RegTypeSpec(plug TypeSpecPlugImpl) {
	x.Plugs.TypeSpec = append(x.Plugs.TypeSpec, plug)
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

func (x *X) RegFunc(plug FuncPlugImpl) {
	x.Plugs.Func = append(x.Plugs.Func, plug)
}

func (x *X) Gen() {
	x.parse()
	x.typeGen()
	x.implGen()
	x.structGen()
	x.typeSpecGen()
	x.funcGen()
	x.callGen()
	return
}

func (x *X) parse() {
	for _, v := range x.Option.Pkg.Syntax {
		astutil.Apply(v, func(c *astutil.Cursor) bool {
			switch t := c.Node().(type) {
			case *ast.CallExpr:
				callDoc := common.GetCommentByTokenPos(x.Option.Pkg, t.Pos())
				if callDoc == nil {
					return false
				}

				doc, err := common.ParseDoc(callDoc.Text())
				if err != nil {
					position := x.Option.Pkg.Fset.Position(t.Pos())
					slog.Error("parse doc error", "err", err, "fileName", position.Filename, "line", position.Line)
					panic(err)
				}

				for _, line := range doc {
					var call CallGoTypeMeta
					call.Name = t.Fun.(*ast.Ident).Name
					call.Doc = doc
					for _, param := range t.Args {
						call.Params = append(call.Params, common.TypeOf(x.Option.Pkg.TypesInfo.TypeOf(param)))
					}
					for _, param := range c.Parent().(*ast.AssignStmt).Lhs {
						call.Results = append(call.Results, common.TypeOf(x.Option.Pkg.TypesInfo.TypeOf(param)))
					}
					slog.Info("parse call", slog.String("name", line.UpFuncName()))
					x.Metas.Call.NameGoTypeMap[line.UpFuncName()] = append(x.Metas.Call.NameGoTypeMap[line.UpFuncName()], call)
				}

			case *ast.FuncDecl:
				if t.Doc.Text() == "" {
					return true
				}

				doc, err := common.ParseDoc(t.Doc.Text())
				if err != nil {
					slog.Error("parse doc error", err, slog.String("doc", t.Doc.Text()))
					panic(err)
				}
				for _, line := range doc {

					var fn FuncGoTypeMeta
					fn.Name = t.Name.Name
					fn.Doc = doc
					for _, param := range t.Type.Params.List {
						fn.Params = append(fn.Params, common.TypeOf(x.Option.Pkg.TypesInfo.TypeOf(param.Type)))
					}

					for _, param := range t.Type.Results.List {
						fn.Results = append(fn.Results, common.TypeOf(x.Option.Pkg.TypesInfo.TypeOf(param.Type)))
					}

					slog.Info("parse func", slog.String("name", line.UpFuncName()))
					x.Metas.Func.NameGoTypeMap[line.UpFuncName()] = append(x.Metas.Func.NameGoTypeMap[line.UpFuncName()], fn)
				}

			case *ast.ImportSpec:
				x.Option.Imports = append(x.Option.Imports, t)
			case *ast.GenDecl:
				if t.Doc.Text() == "" {
					return true
				}
				doc, err := common.ParseDoc(t.Doc.Text())
				if err != nil {
					slog.Error("parse doc error", err, slog.String("doc", t.Doc.Text()))
					panic(err)
				}

				for _, typeSpec := range t.Specs {
					switch st := typeSpec.(type) {
					case *ast.TypeSpec:
						for _, line := range doc {
							slog.Info("parse type", slog.String("name", line.UpFuncName()))
							x.Metas.Type.NameGoTypeMap[line.UpFuncName()] = append(x.Metas.Type.NameGoTypeMap[line.UpFuncName()], TypeGoTypeMeta{
								Doc: doc,
								Obj: x.Option.Pkg.TypesInfo.TypeOf(st.Type),
							})

							slog.Info("parse typeSpec", slog.String("name", line.UpFuncName()))
							x.Metas.TypeSpec.NameGoTypeMap[line.UpFuncName()] = append(x.Metas.TypeSpec.NameGoTypeMap[line.UpFuncName()], TypeSpecGoTypeMeta{
								Doc: doc,
								Obj: st,
							})
						}

						switch st.Type.(type) {
						case *ast.InterfaceType:
							for _, line := range doc {
								slog.Info("parse impl", slog.String("name", line.UpFuncName()))
								x.Metas.Impl.NameGoTypeMap[line.UpFuncName()] = append(x.Metas.Impl.NameGoTypeMap[line.UpFuncName()], InterfaceGoTypeMeta{
									Name: st.Name.Name,
									Doc:  doc,
									Obj:  x.Option.Pkg.TypesInfo.TypeOf(st.Type).(*types.Interface),
								})
							}
						case *ast.StructType:
							for _, line := range doc {
								slog.Info("parse struct", slog.String("name", line.UpFuncName()))
								x.Metas.Struct.NameGoTypeMap[line.UpFuncName()] = append(x.Metas.Struct.NameGoTypeMap[line.UpFuncName()], StructGoTypeMeta{
									Name: st.Name.Name,
									Doc:  doc,
									Obj:  x.Option.Pkg.TypesInfo.TypeOf(st.Type).(*types.Struct),
								})
							}

						}
					}
				}
			default:
				return true
			}
			return true
		}, func(c *astutil.Cursor) bool { return true })
	}
}

func (x *X) implGen() {
	for _, v := range x.Plugs.Impl {
		metas, ok := x.Metas.Impl.NameGoTypeMap[strings.ToUpper(v.Name())]
		if ok {
			err := v.Gen(x.Option, metas)
			if err != nil {
				slog.Error("impl gen error", err, slog.String("name", strings.ToUpper(v.Name())))
			}
		}
	}
}

func (x *X) callGen() {
	for _, v := range x.Plugs.Call {
		metas, ok := x.Metas.Call.NameGoTypeMap[strings.ToUpper(v.Name())]
		if ok {
			slog.Info("call gen start", slog.String("name", strings.ToUpper(v.Name())))
			timeStart := time.Now()
			err := v.Gen(x.Option, metas)
			slog.Info("call gen end", slog.String("name", strings.ToUpper(v.Name())), slog.Duration("time", time.Since(timeStart)))
			if err != nil {
				slog.Error("call gen error", err, slog.String("name", strings.ToUpper(v.Name())))
			}
		}
	}
}

func (x *X) typeGen() {
	for _, v := range x.Plugs.Type {
		metas, ok := x.Metas.Type.NameGoTypeMap[strings.ToUpper(v.Name())]
		if ok {
			err := v.Gen(x.Option, metas)
			if err != nil {
				slog.Error("type gen error", err, slog.String("name", strings.ToUpper(v.Name())))
			}
		}
	}
}

func (x *X) structGen() {
	for _, v := range x.Plugs.Struct {
		metas, ok := x.Metas.Struct.NameGoTypeMap[strings.ToUpper(v.Name())]
		if ok {
			err := v.Gen(x.Option, metas)
			if err != nil {
				slog.Error("struct gen error", err, slog.String("name", strings.ToUpper(v.Name())))
			}
		}
	}
}

func (x *X) funcGen() {
	for _, v := range x.Plugs.Func {
		metas, ok := x.Metas.Func.NameGoTypeMap[strings.ToUpper(v.Name())]
		if ok {
			slog.Info("func gen", slog.String("name", strings.ToUpper(v.Name())))
			spew.Dump(metas)
			err := v.Gen(x.Option, metas)
			if err != nil {
				slog.Error("func gen error", err, slog.String("name", strings.ToUpper(v.Name())))
			}
		}
	}
}

func (x *X) typeSpecGen() {
	for _, v := range x.Plugs.TypeSpec {
		metas, ok := x.Metas.TypeSpec.NameGoTypeMap[strings.ToUpper(v.Name())]
		if ok {
			err := v.Gen(x.Option, metas)
			if err != nil {
				slog.Error("struct gen error", err, slog.String("name", strings.ToUpper(v.Name())))
			}
		}
	}
}

func NewXByPkg(p *packages.Package) (*X, error) {
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
			Struct: StructMeta{
				NameGoTypeMap: make(map[string][]StructGoTypeMeta),
			},
			TypeSpec: TypeSpecMeta{
				NameGoTypeMap: make(map[string][]TypeSpecGoTypeMeta),
			},
			Func: FuncMeta{
				NameGoTypeMap: make(map[string][]FuncGoTypeMeta),
			},
			Call: CallMeta{
				NameGoTypeMap: make(map[string][]CallGoTypeMeta),
			},
		},
		Plugs: Plugs{
			Impl:     make([]InterfacePlugImpl, 0),
			Type:     make([]TypePlugImpl, 0),
			TypeSpec: make([]TypeSpecPlugImpl, 0),
			Struct:   make([]StructPlugImpl, 0),
			Func:     make([]FuncPlugImpl, 0),
			Call:     make([]CallPlugImpl, 0),
		},
	}, nil
}

func NewX(dir string) (*X, error) {
	p, err := common.LoadPkg(dir)
	if err != nil {
		return nil, err
	}
	return NewXByPkg(p)
}

type Option struct {
	Pkg             *packages.Package
	Imports         []*ast.ImportSpec
	MainExtraImport [][]string
}
