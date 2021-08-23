package field

import (
	"go/ast"
	"strings"
)

type Field struct {
	field *ast.Field
	name  *ast.Ident
}

func New(f *ast.Field, name *ast.Ident) *Field {
	return &Field{
		field: f,
		name:  name,
	}
}

func (f *Field) Name() string {
	return f.name.Name
}

func (f *Field) IsExported() bool {
	return f.name.IsExported()
}

func (f *Field) IsOptional() bool {
	if f.field.Doc != nil {
		for _, c := range f.field.Doc.List {
			if isOptionalCommentText(c.Text) {
				return true
			}
		}
	}
	if f.field.Comment != nil {
		for _, c := range f.field.Comment.List {
			if isOptionalCommentText(c.Text) {
				return true
			}
		}
	}
	return false
}

func isOptionalCommentText(s string) bool {
	s = strings.TrimPrefix(s, "//")
	s = strings.TrimSpace(s)
	return strings.HasPrefix(s, "optional")
}
