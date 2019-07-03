package schema

import (
	"fmt"

	"github.com/clear-street/gogen-avro/generator"
	"github.com/clear-street/gogen-avro/imprt"
)

const unionSerializerTemplate = `
func %v(r %v, w io.Writer) error {
	err := writeLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType{
		%v
	}
	return fmt.Errorf("invalid value for %v")
}
`

const unionConstructorTemplate = `
func %v %v {
	return %v{}
}
`

const unionFieldTemplate = `
func (_ %[1]v) DeserializeBoolean(v bool) { panic("Unsupported operation") }
func (_ %[1]v) DeserializeInt(v int32) { panic("Unsupported operation") }
func (_ %[1]v) DeserializeFloat(v float32) { panic("Unsupported operation") }
func (_ %[1]v) DeserializeDouble(v float64) { panic("Unsupported operation") }
func (_ %[1]v) DeserializeBytes(v []byte) { panic("Unsupported operation") }
func (_ %[1]v) DeserializeString(v string) { panic("Unsupported operation") }
func (r %[1]v) DeserializeLong(v int64) { 
	r.UnionType = (%[2]v)(v)
}
func (r %[1]v) Get(i int) types.Field {
	switch (i) {
		%[3]v
	}
	panic("Unknown field index")
}
func (_ %[1]v) SetDefault(i int) { panic("Unsupported operation") }
func (_ %[1]v) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ %[1]v) AppendArray() types.Field { panic("Unsupported operation") }
func (_ %[1]v) Finalize()  { }
`

type UnionField struct {
	name       string
	itemType   []AvroType
	definition []interface{}
}

func NewUnionField(name string, itemType []AvroType, definition []interface{}) *UnionField {
	return &UnionField{
		name:       name,
		itemType:   itemType,
		definition: definition,
	}
}

func (s *UnionField) compositeFieldName() string {
	var UnionFields = "Union"
	for _, i := range s.itemType {
		UnionFields += i.Name()
	}
	return UnionFields
}

func (s *UnionField) Name() string {
	if s.name == "" {
		return generator.ToPublicName(s.compositeFieldName())
	}
	return generator.ToPublicName(s.name)
}

func (s *UnionField) AvroTypes() []AvroType {
	return s.itemType
}

func (s *UnionField) GoType() string {
	return s.Name()
}

func (s *UnionField) unionEnumType() string {
	return fmt.Sprintf("%vType", s.Name())
}

func (s *UnionField) unionEnumDef(p *generator.Package) string {
	var unionTypes string
	for i, item := range s.itemType {
		if ref, ok := item.(*Reference); ok && !Contains(p, ref) {
			name := imprt.UniqName(p.Root(), ref.AvroName().Namespace, item.Name())
			unionTypes += fmt.Sprintf("%v %v = %v\n", s.unionEnumType()+name, s.unionEnumType(), i)
		} else {
			unionTypes += fmt.Sprintf("%v %v = %v\n", s.unionEnumType()+item.Name(), s.unionEnumType(), i)
		}
	}
	return fmt.Sprintf("type %v int\nconst(\n%v)\n", s.unionEnumType(), unionTypes)
}

func (s *UnionField) unionStringerMethodDef(p *generator.Package) string {
	var cases string
	for _, item := range s.itemType {
		name := item.Name()
		if ref, ok := item.(*Reference); ok && !Contains(p, ref) {
			name = imprt.UniqName(p.Root(), ref.AvroName().Namespace, name)
		}

		cases += fmt.Sprintf("case %v:\nreturn %q\n", s.unionEnumType()+name, name)
	}

	return fmt.Sprintf(`
		func (u *%v) Stringify() string {
			switch u.UnionType {
				%v
			default:
				return "unknown"
			}
		}
	`, s.Name(), cases)
}

func (s *UnionField) unionTypeDef(p *generator.Package) string {
	var unionFields string
	for _, i := range s.itemType {
		if ref, ok := i.(*Reference); ok && !Contains(p, ref) {
			unionFields += fmt.Sprintf("%v %v\n", imprt.UniqName(p.Root(), ref.AvroName().Namespace, i.Name()), imprt.Type(p.Root(), ref.AvroName().Namespace, i.GoType()))
		} else {
			unionFields += fmt.Sprintf("%v %v\n", i.Name(), i.GoType())
		}
	}
	unionFields += fmt.Sprintf("UnionType %v", s.unionEnumType())
	return fmt.Sprintf("type %v struct{\n%v\n}\n", s.Name(), unionFields)
}

func (s *UnionField) unionSetMethodDef(p *generator.Package, u AvroType) string {
	t := u.GoType()
	n := u.Name()
	if ref, ok := u.(*Reference); ok && !Contains(p, ref) {
		t = imprt.Type(p.Root(), ref.AvroName().Namespace, t)
		n = imprt.UniqName(p.Root(), ref.AvroName().Namespace, n)
	}

	return fmt.Sprintf(`
		func (u *%v) Set%v(val %v) {
			u.%v = val
			u.UnionType = %v
		}
	`, s.Name(), n, t, n, s.unionEnumType()+n)
}

func (s *UnionField) unionIdentityMethodDef(p *generator.Package, u AvroType) string {
	n := u.Name()
	if ref, ok := u.(*Reference); ok && !Contains(p, ref) {
		n = imprt.UniqName(p.Root(), ref.AvroName().Namespace, n)
	}

	return fmt.Sprintf(`
		func (u *%v) Is%v() bool {
			return u.UnionType == %v
		}
	`, s.Name(), n, s.unionEnumType()+n)
}

