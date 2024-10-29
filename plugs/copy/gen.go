package mapstruct

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"go/types"
	"strings"

	"log/slog"

	"github.com/fitan/genx/common"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/types/typeutil"

	"github.com/fitan/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

const CopyTag = "copy"

const CopyPrefixTag = "@copy-prefix"
const CopyNameTag = "@copy-name"
const CopyMustTag = "@copy-must"
const CopyTargetPathTag = "@copy-target-path"
const CopyTargetNameTag = "@copy-target-name"
const CopyTargetMethodTag = "@copy-target-method"
const CopySoucePathTag = "@copy-source-path"
const CopyAutoCastTag = "@copy-auto-cast"

type Field struct {
	Path []string
	// 别名Path
	AliasPath     []string
	EmbeddedIndex []int
	Name          string
	Type          *common.Type
	ParentDoc     common.Doc
	Doc           common.Doc
	IsNamed       bool

	// doc
	CopyPrefix       string
	CopyName         string
	CopyMust         bool
	CopyTargetPath   string
	CopyTargetName   string
	CopyTargetMethod string
	CopyAutoCast     bool
}

func (f *Field) ParseDoc() {
	var copyPrefix string
	var copyName string
	var copyMust bool
	var copyTargetPath string
	var copyTargetName string
	var copyTargetMethod string
	var copyAutoCast bool

	f.ParentDoc.ByFuncNameAndArgs(CopyPrefixTag, &copyPrefix)

	f.Doc.ByFuncNameAndArgs(CopyPrefixTag, &copyPrefix)
	f.Doc.ByFuncNameAndArgs(CopyNameTag, &copyName)
	copyMust = f.Doc.ByFuncNameAndArgs(CopyMustTag)
	f.Doc.ByFuncNameAndArgs(CopyTargetPathTag, &copyTargetPath)
	f.Doc.ByFuncNameAndArgs(CopyTargetNameTag, &copyTargetName)
	f.Doc.ByFuncNameAndArgs(CopyTargetMethodTag, &copyTargetMethod)
	copyAutoCast = f.Doc.ByFuncNameAndArgs(CopyAutoCastTag)

	f.CopyPrefix = copyPrefix
	f.CopyName = copyName
	f.CopyMust = copyMust
	f.CopyTargetPath = copyTargetPath
	f.CopyTargetName = copyTargetName
	f.CopyTargetMethod = copyTargetMethod
	f.CopyAutoCast = copyAutoCast
}

// func (f Field) CopyMethod() (s1, s2 string) {
// 	var formatPkgName, formatFnName string
// 	has := f.Doc.ByFuncNameAndArgs(structMapFormat, &formatPkgName, &formatFnName)
// 	if has {
// 		return formatPkgName, formatFnName
// 	}

// 	return
// }

func (f Field) SrcIdPath() *jen.Statement {
	path := append([]string{"src"}, f.Path...)
	return jen.Id(strings.Join(path, "."))
}

func (f Field) DestIdPath() *jen.Statement {
	path := append([]string{"dest"}, f.Path...)
	return jen.Id(strings.Join(path, "."))
}

func (f Field) FieldName(s string) string {
	if len(f.Path) == 0 {
		return ""
	}
	hash := sha1.New()
	hash.Write([]byte(s))
	b := hash.Sum(nil)
	return strings.ToLower(f.Path[len(f.Path)-1][0:1]) + f.Path[len(f.Path)-1][1:] + hex.EncodeToString(b)[0:4]
}

type DataFieldMap struct {
	Pkg           *packages.Package
	Name          string
	Type          *common.Type
	NamedMap      OrderMap
	PointerMap    OrderMap
	SliceMap      OrderMap
	MapMap        OrderMap
	BasicMap      OrderMap
	MapMethodList []MapMethod
}

type OrderMap map[string]Field

