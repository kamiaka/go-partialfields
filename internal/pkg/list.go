package pkg

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type List []*packages.Package

func (ls List) FindStruct(name string) *ast.StructType {
	for _, pkg := range ls {
		if s := FindStructInFiles(pkg.Syntax, name); s != nil {
			return s
		}
	}
	return nil
}

func FindStructInFiles(ls []*ast.File, name string) *ast.StructType {
	var ret *ast.StructType
	for _, f := range ls {
		for _, decl := range f.Decls {
			if gen, ok := decl.(*ast.GenDecl); ok {
				ast.Inspect(gen, func(n ast.Node) bool {
					if spec, ok := n.(*ast.TypeSpec); ok {
						if t, ok := spec.Type.(*ast.StructType); ok && spec.Name.Name == name {
							ret = t
							return false
						}
					}
					return true
				})
				if ret != nil {
					return ret
				}
			}
		}
	}
	return nil
}
