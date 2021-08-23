package partialfields

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type IncompleteError struct {
	CompositLit *ast.CompositeLit
	Fields      []string
}

func newIncompleteError(lit *ast.CompositeLit, fields []string) *IncompleteError {
	return &IncompleteError{
		CompositLit: lit,
		Fields:      fields,
	}
}

func (e *IncompleteError) Name() string {
	if s, ok := e.CompositLit.Type.(*ast.SelectorExpr); ok {
		if id, ok := s.X.(*ast.Ident); ok {
			return id.Name + "." + s.Sel.Name
		}
		return s.Sel.Name
	}
	if id, ok := e.CompositLit.Type.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

func (e *IncompleteError) Error() string {
	return fmt.Sprintf("incomplete struct: %s requires %s", e.Name(), strings.Join(e.Fields, ", "))
}

func (e *IncompleteError) Pos() token.Pos {
	return e.CompositLit.Pos()
}

func (e *IncompleteError) End() token.Pos {
	return e.CompositLit.End()
}
