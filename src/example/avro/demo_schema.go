// Code generated by github.com/clear-street/gogen-avro. DO NOT EDIT.
/*
 * SOURCE:
 *     example.avsc
 */

package avro

import (
	"io"

	"github.com/clear-street/gogen-avro/compiler"
	"github.com/clear-street/gogen-avro/container"
	"github.com/clear-street/gogen-avro/vm"
	"github.com/clear-street/gogen-avro/vm/types"
)

type DemoSchema struct {
	IntField    int32
	DoubleField float64
	StringField string
	BoolField   bool
	BytesField  []byte
}

func NewDemoSchemaWriter(writer io.Writer, codec container.Codec, recordsPerBlock int64) (*container.Writer, error) {
	str := &DemoSchema{}
	return container.NewWriter(writer, codec, recordsPerBlock, str.Schema())
}

func DeserializeDemoSchema(r io.Reader) (*DemoSchema, error) {
	t := NewDemoSchema()

	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	return t, err
}

func NewDemoSchema() *DemoSchema {
	return &DemoSchema{}
}

func (r *DemoSchema) Schema() string {
	return "{\"fields\":[{\"name\":\"IntField\",\"type\":\"int\"},{\"name\":\"DoubleField\",\"type\":\"double\"},{\"name\":\"StringField\",\"type\":\"string\"},{\"name\":\"BoolField\",\"type\":\"boolean\"},{\"name\":\"BytesField\",\"type\":\"bytes\"}],\"name\":\"DemoSchema\",\"type\":\"record\"}"
}

func (r *DemoSchema) Serialize(w io.Writer) error {
	return writeDemoSchema(r, w)
}

func (_ *DemoSchema) DeserializeBoolean(v bool)   { panic("Unsupported operation") }
func (_ *DemoSchema) DeserializeInt(v int32)      { panic("Unsupported operation") }
func (_ *DemoSchema) DeserializeLong(v int64)     { panic("Unsupported operation") }
func (_ *DemoSchema) DeserializeFloat(v float32)  { panic("Unsupported operation") }
func (_ *DemoSchema) DeserializeDouble(v float64) { panic("Unsupported operation") }
func (_ *DemoSchema) DeserializeBytes(v []byte)   { panic("Unsupported operation") }
func (_ *DemoSchema) DeserializeString(v string)  { panic("Unsupported operation") }
func (_ *DemoSchema) SetUnionElem(v int64)        { panic("Unsupported operation") }
func (_ *DemoSchema) SetDefault(i int)            { panic("Unsupported operation") }
func (r *DemoSchema) Get(i int) types.Field {
	switch i {
	case 0:
		return (*types.Int)(&r.IntField)
		break
	case 1:
		return (*types.Double)(&r.DoubleField)
		break
	case 2:
		return (*types.String)(&r.StringField)
		break
	case 3:
		return (*types.Boolean)(&r.BoolField)
		break
	case 4:
		return (*types.Bytes)(&r.BytesField)
		break

	}
	panic("Unknown field index")
}
func (_ *DemoSchema) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *DemoSchema) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *DemoSchema) Finalize()                        {}

type DemoSchemaReader struct {
	r io.Reader
	p *vm.Program
}

func NewDemoSchemaReader(r io.Reader) (*DemoSchemaReader, error) {
	containerReader, err := container.NewReader(r)
	if err != nil {
		return nil, err
	}

	t := NewDemoSchema()
	deser, err := compiler.CompileSchemaBytes([]byte(containerReader.AvroContainerSchema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	return &DemoSchemaReader{
		r: containerReader,
		p: deser,
	}, nil
}

func (r *DemoSchemaReader) Read() (*DemoSchema, error) {
	t := NewDemoSchema()
	err := vm.Eval(r.r, r.p, t)
	return t, err
}