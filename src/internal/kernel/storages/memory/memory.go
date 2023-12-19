package memory

import "sync"

// Memory represents an in-memory storage implementation.
// This structure provides thread-safe operations for storing and retrieving data in memory.
//
// Attributes:
// - mu: sync.RWMutex for synchronizing read-write access to the data.
// - data: map[string]any for storing the data with string keys.
type Memory struct {
	mu   sync.RWMutex
	data map[string]any
}

// New creates a new instance of the Memory storage.
// It initializes the storage with an empty map.
//
// Returns:
// - *Memory: A pointer to a newly created Memory instance with initialized storage.
func NewMemory() *Memory {
	return &Memory{
		data: make(map[string]any),
	}
}

// Clear removes all data from the memory storage.
// It locks the storage for write operations, clears the data, and then unlocks it.
func (m *Memory) Clear() {
	m.mu.Lock()
	m.data = make(map[string]any)
	m.mu.Unlock()
}

// Store adds or updates a value in the memory storage with the specified key.
// It locks the storage for write operations, stores/updates the value, and then unlocks it.
//
// Parameters:
// - key: string The key to associate with the value.
// - value: any The value to store, which can be of any type.
// Returns:
// - bool: A boolean indicating whether the key is stored in the storage.
func (m *Memory) Store(key string, value any) bool {
	m.mu.Lock()
	m.data[key] = value
	m.mu.Unlock()
	return true
}

// Read retrieves the value associated with the specified key from the memory storage.
// It locks the storage for read operations, retrieves the value, and then unlocks it.
//
// Parameters:
// - key: string The key for which to retrieve the associated value.
//
// Returns:
// - any: The value associated with the key.
// - bool: A boolean indicating whether the key exists in the storage.
func (m *Memory) Read(key string) (any, bool) {
	m.mu.RLock()
	item, exists := m.data[key]
	m.mu.RUnlock()
	return item, exists
}

// Delete removes the value associated with the specified key from the memory storage.
// It locks the storage for write operations, removes the value, and then unlocks it.
//
// Parameters:
// - key: string The key for which to remove the associated value.
// Returns:
// - bool: A boolean indicating whether the key remove from the storage.
func (m *Memory) Delete(key string) bool {
	m.mu.Lock()
	delete(m.data, key)
	m.mu.Unlock()
	return true
}

// Exists checks if the specified key exists in the memory storage.
// It locks the storage for read operations, checks for key existence, and then unlocks it.
//
// Parameters:
// - key: string The key to check for existence.
//
// Returns:
// - bool: A boolean indicating whether the key exists in the storage.
func (m *Memory) Exists(key string) bool {
	m.mu.RLock()
	_, exists := m.data[key]
	m.mu.RUnlock()
	return exists
}
