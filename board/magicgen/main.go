//go:build magicgen

package main

import (
	"fmt"
	"go/format"
	"os"
	"time"

	"github.com/liamg/chess/board"
)

func main() {
	start := time.Now()
	fmt.Println("Generating magics...")
	data, err := board.GenerateMagics()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error generating magics:", err)
		os.Exit(1)
	}

	formatted, err := format.Source(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error formatting generated code:", err)
		os.Exit(1)
	}

	if err := os.WriteFile("magics_gen.go", formatted, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing file:", err)
		os.Exit(1)
	}
	fmt.Println("Magics generated successfully and written to magics.go in", time.Since(start))
}
