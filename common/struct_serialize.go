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
	Fields      []StructFieldMetaData
}

func NewStructSerialize(pkg *packages.Package) *StructSerialize {
	return &StructSerialize{pkg: pkg, namedRecord: make(map[string]types.Type), Fields: make([]StructFieldMetaData, 0)}
}

func (s *StructSerialize) Parse(t *types.Struct) (res StructMetaData, err error) {
	s.parseType([]string{}, "", t, "", nil)
	res.Fields = s.Fields
	return
}

func (s *StructSerialize) parseType(pre []string, fName string, t types.Type, tag reflect.StructTag, doc Doc) {
	switch tt := t.(type) {
	case *types.Named:
		// 	s.Fields = append(s.Fields, StructFieldMetaData{
		// 	Doc: doc,
		// 	ID:  fName,
		// 	Tag: tag,
		// 	Xtype: TypeOf(t),
		// 	Path: pre,
		// })
		s.parseType(pre, fName, tt.Underlying(), tag, doc)

	case *types.Struct:
		s.Fields = append(s.Fields, StructFieldMetaData{
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

		s.Fields = append(s.Fields, StructFieldMetaData{
			Doc:   doc,
			ID:    fName,
			Tag:   tag,
			Xtype: TypeOf(t),
			Path:  pre,
		})
	case *types.Map:

		s.Fields = append(s.Fields, StructFieldMetaData{
			Doc:   doc,
			ID:    fName,
			Tag:   tag,
			Xtype: TypeOf(t),
			Path:  pre,
		})

	case *types.Array:
		s.Fields = append(s.Fields, StructFieldMetaData{
			Doc:   doc,
			ID:    fName,
			Tag:   tag,
			Xtype: TypeOf(t),
			Path:  pre,
		})
	case *types.Basic:
		s.Fields = append(s.Fields, StructFieldMetaData{
			Doc:   doc,
			ID:    fName,
			Tag:   tag,
			Xtype: TypeOf(t),
			Path:  pre,
		})
	}
}

type StructMetaData struct {
	Fields []StructFieldMetaData
}

type StructFieldMetaData struct {
	Doc   Doc
	ID    string
	Tag   reflect.StructTag
	Xtype *Type
	Path  []string
}
