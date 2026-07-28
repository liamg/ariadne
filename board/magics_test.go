//go:build magicgen

package board

import (
	"go/format"
	"testing"
)

func TestGenerateMagics(t *testing.T) {
	// this is a sanity test that valid go is produced by the generator
	// the values themselves are checked by the magics_test.go file
	data, err := GenerateMagics()
	if err != nil {
		t.Fatalf("GenerateMagics() returned an error: %v", err)
	}

	_, err = format.Source(data) // check that the generated code is valid Go
	if err != nil {
		t.Fatalf("Generated code is not valid Go: %v", err)
	}
}
