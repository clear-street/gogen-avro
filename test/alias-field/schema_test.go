package avro

import (
	"bytes"
	"testing"

	"github.com/clear-street/gogen-avro/compiler"
	evolution "github.com/clear-street/gogen-avro/test/alias-field/evolution"
	"github.com/clear-street/gogen-avro/vm"

	"github.com/stretchr/testify/assert"
)

func TestEvolution(t *testing.T) {
	oldAliasRecord := NewAliasRecord()
	oldAliasRecord.A = "hi"
	oldAliasRecord.C = "bye"

	var buf bytes.Buffer
	err := oldAliasRecord.Serialize(&buf)
	assert.Nil(t, err)

	newAliasRecord := evolution.NewAliasRecord()

	deser, err := compiler.CompileSchemaBytes([]byte(oldAliasRecord.Schema()), []byte(newAliasRecord.Schema()))
	assert.Nil(t, err)

	err = vm.Eval(bytes.NewReader(buf.Bytes()), deser, newAliasRecord)
	assert.Nil(t, err)

	assert.Equal(t, "hi", newAliasRecord.B)
	assert.Equal(t, "bye", newAliasRecord.D)
}
