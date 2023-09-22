package buffer_test

import (
	"testing"

	"github.com/kodmain/kitsune/src/internal/kernel/buffer"
)

func TestSize(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  int
	}{
		{"Size:1", buffer.SIZE_1B, 1},
		{"Size:2", buffer.SIZE_2B, 2},
		{"Size:4", buffer.SIZE_4B, 4},
		{"Size:8", buffer.SIZE_8B, 8},
		{"Size:16", buffer.SIZE_16B, 16},
		{"Size:32", buffer.SIZE_32B, 32},
		{"Size:64", buffer.SIZE_64B, 64},
		{"Size:128", buffer.SIZE_128B, 128},
		{"Size:256", buffer.SIZE_256B, 256},
		{"Size:512", buffer.SIZE_512B, 512},
		{"Size:1024", buffer.SIZE_1KB, 1024},
		{"Size:2048", buffer.SIZE_2KB, 2048},
		{"Size:4096", buffer.SIZE_4KB, 4096},
		{"Size:8192", buffer.SIZE_8KB, 8192},
		{"Size:16384", buffer.SIZE_16KB, 16384},
		{"Size:32768", buffer.SIZE_32KB, 32768},
		{"Size:65536", buffer.SIZE_64KB, 65536},
		{"Size:131072", buffer.SIZE_128KB, 131072},
		{"Size:262144", buffer.SIZE_256KB, 262144},
		{"Size:524288", buffer.SIZE_512KB, 524288},
		{"Size:1048576", buffer.SIZE_1MB, 1048576},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.want {
				t.Errorf("For %s, expected %d, got %d", tt.name, tt.want, tt.value)
			}
		})
	}
}
