package partialfields

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/kamiaka/go-partialfields/internal/pkg"
)

func NewAnalyzer() *analysis.Analyzer {
	var (
		f                          flag.FlagSet
		argExcludesEmptyLiteral    = true
		argExcludesOptionalField   = false
		argExcludesExportedField   = false
		argExcludesUnexportedField = false
		argExcludeFileRegExpPat    = ""
	)

	f.BoolVar(&argExcludesEmptyLiteral, "excludes-empty-literal", argExcludesEmptyLiteral, "excludes struct literal that has no field.")
	f.BoolVar(&argExcludesOptionalField, "excludes-optional-field", argExcludesOptionalField, "excludes field that has comment starts with optional.")
	f.BoolVar(&argExcludesExportedField, "excludes-exported-field", argExcludesExportedField, "excludes exported field")
	f.BoolVar(&argExcludesUnexportedField, "excludes-unexported-field", argExcludesUnexportedField, "excludes unexported field")
	f.StringVar(&argExcludeFileRegExpPat, "exclude-file-pattern", argExcludeFileRegExpPat, `exclude filename pattern (e.g., "_test\\.go$")`)

	return &analysis.Analyzer{
		Name:             "partialfields",
		Doc:              "check for all fields are defined in the struct",
		Flags:            f,
		RunDespiteErrors: false,
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
		ResultType: nil,
		FactTypes:  nil,
		Run: func(pass *analysis.Pass) (interface{}, error) {
			var skipFileRegExp *regexp.Regexp
			if argExcludeFileRegExpPat != "" {
				r, err := regexp.Compile(argExcludeFileRegExpPat)
				if err != nil {
					return nil, fmt.Errorf("invalid ")
				}

				skipFileRegExp = r
			}

			p := pkg.New()

			for _, f := range pass.Files {
				if isGeneratedFile(f) {
					continue
				}

				if err := p.SetImports(f.Imports); err != nil {
					return nil, err
				}

				v := &visitor{
					excludesExportedField:   argExcludesExportedField,
					excludesUnexportedField: argExcludesUnexportedField,
					excludeFilePattern:      skipFileRegExp,
					errors:                  nil,
					pkg:                     p,
					files:                   pass.Files,
					cmap:                    ast.NewCommentMap(pass.Fset, f, f.Comments),
					prev:                    nil,
					prevs:                   map[ast.Node]ast.Node{},

					pass: pass,
				}

				ast.Walk(v, f)

				// for k, ls := range v.cmap {
				// 	fmt.Println("")
				// 	for _, cg := range ls {
				// 		for _, c := range cg.List {
				// 			fmt.Println(c.Text)
				// 		}
				// 	}
				// 	fmt.Printf("%#v\n", k)
				// }

				// fmt.Println("===========================")
				// for k, ls := range v.cmap {
				// 	fmt.Printf("\n%#v\n", k)
				// 	if id, ok := k.(*ast.Ident); ok && id.Obj != nil {
				// 		fmt.Printf("  Decl: %#v\n", id.Obj.Decl)
				// 	}
				// 	for _, cg := range ls {
				// 		for _, c := range cg.List {
				// 			fmt.Printf("  %v\n", c.Text)
				// 		}
				// 	}
				// }
				// fmt.Println("---------------------------")

				// ast.Inspect(f, func(n ast.Node) bool {
				// 	switch n.(type) {
				// 	case *ast.StructType:
				// 		fmt.Printf("\nStruct")
				// 	default:
				// 		fmt.Printf("\nX: %#v\n", n)

				// 	}
				// 	return true
				// })
				fmt.Println("\n-------------------------")

				for _, err := range v.errors {
					var e *IncompleteError
					if !errors.As(err, &e) {
						return nil, err
					}
					pass.Report(analysis.Diagnostic{
						Pos:     e.Pos(),
						End:     e.End(),
						Message: e.Error(),
					})
				}
			}
			return nil, nil
		},
	}
}

var generatedCommentPat = regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`)

func isGeneratedFile(f *ast.File) bool {
	for _, g := range f.Comments {
		for _, c := range g.List {
			if generatedCommentPat.MatchString(c.Text) {
				return true
			}
		}
	}
	return false
}
