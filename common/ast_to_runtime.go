package common

import (
	"go/types"
	"golang.org/x/exp/slog"
	"reflect"
)

func TypesStruct2RuntimeStruct(t *types.Struct) reflect.Type {
	var fields []reflect.StructField
	for i := 0; i < t.NumFields(); i++ {
		field := t.Field(i)
		fieldName := field.Name()
		fieldType := field.Type()
		fieldTag := t.Tag(i)
		var pkgPath string
		if !field.Anonymous() {
			pkgPath = field.Pkg().Path()
		}
		fields = append(fields, reflect.StructField{
			Name:      fieldName,
			PkgPath:   pkgPath,
			Type:      TypesType2ReflectType(fieldType),
			Tag:       reflect.StructTag(fieldTag),
			Offset:    0,
			Index:     nil,
			Anonymous: field.Anonymous(),
		})
	}

	return reflect.StructOf(fields)
}

func TypesMap2RuntimeMap(t *types.Map) reflect.Type {
	return reflect.MapOf(TypesType2ReflectType(t.Key()), TypesType2ReflectType(t.Elem()))
}

func TypesSlice2RuntimeSlice(t *types.Slice) reflect.Type {
	return reflect.SliceOf(TypesType2ReflectType(t.Elem()))
}

func TypesArray2RuntimeArray(t *types.Array) reflect.Type {
	return reflect.ArrayOf(int(t.Len()), TypesType2ReflectType(t.Elem()))
}

func TypesPoint2RuntimePoint(t *types.Pointer) reflect.Type {
	return reflect.PtrTo(TypesType2ReflectType(t.Elem()))
}

func TypesType2ReflectType(p types.Type) reflect.Type {
	slog.Info("TypesType2ReflectType", slog.Any("p", p))

	switch pt := p.(type) {
	case *types.Named:
		return TypesType2ReflectType(pt.Underlying())
	case *types.Struct:
		return TypesStruct2RuntimeStruct(pt)
	case *types.Map:
		return TypesMap2RuntimeMap(pt)
	case *types.Slice:
		return TypesSlice2RuntimeSlice(pt)
	case *types.Array:
		return TypesArray2RuntimeArray(pt)
	case *types.Pointer:
		return TypesPoint2RuntimePoint(pt)
	case *types.Basic:
		switch pt.Kind() {
		case types.Bool:
			return reflect.TypeOf(false)
		case types.Int:
			return reflect.TypeOf(int(0))
		case types.Int8:
			return reflect.TypeOf(int8(0))
		case types.Int16:
			return reflect.TypeOf(int16(0))
		case types.Int32:
			return reflect.TypeOf(int32(0))
		case types.Int64:
			return reflect.TypeOf(int64(0))
		case types.Uint:
			return reflect.TypeOf(uint(0))
		case types.Uint8:
			return reflect.TypeOf(uint8(0))
		case types.Uint16:
			return reflect.TypeOf(uint16(0))
		case types.Uint32:
			return reflect.TypeOf(uint32(0))
		case types.Uint64:
			return reflect.TypeOf(uint64(0))
		case types.Uintptr:
			return reflect.TypeOf(uintptr(0))
		case types.Float32:
			return reflect.TypeOf(float32(0))
		case types.Float64:
			return reflect.TypeOf(float64(0))
		case types.Complex64:
			return reflect.TypeOf(complex64(0))
		case types.Complex128:
			return reflect.TypeOf(complex128(0))
		case types.String:
			return reflect.TypeOf("")
		case types.UnsafePointer:
			return reflect.TypeOf(uintptr(0))
		default:
			slog.Error("TypesType2ReflectType Default", slog.Any("pt.Kind()", pt.Kind()))
		}
	default:
		slog.Error("TypesType2ReflectType Default", slog.Any("pt", pt))

	}

	panic(slog.Any("p", p.String()))

	return nil
}
