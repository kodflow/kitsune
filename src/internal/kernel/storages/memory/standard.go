package memory

// instance is a global variable that holds a pointer to a Memory instance.
var instance *Memory = nil

// standard returns a pointer to a Memory instance.
// If the instance is nil, it creates a new instance using the New() function.
// It then returns the instance.
func standard() *Memory {
	if instance == nil {
		instance = New()
	}
	return instance
}

// Clear clears the memory storage.
func Clear() {
	standard().Clear()
}

// Store stores a value in the memory storage.
// It takes a key string and a value interface{} as parameters.
// The key is used to identify the stored value.
// The value can be of any type.
func Store(key string, value interface{}) {
	standard().Store(key, value)
}

// Read retrieves a value from the memory storage based on the given key.
// It returns the value and a boolean indicating whether the key exists in the storage.
func Read(key string) (interface{}, bool) {
	return standard().Read(key)
}

// Delete removes a value from the memory storage based on the given key.
func Delete(key string) {
	standard().Delete(key)
}

// Exists checks if a key exists in the memory storage.
// It returns true if the key exists, false otherwise.
func Exists(key string) bool {
	return standard().Exists(key)
}
