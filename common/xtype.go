package common

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"go/types"
	"log"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/fitan/jennifer/jen"
)

// ThisVar is used as name for the reference to the converter interface.
const ThisVar = "c"

// Signature represents a signature for conversion.
type Signature struct {
	Source string
	Target string
}

// Type is a helper wrapper for types.Type.
type Type struct {
	T             types.Type
	Interface     bool
	InterfaceType *types.Interface
	Struct        bool
	StructType    *types.Struct
	Named         bool
	NamedType     *types.Named
	Pointer       bool
	PointerType   *types.Pointer
	PointerInner  *Type
	List          bool
	ListFixed     bool
	ListInner     *Type
	Map           bool
	MapType       *types.Map
	MapKey        *Type
	MapValue      *Type
	Basic         bool
	BasicType     *types.Basic
	Func          bool
	FuncType      *types.Func
}

// StructField holds the type of a struct field and its name.
type StructField struct {
	Name string
	Type *Type
}

// StructField returns the type of a struct field and its name upon successful match or
// an error if it is not found. This method will also return a detailed error if matchIgnoreCase
// is enabled and there are multiple non-exact matches.
func (t Type) StructField(name string, ignoreCase bool, ignore map[string]struct{}) (*StructField, error) {
	if !t.Struct {
		panic("trying to get field of non struct")
	}

	var ambMatches []*StructField
	for y := 0; y < t.StructType.NumFields(); y++ {
		m := t.StructType.Field(y)
		if _, ignored := ignore[m.Name()]; ignored {
			continue
		}
		if m.Name() == name {
			// exact match takes precedence over case-insensitive match
			return &StructField{Name: m.Name(), Type: TypeOf(m.Type())}, nil
		}
		if ignoreCase && strings.EqualFold(m.Name(), name) {
			ambMatches = append(ambMatches, &StructField{Name: m.Name(), Type: TypeOf(m.Type())})
			// keep going to ensure struct does not have another case-insensitive match
		}
	}

	switch len(ambMatches) {
	case 0:
		return nil, fmt.Errorf("%q does not exist", name)
	case 1:
		return ambMatches[0], nil
	default:
		ambNames := make([]string, 0, len(ambMatches))
		for _, m := range ambMatches {
			ambNames = append(ambNames, m.Name)
		}
		return nil, ambiguousMatchError(name, ambNames)
	}
}

// JenID a jennifer code wrapper with extra infos.
type JenID struct {
	Code     *jen.Statement
	Variable bool
}

// VariableID is used, when the ID can be referenced. F.ex it is not a function call.
func VariableID(code *jen.Statement) *JenID {
	return &JenID{Code: code, Variable: true}
}

// OtherID is used, when the ID isn't a variable id.
func OtherID(code *jen.Statement) *JenID {
	return &JenID{Code: code, Variable: false}
}

// TypeOf creates a Type.
func TypeOf(t types.Type) *Type {
	rt := &Type{}
	rt.T = t
	switch value := t.(type) {
	case *types.Pointer:
		rt.Pointer = true
		rt.PointerType = value
		rt.PointerInner = TypeOf(value.Elem())
	case *types.Basic:
		rt.Basic = true
		rt.BasicType = value
	case *types.Map:
		rt.Map = true
		rt.MapType = value
		rt.MapKey = TypeOf(value.Key())
		rt.MapValue = TypeOf(value.Elem())
	case *types.Slice:
		rt.List = true
		rt.ListInner = TypeOf(value.Elem())
	case *types.Array:
		rt.List = true
		rt.ListFixed = true
		rt.ListInner = TypeOf(value.Elem())
	case *types.Named:
		underlying := TypeOf(value.Underlying())
		underlying.T = value
		underlying.Named = true
		underlying.NamedType = value
		return underlying
	case *types.Struct:
		rt.Struct = true
		rt.StructType = value
	case *types.Interface:
		rt.Interface = true
		rt.InterfaceType = value
	default:
		log.Println("unknown type", t)
		//panic("unknown types.Type " + t.String())
	}
	return rt
}

// ID returns a deteministically generated id that may be used as variable.
func (t *Type) ID() string {
	return t.asID(true, true)
}