func (o OrderMap) GetByField(f Field) (res Field, has bool) {
	fName := lo.Ternary(lo.IsEmpty(f.CopyName), f.Name, f.CopyName)
	fName = lo.Ternary(lo.IsEmpty(f.CopyPrefix), fName, strings.TrimPrefix(fName, f.CopyPrefix))

	slog.Debug("getByField", "fName", fName, "path", strings.Join(f.Path, "."), "aliasPath", strings.Join(f.AliasPath, "."))
	// 通过路径查找 最高优先级
	if lo.IsNotEmpty(f.CopyTargetPath) {
		for _, v := range o {
			// slog.Info("copyTargetPath", "copyTargetPath", f.CopyTargetPath, "path", strings.Join(v.Path, "."))
			if strings.Join(v.Path, ".") == f.CopyTargetPath {
				res = v
				has = true
				return
			}
		}
	}

	if lo.IsNotEmpty(f.CopyTargetName) {
		for _, v := range o {
			// slog.Info("copyTargetName", "copyTargetName", f.CopyTargetName, "name", v.Name)
			if f.CopyTargetName == v.Name {
				res = v
				has = true
				return
			}
		}
	}

	// 通过名字查找 如果路径一样最高优先级 如果路径不一样名字一样也行
	for _, v := range o {
		vName := lo.Ternary(lo.IsEmpty(v.CopyName), v.Name, v.CopyName)
		// slog.Info("getByField eq name", "fName", fName, "vName", vName, "sourceName", f.Name, "copyName", f.CopyName, "doc", f.Doc)
		if fName == vName {
			if slices.Equal(lo.DropRight(v.AliasPath, 1), lo.DropRight(f.AliasPath, 1)) {
				// fmt.Println("fName == vName", fName, vName, lo.DropRight(v.AliasPath, 1), lo.DropRight(f.AliasPath, 1))
				res = v
				has = true
				// 最高优先级直接返回
				return
			}
		}
	}

	for _, of := range o {
		if DepthFind(of, 0, f, 0) {
			return of, true
		}
	}

	if f.CopyMust {
		slog.Error("字段未找到", "name", f.Name, "path", strings.Join(f.Path, "."))
		panic("字段未找到")
	}
	return Field{}, false
}

func DepthFind(dest Field, destIndex int, src Field, srcIndex int) bool {
	/* if destIndex == 0 {
		fmt.Println("first", "destIndex", dest.AliasPath, dest.EmbeddedIndex, "srcIndex", src.AliasPath, src.EmbeddedIndex)
	} */
	if destIndex > len(dest.AliasPath)-1 && srcIndex > len(src.AliasPath)-1 {
		return true
	}
	if destIndex > len(dest.AliasPath)-1 || srcIndex > len(src.AliasPath)-1 {
		return false
	}
	if len(src.AliasPath)-1 == srcIndex && len(dest.AliasPath)-1 == destIndex && src.AliasPath[srcIndex] == dest.AliasPath[destIndex] {
		return true
	}

	// 相应的节点
	var eq0 bool
	// src 隐藏
	var eq1 bool
	// dest 隐藏
	var eq2 bool
	// src dest 隐藏
	var eq3 bool
	if src.AliasPath[srcIndex] == dest.AliasPath[destIndex] {
		eq0 = DepthFind(dest, destIndex+1, src, srcIndex+1)
	}

	if srcIndex+1 < len(src.AliasPath) && lo.Contains(src.EmbeddedIndex, srcIndex) && src.AliasPath[srcIndex+1] == dest.AliasPath[destIndex] {
		eq1 = DepthFind(dest, destIndex+1, src, srcIndex+2)
	}

	if destIndex+1 < len(dest.AliasPath) && lo.Contains(dest.EmbeddedIndex, destIndex) && dest.AliasPath[destIndex+1] == src.AliasPath[srcIndex] {
		eq2 = DepthFind(dest, destIndex+2, src, srcIndex+1)
	}

	if destIndex+1 < len(dest.AliasPath) && srcIndex+1 < len(src.AliasPath) && lo.Contains(dest.EmbeddedIndex, destIndex) && lo.Contains(src.EmbeddedIndex, srcIndex) && src.AliasPath[srcIndex+1] == dest.AliasPath[destIndex+1] {
		eq3 = DepthFind(dest, destIndex+2, src, srcIndex+2)
	}

	return eq0 || eq1 || eq2 || eq3
}

// named 比较特殊要先从段的路径开始排序， 路径长的可能是子路径会影响赋值
// 先根据path长度升序，同长度的再根据名字排序

