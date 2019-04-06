package types

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
	return fmt.Errorf("Invalid value for %v")
}
`

const unionDeserializerTemplate = `
func %v(r io.Reader) (%v, error) {
	field, err := readLong(r)
	var unionStr %v
	if err != nil {
		return unionStr, err
	}
	unionStr.UnionType = %v(field)
	switch unionStr.UnionType {
		%v
	default:	
		return unionStr, fmt.Errorf("Invalid value for %v")
	}
	return unionStr, nil
}
`

type unionField struct {
	name       string
	itemType   []AvroType
	definition []interface{}
}

func NewUnionField(name string, itemType []AvroType, definition []interface{}) *unionField {
	return &unionField{
		name:       name,
		itemType:   itemType,
		definition: definition,
	}
}

func (s *unionField) compositeFieldName() string {
	var unionFields = "Union"
	for _, i := range s.itemType {
		unionFields += i.Name()
	}
	return unionFields
}

func (s *unionField) Name() string {
	return s.GoType()
}

func (s *unionField) GoType() string {
	if s.name == "" {
		return generator.ToPublicName(s.compositeFieldName())
	}
	return generator.ToPublicName(s.name)
}

func (s *unionField) unionEnumType() string {
	return fmt.Sprintf("%vType", s.Name())
}

func (s *unionField) unionEnumDef(p *generator.Package) string {
	var unionTypes string
	for i, item := range s.itemType {
		if ref, ok := item.(*Reference); ok && ref.AvroName().Namespace != p.Name() {
			name := imprt.UniqName(p.Root(), ref.AvroName().Namespace, item.Name())
			unionTypes += fmt.Sprintf("%v %v = %v\n", s.unionEnumType()+name, s.unionEnumType(), i)
		} else {
			unionTypes += fmt.Sprintf("%v %v = %v\n", s.unionEnumType()+item.Name(), s.unionEnumType(), i)
		}
	}
	return fmt.Sprintf("type %v int\nconst(\n%v)\n", s.unionEnumType(), unionTypes)
}

func (s *unionField) unionStringerMethodDef(p *generator.Package) string {
	var cases string
	for _, item := range s.itemType {
		name := item.Name()
		if ref, ok := item.(*Reference); ok && ref.AvroName().Namespace != p.Name() {
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

func (s *unionField) unionTypeDef(p *generator.Package) string {
	var unionFields string
	for _, i := range s.itemType {
		if ref, ok := i.(*Reference); ok && ref.AvroName().Namespace != p.Name() {
			unionFields += fmt.Sprintf("%v %v\n", imprt.UniqName(p.Root(), ref.AvroName().Namespace, i.Name()), imprt.Type(p.Root(), ref.AvroName().Namespace, i.GoType()))
		} else {
			unionFields += fmt.Sprintf("%v %v\n", i.Name(), i.GoType())
		}
	}
	unionFields += fmt.Sprintf("UnionType %v", s.unionEnumType())
	return fmt.Sprintf("type %v struct{\n%v\n}\n", s.Name(), unionFields)
}

func (s *unionField) unionSetMethodDef(p *generator.Package, u AvroType) string {
	t := u.GoType()
	n := u.Name()
	if ref, ok := u.(*Reference); ok && ref.AvroName().Namespace != p.Name() {
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

func (s *unionField) unionIdentityMethodDef(p *generator.Package, u AvroType) string {
	n := u.Name()
	if ref, ok := u.(*Reference); ok && ref.AvroName().Namespace != p.Name() {
		n = imprt.UniqName(p.Root(), ref.AvroName().Namespace, n)
	}

	return fmt.Sprintf(`
		func (u *%v) Is%v() bool {
			return u.UnionType == %v
		}
	`, s.Name(), n, s.unionEnumType()+n)
}

func (s *unionField) unionSerializer(p *generator.Package) string {
	switchCase := ""
	for _, t := range s.itemType {
		n := t.Name()
		if ref, ok := t.(*Reference); ok && ref.AvroName().Namespace != p.Name() {
			n = imprt.UniqName(p.Root(), ref.AvroName().Namespace, n)
		}
		switchCase += fmt.Sprintf("case %v:\nreturn %v(r.%v, w)\n", s.unionEnumType()+n, t.SerializerMethod(p), n)
	}
	return fmt.Sprintf(unionSerializerTemplate, s.SerializerMethod(p), s.GoType(), switchCase, s.GoType())
}

func (s *unionField) unionDeserializer(p *generator.Package) string {
	switchCase := ""
	for _, t := range s.itemType {
		n := t.Name()
		if ref, ok := t.(*Reference); ok && ref.AvroName().Namespace != p.Name() {
			n = imprt.UniqName(p.Root(), ref.AvroName().Namespace, n)
		}
		switchCase += fmt.Sprintf("case %v:\nval, err :=  %v(r)\nif err != nil {return unionStr, err}\nunionStr.%v = val\n", s.unionEnumType()+n, t.DeserializerMethod(p), n)
	}
	return fmt.Sprintf(unionDeserializerTemplate, s.DeserializerMethod(p), s.GoType(), s.GoType(), s.unionEnumType(), switchCase, s.GoType())
}

func (s *unionField) filename() string {
	return generator.ToSnake(s.GoType()) + ".go"
}

func (s *unionField) SerializerMethod(p *generator.Package) string {
	return fmt.Sprintf("write%v", s.Name())
}

func (s *unionField) DeserializerMethod(p *generator.Package) string {
	return fmt.Sprintf("read%v", s.Name())
}

func (s *unionField) AddStruct(p *generator.Package, containers bool) error {
	p.AddStruct(s.filename(), s.unionEnumType(), s.unionEnumDef(p))
	p.AddStruct(s.filename(), s.Name(), s.unionTypeDef(p))
	p.AddFunction(s.filename(), s.Name(), "stringer", s.unionStringerMethodDef(p))
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
		if ref, ok := f.(*Reference); ok && ref.AvroName().Namespace != p.Name() {
			p.AddImport(s.filename(), imprt.Path(p.Root(), ref.AvroName().Namespace))
		}
	}

	return nil
}

func (s *unionField) AddSerializer(p *generator.Package) {
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

func (s *unionField) AddDeserializer(p *generator.Package) {
	p.AddImport(UTIL_FILE, "fmt")
	p.AddFunction(UTIL_FILE, "", s.DeserializerMethod(p), s.unionDeserializer(p))
	p.AddFunction(UTIL_FILE, "", "readLong", readLongMethod)
	p.AddImport(UTIL_FILE, "io")
	for _, f := range s.itemType {
		f.AddDeserializer(p)
	}
}

func (s *unionField) ResolveReferences(n *Namespace) error {
	var err error
	for _, f := range s.itemType {
		err = f.ResolveReferences(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *unionField) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	var err error
	for i, item := range s.itemType {
		s.definition[i], err = item.Definition(scope)
		if err != nil {
			return nil, err
		}
	}
	return s.definition, nil
}

func (s *unionField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	lvalue = fmt.Sprintf("%v.%v", lvalue, s.itemType[0].Name())
	return s.itemType[0].DefaultValue(lvalue, rvalue)
}
