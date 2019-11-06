package schema

import (
	"fmt"
	"strings"

	"github.com/clear-street/gogen-avro/generator"
	"github.com/clear-street/gogen-avro/imprt"
)

const enumTypeDef = `
%v
type %v int32

const (
%v
)
`

const enumTypeStringer = `
func (e %v) String() string {
	switch e {
%v
	}
	return "unknown"
}
`

const enumTypeParser = `
func Parse%v(val string) %v {
	switch val {
		%v
	}
	panic("unknown value: " + val)
}
`

const enumTypeIs = `
func Is%v(val string) bool {
	switch val {
		%v
	}
	return false
}
`

const enumSerializerDef = `
func %v(r %v, w io.Writer) error {
	return writeInt(int32(r), w)
}
`

type EnumDefinition struct {
	name       QualifiedName
	aliases    []QualifiedName
	symbols    []string
	doc        string
	definition map[string]interface{}
}

func NewEnumDefinition(name QualifiedName, aliases []QualifiedName, symbols []string, doc string, definition map[string]interface{}) *EnumDefinition {
	return &EnumDefinition{
		name:       name,
		aliases:    aliases,
		symbols:    symbols,
		doc:        doc,
		definition: definition,
	}
}

func (e *EnumDefinition) Name() string {
	return e.GoType()
}

func (e *EnumDefinition) SimpleName() string {
	return e.name.Name
}

func (e *EnumDefinition) AvroName() QualifiedName {
	return e.name
}

func (e *EnumDefinition) Aliases() []QualifiedName {
	return e.aliases
}

func (e *EnumDefinition) GoType() string {
	return generator.ToPublicName(e.name.Name)
}

func (e *EnumDefinition) typeList() string {
	typeStr := ""
	for i, t := range e.symbols {
		typeStr += fmt.Sprintf("%v %v = %v\n", generator.ToPublicName(e.GoType()+strings.Title(t)), e.GoType(), i)
	}
	return typeStr
}

func (e *EnumDefinition) stringerList() string {
	stringerStr := ""
	for _, t := range e.symbols {
		stringerStr += fmt.Sprintf("case %v:\n return %q\n", generator.ToPublicName(e.GoType()+strings.Title(t)), t)
	}
	return stringerStr
}

func (e *EnumDefinition) parserList() string {
	parserStr := ""
	for i, t := range e.symbols {
		parserStr += fmt.Sprintf("case %q:\n return %v\n", t, i)
	}
	return parserStr
}

func (e *EnumDefinition) isList() string {
	parserStr := ""
	for _, t := range e.symbols {
		parserStr += fmt.Sprintf("case %q:\n return true\n", t)
	}
	return parserStr
}

func (e *EnumDefinition) structDef() string {
	var doc string
	if e.doc != "" {
		doc = fmt.Sprintf("// %v", e.doc)
	}
	return fmt.Sprintf(enumTypeDef, doc, e.GoType(), e.typeList())
}

func (e *EnumDefinition) stringerDef() string {
	return fmt.Sprintf(enumTypeStringer, e.GoType(), e.stringerList())
}

func (e *EnumDefinition) parserDef() string {
	return fmt.Sprintf(enumTypeParser, e.GoType(), e.GoType(), e.parserList())
}

func (e *EnumDefinition) isDef() string {
	return fmt.Sprintf(enumTypeIs, e.GoType(), e.isList())
}

func (e *EnumDefinition) serializerMethodDef(p *generator.Package) string {
	return fmt.Sprintf(enumSerializerDef, e.SerializerMethod(p), e.GoType())
}

func (e *EnumDefinition) SerializerMethod(p *generator.Package) string {
	if !Contains(p, e) {
		pkg := imprt.Pkg(p.Root(), e.AvroName().Namespace)
		return fmt.Sprintf("%s.Write%s", pkg, e.GoType())
	}
	return "Write" + e.GoType()
}

func (e *EnumDefinition) filename() string {
	return generator.ToSnake(e.GoType()) + ".go"
}

func (e *EnumDefinition) AddStruct(p *generator.Package, _ bool) error {
	if p.Name() != e.AvroName().Namespace {
		return nil
	}

	p.AddStruct(e.filename(), e.GoType(), e.structDef())
	p.AddFunction(e.filename(), e.GoType(), "String", e.stringerDef())
	p.AddFunction(e.filename(), e.GoType(), "Parse", e.parserDef())
	p.AddFunction(e.filename(), e.GoType(), "Is", e.isDef())
	return nil
}

func (e *EnumDefinition) AddSerializer(p *generator.Package) {
	if !Contains(p, e) {
		p.AddImport(UTIL_FILE, imprt.Path(p.Root(), e.AvroName().Namespace))
		return
	}

	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeInt", writeIntMethod)
	p.AddFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.AddFunction(UTIL_FILE, "", e.SerializerMethod(p), e.serializerMethodDef(p))
	p.AddImport(UTIL_FILE, "io")
}

func (s *EnumDefinition) ResolveReferences(n *Namespace) error {
	return nil
}

func (s *EnumDefinition) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	if _, ok := scope[s.name]; ok {
		return s.name.String(), nil
	}
	scope[s.name] = 1
	return s.definition, nil
}

func (s *EnumDefinition) DefaultValue(p *generator.Package, lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(string); !ok {
		return "", fmt.Errorf("Expected string as default for field %v, got %q", lvalue, rvalue)
	}
	namespace := ""
	if p.Name() != s.AvroName().Namespace {
		lastDot := strings.LastIndex(s.AvroName().Namespace, ".")
		if lastDot >= 0 {
			namespace = s.AvroName().Namespace[lastDot+1:] + "."
		}
	}
	return fmt.Sprintf("%v = %v", lvalue, namespace+generator.ToPublicName(s.GoType()+strings.Title(rvalue.(string)))), nil
}

func (s *EnumDefinition) IsReadableBy(d Definition) bool {
	otherEnum, ok := d.(*EnumDefinition)
	return ok && otherEnum.name == s.name
}

func (s *EnumDefinition) WrapperType() string {
	return "types.Int"
}
