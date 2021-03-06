package types

type String string

func (b *String) DeserializeBoolean(v bool) {
	panic("Unable to assign boolean to string field")
}

func (b *String) DeserializeInt(v int32) {
	panic("Unable to assign int to string field")
}

func (b *String) DeserializeLong(v int64) {
	panic("Unable to assign long to string field")
}

func (b *String) DeserializeFloat(v float32) {
	panic("Unable to assign float to string field")
}

func (b *String) SetUnionElem(v int64) {
	panic("Unable to assign union elem to string field")
}

func (b *String) DeserializeDouble(v float64) {
	panic("Unable to assign double to string field")
}

func (b *String) DeserializeBytes(v []byte) {
	*(*string)(b) = string(v)
}

func (b *String) DeserializeString(v string) {
	*(*string)(b) = v
}

func (b *String) Get(i int) Field {
	panic("Unable to get field from string field")
}

func (b *String) SetDefault(i int) {
	panic("Unable to set default on string field")
}

func (b *String) AppendMap(key string) Field {
	panic("Unable to append map key to from string field")
}

func (b *String) AppendArray() Field {
	panic("Unable to append array element to from string field")
}

func (b *String) Finalize() {}