func (o OrderMap) OrderKeys() []string {
	var keys []string
	mapKey := make(map[int][]string, 0)
	for k, v := range o {
		pathLen := len(v.Path)
		if mapKey[pathLen] == nil {
			mapKey[pathLen] = make([]string, 0)
		}
		mapKey[pathLen] = append(mapKey[pathLen], k)
	}

	mapIndex := lo.Keys(mapKey)
	slices.Sort(mapIndex)
	for _, v := range mapIndex {
		k := mapKey[v]
		slices.Sort(k)
		keys = append(keys, k...)
	}

	return keys

}

type MapMethod struct {
	Name     string
	ParamID  string
	ResultID string
	IsError  bool
}

type StructFields struct {
	NamedMap   OrderMap
	PointerMap OrderMap
	SliceMap   OrderMap
	MapMap     OrderMap
	BasicMap   OrderMap
	StructMap  OrderMap
}

func NewDataFieldMap(pkg *packages.Package, path []string, aliasPath []string, name string, commonType *common.Type) *DataFieldMap {
	m := &DataFieldMap{
		Pkg:           pkg,
		Name:          name,
		Type:          commonType,
		NamedMap:      map[string]Field{},
		PointerMap:    map[string]Field{},
		SliceMap:      map[string]Field{},
		MapMap:        map[string]Field{},
		BasicMap:      map[string]Field{},
		MapMethodList: []MapMethod{},
	}
	m.Parse(Field{
		Path:          path,
		AliasPath:     aliasPath,
		EmbeddedIndex: []int{},
		Name:          "_root",
		Type:          commonType,
		Doc:           nil,
	})

	methodSet := typeutil.IntuitiveMethodSet(commonType.T, nil)
	lo.ForEach(methodSet, func(item *types.Selection, index int) {
		sig := item.Type().(*types.Signature)
		var paramList []string
		var resultList []string
		for i := 0; i < sig.Params().Len(); i++ {
			paramList = append(paramList, sig.Params().At(i).Type().String())
		}

		for i := 0; i < sig.Results().Len(); i++ {
			resultList = append(resultList, sig.Results().At(i).Type().String())
		}

		mapMethod := MapMethod{}
		mapMethod.Name = item.Obj().Name()
		if len(paramList) > 0 {
			mapMethod.ParamID = paramList[0]
		}
		if lo.Contains(resultList, "error") {
			mapMethod.IsError = true
		}

		m.MapMethodList = append(m.MapMethodList, mapMethod)
		// for i := 0; i < sig.RecvTypeParams().Len(); i++ {
		// 	rtp := sig.RecvTypeParams().At(i)
		// 	slog.Info("rtp", "string", rtp.String())
		// 	slog.Info("rtp", "constraint string", rtp.Constraint().String())
		// 	spew.Dump(rtp.Obj())
		// }
	})
	return m
}

func (d *DataFieldMap) saveField(m map[string]Field, name string, field Field) {
	//var oldField Field
	//var ok bool
	mapID := strings.Join(field.Path, ".")
	if oldField, ok := m[mapID]; !ok {
		m[mapID] = field
		return
	} else {
		slog.Error("MapID 重复.", "Path", strings.Join(field.Path, "."), "oldPath", strings.Join(oldField.Path, "."))
		panic("MapID 重复.")
	}

	//fmt.Printf("作用域内重复定义: %s. src.DestIdPath: %s, src.SrcIdPath: %s, dest.DestIdPath: %s, dest.SrcIdPath: %s \n", name, oldField.DestIdPath().GoString(), oldField.SrcIdPath().GoString(), field.DestIdPath().GoString(), field.SrcIdPath().GoString())
	//if len(oldField.Path) > len(field.Path) {
	//	m[name] = field
	//}

	return
}

