package schema

import (
	"fmt"

	"github.com/clear-street/gogen-avro/generator"
)

/*
  A named Reference to a user-defined type (fixed, enum, record). Just a wrapper with a name around a Definition.
*/

type Reference struct {
	TypeName QualifiedName
	Def      Definition
}

func NewReference(typeName QualifiedName) *Reference {
	return &Reference{
		TypeName: typeName,
	}
}

func (s *Reference) Name() string {
	return s.Def.Name()
}

func (s *Reference) AvroName() QualifiedName {
	return s.TypeName
}

func (s *Reference) GoType() string {
	return s.Def.GoType()
}

func (s *Reference) SerializerMethod(p *generator.Package) string {
	return s.Def.SerializerMethod(p)
}

func (s *Reference) SimpleName() string {
	return s.Def.SimpleName()
}

func (s *Reference) AddStruct(p *generator.Package, containers bool) error {
	return s.Def.AddStruct(p, containers)
}

func (s *Reference) AddSerializer(p *generator.Package) {
	s.Def.AddSerializer(p)
}

func (s *Reference) ResolveReferences(n *Namespace) error {
	if s.Def == nil {
		var ok bool
		if s.Def, ok = n.Definitions[s.TypeName]; !ok {
			t := QualifiedName{
				Name:      s.TypeName.Name,
				Namespace: "",
			}
			if s.Def, ok = n.Definitions[t]; !ok {
				return fmt.Errorf("Unable to resolve definition of type %v (%v,%v)\n", s.TypeName, s.TypeName.Namespace, s.TypeName.Name)
			}
		}
		return s.Def.ResolveReferences(n)
	}
	return nil
}

func (s *Reference) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	return s.Def.Definition(scope)
}

func (s *Reference) DefaultValue(p *generator.Package, lvalue string, rvalue interface{}) (string, error) {
	return s.Def.DefaultValue(p, lvalue, rvalue)
}

func (s *Reference) WrapperType() string {
	return s.Def.WrapperType()
}

func (s *Reference) IsReadableBy(f AvroType) bool {
	if reader, ok := f.(*Reference); ok {
		return s.Def.IsReadableBy(reader.Def)
	}
	return false
}
