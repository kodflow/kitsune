package memory_test

import (
	"testing"

	"github.com/kodmain/kitsune/src/internal/kernel/storages/memory"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name string
}

func TestStandard(t *testing.T) {
	data := &testStruct{Name: "test"}

	memory.Store(data.Name, data)
	exist := memory.Exists(data.Name)
	assert.True(t, exist, "Should return true if the key exists")

	recoveredData, exist := memory.Read(data.Name)
	assert.True(t, exist, "Should not return an error")
	assert.Equal(t, data, recoveredData, "Should return the same data")

	memory.Delete(data.Name)

	recoveredData, exist = memory.Read(data.Name)
	assert.False(t, exist, "Should not return an error")
	assert.Nil(t, recoveredData, "Should return nil")
}