// func (d *DataFieldMap) Parse(prefix []string, name string, t types.Type, doc *ast.CommentGroup) {
func (d *DataFieldMap) Parse(f Field) {
	(&f).ParseDoc()
	switch v := f.Type.T.(type) {
	case *types.Pointer:
		d.Parse(Field{
			Name:          f.Name,
			Type:          common.TypeOf(f.Type.PointerType.Elem()),
			Path:          f.Path,
			AliasPath:     f.AliasPath,
			Doc:           f.Doc,
			EmbeddedIndex: f.EmbeddedIndex,
		})
	case *types.Basic:
		d.saveField(d.BasicMap, f.Name, f)
		return
	case *types.Map:
		d.saveField(d.MapMap, f.Name, f)
		return
	case *types.Slice:
		d.saveField(d.SliceMap, f.Name, f)
		return
	case *types.Array:
	case *types.Named:
		d.saveField(d.NamedMap, f.Name, f)
		d.Parse(Field{
			Name:          f.Name,
			Type:          common.TypeOf(v.Underlying()),
			Path:          f.Path,
			AliasPath:     f.AliasPath,
			Doc:           f.Doc,
			IsNamed:       true,
			EmbeddedIndex: f.EmbeddedIndex,
		})
		return
	case *types.Struct:
		for i := 0; i < v.NumFields(); i++ {
			field := v.Field(i)

			if !field.Exported() {
				continue
			}
			indexName := field.Name()
			fieldDoc := common.GetCommentByTokenPos(d.Pkg, field.Pos())
			parseDoc, err := common.ParseDoc(fieldDoc.Text())
			if err != nil {
				slog.Error("parseDoc err", "err", err, "doc", parseDoc)
				panic(err)
			}
			var aliasName string
			parseDoc.ByFuncNameAndArgs(CopyNameTag, &aliasName)
			if lo.IsEmpty(f.CopyName) {
				aliasName = field.Name()
			}
			var path []string
			var aliasPath []string

			// if !field.Embedded() {
			path = append(f.Path[0:], field.Name())
			aliasPath = append(f.AliasPath[0:], aliasName)
			// }
			d.Parse(Field{
				Path:          path,
				AliasPath:     aliasPath,
				EmbeddedIndex: lo.Ternary(field.Embedded(), append(f.EmbeddedIndex[0:], len(path)-1), f.EmbeddedIndex),
				Name:          indexName,
				Type:          common.TypeOf(field.Type()),
				ParentDoc:     f.Doc,
				Doc:           parseDoc,
			})
		}
		return
		//default:
		//	panic("unknown types.Type " + f.Type.T.String())
	}
}

type Recorder struct {
	m map[string]struct{}
}

func NewRecorder() *Recorder {
	return &Recorder{m: map[string]struct{}{}}
}

func (r *Recorder) Reg(name string) {
	r.m[name] = struct{}{}
}

func (r *Recorder) Lookup(name string) bool {
	_, ok := r.m[name]
	return ok
}

//func NewResponse(pkg *packages.Package, f *types.Func, responseName string) *Copy {
//	fnName := f.Id()
//	src := f.Type().(*types.Signature).Results().At(0)
//	srcType := src.Type()
//	_, typeFile := path.Split(src.Type().String())
//	srcTypeName := strings.TrimPrefix(strings.TrimPrefix(typeFile, src.Pkg().Name()), ".")
//	fmt.Println("name: ", src.Name(), "id: ", src.Id(), "typestring", src.Type(), "pkg: ", src.Pkg().Name(), "srctypename: ", srcTypeName)
//	//srcName := fnType.Results.List[0].Names[0].Name
//	//spew.Dump(pkg.Types.Scope())
//	//fmt.Println("srcName: ", srcName)
//	//srcType := pkg.TypesInfo.TypeOf(fnType.Results.List[0].Type)
//	//srcType := pkg.TypesInfo.Types[fnType]
//	//fmt.Println("names: ", pkg.Types.Scope().Names(), "path: ", pkg.Types.Path())
//	destType := pkg.Types.Scope().Lookup(responseName)
//
//	jenF := jen.NewFile("Copy")
//	jenF.Add(jen.Type().Id(fnName + "Copy").Struct())
//
//	dto := Copy{
//		Pkg:            pkg,
//		JenF:           jenF,
//		Recorder:       NewRecorder(),
//		SrcParentPath:  []string{},
//		SrcPath:        []string{},
//		Src:            NewDataFieldMap(pkg, []string{}, "src", common.TypeOf(srcType), srcType),
//		DestParentPath: []string{},
//		DestPath:       []string{},
//		Dest:           NewDataFieldMap(pkg, []string{}, "dest", common.TypeOf(destType.Type()), destType.Type()),
//		DefaultFn: jen.Func().Params(jen.Id("d").Id("*" + fnName + "Copy")).
//			Id("Copy").Params(jen.Id("src").Id(srcTypeName)).Params(jen.Id("dest").Id(responseName)),
//		StructName: fnName,
//	}
//	dto.Gen()
//	return &dto
//}

