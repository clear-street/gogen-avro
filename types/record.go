package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/clear-street/gogen-avro/generator"
	"github.com/clear-street/gogen-avro/imprt"
)

const recordStructDefTemplate = `
%v
type %v struct {
%v
}
`

const recordSchemaTemplate = `func (r %v) Schema() string {
 return %v
}
`

const recordQualifiedName = `func (r %v) QualifiedName() (string, string) {
	return %q, %q
}
`

const recordConstructorTemplate = `
	func %v %v {
		v := &%v{
			%v
		}
		%v
		return v
	}
`

const recordStructPublicSerializerTemplate = `
func (r %v) Serialize(w io.Writer) error {
	return %v(r, w)
}
`

const recordStructDeserializerTemplate = `
func %v(r io.Reader) (%v, error) {
	var str = &%v{}
	var err error
	%v
	return str, nil
}
`

const recordStructPublicDeserializerTemplate = `
func %v(r io.Reader) (%v, error) {
	return %v(r)
}
`

const recordWriterTemplate = `
func %v(writer io.Writer, codec container.Codec, recordsPerBlock int64) (*container.Writer, error) {
	str := &%v{}
	return container.NewWriter(writer, codec, recordsPerBlock, str.Schema())
}
`

type RecordDefinition struct {
	name     QualifiedName
	aliases  []QualifiedName
	fields   []*Field
	doc      string
	metadata map[string]interface{}
}

func NewRecordDefinition(name QualifiedName, aliases []QualifiedName, fields []*Field, doc string, metadata map[string]interface{}) *RecordDefinition {
	return &RecordDefinition{
		name:     name,
		aliases:  aliases,
		fields:   fields,
		doc:      doc,
		metadata: metadata,
	}
}

func (r *RecordDefinition) AvroName() QualifiedName {
	return r.name
}

func (r *RecordDefinition) Name() string {
	return generator.ToPublicName(r.name.Name)
}

func (r *RecordDefinition) GoType() string {
	return fmt.Sprintf("*%v", r.Name())
}

func (r *RecordDefinition) Aliases() []QualifiedName {
	return r.aliases
}

func (r *RecordDefinition) structFields(p *generator.Package) string {
	var fieldDefinitions string
	for _, f := range r.fields {
		if f.Doc() != "" {
			fieldDefinitions += fmt.Sprintf("\n// %v\n", f.Doc())
		}

		if ref, ok := f.avroType.(*Reference); ok && ref.AvroName().Namespace != r.AvroName().Namespace {
			t := imprt.Type(p.Root(), ref.AvroName().Namespace, f.Type().GoType())
			fieldDefinitions += fmt.Sprintf("%v %v\n", f.GoName(), t)
		} else {
			fieldDefinitions += fmt.Sprintf("%v %v\n", f.GoName(), f.Type().GoType())
		}
	}
	return fieldDefinitions
}

func (r *RecordDefinition) fieldSerializers(p *generator.Package) string {
	serializerMethods := "var err error\n"
	for _, f := range r.fields {
		serializerMethods += fmt.Sprintf("err = %v(r.%v, w)\nif err != nil {return err}\n", f.Type().SerializerMethod(p), f.GoName())
	}
	return serializerMethods
}

func (r *RecordDefinition) fieldDeserializers(p *generator.Package) string {
	deserializerMethods := ""
	for _, f := range r.fields {
		deserializerMethods += fmt.Sprintf("str.%v, err = %v(r)\nif err != nil {return nil, err}\n", f.GoName(), f.Type().DeserializerMethod(p))
	}
	return deserializerMethods
}

func (r *RecordDefinition) structDefinition(p *generator.Package) string {
	var doc string
	if r.doc != "" {
		doc = fmt.Sprintf("// %v", r.doc)
	}
	return fmt.Sprintf(recordStructDefTemplate, doc, r.Name(), r.structFields(p))
}

func (r *RecordDefinition) serializerMethodDef(p *generator.Package) string {
	return fmt.Sprintf("func %v(r %v, w io.Writer) error {\n%v\nreturn nil\n}", r.SerializerMethod(p), r.GoType(), r.fieldSerializers(p))
}

