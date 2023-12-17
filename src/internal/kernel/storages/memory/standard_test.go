package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name string
}

func TestStandard(t *testing.T) {
	data := &testStruct{Name: "test"}
	clear := &testStruct{Name: "clear"}

	Store(data.Name, data)
	Store(clear.Name, clear)

	exist := Exists(data.Name)
	assert.True(t, exist, "should return true if the key exists")

	recoveredData, exist := Read(data.Name)
	assert.True(t, exist, "should not return an error")
	assert.Equal(t, data, recoveredData, "should return the same data")

	Delete(data.Name)

	recoveredData, exist = Read(data.Name)
	assert.False(t, exist, "should not return an error")
	assert.Nil(t, recoveredData, "Should return nil")
	recoveredData, exist = Read(clear.Name)
	assert.True(t, exist, "should not return an error")
	assert.NotNil(t, recoveredData, "should return nil")

	Clear()

	recoveredData, exist = Read(clear.Name)
	assert.False(t, exist, "should not return an error")
	assert.Nil(t, recoveredData, "should return nil")
}