type Copy struct {
	Pkg            *packages.Package
	JenF           *jen.File
	Recorder       *Recorder
	SrcParentPath  []string
	SrcPath        []string
	Src            *DataFieldMap
	DestParentPath []string
	DestPath       []string
	Dest           *DataFieldMap
	Nest           []*Copy
	DefaultFn      *jen.Statement
	StructName     string
	Head           bool
	// namedEq: 源字段名与目标字段名相同 它的子路径不需要再次拷贝
	NamedEq []string
}

func (d *Copy) FnName() string {
	return d.Src.Type.ID() + "To" + common.UpFirst(d.Dest.Type.ID())
}

func (d *Copy) SumPath() string {
	return strings.Join(d.SrcPath, ".") + ":" + strings.Join(d.DestPath, ".")
}

func (d *Copy) Doc() *jen.Statement {
	code := make(jen.Statement, 0)
	// code = append(code, jen.Comment("parentPath: "+strings.Join(d.SrcParentPath, ".")+":"+strings.Join(d.DestParentPath, ".")))
	// code = append(code, jen.Comment("path: "+strings.Join(d.SrcPath, ".")+":"+strings.Join(d.DestPath, ".")))
	return &code
}

func (d *Copy) SumParentPath() string {
	return strings.Join(d.SrcParentPath, ".") + ":" + strings.Join(d.DestParentPath, ".")
}

func (d *Copy) Gen() {
	slog.Debug("gen fn name: ", "fnName", d.FnName())
	slog.Debug("basic map: ", "srcKeys", d.Src.BasicMap.OrderKeys(), "descKeys", d.Dest.BasicMap.OrderKeys())
	slog.Debug("map map: ", "srcKeys", d.Src.MapMap.OrderKeys(), "descKeys", d.Dest.MapMap.OrderKeys())
	slog.Debug("slice map: ", "srcKeys", d.Src.SliceMap.OrderKeys(), "descKeys", d.Dest.SliceMap.OrderKeys())
	slog.Debug("pointer map: ", "srcKeys", d.Src.PointerMap.OrderKeys(), "descKeys", d.Dest.PointerMap.OrderKeys())
	has, fn := d.GenFn(d.FnName(), d.Src.Type.TypeAsJenComparePkgName(d.Pkg), d.Dest.Type.TypeAsJenComparePkgName(d.Pkg))
	if has {
		return
	}
	bind := make(jen.Statement, 0)
	if d.Src.Type.Pointer {
		bind = append(bind, jen.Id("if src == nil { return }"))
	}

	// 为了不改变外面的指针所以不需要New
	if !d.Head && d.Dest.Type.Pointer {
		bind = append(bind, jen.Id("dest = new").Call(common.TypeOf(d.Dest.Type.PointerType.Elem()).TypeAsJenComparePkgName(d.Pkg)))
	}
	// if d.Dest.Type.Pointer {
	// bind = append(bind, jen.Id("dest = new").Call(common.TypeOf(d.Dest.Type.PointerType.Elem()).TypeAsJenComparePkgName(d.Pkg)))
	// }
	bind = append(bind, jen.Comment("named map"))
	bind = append(bind, d.GenNamed()...)
	bind = append(bind, jen.Comment("basic map"))
	bind = append(bind, d.GenBasic()...)
	bind = append(bind, jen.Comment("slice map"))
	bind = append(bind, d.GenSlice()...)
	bind = append(bind, jen.Comment("map map"))
	bind = append(bind, d.GenMap()...)
	bind = append(bind, jen.Comment("pointer map"))
	bind = append(bind, d.GenPointer()...)
	bind = append(bind, jen.Comment("method map"))
	bind = append(bind, d.GenMapMethod()...)
	bind = append(bind, jen.Return())

	fn.Block(bind...)
	//d.JenF.Add(d.Doc())
	d.JenF.Add(fn)
	for _, v := range d.Nest {
		v.Gen()
	}
}

