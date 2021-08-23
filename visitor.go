package partialfields

import (
	"fmt"
	"go/ast"
	"regexp"

	"github.com/kamiaka/go-partialfields/internal/field"
	"github.com/kamiaka/go-partialfields/internal/pkg"
	"golang.org/x/tools/go/analysis"
)

type visitor struct {
	excludesEmptyLiteral    bool
	excludesOptionalField   bool
	excludesExportedField   bool
	excludesUnexportedField bool
	excludeFilePattern      *regexp.Regexp
	errors                  []error
	pkg                     *pkg.Pkg
	files                   []*ast.File
	cmap                    ast.CommentMap
	prev                    ast.Node
	prevs                   map[ast.Node]ast.Node

	pass *analysis.Pass
}

func (v *visitor) print(indent string, expr ast.Node) {
	fmt.Println("")

	if ls, ok := v.cmap[expr]; ok {
		for _, cg := range ls {
			for _, c := range cg.List {
				fmt.Println(indent + c.Text)
			}
		}
	}
	switch stmt := expr.(type) {
	case *ast.Ident:
		fmt.Println(indent + fmt.Sprintf("Ident: %v", stmt.Name))
		obj := stmt.Obj
		if obj != nil {
			fmt.Println(indent, fmt.Sprintf("  .Obj.Kind: %v", obj.Kind.String()))
			fmt.Println(indent, fmt.Sprintf("  .Obj.Decl: %#v", obj.Decl))
		}
	default:
		fmt.Println(indent + fmt.Sprintf("? %#v", stmt))
	}
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	v.prevs[node] = v.prev
	v.prev = node

	switch stmt := node.(type) {
	case *ast.CompositeLit:
		v.print("", node)
		if err := v.partialfieldsErr(stmt); err != nil {
			v.addError(err)
		}
	// case *ast.CommentGroup:
	// 	fmt.Printf("\n// %#v, %d\n", stmt.List, len(stmt.List))
	// 	for i, c := range stmt.List {
	// 		fmt.Printf("  - %d: %#v\n", i, c)
	// 	}
	// case *ast.Ident:
	// 	fmt.Printf("\nID: %#v\n", stmt)
	// 	if stmt.Obj != nil {
	// 		fmt.Printf("  Decl: %#v\n", stmt.Obj.Decl)
	// 	}
	// case *ast.ValueSpec:
	// 	fmt.Printf("\nValueSpec: %#v\n", stmt)
	// 	if stmt.Doc != nil {
	// 		for _, c := range stmt.Doc.List {
	// 			fmt.Printf("  Doc: %v\n", c.Text)
	// 		}
	// 	}
	// 	if stmt.Comment != nil {
	// 		for _, c := range stmt.Comment.List {
	// 			fmt.Printf("  Comment: %v\n", c.Text)
	// 		}
	// 	}

	// case *ast.KeyValueExpr:
	// 	fmt.Printf("\nKV.Key: %#v\n", stmt.Key)
	// 	fmt.Printf("  : %#v\n", stmt.Value)
	default:
		if node != nil {
			v.print("", node)
			ast.Print(v.pass.Fset, node)
		}
	}
	return v
}

func (v *visitor) partialfieldsErr(lit *ast.CompositeLit) error {
	if v.excludesEmptyLiteral && len(lit.Elts) == 0 {
		return nil
	}

	s, isImported, err := v.structType(lit)
	if err != nil {
		return err
	}
	if s == nil {
		return nil
	}

	var fields []*field.Field
	for _, f := range s.Fields.List {
		for _, n := range f.Names {
			fields = append(fields, field.New(f, n))
		}
	}

	if len(fields) == len(lit.Elts) {
		return nil
	}

	decls := map[string]bool{}
	for _, el := range lit.Elts {
		if kv, ok := el.(*ast.KeyValueExpr); ok {
			if k, ok := kv.Key.(*ast.Ident); ok {
				decls[k.Name] = true
			}
		}
	}

	var undecls []string
	for _, field := range fields {
		if decls[field.Name()] ||
			(v.excludesExportedField && field.IsExported()) ||
			((isImported || v.excludesUnexportedField) && !field.IsExported()) ||
			(v.excludesOptionalField && field.IsOptional()) {
			continue
		}

		undecls = append(undecls, field.Name())
	}

	if len(undecls) == 0 {
		return nil
	}

	return newIncompleteError(lit, undecls)
}

func (v *visitor) structType(lit *ast.CompositeLit) (t *ast.StructType, isImported bool, e error) {
	switch stmt := lit.Type.(type) {
	case *ast.Ident:
		return pkg.FindStructInFiles(v.files, stmt.Name), false, nil
	case *ast.SelectorExpr:
		s, err := v.selectorStructType(stmt)
		if err != nil {
			return nil, false, err
		}
		return s, true, nil
	default:
		return nil, false, nil
	}
}

func (v *visitor) selectorStructType(s *ast.SelectorExpr) (*ast.StructType, error) {
	if id, ok := s.X.(*ast.Ident); ok {
		p, err := v.pkg.List(id.Name)
		if err != nil {
			return nil, err
		}

		return p.FindStruct(s.Sel.Name), nil
	}
	return nil, nil
}

func (v *visitor) addError(err error) {
	v.errors = append(v.errors, err)
}