func (s *UnionField) unionSerializer(p *generator.Package) string {
	switchCase := ""
	for _, t := range s.itemType {
		n := t.Name()
		if ref, ok := t.(*Reference); ok && !Contains(p, ref) {
			n = imprt.UniqName(p.Root(), ref.AvroName().Namespace, n)
		}
		switchCase += fmt.Sprintf("case %v:\nreturn %v(r.%v, w)\n", s.unionEnumType()+n, t.SerializerMethod(p), n)
	}
	return fmt.Sprintf(unionSerializerTemplate, s.SerializerMethod(p), s.GoType(), switchCase, s.GoType())
}

func (s *UnionField) FieldsMethodDef(p *generator.Package) string {
	getBody := ""
	for i, f := range s.itemType {
		name := f.Name()
		if ref, ok := f.(*Reference); ok && !Contains(p, ref) {
			name = imprt.UniqName(p.Root(), ref.AvroName().Namespace, name)
		}
		getBody += fmt.Sprintf("case %v:\n", i)
		if constructor, ok := getConstructableForType(f); ok {
			getBody += fmt.Sprintf("r.%v = %v\n", name, constructor.ConstructorMethod(p))
		}
		if f.WrapperType() == "" {
			getBody += fmt.Sprintf("return r.%v", name)
		} else {
			getBody += fmt.Sprintf("return (*%v)(&r.%v)", f.WrapperType(), name)
		}
		getBody += "\nbreak\n"
	}
	return fmt.Sprintf(unionFieldTemplate, s.GoType(), s.unionEnumType(), getBody)
}

func (s *UnionField) filename() string {
	return generator.ToSnake(s.Name()) + ".go"
}

func (s *UnionField) SerializerMethod(p *generator.Package) string {
	return fmt.Sprintf("write%v", s.Name())
}

func (s *UnionField) AddStruct(p *generator.Package, containers bool) error {
	p.AddStruct(s.filename(), s.unionEnumType(), s.unionEnumDef(p))
	p.AddStruct(s.filename(), s.Name(), s.unionTypeDef(p))
	p.AddFunction(s.filename(), s.Name(), "stringer", s.unionStringerMethodDef(p))
	p.AddFunction(s.filename(), s.GoType(), s.ConstructorMethod(), s.constructorMethodDef())
	for _, f := range s.itemType {
		err := f.AddStruct(p, containers)
		if err != nil {
			return err
		}
	}
	for _, f := range s.itemType {
		set := s.unionSetMethodDef(p, f)
		p.AddFunction(s.filename(), "set", set, set)

		ident := s.unionIdentityMethodDef(p, f)
		p.AddFunction(s.filename(), "identity", ident, ident)
	}
	for _, f := range s.itemType {
		if ref, ok := f.(*Reference); ok && !Contains(p, ref) {
			p.AddImport(s.filename(), imprt.Path(p.Root(), ref.AvroName().Namespace))
		}
	}
	p.AddImport(s.filename(), "github.com/clear-street/gogen-avro/vm/types")
	p.AddFunction(s.filename(), s.GoType(), "fieldTemplate", s.FieldsMethodDef(p))

	return nil
}

func (s *UnionField) AddSerializer(p *generator.Package) {
	p.AddImport(UTIL_FILE, "fmt")
	p.AddFunction(UTIL_FILE, "", s.SerializerMethod(p), s.unionSerializer(p))
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeLong", writeLongMethod)
	p.AddFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.AddImport(UTIL_FILE, "io")
	for _, f := range s.itemType {
		f.AddSerializer(p)
	}
}

func (s *UnionField) ResolveReferences(n *Namespace) error {
	var err error
	for _, f := range s.itemType {
		err = f.ResolveReferences(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *UnionField) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	var err error
	for i, item := range s.itemType {
		s.definition[i], err = item.Definition(scope)
		if err != nil {
			return nil, err
		}
	}
	return s.definition, nil
}

func (s *UnionField) DefaultValue(p *generator.Package, lvalue string, rvalue interface{}) (string, error) {
	defaultType := s.itemType[0]
	init := fmt.Sprintf("%v = %v\n", lvalue, s.ConstructorMethod())
	lvalue = fmt.Sprintf("%v.%v", lvalue, defaultType.Name())
	constructorCall := ""
	if constructor, ok := getConstructableForType(defaultType); ok {
		constructorCall = fmt.Sprintf("%v = %v\n", lvalue, constructor.ConstructorMethod(p))
	}
	assignment, err := defaultType.DefaultValue(p, lvalue, rvalue)
	return init + constructorCall + assignment, err
}

func (s *UnionField) WrapperType() string {
	return ""
}

func (s *UnionField) IsReadableBy(f AvroType) bool {
	// Report if *any* writer type could be deserialized by the reader
	for _, t := range s.AvroTypes() {
		if readerUnion, ok := f.(*UnionField); ok {
			for _, rt := range readerUnion.AvroTypes() {
				if t.IsReadableBy(rt) {
					return true
				}
			}
		} else {
			if t.IsReadableBy(f) {
				return true
			}
		}
	}
	return false
}

func (s *UnionField) ConstructorMethod() string {
	return fmt.Sprintf("New%v()", s.Name())
}

func (s *UnionField) constructorMethodDef() string {
	return fmt.Sprintf(unionConstructorTemplate, s.ConstructorMethod(), s.GoType(), s.Name())
}

func (s *UnionField) Equals(reader *UnionField) bool {
	if len(reader.AvroTypes()) != len(s.AvroTypes()) {
		return false
	}

	for i, t := range s.AvroTypes() {
		readerType := reader.AvroTypes()[i]
		if writerRef, ok := t.(*Reference); ok {
			if readerRef, ok := readerType.(*Reference); ok {
				if readerRef.TypeName != writerRef.TypeName {
					return false
				}
			} else {
				return false
			}
		} else if t != readerType {
			return false
		}
	}
	return true
}

func (s *UnionField) SimpleName() string {
	return s.GoType()
}