func (d *Copy) GenExtraCopyMethod(bind *jen.Statement, destV, srcV Field) (has bool) {
	// pkgName, methodName := destV.CopyMethod()
	// if pkgName == "" && methodName == "" {
	// 	return false
	// }

	// bind.Add(destV.DestIdPath().Op("=").Add(jen.Qual(pkgName, methodName).Call(srcV.SrcIdPath())))
	// return true
	return false

}
func (d *Copy) GenNamed() jen.Statement {
	bind := make(jen.Statement, 0)
	for _, k := range d.Dest.NamedMap.OrderKeys() {
		v := d.Dest.NamedMap[k]
		if d.CheckNamedEq(v) {
			continue
		}

		if !lo.IsEmpty(v.CopyTargetMethod) {
			bind.Add(v.DestIdPath().Op("=").Id("src." + v.CopyTargetMethod + "()"))
			continue
		}
		srcV, ok := d.Src.NamedMap.GetByField(v)
		if !ok {
			continue
		}

		if d.GenExtraCopyMethod(&bind, v, srcV) {
			continue
		}

		if v.Type.ID() == srcV.Type.ID() {
			d.NamedEq = append(d.NamedEq, v.DestIdPath().GoString())
			bind.Add(v.DestIdPath().Op("=").Add(srcV.SrcIdPath()))
		}

	}

	return bind
}

func (d *Copy) CheckNamedEq(f Field) bool {
	for _, v := range d.NamedEq {
		if strings.HasPrefix(f.DestIdPath().GoString(), v) {
			return true
		}
	}
	return false
}

func (d *Copy) GenBasic() jen.Statement {
	bind := make(jen.Statement, 0)
	for _, k := range d.Dest.BasicMap.OrderKeys() {
		v := d.Dest.BasicMap[k]

		if d.CheckNamedEq(v) {
			continue
		}

		if !lo.IsEmpty(v.CopyTargetMethod) {
			// bind = append(bind, jen.Comment("source: "+v.SrcIdPath().GoString()+" target: "+v.CopyTargetMethod+"()"))
			bind.Add(v.DestIdPath().Op("=").Id("src." + v.CopyTargetMethod + "()"))
			continue
		}

		srcV, ok := d.Src.BasicMap.GetByField(v)
		if !ok {
			continue
		}

		if d.GenExtraCopyMethod(&bind, v, srcV) {
			continue
		}
		//dtoMethod := v.CopyMethod()
		//if dtoMethod != nil {
		//	bind.Add(v.DestIdPath().Op("=").Add(dtoMethod.Call(srcV.SrcIdPath())))
		//	continue
		//}
		// bind = append(bind, jen.Comment("source: "+v.SrcIdPath().GoString()+" target: "+v.DestIdPath().GoString()))
		bind.Add(v.DestIdPath().Op("=").Add(srcV.SrcIdPath()))
	}
	return bind
}

