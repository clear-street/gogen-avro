package schema

import (
	"github.com/clear-street/gogen-avro/generator"
)

type Constructable interface {
	ConstructorMethod(p *generator.Package) string
}

func getConstructableForType(t AvroType) (Constructable, bool) {
	if c, ok := t.(Constructable); ok {
		return c, true
	}
	if ref, ok := t.(*Reference); ok {
		if c, ok := ref.Def.(Constructable); ok {
			return c, true
		}
	}
	return nil, false
}
