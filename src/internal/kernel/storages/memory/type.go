package memory

import (
	"sync"
	"time"
)

type Item struct {
	Value     interface{}
	CreatedAt time.Time
}

type Memory struct {
	mu   sync.RWMutex
	data map[string]interface{}
}