func (d *Copy) GenMap() jen.Statement {
	bind := make(jen.Statement, 0)
	for _, key := range d.Dest.MapMap.OrderKeys() {
		v := d.Dest.MapMap[key]
		if d.CheckNamedEq(v) {
			continue
		}
		srcV, ok := d.Src.MapMap.GetByField(v)
		if !ok {
			// fmt.Printf("not found %s in %s\n", v.Name, d.SumPath())
			continue
		}
		//if v.Doc != nil {
		//	bind.Add(jen.Comment(v.Doc.Text()))
		//}

		if d.GenExtraCopyMethod(&bind, v, srcV) {
			continue
		}

		// 类型一样 直接=
		if v.Type.T.String() == srcV.Type.T.String() {
			bind.Add(v.DestIdPath().Op("=").Add(srcV.SrcIdPath()))
			continue
		}

		bind.Add(v.DestIdPath().Op("=").Make(v.Type.TypeAsJenComparePkgName(d.Pkg), jen.Id("len").Call(srcV.SrcIdPath())))
		block := v.DestIdPath().Index(jen.Id("key")).Op("=").Add(srcV.SrcIdPath()).Index(jen.Id("value"))
		if !v.Type.MapValue.Basic {
			srcMapValue := srcV.Type.MapValue
			destMapValue := v.Type.MapValue
			//srcName := destMapValue.HashID(d.SumPath())
			//destName := destMapValue.HashID(d.SumPath())
			srcName := srcV.FieldName(d.SumPath())
			destName := v.FieldName(d.SumPath())
			nestCopy := &Copy{
				Pkg:            d.Pkg,
				JenF:           d.JenF,
				Recorder:       d.Recorder,
				SrcParentPath:  append(d.SrcParentPath, srcV.Path...),
				SrcPath:        append([]string{}, srcV.Path...),
				Src:            NewDataFieldMap(d.Pkg, []string{}, []string{}, srcName, srcMapValue),
				DestParentPath: append(d.DestParentPath, v.Path...),
				DestPath:       append([]string{}, v.Path...),
				Dest:           NewDataFieldMap(d.Pkg, []string{}, []string{}, destName, destMapValue),
				Nest:           make([]*Copy, 0),
				StructName:     d.StructName,
			}
			d.Nest = append(d.Nest, nestCopy)

			block = v.DestIdPath().Index(jen.Id("key")).Op("=").Id("d." + nestCopy.FnName()).Call(jen.Id("value"))
		}
		bind.Add(jen.For(
			jen.List(jen.Id("key"), jen.Id("value")).
				Op(":=").Range().Add(srcV.SrcIdPath()).Block(
				block,
			)))
	}
	return bind
}

func (d *Copy) GenMapMethod() jen.Statement {
	bind := make(jen.Statement, 0)
	lo.ForEach(d.Dest.MapMethodList, func(item MapMethod, index int) {
		if item.ParamID == d.Src.Type.T.String() {
			bind.Add(jen.Id("(&dest)").Dot(item.Name).Call(jen.Id("src")).Line())
		}
	})
	return bind
}

func (d *Copy) GenPointer() jen.Statement {
	bind := make(jen.Statement, 0)
	for _, key := range d.Dest.PointerMap.OrderKeys() {
		v := d.Dest.PointerMap[key]
		if d.CheckNamedEq(v) {
			continue
		}
		srcV, ok := d.Src.PointerMap.GetByField(v)
		if !ok {
			continue
		}

		//if v.Doc != nil {
		//	bind.Add(jen.Comment(v.Doc.Text()))
		//}

		if d.GenExtraCopyMethod(&bind, v, srcV) {
			continue
		}

		if v.Type.T.String() == srcV.Type.T.String() {
			bind.Add(v.DestIdPath().Op("=").Add(srcV.SrcIdPath()))
			continue
		}
		if v.Type.PointerInner.Basic {
			bind.Add(v.DestIdPath().Op("=").Add(srcV.SrcIdPath()))
		} else {
			srcLiner := srcV.Type.PointerInner
			destLiner := v.Type.PointerInner
			srcName := srcV.FieldName(d.SumPath())
			destName := v.FieldName(d.SumPath())
			//destName := srcLiner.HashID(d.SumPath())
			nestCopy := &Copy{
				Pkg:            d.Pkg,
				JenF:           d.JenF,
				Recorder:       d.Recorder,
				SrcParentPath:  append(d.SrcParentPath, srcV.Path...),
				SrcPath:        append([]string{}, srcV.Path...),
				Src:            NewDataFieldMap(d.Pkg, []string{}, []string{}, srcName, srcLiner),
				DestParentPath: append(d.DestParentPath, v.Path...),
				DestPath:       append([]string{}, v.Path...),
				Dest:           NewDataFieldMap(d.Pkg, []string{}, []string{}, destName, destLiner),
				Nest:           make([]*Copy, 0),
				StructName:     d.StructName,
			}
			d.Nest = append(d.Nest, nestCopy)

			bind.Add(
				jen.If(srcV.SrcIdPath().Op("!=").Nil()).Block(
					jen.Id("v").Op(":=").Id("d."+nestCopy.FnName()).Call(jen.Id("*").Add(srcV.SrcIdPath())),
					v.DestIdPath().Op("=").Id("&v"),
				),
			)
		}
	}
	return bind
}

