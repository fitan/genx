package gen

import (
	"embed"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/fitan/genx/common"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

type TypeName int

type X struct {
	Option Option
	Metas  Metas
	Plugs  Plugs
	WG     *conc.WaitGroup
	TUI    *Model
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
	x.typeGen()
	x.implGen()
	x.structGen()
	x.typeSpecGen()
	x.funcGen()
	x.callGen()
	x.WG.Wait()

	x.TUI.PkgEnd(UpdateTreeReq{
		PkgName: x.Option.Pkg.PkgPath,
	})
	return
}

func (x *X) Parse() {
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
					parentAssignStmt, ok := c.Parent().(*ast.AssignStmt)
					if ok {
						for _, param := range parentAssignStmt.Lhs {
							call.Results = append(call.Results, common.TypeOf(x.Option.Pkg.TypesInfo.TypeOf(param)))
						}
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
					if t.Type.Params != nil {
						for _, param := range t.Type.Params.List {
							fn.Params = append(fn.Params, common.TypeOf(x.Option.Pkg.TypesInfo.TypeOf(param.Type)))
						}
					}

					if t.Type.Results != nil {
						for _, param := range t.Type.Results.List {
							fn.Results = append(fn.Results, common.TypeOf(x.Option.Pkg.TypesInfo.TypeOf(param.Type)))
						}
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
									Name:   st.Name.Name,
									Doc:    doc,
									RawDoc: t.Doc,
									Obj:    x.Option.Pkg.TypesInfo.TypeOf(st.Type).(*types.Interface),
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

func (x *X) implByName(name string) ([]InterfaceGoTypeMeta, bool) {
	meta, ok := x.Metas.Impl.NameGoTypeMap[strings.ToUpper(name)]
	return meta, ok
}

func (x *X) implGen() {
	for _, v := range x.Plugs.Impl {
		metas, ok := x.implByName(v.Name())
		if ok {
			modelName := v.Name()
			x.WG.Go(func() {
				x.UpdateTUI(modelName, func() (gens []GenResult, err error) { return v.Gen(x.Option, metas) })
			})
		}
	}
}

func (x *X) callByName(name string) ([]CallGoTypeMeta, bool) {
	meta, ok := x.Metas.Call.NameGoTypeMap[strings.ToUpper(name)]
	return meta, ok
}

func (x *X) callGen() {
	for _, v := range x.Plugs.Call {
		metas, ok := x.callByName(v.Name())
		if ok {
			modelName := v.Name()
			x.WG.Go(func() {
				x.UpdateTUI(modelName, func() (gens []GenResult, err error) { return v.Gen(x.Option, metas) })
			})
		}
	}
}

func (x *X) typeByName(name string) ([]TypeGoTypeMeta, bool) {
	meta, ok := x.Metas.Type.NameGoTypeMap[strings.ToUpper(name)]
	return meta, ok
}

func (x *X) typeGen() {
	for _, v := range x.Plugs.Type {
		metas, ok := x.typeByName(v.Name())
		if ok {
			modelName := v.Name()
			x.WG.Go(func() {
				x.UpdateTUI(modelName, func() (gens []GenResult, err error) { return v.Gen(x.Option, metas) })
			})
		}
	}
}

func (x *X) structByName(name string) ([]StructGoTypeMeta, bool) {
	meta, ok := x.Metas.Struct.NameGoTypeMap[strings.ToUpper(name)]
	return meta, ok
}

func (x *X) structGen() {
	for _, v := range x.Plugs.Struct {
		metas, ok := x.structByName(v.Name())
		if ok {
			modelName := v.Name()
			x.WG.Go(func() {
				x.UpdateTUI(modelName, func() (gens []GenResult, err error) { return v.Gen(x.Option, metas) })
			})
		}
	}
}

func (x *X) funcByName(name string) ([]FuncGoTypeMeta, bool) {
	meta, ok := x.Metas.Func.NameGoTypeMap[strings.ToUpper(name)]
	return meta, ok
}

func (x *X) funcGen() {
	for _, v := range x.Plugs.Func {
		metas, ok := x.funcByName(v.Name())
		if ok {
			modelName := v.Name()
			x.WG.Go(func() {
				x.UpdateTUI(modelName, func() (gens []GenResult, err error) { return v.Gen(x.Option, metas) })
			})
		}
	}
}

func (x *X) typeSpecByName(name string) ([]TypeSpecGoTypeMeta, bool) {
	meta, ok := x.Metas.TypeSpec.NameGoTypeMap[strings.ToUpper(name)]
	return meta, ok
}

func (x *X) typeSpecGen() {
	for _, v := range x.Plugs.TypeSpec {
		metas, ok := x.typeSpecByName(v.Name())
		if ok {
			modelName := v.Name()
			x.WG.Go(func() {
				x.UpdateTUI(modelName, func() (gens []GenResult, err error) { return v.Gen(x.Option, metas) })
			})
		}
	}
}

func (x *X) UpdateTUI(plugName string, f func() (gens []GenResult, err error)) {

	x.TUI.PlugStart(UpdateTreeReq{
		PkgName:  x.Option.Pkg.PkgPath,
		PlugName: plugName,
	})

	gens, err := f()
	if err != nil {
		x.TUI.PlugEnd(UpdateTreeReq{
			PkgName:  x.Option.Pkg.PkgPath,
			PlugName: plugName,
			FileName: "",
			Status:   2,
			Err:      err.Error(),
		})
		return
	}

	gw := conc.NewWaitGroup()

	for _, gen := range gens {
		gw.Go(func() {
			x.TUI.FileStart(UpdateTreeReq{
				PkgName:  x.Option.Pkg.PkgPath,
				PlugName: plugName,
				FileName: gen.FileName,
				Status:   0,
				Err:      "",
			})

			cover, err := common.WriteGoWithOpt(gen.FileName, gen.FileStr, common.WriteOpt{
				Cover: gen.Cover,
			})

			x.TUI.FileEnd(UpdateTreeReq{
				PkgName:  x.Option.Pkg.PkgPath,
				PlugName: plugName,
				FileName: gen.FileName,
				Status:   lo.Ternary(cover, 1, 3),
				Err: lo.TernaryF(err != nil, func() string {
					return err.Error()
				}, func() string {
					return ""
				}),
			})

		})
	}

	gw.Wait()

	x.TUI.PlugEnd(UpdateTreeReq{
		PkgName:  x.Option.Pkg.PkgPath,
		PlugName: plugName,
		Status:   1,
	})

	return
}

func NewXByPkg(static embed.FS, p *packages.Package, tui *Model, config *Config, preload map[string][]*packages.Package) (*X, error) {
	x := &X{
		WG: conc.NewWaitGroup(),
		Option: Option{
			Static:          static,
			Pkg:             p,
			Dir:             filepath.Dir(p.GoFiles[0]),
			Imports:         make([]*ast.ImportSpec, 0),
			MainExtraImport: make([][]string, 0),
			Config:          config,
			PreloadPkg:      preload,
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
		TUI: tui,
	}

	x.Parse()
	return x, nil
}

func NewX(static embed.FS, dir string, tui *Model) (res []*X, err error) {
	config := findGenXConfig()

	preloadMap := make(map[string][]*packages.Package, 0)

	for _, preload := range config.Preloads {
		d, err := common.GetPkgAbsPath(preload.Path)
		if err != nil {
			err = fmt.Errorf("preload %s err: %s", preload.Path, err.Error())
			return nil, err
		}
		pkgs, err := common.LoadPkg(d)

		if err != nil {
			err = fmt.Errorf("preload %s err: %s", preload.Path, err.Error())
			return nil, err
		}

		if preload.Alias == "" {
			_, file := filepath.Split(preload.Path)
			preloadMap[file] = pkgs
		} else {
			preloadMap[preload.Alias] = pkgs
		}
	}

	ps, err := common.LoadPkg(dir)
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		x, err := NewXByPkg(static, p, tui, config, preloadMap)
		if err != nil {
			return nil, err
		}
		res = append(res, x)
	}

	return
}

type Option struct {
	Static          embed.FS
	Pkg             *packages.Package
	Dir             string
	Imports         []*ast.ImportSpec
	MainExtraImport [][]string
	Config          *Config
	PreloadPkg      map[string][]*packages.Package
}

func findGenXConfig() *Config {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	currentDir := ""
	tmpFilePath := ""

	for {
		tmpFilePath = filepath.Join(pwd, "genx.yaml")
		if currentDir == tmpFilePath {
			panic("genx.yaml not found")
		}
		if _, err := os.Stat(tmpFilePath); os.IsNotExist(err) {
			currentDir = tmpFilePath
			pwd = filepath.Dir(pwd)
			continue
		}
		break
	}

	k := koanf.New(".")

	err = k.Load(file.Provider(tmpFilePath), yaml.Parser())
	if err != nil {
		panic(err)
	}

	config := &Config{}
	k.Unmarshal("", config)
	return config
}
