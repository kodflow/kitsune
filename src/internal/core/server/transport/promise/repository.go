// Package promise handles functionalities to create, find, and manage asynchronous promises.
package promise

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kodflow/kitsune/src/config"
	"github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger"
)

// Repository manages a set of promises.
// It provides concurrent-safe operations to store and retrieve promises by their ID.
type Repository struct {
	promises map[string]*Promise // Stores promises by their ID.
	mu       sync.Mutex          // Mutex to ensure safe concurrent access to the promises map.
}

// repository is an instance of Repository that stores active promises.
var repository = &Repository{
	promises: map[string]*Promise{},
}

// new creates a new Promise instance.
// It generates a unique ID using the UUID library and initializes a Promise.
//
// Returns:
// - *Promise: A pointer to a newly created Promise with a unique ID.
// - error: An error if UUID generation fails.
func new() (*Promise, error) {
	v4, err := uuid.NewRandom() // Generate a new random UUID.

	if err != nil {
		return nil, err
	}

	p := &Promise{
		Id:        v4.String(), // Set the Promise ID to the UUID string.
		responses: []*generated.Response{},
	}

	return p, nil
}

// Create creates a new promise with a callback function.
// It stores the promise in the repository and starts a timer to remove the promise after a timeout.
//
// Parameters:
// - callback: func(...*generated.Response) Function to call when the promise is resolved.
//
// Returns:
// - *Promise: A pointer to the newly created Promise.
// - error: An error if the Promise creation fails.
func Create(callback func(...*generated.Response)) (*Promise, error) {
	promise, err := new()
	promise.callback = callback // Set the callback function for the promise.
	if err != nil {
		return promise, err
	}

	repository.mu.Lock()
	repository.promises[promise.Id] = promise // Store the promise in the repository.
	repository.mu.Unlock()

	// Start a timer to remove the promise after a predefined timeout duration.
	go func(promise *Promise) {
		<-time.After(time.Second * config.DEFAULT_TIMEOUT)
		if !promise.Closed {
			// Delete the promise after timeout.
			promise.Close()
			logger.Error(fmt.Errorf("promise %s timed out", promise.Id))
		}
	}(promise)

	return promise, nil
}

// Find locates a promise in the repository by its ID.
// It returns the promise if found, or an error if the promise ID is empty or not found.
//
// Parameters:
// - promiseId: string The ID of the promise to locate.
//
// Returns:
// - *Promise: The promise associated with the given ID, if found.
// - error: An error if the promise ID is empty or the promise is not found.
func Find(promiseId string) (*Promise, error) {
	if promiseId == "" {
		return nil, fmt.Errorf("promise no pid")
	}

	repository.mu.Lock()
	defer repository.mu.Unlock()
	if promise, ok := repository.promises[promiseId]; ok {
		return promise, nil
	}

	return nil, fmt.Errorf("promise %s not found", promiseId)
}