func (t *Type) HashID(s string) string {
	hash := sha1.New()
	hash.Write([]byte(s))
	b := hash.Sum(nil)
	return t.ID() + hex.EncodeToString(b)[0:4]
}

// UnescapedID returns a deteministically generated id that may be used as variable
// reserved keywords aren't escaped.
func (t *Type) UnescapedID() string {
	return t.asID(true, false)
}

func (t *Type) asID(seeNamed, escapeReserved bool) string {
	if seeNamed && t.Named {
		pkgName := t.NamedType.Obj().Pkg().Name()
		name := pkgName + t.NamedType.Obj().Name()
		return name
	}
	if t.List {
		return t.ListInner.asID(true, false) + "List"
	}
	if t.Basic {
		if escapeReserved {
			return "x" + t.BasicType.String()
		}
		return t.BasicType.String()
	}
	if t.Pointer {
		return "p" + strings.Title(t.PointerInner.asID(true, false))
	}
	if t.Map {
		return "map" + strings.Title(t.MapKey.asID(true, false)+strings.Title(t.MapValue.asID(true, false)))
	}
	if t.Struct {
		if escapeReserved {
			return "xstruct"
		}
		return "struct"
	}
	return "unknown"
}

// TypeAsJen returns a jen representation of the type.
func (t Type) TypeAsJen() *jen.Statement {
	if t.Named {
		return toCode(t.NamedType, &jen.Statement{})
	}
	return toCode(t.T, &jen.Statement{})
}

func (t Type) TypeAsJenComparePkgNameString(pkg *packages.Package) string {
	if t.Named {
		return toCodeComparePkgNameString(pkg, t.NamedType, "")
	}
	return toCodeComparePkgNameString(pkg, t.T, "")
}

// ToCodeComparePkgNameString 导出版本用于测试
func ToCodeComparePkgNameString(pkg *packages.Package, t types.Type, s string) string {
	return toCodeComparePkgNameString(pkg, t, s)
}

func toCodeComparePkgNameString(pkg *packages.Package, t types.Type, s string) string {
	switch cast := t.(type) {
	case *types.Named:
		if cast.Obj().Pkg() == nil || pkg.Name == cast.Obj().Pkg().Name() {
			return s + cast.Obj().Name()
		}
		return s + cast.Obj().Pkg().Name() + "." + cast.Obj().Name()
	case *types.Map:
		key := toCodeComparePkgNameString(pkg, cast.Key(), "")
		return s + "map[" + key + "]" + toCodeComparePkgNameString(pkg, cast.Elem(), "")
	case *types.Slice:
		nest := toCodeComparePkgNameString(pkg, cast.Elem(), "")
		if strings.HasPrefix("...", nest) {
			return s + nest
		}
		return s + "[]" + nest
	case *types.Array:
		n := strconv.FormatInt(cast.Len(), 10)
		return s + "[" + n + "]" + toCodeComparePkgNameString(pkg, cast.Elem(), "")
	case *types.Pointer:
		return s + "*" + toCodeComparePkgNameString(pkg, cast.Elem(), "")
	case *types.Basic:
		return s + cast.String()
	case *types.Struct:
		return s + t.String()
	case *types.Signature:
		// 手动构建函数签名，正确处理可变参数，并添加 func 关键字
		var params []string
		for i := 0; i < cast.Params().Len(); i++ {
			param := cast.Params().At(i)
			paramType := toCodeComparePkgNameString(pkg, param.Type(), "")

			// 如果是最后一个参数且函数是可变参数，则添加...前缀
			if i == cast.Params().Len()-1 && cast.Variadic() {
				// 移除[]前缀，添加...前缀
				if strings.HasPrefix(paramType, "[]") {
					paramType = "..." + paramType[2:]
				}
			}

			if param.Name() != "" {
				params = append(params, param.Name()+" "+paramType)
			} else {
				params = append(params, paramType)
			}
		}

		var results []string
		for i := 0; i < cast.Results().Len(); i++ {
			result := cast.Results().At(i)
			resultType := toCodeComparePkgNameString(pkg, result.Type(), "")

			if result.Name() != "" {
				results = append(results, result.Name()+" "+resultType)
			} else {
				results = append(results, resultType)
			}
		}

		signature := "(" + strings.Join(params, ", ") + ")"
		if len(results) > 0 {
			if len(results) == 1 && cast.Results().At(0).Name() == "" {
				signature += " " + results[0]
			} else {
				signature += " (" + strings.Join(results, ", ") + ")"
			}
		}

		return s + "func" + signature
	default:
		return s + t.String()
	}
	//panic("unsupported type " + t.String())
}

