package main

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"os"

	"golang.org/x/tools/imports"
)

func fmtFile(fileName string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// gofmt the modified file
	var buf bytes.Buffer
	if err = format.Node(&buf, fset, node); err != nil {
		return err
	}

	formattedCode, err := imports.Process(fileName, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("failed to perform goimports: %w", err)
	}

	if err := os.WriteFile(fileName, formattedCode, 0o600); err != nil {
		return err
	}

	return nil
}
