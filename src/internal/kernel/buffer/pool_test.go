package buffer_test

import (
	"testing"

	"github.com/kodmain/kitsune/src/internal/kernel/buffer"
)

func TestPool(t *testing.T) {
	t.Run("Exist:Failure(should not exist)", func(t *testing.T) {
		pool := buffer.ExistPool(buffer.SIZE_1B)

		if pool != nil {
			t.Errorf("Expected should not exist, got pool")
		}
	})

	t.Run("Create", func(t *testing.T) {
		pool := buffer.GetPool(buffer.SIZE_1B)

		if pool == nil {
			t.Errorf("Expected pool to be created, got nil")
		}
	})

	t.Run("Exist:Success", func(t *testing.T) {
		pool := buffer.ExistPool(buffer.SIZE_1B)

		if pool == nil {
			t.Errorf("Expected should exist, got nil")
		}
	})

	t.Run("Get", func(t *testing.T) {
		pool := buffer.GetPool(buffer.SIZE_1B)

		buf := pool.Get()

		if buf == nil {
			t.Errorf("Expected buffer to be allocated, got nil")
		}

		if len(*buf) != buffer.SIZE_1B {
			t.Errorf("Expected buffer of size %d, got %d", buffer.SIZE_1B, len(*buf))
		}
	})

	t.Run("Put", func(t *testing.T) {
		pool := buffer.GetPool(buffer.SIZE_1B)
		buf := pool.Get()

		pool.Put(buf)

		newBuf := pool.Get()
		if buf != newBuf {
			t.Errorf("Expected the buffer to be reused")
		}
	})
}
