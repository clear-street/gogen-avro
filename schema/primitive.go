package schema

import (
	"github.com/clear-street/gogen-avro/generator"
)

// Common methods for all primitive types
type PrimitiveField struct {
	definition       interface{}
	name             string
	goType           string
	serializerMethod string
}

func (s *PrimitiveField) Name() string {
	return s.name
}

func (s *PrimitiveField) GoType() string {
	return s.goType
}

func (s *PrimitiveField) SerializerMethod(p *generator.Package) string {
	return s.serializerMethod
}

func (s *PrimitiveField) AddStruct(p *generator.Package, _ bool) error {
	return nil
}

func (s *PrimitiveField) ResolveReferences(n *Namespace) error {
	return nil
}

func (s *PrimitiveField) Definition(_ map[QualifiedName]interface{}) (interface{}, error) {
	return s.definition, nil
}

func (s *PrimitiveField) SimpleName() string {
	return s.name
}
