package schema

import (
	"fmt"

	"github.com/clear-street/gogen-avro/generator"
)

const writeFixedMethod = `
func %v(r %v, w io.Writer) error {
	_, err := w.Write(r[:])
	return err
}
`

const fixedFieldTemplate = `
type %[1]v %[2]v

func (_ *%[1]v) DeserializeBoolean(v bool) { panic("Unsupported operation") }
func (_ *%[1]v) DeserializeInt(v int32) { panic("Unsupported operation") }
func (_ *%[1]v) DeserializeLong(v int64) { panic("Unsupported operation") }
func (_ *%[1]v) DeserializeFloat(v float32) { panic("Unsupported operation") }
func (_ *%[1]v) DeserializeDouble(v float64) { panic("Unsupported operation") }
func (r *%[1]v) DeserializeBytes(v []byte) { 
	copy((*r)[:], v)
}
func (_ *%[1]v) DeserializeString(v string) { panic("Unsupported operation") }
func (_ *%[1]v) SetUnionElem(v int64) { panic("Unsupported operation") }
func (_ *%[1]v) Get(i int) types.Field { panic("Unsupported operation") }
func (_ *%[1]v) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *%[1]v) AppendArray() types.Field { panic("Unsupported operation") }
func (_ *%[1]v) Finalize() { }
func (_ *%[1]v) SetDefault(i int) { panic("Unsupported operation") }
`

type FixedDefinition struct {
	name       QualifiedName
	aliases    []QualifiedName
	sizeBytes  int
	definition map[string]interface{}
}

func NewFixedDefinition(name QualifiedName, aliases []QualifiedName, sizeBytes int, definition map[string]interface{}) *FixedDefinition {
	return &FixedDefinition{
		name:       name,
		aliases:    aliases,
		sizeBytes:  sizeBytes,
		definition: definition,
	}
}

func (s *FixedDefinition) Name() string {
	return s.GoType()
}

func (s *FixedDefinition) SimpleName() string {
	return generator.ToPublicSimpleName(s.name.Name)
}

func (s *FixedDefinition) AvroName() QualifiedName {
	return s.name
}

func (s *FixedDefinition) Aliases() []QualifiedName {
	return s.aliases
}

func (s *FixedDefinition) GoType() string {
	return generator.ToPublicName(s.name.Name)
}

func (s *FixedDefinition) SizeBytes() int {
	return s.sizeBytes
}

func (s *FixedDefinition) serializerMethodDef(p *generator.Package) string {
	return fmt.Sprintf(writeFixedMethod, s.SerializerMethod(p), s.GoType())
}

func (s *FixedDefinition) typeDef() string {
	return fmt.Sprintf("type %v [%v]byte\n", s.GoType(), s.sizeBytes)
}

func (s *FixedDefinition) filename() string {
	return generator.ToSnake(s.GoType()) + ".go"
}

func (s *FixedDefinition) SerializerMethod(p *generator.Package) string {
	return fmt.Sprintf("write%v", s.GoType())
}

func (s *FixedDefinition) AddStruct(p *generator.Package, _ bool) error {
	p.AddStruct(s.filename(), s.GoType(), s.typeDef())
	return nil
}

func (s *FixedDefinition) AddSerializer(p *generator.Package) {
	p.AddImport(UTIL_FILE, "io")
	p.AddImport(UTIL_FILE, "github.com/clear-street/gogen-avro/vm/types")
	p.AddFunction(UTIL_FILE, "", s.SerializerMethod(p), s.serializerMethodDef(p))
	p.AddFunction(UTIL_FILE, s.GoType(), "fieldTemplate", s.FieldsMethodDef())
}

func (s *FixedDefinition) ResolveReferences(n *Namespace) error {
	return nil
}

func (s *FixedDefinition) FieldsMethodDef() string {
	return fmt.Sprintf(fixedFieldTemplate, s.WrapperType(), s.GoType(), s.sizeBytes)
}

func (s *FixedDefinition) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	if _, ok := scope[s.name]; ok {
		return s.name.String(), nil
	}
	scope[s.name] = 1
	return s.definition, nil
}

func (s *FixedDefinition) DefaultValue(p *generator.Package, lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(string); !ok {
		return "", fmt.Errorf("Expected string as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = []byte(%q)", lvalue, rvalue), nil
}

func (s *FixedDefinition) IsReadableBy(d Definition) bool {
	if fixed, ok := d.(*FixedDefinition); ok {
		return fixed.sizeBytes == s.sizeBytes && fixed.name == s.name
	}
	return false
}

func (s *FixedDefinition) WrapperType() string {
	return fmt.Sprintf("%vWrapper", s.GoType())
}