func (r *RecordDefinition) deserializerMethodDef(p *generator.Package) string {
	return fmt.Sprintf(recordStructDeserializerTemplate, r.DeserializerMethod(p), r.GoType(), r.Name(), r.fieldDeserializers(p))
}

func (r *RecordDefinition) SerializerMethod(p *generator.Package) string {
	if p.Name() != r.AvroName().Namespace {
		pkg := imprt.Pkg(p.Root(), r.AvroName().Namespace)
		return fmt.Sprintf("%s.Write%s", pkg, r.Name())
	}

	return fmt.Sprintf("Write%v", r.Name())
}

func (r *RecordDefinition) DeserializerMethod(p *generator.Package) string {
	if p.Name() != r.AvroName().Namespace {
		pkg := imprt.Pkg(p.Root(), r.AvroName().Namespace)
		return fmt.Sprintf("%s.Read%s", pkg, r.Name())
	}

	return fmt.Sprintf("Read%v", r.Name())
}

func (r *RecordDefinition) publicDeserializerMethod() string {
	return fmt.Sprintf("Deserialize%v", r.Name())
}

func (r *RecordDefinition) recordWriterMethod() string {
	return fmt.Sprintf("New%vWriter", r.Name())
}

func (r *RecordDefinition) recordWriterMethodDef() string {
	return fmt.Sprintf(recordWriterTemplate, r.recordWriterMethod(), r.Name())
}

func (r *RecordDefinition) publicSerializerMethodDef(p *generator.Package) string {
	return fmt.Sprintf(recordStructPublicSerializerTemplate, r.GoType(), r.SerializerMethod(p))
}

func (r *RecordDefinition) publicDeserializerMethodDef(p *generator.Package) string {
	return fmt.Sprintf(recordStructPublicDeserializerTemplate, r.publicDeserializerMethod(), r.GoType(), r.DeserializerMethod(p))
}

func (r *RecordDefinition) filename() string {
	return generator.ToSnake(r.Name()) + ".go"
}

func (r *RecordDefinition) schemaMethodDef() (string, error) {
	def, err := r.Definition(make(map[QualifiedName]interface{}))
	if err != nil {
		return "", err
	}

	schemaJson, _ := json.Marshal(def)
	return fmt.Sprintf(recordSchemaTemplate, r.GoType(), strconv.Quote(string(schemaJson))), nil
}

func (r *RecordDefinition) qualifiedNameMethodDef() (string, error) {
	avroName := r.AvroName()
	return fmt.Sprintf(recordQualifiedName, r.GoType(), avroName.Namespace, avroName.Name), nil
}

func (r *RecordDefinition) AddStruct(p *generator.Package, containers bool) error {
	if p.Name() != r.AvroName().Namespace {
		return nil
	}

	// Import guard, to avoid circular dependencies
	if !p.HasStruct(r.filename(), r.GoType()) {
		p.AddStruct(r.filename(), r.GoType(), r.structDefinition(p))
		schemaDef, err := r.schemaMethodDef()
		if err != nil {
			return err
		}

		p.AddFunction(r.filename(), r.GoType(), "Schema", schemaDef)
		constructorMethodDef, err := r.ConstructorMethodDef()
		if err != nil {
			return err
		}

		for _, f := range r.fields {
			ref, ok := f.avroType.(*Reference)
			if !ok || ref.AvroName().Namespace == r.AvroName().Namespace {
				continue
			}
			p.AddImport(r.filename(), imprt.Path(p.Root(), ref.AvroName().Namespace))
		}

		qnDef, err := r.qualifiedNameMethodDef()
		if err != nil {
			return err
		}

		p.AddFunction(r.filename(), r.GoType(), "QualifiedName", qnDef)

		if containers {
			p.AddImport(r.filename(), "github.com/clear-street/gogen-avro/container")
			p.AddFunction(r.filename(), "", r.recordWriterMethod(), r.recordWriterMethodDef())
		}

		p.AddFunction(r.filename(), r.GoType(), r.ConstructorMethod(), constructorMethodDef)
		for _, f := range r.fields {
			f.Type().AddStruct(p, containers)
		}
	}
	return nil
}

