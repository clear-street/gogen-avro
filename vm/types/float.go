package types

type Float float32

func (b *Float) DeserializeBoolean(v bool) {
	panic("Unable to assign boolean to float field")
}

func (b *Float) DeserializeInt(v int32) {
	*(*float32)(b) = float32(v)
}

func (b *Float) DeserializeLong(v int64) {
	*(*float32)(b) = float32(v)
}

func (b *Float) DeserializeFloat(v float32) {
	*(*float32)(b) = v
}

func (b *Float) SetUnionElem(v int64) {
	panic("Unable to assign union elem to float field")
}

func (b *Float) DeserializeDouble(v float64) {
	panic("Unable to assign double to float field")
}

func (b *Float) DeserializeBytes(v []byte) {
	panic("Unable to assign double to float field")
}

func (b *Float) DeserializeString(v string) {
	panic("Unable to assign double to float field")
}

func (b *Float) Get(i int) Field {
	panic("Unable to get field from float field")
}

func (b *Float) SetDefault(i int) {
	panic("Unable to set default on float field")
}

func (b *Float) AppendMap(key string) Field {
	panic("Unable to append map key to from float field")
}

func (b *Float) AppendArray() Field {
	panic("Unable to append array element to from float field")
}

func (b *Float) Finalize() {}
