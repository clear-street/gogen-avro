package schema

import (
	"github.com/clear-street/gogen-avro/generator"
)

const writeNullMethod = `
func writeNull(_ interface{}, _ io.Writer) error {
	return nil
}
`

type NullField struct {
	PrimitiveField
}

func NewNullField(definition interface{}) *NullField {
	return &NullField{PrimitiveField{
		definition:       definition,
		name:             "Null",
		goType:           "*types.NullVal",
		serializerMethod: "writeNull",
	}}
}

func (s *NullField) AddSerializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", "writeNull", writeNullMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *NullField) DefaultValue(p *generator.Package, lvalue string, rvalue interface{}) (string, error) {
	return "", nil
}

func (s *NullField) WrapperType() string {
	return ""
}

func (s *NullField) IsReadableBy(f AvroType) bool {
	_, ok := f.(*NullField)
	return ok
}
