package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
)

var (
	src         = flag.String("s", "", "src file")
	dest        = flag.String("d", "", "dest file")
	packageName = flag.String("p", "", "package name")
)

func main() {
	flag.Parse()

	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, *src, nil, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	ast.Print(fset, f)
}
