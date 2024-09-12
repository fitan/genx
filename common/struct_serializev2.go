package common

import (
	"go/types"
	"log/slog"
	"reflect"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

type StructSerializeV2 struct {
	pkg    *packages.Package
	t      *Type
	Fields *[]StructFieldMetaDataV2 `json:"fields"`
}

type StructFieldMetaDataV2 struct {
	Name       string            `json:"name"`
	Tag        reflect.StructTag `json:"tag"`
	TypeX      *Type
	Named      bool                     `json:"named"`
	Ptr        bool                     `json:"ptr"`
	NextFields *[]StructFieldMetaDataV2 `json:"nextFields"`
}

func NewStructSerializeV2(pkg *packages.Package, t types.Type) *StructSerializeV2 {
	obj := &StructSerializeV2{pkg: pkg, t: TypeOf(t)}
	obj.run()
	return obj
}

type StructMetaDataV2ParseParams struct {
	FieldName string
	T         types.Type
	Fields    *[]StructFieldMetaDataV2
	Ptr       bool
	Named     bool
	Tag       reflect.StructTag
	Doc       Doc
}

func (s *StructSerializeV2) run() {
	fields := []StructFieldMetaDataV2{}
	s.Fields = &fields
	s.parse(StructMetaDataV2ParseParams{
		FieldName: "",
		T:         s.t.T,
		Fields:    &fields,
		Ptr:       false,
		Named:     false,
		Tag:       "",
		Doc:       []DocLine{},
	})
}

func (s *StructSerializeV2) parse(v StructMetaDataV2ParseParams) {
	slog.Info("parse struct", "t", v.T.String())
	switch tt := v.T.(type) {
	case *types.Struct:
		nextFields := []StructFieldMetaDataV2{}
		*v.Fields = append(*v.Fields, StructFieldMetaDataV2{
			Name:       v.FieldName,
			Tag:        v.Tag,
			TypeX:      TypeOf(v.T),
			Named:      v.Named,
			Ptr:        v.Ptr,
			NextFields: &nextFields,
		})

		for i := 0; i < tt.NumFields(); i++ {
			field := tt.Field(i)
			if !field.Exported() {
				continue
			}
			tag := reflect.StructTag(tt.Tag(i))

			docStr := GetCommentByTokenPos(s.pkg, field.Pos()).Text()
			doc, err := ParseDoc(docStr)
			if err != nil {
				err = errors.Wrapf(err, "field name %s parse doc", field.Name())
				slog.Error("parse doc", "err", err)
				panic(err)
			}

			paramV := StructMetaDataV2ParseParams{
				FieldName: field.Name(),
				T:         field.Type(),
				Fields:    &nextFields,
				Ptr:       false,
				Named:     false,
				Tag:       tag,
				Doc:       doc,
			}

			s.parse(paramV)
		}

	case *types.Basic:
		*v.Fields = append(*v.Fields, StructFieldMetaDataV2{
			Name:       v.FieldName,
			Named:      v.Named,
			Tag:        v.Tag,
			TypeX:      TypeOf(tt),
			Ptr:        v.Ptr,
			NextFields: &[]StructFieldMetaDataV2{},
		})
	case *types.Named:
		v.T = tt.Underlying()
		v.Named = true
		s.parse(v)
	case *types.Pointer:
		v.Ptr = true
		v.T = tt.Elem()
		s.parse(v)
	case *types.Slice:
		*v.Fields = append(*v.Fields, StructFieldMetaDataV2{
			Name:       v.FieldName,
			Named:      v.Named,
			Tag:        v.Tag,
			TypeX:      TypeOf(v.T),
			Ptr:        v.Ptr,
			NextFields: &[]StructFieldMetaDataV2{},
		})
	case *types.Map:
		*v.Fields = append(*v.Fields, StructFieldMetaDataV2{
			Name:       v.FieldName,
			Named:      v.Named,
			Tag:        v.Tag,
			TypeX:      TypeOf(v.T),
			Ptr:        v.Ptr,
			NextFields: &[]StructFieldMetaDataV2{},
		})
	case *types.Array:
		*v.Fields = append(*v.Fields, StructFieldMetaDataV2{
			Name:       v.FieldName,
			Named:      v.Named,
			Tag:        v.Tag,
			TypeX:      TypeOf(v.T),
			Ptr:        v.Ptr,
			NextFields: &[]StructFieldMetaDataV2{},
		})
	case *types.Interface:
		*v.Fields = append(*v.Fields, StructFieldMetaDataV2{
			Name:       v.FieldName,
			Named:      v.Named,
			Tag:        v.Tag,
			TypeX:      TypeOf(v.T),
			Ptr:        v.Ptr,
			NextFields: &[]StructFieldMetaDataV2{},
		})
	}
}
