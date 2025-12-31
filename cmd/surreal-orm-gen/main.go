package main

import (
	"flag"
	"log"

	"github.com/yaroher/surrealdb.go.orm/internal/codegen"
)

func main() {
	var dir string
	flag.StringVar(&dir, "dir", ".", "package directory to scan")
	flag.Parse()

	if err := codegen.Generate(codegen.Options{Dirs: []string{dir}}); err != nil {
		log.Fatal(err)
	}
}
