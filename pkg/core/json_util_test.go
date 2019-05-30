package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testJsonStruct struct {
	Foo string `json:"foo"`
}

func Test_jsonbSqlUnmarshal(t *testing.T) {
	jsonb := []byte(`{"foo": "val"}`)
	obj := &testJsonStruct{}
	err := jsonbSqlUnmarshal(jsonb, obj)

	assert.NoError(t, err)
	assert.Equal(t, "val", obj.Foo)
}

func Test_jsonbSqlUnmarshal_NotByteThrows(t *testing.T) {
	err := jsonbSqlUnmarshal("not[]byte", nil)

	assert.Equal(t, "unable to read jsonb field as []byte", err.Error())
}

func Test_jsonbSqlUnmarshal_IgnoresNil(t *testing.T) {
	obj := &testJsonStruct{}

	err := jsonbSqlUnmarshal(nil, obj)

	assert.NoError(t, err)
	assert.Empty(t, obj)
}
