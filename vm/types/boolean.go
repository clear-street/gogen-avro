package types

type Boolean bool

func (b *Boolean) DeserializeBoolean(v bool) {
	*(*bool)(b) = v
}

func (b *Boolean) DeserializeInt(v int32) {
	panic("Unable to assign int to boolean field")
}

func (b *Boolean) DeserializeLong(v int64) {
	panic("Unable to assign long to boolean field")
}

func (b *Boolean) DeserializeFloat(v float32) {
	panic("Unable to assign float to boolean field")
}

func (b *Boolean) DeserializeDouble(v float64) {
	panic("Unable to assign double to boolean field")
}

func (b *Boolean) DeserializeBytes(v []byte) {
	panic("Unable to assign bytes to boolean field")
}

func (b *Boolean) DeserializeString(v string) {
	panic("Unable to assign string to boolean field")
}

func (b *Boolean) SetUnionElem(v int64) {
	panic("Unable to assign union elem to boolean field")
}

func (b *Boolean) Get(i int) Field {
	panic("Unable to get field from boolean field")
}

func (b *Boolean) SetDefault(i int) {
	panic("Unable to set default on boolean field")
}

func (b *Boolean) AppendMap(key string) Field {
	panic("Unable to append map key to from boolean field")
}

func (b *Boolean) AppendArray() Field {
	panic("Unable to append array element to from boolean field")
}

func (b *Boolean) Finalize() {}
