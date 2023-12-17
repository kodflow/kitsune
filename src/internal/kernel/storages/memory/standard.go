package memory

// instance is a global variable that holds a pointer to a Memory instance.
// It is used to implement a singleton pattern for the Memory type, ensuring that only one instance of Memory is used across the application.
var instance *Memory = nil

// standard returns a pointer to a Memory instance.
// It ensures that only one instance of Memory is created and used throughout the application (singleton pattern).
// If the instance is nil, it creates a new instance using the New() function, then returns the instance.
//
// Returns:
// - *Memory: A pointer to the singleton Memory instance.
func standard() *Memory {
	if instance == nil {
		instance = New()
	}
	return instance
}

// Clear clears the memory storage.
// It invokes the Clear method on the singleton Memory instance, resetting the storage to its initial state.
func Clear() {
	standard().Clear()
}

// Store stores a value in the memory storage.
// It saves the value in the singleton Memory instance, accessible via the provided key.
//
// Parameters:
// - key: string The key to store the value under.
// - value: interface{} The value to store, can be of any type.
func Store(key string, value interface{}) {
	standard().Store(key, value)
}

// Read retrieves a value from the memory storage based on the given key.
// It returns the value associated with the key and a boolean indicating whether the key exists in the storage.
//
// Parameters:
// - key: string The key corresponding to the value to retrieve.
//
// Returns:
// - interface{}: The value associated with the key.
// - bool: A boolean indicating whether the key exists in the storage.
func Read(key string) (interface{}, bool) {
	return standard().Read(key)
}

// Delete removes a value from the memory storage based on the given key.
// It deletes the value associated with the key from the singleton Memory instance.
//
// Parameters:
// - key: string The key corresponding to the value to be removed.
func Delete(key string) {
	standard().Delete(key)
}

// Exists checks if a key exists in the memory storage.
// It checks the presence of a key in the singleton Memory instance.
//
// Parameters:
// - key: string The key to check in the storage.
//
// Returns:
// - bool: True if the key exists in the storage, false otherwise.
func Exists(key string) bool {
	return standard().Exists(key)
}
