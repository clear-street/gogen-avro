package types

import (
	"github.com/clear-street/gogen-avro/generator"
	"github.com/clear-street/gogen-avro/imprt"
)

type HasAvroName interface {
	AvroName() QualifiedName
}

/*
  The definition of a record, fixed or enum satisfies this interface.
*/

type Definition interface {
	HasAvroName
	Aliases() []QualifiedName

	// A user-friendly name that can be built into a Go string (for unions, mostly)
	Name() string

	GoType() string

	SerializerMethod(*generator.Package) string
	DeserializerMethod(*generator.Package) string

	// Add the imports and struct for the definition of this type to the generator.Package
	AddStruct(*generator.Package, bool) error
	AddSerializer(*generator.Package)
	AddDeserializer(*generator.Package)

	// Resolve references to user-defined types
	ResolveReferences(*Namespace) error

	// A JSON object defining this object, for writing the schema back out
	Definition(scope map[QualifiedName]interface{}) (interface{}, error)
	DefaultValue(lvalue string, rvalue interface{}) (string, error)
}

//Contains returns whether the package contains the definition
func Contains(p *generator.Package, def HasAvroName) bool {
	pkgPath := imprt.Path(p.Root(), p.Name())
	defPath := imprt.Path(p.Root(), def.AvroName().Namespace)
	return pkgPath == defPath
}
