package pkg

import (
	"fmt"
	"go/ast"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Pkg struct {
	pathByName    map[string]string
	pkgListByPath map[string]List
}

func New() *Pkg {
	return &Pkg{
		pathByName:    map[string]string{},
		pkgListByPath: map[string]List{},
	}
}

func (p *Pkg) SetImports(specs []*ast.ImportSpec) error {
	for _, spec := range specs {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return err
		}

		var name string
		if spec.Name != nil {
			name = spec.Name.Name
		} else {
			name = path[strings.LastIndexByte(path, '/')+1:]
		}
		p.pathByName[name] = path
	}
	return nil
}

func (p *Pkg) List(name string) (List, error) {
	path, ok := p.pathByName[name]
	if !ok {
		return nil, fmt.Errorf("unknown package name: %s", name)
	}
	if pkgs, ok := p.pkgListByPath[path]; ok {
		return pkgs, nil
	}

	pkgs, err := p.load(path)
	if err != nil {
		return nil, err
	}
	p.pkgListByPath[path] = pkgs

	return pkgs, nil
}

func (p *Pkg) load(path string) (List, error) {
	// partial:packages.Config
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedSyntax,
	}, path)
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}
	for _, pkg := range pkgs {
		if len(pkg.Errors) != 0 {
			for _, err := range pkg.Errors {
				fmt.Printf("err: %#v\n", err)
			}
			// return nil, fmt.Errorf("ERR: %#v\n", pkg.Errors)
		}
	}
	return pkgs, nil
}
