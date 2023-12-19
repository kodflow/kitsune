package memory

import (
	"hash/fnv"
)

const shardCount = 1024

// Shared represents a shared memory storage.
type Shared struct {
	shards [shardCount]*Memory
}

// NewSharedMemory creates a new instance of Shared.
// It initializes the shards with empty memory.
//
// Returns:
// - s: *Shared - The new instance of Shared.
func NewSharedMemory() *Shared {
	s := &Shared{}
	for i := 0; i < shardCount; i++ {
		s.shards[i] = &Memory{data: make(map[string]interface{})}
	}
	return s
}

// getShard returns the shard associated with the given key.
//
// Parameters:
// - key: string - The key used to determine the shard.
//
// Returns:
// - shard: *Memory - The shard associated with the key.
func (s *Shared) getShard(key string) *Memory {
	hasher := fnv.New32()
	hasher.Write([]byte(key))
	return s.shards[hasher.Sum32()%shardCount]
}

// Store stores a value in the shared memory storage.
//
// Parameters:
// - key: string - The key used to store the value.
// - value: interface{} - The value to be stored.
//
// Returns:
// - success: bool - True if the value was successfully stored, false otherwise.
func (s *Shared) Store(key string, value interface{}) bool {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	shard.data[key] = value
	return true
}

// Read retrieves a value from the shared memory storage.
//
// Parameters:
// - key: string - The key used to retrieve the value.
//
// Returns:
// - value: interface{} - The retrieved value.
// - exists: bool - True if the value exists, false otherwise.
func (s *Shared) Read(key string) (interface{}, bool) {
	shard := s.getShard(key)
	shard.mu.RLock()
	value, exists := shard.data[key]
	shard.mu.RUnlock()
	return value, exists
}

// Delete removes a value from the shared memory storage.
//
// Parameters:
// - key: string - The key used to delete the value.
// Delete deletes the value associated with the given key from the shared memory storage.
func (s *Shared) Delete(key string) {
	shard := s.getShard(key)
	shard.mu.Lock()
	delete(shard.data, key)
	shard.mu.Unlock()
}

// Exists checks if a key exists in the shared memory storage.
// It returns true if the key exists, otherwise it returns false.
//
// Parameters:
// - key: string - The key to check for existence in the shared memory storage.
//
// Returns:
// - exists: bool - True if the key exists, false otherwise.
func (s *Shared) Exists(key string) bool {
	shard := s.getShard(key)
	shard.mu.RLock()
	_, exists := shard.data[key]
	shard.mu.RUnlock()
	return exists
}
