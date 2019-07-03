// Wrappers for Avro primitive types implementing the methods required by GADGT
package types

// The interface neeed by GADGT to enter and set fields on a type
// Most types only need to implement a subset
type Field interface {
	// Assign a primitive field
	DeserializeBoolean(v bool)
	DeserializeInt(v int32)
	DeserializeLong(v int64)
	DeserializeFloat(v float32)
	DeserializeDouble(v float64)
	DeserializeBytes(v []byte)
	DeserializeString(v string)

	// Get a nested field
	Get(i int) Field
	// Set the default value for a given field
	SetDefault(i int)

	// Append a new value to a map or array and enter it
	AppendMap(key string) Field
	AppendArray() Field

	// Finalize a field if necessary
	Finalize()
}
