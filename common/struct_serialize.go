package common

import (
	"go/types"
	"reflect"

	"golang.org/x/exp/slog"
	"golang.org/x/tools/go/packages"
)

type StructSerialize struct {
	pkg         *packages.Package
	namedRecord map[string]types.Type
	Fields      []Field
}

func (s *StructSerialize) Parse(t *types.Struct) (StructMetaDate, error) {
	return StructMetaDate{}, nil
}

func (s *StructSerialize) parseStruct(t *types.Struct) {
	s.parseType([]string{}, "", t, "", nil)
}

func (s *StructSerialize) parseType(pre []string, fName string, t types.Type, tag reflect.StructTag, doc *Doc) {
	switch tt := t.(type) {
	case *types.Named:
		// 	s.Fields = append(s.Fields, Field{
		// 	Doc: doc,
		// 	ID:  fName,
		// 	Tag: tag,
		// 	Xtype: TypeOf(t),
		// 	Path: pre,
		// })
		s.parseType(pre, fName, tt.Underlying(), tag, doc)

	case *types.Struct:
		s.Fields = append(s.Fields, Field{
			Doc:   doc,
			ID:    fName,
			Tag:   tag,
			Xtype: TypeOf(t),
			Path:  pre,
		})
		for i := 0; i < tt.NumFields(); i++ {
			field := tt.Field(i)
			fieldName := field.Name()
			fieldType := field.Type()
			fTag := reflect.StructTag(tt.Tag(i))
			docStr := GetCommentByTokenPos(s.pkg, field.Pos()).Text()
			doc, err := ParseDoc(docStr)
			if err != nil {
				slog.Error("ParseDoc", err, slog.String("comment", docStr))
				panic(err)
			}
			s.parseType(append(pre, fieldName), fieldName, fieldType, fTag, doc)
		}
	case *types.Pointer:
		s.parseType(pre, fName, tt.Elem(), tag, doc)
	case *types.Slice:

		s.Fields = append(s.Fields, Field{
			Doc:   doc,
			ID:    fName,
			Tag:   tag,
			Xtype: TypeOf(t),
			Path:  pre,
		})
	case *types.Map:

		s.Fields = append(s.Fields, Field{
			Doc:   doc,
			ID:    fName,
			Tag:   tag,
			Xtype: TypeOf(t),
			Path:  pre,
		})

	case *types.Array:
		s.Fields = append(s.Fields, Field{
			Doc:   doc,
			ID:    fName,
			Tag:   tag,
			Xtype: TypeOf(t),
			Path:  pre,
		})
	case *types.Basic:
		s.Fields = append(s.Fields, Field{
			Doc:   doc,
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