func (r *RecordDefinition) AddSerializer(p *generator.Package) {
	if p.Name() != r.AvroName().Namespace {
		p.AddImport(UTIL_FILE, imprt.Path(p.Root(), r.AvroName().Namespace))
		return
	}

	// Import guard, to avoid circular dependencies
	if !p.HasFunction(UTIL_FILE, "", r.SerializerMethod(p)) {
		p.AddImport(r.filename(), "io")
		p.AddFunction(UTIL_FILE, "", r.SerializerMethod(p), r.serializerMethodDef(p))
		p.AddFunction(r.filename(), r.GoType(), "Serialize", r.publicSerializerMethodDef(p))
		for _, f := range r.fields {
			f.Type().AddSerializer(p)
		}
	}
}

func (r *RecordDefinition) AddDeserializer(p *generator.Package) {
	if p.Name() != r.AvroName().Namespace {
		p.AddImport(UTIL_FILE, imprt.Path(p.Root(), r.AvroName().Namespace))
		return
	}

	// Import guard, to avoid circular dependencies
	if !p.HasFunction(UTIL_FILE, "", r.DeserializerMethod(p)) {
		p.AddImport(r.filename(), "io")
		p.AddFunction(UTIL_FILE, "", r.DeserializerMethod(p), r.deserializerMethodDef(p))
		p.AddFunction(r.filename(), "", r.publicDeserializerMethod(), r.publicDeserializerMethodDef(p))
		for _, f := range r.fields {
			f.Type().AddDeserializer(p)
		}
	}
}

func (r *RecordDefinition) ResolveReferences(n *Namespace) error {
	var err error
	for _, f := range r.fields {
		err = f.Type().ResolveReferences(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RecordDefinition) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	if _, ok := scope[r.name]; ok {
		return r.name.String(), nil
	}
	scope[r.name] = 1
	fields := make([]map[string]interface{}, 0)
	for _, f := range r.fields {
		def, err := f.Definition(scope)
		if err != nil {
			return nil, err
		}
		fields = append(fields, def)
	}

	r.metadata["fields"] = fields
	return r.metadata, nil
}

func (r *RecordDefinition) ConstructorMethod() string {
	return fmt.Sprintf("New%v()", r.Name())
}

func (r *RecordDefinition) fieldConstructors() (string, error) {
	constructors := ""
	for _, f := range r.fields {
		if constructor, ok := getConstructableForType(f.Type()); ok {
			constructors += fmt.Sprintf("%v: %v,\n", f.GoName(), constructor.ConstructorMethod())
		}
	}
	return constructors, nil
}

func (r *RecordDefinition) defaultValues() (string, error) {
	defaults := ""
	for _, f := range r.fields {
		if f.hasDef {
			def, err := f.Type().DefaultValue(fmt.Sprintf("v.%v", f.GoName()), f.Default())
			if err != nil {
				return "", err
			}
			defaults += def + "\n"
		}
	}
	return defaults, nil
}

func (r *RecordDefinition) ConstructorMethodDef() (string, error) {
	defaults, err := r.defaultValues()
	if err != nil {
		return "", err
	}

	fieldConstructors, err := r.fieldConstructors()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(recordConstructorTemplate, r.ConstructorMethod(), r.GoType(), r.Name(), fieldConstructors, defaults), nil
}

func (r *RecordDefinition) FieldByName(name string) *Field {
	for _, f := range r.fields {
		if f.Name() == name {
			return f
		}
	}
	return nil
}

func (r *RecordDefinition) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	items := rvalue.(map[string]interface{})
	fieldSetters := ""
	for k, v := range items {
		field := r.FieldByName(k)
		fieldSetter, err := field.Type().DefaultValue(fmt.Sprintf("%v.%v", lvalue, field.GoName()), v)
		if err != nil {
			return "", err
		}

		fieldSetters += fieldSetter + "\n"
	}
	return fieldSetters, nil
}
