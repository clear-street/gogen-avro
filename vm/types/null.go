package types

type NullVal struct{}

func (b *NullVal) DeserializeBoolean(v bool) {
	panic("Unable to assign boolean to null field")
}

func (b *NullVal) DeserializeInt(v int32) {
	panic("Unable to assign boolean to null field")
}

func (b *NullVal) DeserializeLong(v int64) {
	panic("Unable to assign long to null field")
}

func (b *NullVal) DeserializeFloat(v float32) {
	panic("Unable to assign float to null field")
}

func (b *NullVal) SetUnionElem(v int64) {
	panic("Unable to assign union elem to null field")
}

func (b *NullVal) DeserializeDouble(v float64) {
	panic("Unable to assign double to null field")
}

func (b *NullVal) DeserializeBytes(v []byte) {
	panic("Unable to assign bytes to null field")
}

func (b *NullVal) DeserializeString(v string) {
	panic("Unable to assign string to null field")
}

func (b *NullVal) Get(i int) Field {
	panic("Unable to get field from null field")
}

func (b *NullVal) SetDefault(i int) {
	panic("Unable to set default on null field")
}

func (b *NullVal) AppendMap(key string) Field {
	panic("Unable to append map key to from null field")
}

func (b *NullVal) AppendArray() Field {
	panic("Unable to append array element to from null field")
}

func (b *NullVal) Finalize() {}
