package common

import (
	"go/types"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/go/packages"
	"reflect"
)

type StructSerialize struct {
	pkg         *packages.Package
	namedRecord map[string]types.Type
	Fields      []Field
}

func (s *StructSerialize) Parse(t *types.Struct) (StructMetaDate, error) {
	typeof
}

func (s *StructSerialize) parseStruct(t *types.Struct) {

}

func (s *StructSerialize) parseType(pre []string, fName string, t types.Type, tag reflect.StructTag, doc string) {
	switch tt := t.(type) {
	case *types.Named:

	case *types.Struct:
		for i := 0; i < tt.NumFields(); i++ {
			field := tt.Field(i)
			fieldName := field.Name()
			fieldType := field.Type()
			fTag := reflect.StructTag(tt.Tag(i))
			s.parseType(append(pre, fieldName), fieldName, fieldType, fTag, GetCommentByTokenPos(s.pkg, field.Pos()).Text())
		}
	case *types.Pointer:
	case *types.Slice:
	case *types.Map:
	case *types.Array:
	case *types.Basic:
		d, err := ParseDoc(doc)
		if err != nil {
			slog.Error("ParseDoc", err, slog.String("comment", doc))
		}
		s.Fields = append(s.Fields, Field{
			Doc:   d,
			ID:    fName,
			Tag:   tag,
			Xtype: TypeOf(t),
			Path:  pre,
		})
	}
}

type StructMetaDate struct {
	Fields []Field
}

type Field struct {
	Doc   *Doc
	ID    string
	Tag   reflect.StructTag
	Xtype *Type
	Path  []string
}
