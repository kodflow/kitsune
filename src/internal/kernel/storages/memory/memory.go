package memory

import "sync"

// Memory represents an in-memory storage implementation.
// This structure provides thread-safe operations for storing and retrieving data in memory.
//
// Attributes:
// - mu: sync.RWMutex for synchronizing read-write access to the data.
// - data: map[string]interface{} for storing the data with string keys.
type Memory struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

// New creates a new instance of the Memory storage.
// It initializes the storage with an empty map.
//
// Returns:
// - *Memory: A pointer to a newly created Memory instance with initialized storage.
func New() *Memory {
	return &Memory{
		data: make(map[string]interface{}),
	}
}

// Clear removes all data from the memory storage.
// It locks the storage for write operations, clears the data, and then unlocks it.
func (m *Memory) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string]interface{})
}

// Store adds or updates a value in the memory storage with the specified key.
// It locks the storage for write operations, stores/updates the value, and then unlocks it.
//
// Parameters:
// - key: string The key to associate with the value.
// - value: interface{} The value to store, which can be of any type.
func (m *Memory) Store(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value
}

// Read retrieves the value associated with the specified key from the memory storage.
// It locks the storage for read operations, retrieves the value, and then unlocks it.
//
// Parameters:
// - key: string The key for which to retrieve the associated value.
//
// Returns:
// - interface{}: The value associated with the key.
// - bool: A boolean indicating whether the key exists in the storage.
func (m *Memory) Read(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.data[key]
	return item, exists
}

// Delete removes the value associated with the specified key from the memory storage.
// It locks the storage for write operations, removes the value, and then unlocks it.
//
// Parameters:
// - key: string The key for which to remove the associated value.
func (m *Memory) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
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
	defer m.mu.RUnlock()

	_, exists := m.data[key]
	return exists
}
