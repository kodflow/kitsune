package memory

import (
	"sync"
	"time"
)

// Item represents an item stored in memory storage.
type Item struct {
	Value     interface{} // Value is the value of the item.
	CreatedAt time.Time   // CreatedAt is the timestamp when the item was created.
}

// Memory represents an in-memory storage implementation.
type Memory struct {
	mu   sync.RWMutex
	data map[string]interface{}
}
