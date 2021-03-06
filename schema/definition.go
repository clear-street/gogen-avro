package schema

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
	SimpleName() string

	GoType() string

	SerializerMethod(*generator.Package) string

	// Add the imports and struct for the definition of this type to the generator.Package
	AddStruct(*generator.Package, bool) error
	AddSerializer(*generator.Package)

	// Resolve references to user-defined types
	ResolveReferences(*Namespace) error

	// A JSON object defining this object, for writing the schema back out
	Definition(scope map[QualifiedName]interface{}) (interface{}, error)
	DefaultValue(p *generator.Package, lvalue string, rvalue interface{}) (string, error)
	IsReadableBy(f Definition) bool
	WrapperType() string
}

//Contains returns whether the package contains the definition
func Contains(p *generator.Package, def HasAvroName) bool {
	pkgPath := imprt.Path(p.Root(), p.Name())
	defPath := imprt.Path(p.Root(), def.AvroName().Namespace)
	return pkgPath == defPath
}
