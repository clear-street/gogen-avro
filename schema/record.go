package schema

import (
	"encoding/json"
	"fmt"
	"sort"
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
const recordQualifiedName = `func (r %v) QualifiedName() string {
	return %q + "." + %q
}
`

const recordSchemaNameTemplate = `func (r %v) SchemaName() string {
 return %v
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

const recordStructPublicDeserializerTemplate = `
func %v(r io.Reader) (%v, error) {
	t := %v

	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
        if err != nil {
		return nil, err
	}

        err = vm.Eval(r, deser, t)
	return t, err
}
`

const recordWriterTemplate = `
func %v(writer io.Writer, codec container.Codec, recordsPerBlock int64) (*container.Writer, error) {
	str := &%v{}
	return container.NewWriter(writer, codec, recordsPerBlock, str.Schema())
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

const recordFieldTemplate = `
func (_ %[1]v) DeserializeBoolean(v bool) { panic("Unsupported operation") }
func (_ %[1]v) DeserializeInt(v int32) { panic("Unsupported operation") }
func (_ %[1]v) DeserializeLong(v int64) { panic("Unsupported operation") }
func (_ %[1]v) DeserializeFloat(v float32) { panic("Unsupported operation") }
func (_ %[1]v) DeserializeDouble(v float64) { panic("Unsupported operation") }
func (_ %[1]v) DeserializeBytes(v []byte) { panic("Unsupported operation") }
func (_ %[1]v) DeserializeString(v string) { panic("Unsupported operation") }
func (_ %[1]v) SetUnionElem(v int64) { panic("Unsupported operation") }
func (r %[1]v) Get(i int) types.Field {
	switch (i) {
		%[2]v
	}
	panic("Unknown field index")
}
func (r %[1]v) SetDefault(i int) {
	switch (i) {
		%[3]v
	}
	panic("Unknown field index")
}
func (_ %[1]v) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ %[1]v) AppendArray() types.Field { panic("Unsupported operation") }
func (_ %[1]v) Finalize() { }
`

const recordReaderTemplate = `
type %[1]v struct {
	r io.Reader
	p *vm.Program
}

func New%[1]v(r io.Reader) (*%[1]v, error){
	containerReader, err := container.NewReader(r)
	if err != nil {
		return nil, err
	}

	t := %[3]v
	deser, err := compiler.CompileSchemaBytes([]byte(containerReader.AvroContainerSchema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	return &%[1]v{
		r: containerReader,
		p: deser,
	}, nil
}

func (r *%[1]v) Read() (%[2]v, error) {
	t := %[3]v
        err := vm.Eval(r.r, r.p, t)
	return t, err
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
	return generator.ToPublicName(r.name.String())
}

func (r *RecordDefinition) SimpleName() string {
	return generator.ToPublicName(r.name.Name)
}

func (r *RecordDefinition) GoType() string {
	return fmt.Sprintf("*%v", r.Name())
}

func (r *RecordDefinition) Aliases() []QualifiedName {
	return r.aliases
}

func (r *RecordDefinition) structFields(p *generator.Package) string {
	var definitions string
	for _, f := range r.fields {
		var field string

		// Prepend doc if exists
		if f.Doc() != "" {
			field += fmt.Sprintf("\n// %v\n", f.Doc())
		}

		field += fmt.Sprintf("%v %v", f.SimpleName(), f.Type().GoType())

		if f.Tags() != "" {
			field += " `" + f.Tags() + "`"
		}

		if ref, ok := f.avroType.(*Reference); ok && !Contains(p, ref) {
			t := imprt.Type(p.Root(), ref.AvroName().Namespace, f.Type().GoType())
			definitions += fmt.Sprintf("%v %v\n", f.GoName(), t)
		} else {
			definitions += fmt.Sprintf("%v %v\n", f.GoName(), f.Type().GoType())
		}
	}

	return definitions
}

func (r *RecordDefinition) fieldSerializers(p *generator.Package) string {
	if r.fields == nil || len(r.fields) == 0 {
		//in case the record has no fields just return empty fieldSerializers
		return ""
	}
	serializerMethods := "var err error\n"
	for _, f := range r.fields {
		serializerMethods += fmt.Sprintf("err = %v(r.%v, w)\nif err != nil {return err}\n", f.Type().SerializerMethod(p), f.GoName())
	}
	return serializerMethods
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

func (r *RecordDefinition) SerializerMethod(p *generator.Package) string {
	if !Contains(p, r) {
		pkg := imprt.Pkg(p.Root(), r.AvroName().Namespace)
		return fmt.Sprintf("%s.Write%s", pkg, r.Name())
	}

	return fmt.Sprintf("Write%v", r.Name())
}

func (r *RecordDefinition) DeserializerMethod(p *generator.Package) string {
	if !Contains(p, r) {
		pkg := imprt.Pkg(p.Root(), r.AvroName().Namespace)
		return fmt.Sprintf("%s.Read%s", pkg, r.Name())
	}

	return fmt.Sprintf("Read%v", r.Name())
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

func (r *RecordDefinition) publicDeserializerMethod() string {
	return fmt.Sprintf("Deserialize%v", r.Name())
}

func (r *RecordDefinition) publicDeserializerMethodDef(p *generator.Package) string {
	return fmt.Sprintf(recordStructPublicDeserializerTemplate, r.publicDeserializerMethod(), r.GoType(), r.ConstructorMethod(p))
}

func (r *RecordDefinition) schemaNameMethodDef() (string, error) {
	return fmt.Sprintf(recordSchemaNameTemplate, r.GoType(), strconv.Quote(r.name.String())), nil
}

func (r *RecordDefinition) AddStruct(p *generator.Package, containers bool) error {
	if !Contains(p, r) {
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
		constructorMethodDef, err := r.ConstructorMethodDef(p)
		if err != nil {
			return err
		}

		for _, f := range r.fields {
			ref, ok := f.avroType.(*Reference)
			if !ok || Contains(p, ref) {
				continue
			}
			p.AddImport(r.filename(), imprt.Path(p.Root(), ref.AvroName().Namespace))
		}

		qnDef, err := r.qualifiedNameMethodDef()
		if err != nil {
			return err
		}

		p.AddFunction(r.filename(), r.GoType(), "QualifiedName", qnDef)
		p.AddConstant(r.filename(), r.AvroName().Name+"QualifiedName", r.AvroName().String())
		p.AddConstant(r.filename(), r.AvroName().Name+"Name", r.AvroName().Name)
		p.AddConstant(r.filename(), r.AvroName().Name+"Namespace", r.AvroName().Namespace)
		schemaNameDef, err := r.schemaNameMethodDef()
		if err != nil {
			return err
		}

		p.AddFunction(r.filename(), r.GoType(), "SchemaName", schemaNameDef)

		if containers {
			p.AddImport(r.filename(), "github.com/clear-street/gogen-avro/container")
			p.AddFunction(r.filename(), "", r.recordWriterMethod(), r.recordWriterMethodDef())
		}

		p.AddImport(r.filename(), "github.com/clear-street/gogen-avro/vm/types")
		p.AddImport(r.filename(), "github.com/clear-street/gogen-avro/vm")
		p.AddImport(r.filename(), "github.com/clear-street/gogen-avro/compiler")
		p.AddFunction(r.filename(), r.GoType(), "fieldTemplate", r.FieldsMethodDef(p))
		p.AddFunction(r.filename(), r.GoType(), "recordReader", r.recordReaderDef(p))
		p.AddFunction(r.filename(), r.GoType(), r.ConstructorMethod(p), constructorMethodDef)
		p.AddFunction(r.filename(), r.GoType(), r.publicDeserializerMethod(), r.publicDeserializerMethodDef(p))
		for _, f := range r.fields {
			f.Type().AddStruct(p, containers)
		}
	}
	return nil
}

func (r *RecordDefinition) AddSerializer(p *generator.Package) {
	if !Contains(p, r) {
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

func (r *RecordDefinition) ConstructorMethod(p *generator.Package) string {
	pkg := ""
	if !Contains(p, r) {
		pkg = fmt.Sprintf("%v.", imprt.Pkg(p.Root(), r.AvroName().Namespace))
	}
	return fmt.Sprintf("%vNew%v()", pkg, r.Name())
}

func (r *RecordDefinition) fieldConstructors(p *generator.Package) (string, error) {
	constructors := ""
	for _, f := range r.fields {
		if constructor, ok := getConstructableForType(f.Type()); ok {
			constructors += fmt.Sprintf("%v: %v,\n", f.GoName(), constructor.ConstructorMethod(p))
		}
	}
	return constructors, nil
}

func (r *RecordDefinition) defaultMethodDef(p *generator.Package) (string, error) {
	defaults := ""
	for i, f := range r.fields {
		if f.hasDef {
			defaults += fmt.Sprintf("case %v:\n", i)
			def, err := f.Type().DefaultValue(p, fmt.Sprintf("r.%v", f.GoName()), f.Default())
			if err != nil {
				return "", err
			}
			defaults += def + "\nreturn\n"
		}
	}
	return defaults, nil
}

func (r *RecordDefinition) getMethodDef(p *generator.Package) string {
	getBody := ""
	for i, f := range r.fields {
		getBody += fmt.Sprintf("case %v:\n", i)
		if constructor, ok := getConstructableForType(f.Type()); ok {
			getBody += fmt.Sprintf("r.%v = %v\n", f.GoName(), constructor.ConstructorMethod(p))
		}
		if f.Type().WrapperType() == "" {
			pointer := "&"
			if _, ok := f.Type().(*Reference); ok {
				pointer = ""
			}
			getBody += fmt.Sprintf("return %vr.%v\n", pointer, f.GoName())
		} else {
			getBody += fmt.Sprintf("return (*%v)(&r.%v)\n", f.Type().WrapperType(), f.GoName())
		}
	}
	return getBody
}

func (r *RecordDefinition) FieldsMethodDef(p *generator.Package) string {
	getBody := r.getMethodDef(p)
	defaultBody, _ := r.defaultMethodDef(p)
	return fmt.Sprintf(recordFieldTemplate, r.GoType(), getBody, defaultBody)
}

func (r *RecordDefinition) defaultValues(p *generator.Package) (string, error) {
	defaults := ""
	for _, f := range r.fields {
		if f.hasDef {
			def, err := f.Type().DefaultValue(p, fmt.Sprintf("v.%v", f.GoName()), f.Default())
			if err != nil {
				return "", err
			}
			defaults += def + "\n"
		}
	}
	return defaults, nil
}

func (r *RecordDefinition) ConstructorMethodDef(p *generator.Package) (string, error) {
	defaults, err := r.defaultValues(p)
	if err != nil {
		return "", err
	}

	fieldConstructors, err := r.fieldConstructors(p)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(recordConstructorTemplate, r.ConstructorMethod(p), r.GoType(), r.Name(), fieldConstructors, defaults), nil
}

func (r *RecordDefinition) recordReaderTypeName() string {
	return r.Name() + "Reader"
}

func (r *RecordDefinition) recordReaderDef(p *generator.Package) string {
	return fmt.Sprintf(recordReaderTemplate, r.recordReaderTypeName(), r.GoType(), r.ConstructorMethod(p))
}

func (r *RecordDefinition) FieldByName(name string) *Field {
	for _, f := range r.fields {
		if f.Name() == name {
			return f
		}
	}
	return nil
}

func (r *RecordDefinition) DefaultValue(p *generator.Package, lvalue string, rvalue interface{}) (string, error) {
	items := rvalue.(map[string]interface{})
	fieldSetters := ""
	sorted := make([]string, 0, len(items))
	for k, _ := range items {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	for _, k := range sorted {
		v := items[k]
		field := r.FieldByName(k)
		fieldSetter, err := field.Type().DefaultValue(p, fmt.Sprintf("%v.%v", lvalue, field.GoName()), v)
		if err != nil {
			return "", err
		}

		fieldSetters += fieldSetter + "\n"
	}
	return fieldSetters, nil
}

func (r *RecordDefinition) Fields() []*Field {
	return r.fields
}

func (s *RecordDefinition) IsReadableBy(d Definition) bool {
	reader, ok := d.(*RecordDefinition)
	return ok && reader.name == s.name
}

func (s *RecordDefinition) WrapperType() string {
	return ""
}
