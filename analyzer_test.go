package partialfields

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_a(t *testing.T) {
	testdata := analysistest.TestData()
	a := NewAnalyzer()

	analysistest.Run(t, testdata, a, "tmp")
}