func (d *Copy) GenSlice() jen.Statement {
	bind := make(jen.Statement, 0)
	for _, key := range d.Dest.SliceMap.OrderKeys() {
		v := d.Dest.SliceMap[key]

		if d.CheckNamedEq(v) {
			continue
		}

		srcV, ok := d.Src.SliceMap.GetByField(v)
		if !ok {
			continue
		}
		//if v.Doc != nil {
		//	bind.Add(jen.Comment(v.Doc.Text()))
		//}
		//fmt.Println("common", "ttype", "slice", "id", v.Type.ID(), "unescapedid", v.Type.UnescapedID(), "jen", v.Type.TypeAsJenComparePkgName().Render(os.Stdout))

		if d.GenExtraCopyMethod(&bind, v, srcV) {
			continue
		}

		if v.Type.T.String() == srcV.Type.T.String() {
			bind.Add(v.DestIdPath().Op("=").Add(srcV.SrcIdPath()))
			continue
		}
		bind.Add(v.DestIdPath().Op("=").Make(v.Type.TypeAsJenComparePkgName(d.Pkg), jen.Id("0"), jen.Id("len").Call(srcV.SrcIdPath())))
		block := v.DestIdPath().Index(jen.Id("i")).Op("=").Add(srcV.SrcIdPath()).Index(jen.Id("i"))
		if !v.Type.ListInner.Basic {
			srcLiner := srcV.Type.ListInner
			destLiner := v.Type.ListInner
			//fmt.Println("listInner", srcLiner.TypeAsJen().GoString())
			//srcName := srcLiner.HashID(d.SumPath())
			srcName := srcV.FieldName(d.SumPath())
			destName := v.FieldName(d.SumPath())
			//destName := srcLiner.HashID(d.SumPath())
			nestCopy := &Copy{
				Pkg:           d.Pkg,
				JenF:          d.JenF,
				Recorder:      d.Recorder,
				SrcParentPath: append(d.SrcParentPath, srcV.Path...),
				//SrcPath:  append([]string{}, srcV.Path...),
				SrcPath:        d.SrcPath[0:],
				Src:            NewDataFieldMap(d.Pkg, []string{}, []string{}, srcName, srcLiner),
				DestParentPath: append(d.DestParentPath, v.Path...),
				//DestPath: append([]string{}, v.Path...),
				DestPath:   d.DestPath[0:],
				Dest:       NewDataFieldMap(d.Pkg, []string{}, []string{}, destName, destLiner),
				Nest:       make([]*Copy, 0),
				StructName: d.StructName,
			}
			d.Nest = append(d.Nest, nestCopy)

			block = v.DestIdPath().Index(jen.Id("i")).Op("=").Id("d." + nestCopy.FnName()).Call(srcV.SrcIdPath().Index(jen.Id("i")))
		}
		bind.Add(jen.For(
			jen.Id("i").Op(":=").Lit(0),
			jen.Id("i").Op("<").Id("len").Call(srcV.SrcIdPath()),
			jen.Id("i").Op("++")).
			Block(
				block,
			))
	}
	return bind
}

func (d *Copy) GenFn(funcName string, srcId, destId jen.Code) (has bool, fn *jen.Statement) {
	if d.DefaultFn != nil {
		return false, d.DefaultFn
	}
	srcType := jen.Type().Id("src").Add(srcId)
	destType := jen.Type().Id("dest").Add(destId)

	funcKey := fmt.Sprintf("%s_%s_%s", funcName, srcType.GoString(), destType.GoString())
	//fmt.Printf("funcName: %s, srcpath: %#v, destpath %#v \n", funcName,srcType, destType)

	has = d.Recorder.Lookup(funcKey)
	if has {
		return has, nil
	}
	d.Recorder.Reg(funcKey)

	return false, jen.Func().Params(jen.Id("d").Id(d.StructName)).
		Id(funcName).Params(jen.Id("src").Add(srcId)).Params(jen.Id("dest").Add(destId))
}
