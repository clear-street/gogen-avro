package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/clear-street/gogen-avro/generator"
	"github.com/clear-street/gogen-avro/imprt"
	"github.com/clear-street/gogen-avro/types"
)

func main() {
	packageName := flag.String("package", "avro", "Root package")
	containers := flag.Bool("containers", false, "Whether to generate container writer methods")
	shortUnions := flag.Bool("short-unions", false, "Whether to use shorter names for Union types")

	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "Usage: gogen-avro [--short-unions] [--package=<root package>] [--containers] <target directory> <schema files>\n")
		os.Exit(1)
	}
	fmt.Println(packageName)

	targetDir := flag.Arg(0)
	files := flag.Args()[1:]
	namespace := types.NewNamespace(*shortUnions)

	for _, fileName := range files {
		schema, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %q - %v\n", fileName, err)
			os.Exit(2)
		}

		_, err = namespace.TypeForSchema(schema)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding schema for file %q - %v\n", fileName, err)
			os.Exit(3)
		}
	}

	for _, v := range namespace.Schemas {
		if err := v.Root.ResolveReferences(namespace); err != nil {
			panic(err)
		}
	}

	pkgs := map[string]*generator.Package{}
	for k, v := range namespace.Definitions {
		pkg, ok := pkgs[k.Namespace]
		if !ok {
			pkg = generator.NewPackage(*packageName, k.Namespace)
			pkgs[k.Namespace] = pkg
		}

		v.AddStruct(pkg, *containers)
		v.AddSerializer(pkg)
		v.AddDeserializer(pkg)
	}

	if err := os.RemoveAll(filepath.Join(targetDir, *packageName)); err != nil {
		panic(err)
	}

	for k, v := range pkgs {
		path := filepath.Join(targetDir, imprt.Path(*packageName, k))
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			panic(err)
		}
		err := v.WriteFiles(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing source files to directory %q - %v\n", path, err)
			os.Exit(4)
		}
	}
}

// codegenComment generates a comment informing readers they are looking at
// generated code and lists the source avro files used to generate the code
//
// invariant: sources > 0
func codegenComment(sources []string) string {
	const fileComment = `// Code generated by github.com/clear-street/gogen-avro. DO NOT EDIT.
/*
 * %s
 */`
	var sourceBlock []string
	if len(sources) == 1 {
		sourceBlock = append(sourceBlock, "SOURCE:")
	} else {
		sourceBlock = append(sourceBlock, "SOURCES:")
	}

	for _, source := range sources {
		_, fName := filepath.Split(source)
		sourceBlock = append(sourceBlock, fmt.Sprintf(" *     %s", fName))
	}

	return fmt.Sprintf(fileComment, strings.Join(sourceBlock, "\n"))
}
