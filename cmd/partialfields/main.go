package main

import (
	partialfields "github.com/kamiaka/go-partialfields"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(
		partialfields.NewAnalyzer(),
	)
}
