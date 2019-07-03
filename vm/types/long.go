package types

type Long int64

func (b *Long) DeserializeBoolean(v bool) {
	panic("Unable to assign boolean to long field")
}

func (b *Long) DeserializeInt(v int32) {
	*(*int64)(b) = int64(v)
}

func (b *Long) DeserializeLong(v int64) {
	*(*int64)(b) = v
}

func (b *Long) DeserializeFloat(v float32) {
	panic("Unable to assign float to long field")
}

func (b *Long) SetUnionElem(v int64) {
	panic("Unable to assign union elem to long field")
}

func (b *Long) DeserializeDouble(v float64) {
	panic("Unable to assign double to long field")
}

func (b *Long) DeserializeBytes(v []byte) {
	panic("Unable to assign bytes to long field")
}

func (b *Long) DeserializeString(v string) {
	panic("Unable to assign string to long field")
}

func (b *Long) Get(i int) Field {
	panic("Unable to get field from long field")
}

func (b *Long) SetDefault(i int) {
	panic("Unable to set default on long field")
}

func (b *Long) AppendMap(key string) Field {
	panic("Unable to append map key to from long field")
}

func (b *Long) AppendArray() Field {
	panic("Unable to append array element to from long field")
}

func (b *Long) Finalize() {}