func (t Type) TypeAsJenComparePkgName(pkg *packages.Package) *jen.Statement {
	if t.Named {
		return toCodeComparePkgName(pkg, t.NamedType, &jen.Statement{})
	}
	return toCodeComparePkgName(pkg, t.T, &jen.Statement{})
}

func toCodeComparePkgName(pkg *packages.Package, t types.Type, st *jen.Statement) *jen.Statement {
	switch cast := t.(type) {
	case *types.Named:
		if cast.Obj().Pkg() == nil || pkg.Name == cast.Obj().Pkg().Name() {
			return st.Id(cast.Obj().Name())
		}
		return st.Qual(cast.Obj().Pkg().Path(), cast.Obj().Name())
	case *types.Map:
		key := toCodeComparePkgName(pkg, cast.Key(), &jen.Statement{})
		return toCodeComparePkgName(pkg, cast.Elem(), st.Map(key))
	case *types.Slice:
		return toCodeComparePkgName(pkg, cast.Elem(), st.Index())
	case *types.Array:
		return toCodeComparePkgName(pkg, cast.Elem(), st.Index(jen.Lit(int(cast.Len()))))
	case *types.Pointer:
		return toCodeComparePkgName(pkg, cast.Elem(), st.Op("*"))
	case *types.Basic:
		return toCodeBasic(cast.Kind(), st)
	case *types.Struct:
		return st.Id(t.String())
	case *types.Signature:
		// 对于函数签名类型，返回func关键字加上签名
		// 注意：这里不处理可变参数，因为Jennifer会在更高层处理
		return st.Func().Params().Params()
	}
	panic("unsupported type " + t.String())
}

func toCode(t types.Type, st *jen.Statement) *jen.Statement {
	switch cast := t.(type) {
	case *types.Named:
		if cast.Obj().Pkg() == nil {
			return st.Id(cast.Obj().Name())
		}
		return st.Qual(cast.Obj().Pkg().Path(), cast.Obj().Name())
	case *types.Map:
		key := toCode(cast.Key(), &jen.Statement{})
		return toCode(cast.Elem(), st.Map(key))
	case *types.Slice:
		return toCode(cast.Elem(), st.Index())
	case *types.Array:
		return toCode(cast.Elem(), st.Index(jen.Lit(int(cast.Len()))))
	case *types.Pointer:
		return toCode(cast.Elem(), st.Op("*"))
	case *types.Basic:
		return toCodeBasic(cast.Kind(), st)
	case *types.Struct:
		return st.Id(t.String())
	}
	panic("unsupported type " + t.String())
}

func toCodeBasic(t types.BasicKind, st *jen.Statement) *jen.Statement {
	switch t {
	case types.String:
		return st.String()
	case types.Int:
		return st.Int()
	case types.Int8:
		return st.Int8()
	case types.Int16:
		return st.Int16()
	case types.Int32:
		return st.Int32()
	case types.Int64:
		return st.Int64()
	case types.Uint:
		return st.Uint()
	case types.Uint8:
		return st.Uint8()
	case types.Uint16:
		return st.Uint16()
	case types.Uint32:
		return st.Uint32()
	case types.Uint64:
		return st.Uint64()
	case types.Bool:
		return st.Bool()
	case types.Complex128:
		return st.Complex128()
	case types.Complex64:
		return st.Complex64()
	case types.Float32:
		return st.Float32()
	case types.Float64:
		return st.Float64()
	default:
		return st.Id("unsupported")
		// panic(fmt.Sprintf("unsupported type %T", t))
	}
}

func ambiguousMatchError(name string, ambNames []string) error {
	return fmt.Errorf(`multiple matches found for %q. Possible matches: %s.

Explicitly define the mapping via goverter:map. Example:

    goverter:map %s %s

See https://github.com/jmattheis/goverter#struct-field-mapping`, name, strings.Join(ambNames, ", "), ambNames[0], name)
}
