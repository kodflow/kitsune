package memory

// New creates a new instance of the Memory storage.
func New() *Memory {
	return &Memory{
		data: make(map[string]interface{}),
	}
}

// Clear removes all data from the memory storage.
func (m *Memory) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string]interface{})
}

// Store adds or updates a value in the memory storage with the specified key.
func (m *Memory) Store(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value
}

// Read retrieves the value associated with the specified key from the memory storage.
// It returns the value and a boolean indicating whether the key exists in the storage.
func (m *Memory) Read(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.data[key]
	return item, exists
}

// Delete removes the value associated with the specified key from the memory storage.
func (m *Memory) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
}

// Exists checks if the specified key exists in the memory storage.
// It returns a boolean indicating whether the key exists.
func (m *Memory) Exists(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.data[key]
	return exists
}
