package types

type Double float64

func (b *Double) DeserializeBoolean(v bool) {
	panic("Unable to assign boolean to double field")
}

func (b *Double) DeserializeInt(v int32) {
	*(*float64)(b) = float64(v)
}

func (b *Double) DeserializeLong(v int64) {
	*(*float64)(b) = float64(v)
}

func (b *Double) DeserializeFloat(v float32) {
	*(*float64)(b) = float64(v)
}

func (b *Double) DeserializeDouble(v float64) {
	*(*float64)(b) = v
}

func (b *Double) SetUnionElem(v int64) {
	panic("Unable to assign union elem to double field")
}

func (b *Double) DeserializeBytes(v []byte) {
	panic("Unable to assign bytes to double field")
}

func (b *Double) DeserializeString(v string) {
	panic("Unable to assign string to double field")
}

func (b *Double) Get(i int) Field {
	panic("Unable to get field from double field")
}

func (b *Double) SetDefault(i int) {
	panic("Unable to set default on double field")
}

func (b *Double) AppendMap(key string) Field {
	panic("Unable to append map key to from double field")
}

func (b *Double) AppendArray() Field {
	panic("Unable to append array element to from double field")
}

func (b *Double) Finalize() {}
