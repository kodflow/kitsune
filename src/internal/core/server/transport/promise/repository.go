// Package promise handles functionalities to create, find, and manage asynchronous promises.
package promise

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kodmain/kitsune/src/config"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger"
)

// Repository manages a set of promises.
type Repository struct {
	promises map[string]*Promise // Stores promises by their ID
	mu       sync.Mutex          // Mutex to ensure safe concurrent access
}

// Repository instance that stores active promises.
var repository = &Repository{
	promises: map[string]*Promise{},
}

// new creates a new Promise instance.
func new() (*Promise, error) {
	v4, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	p := &Promise{
		Id:        v4.String(),
		responses: []*transport.Response{},
	}

	return p, nil
}

// Create creates a new promise with a callback function.
// callback: Function to call when the promise is resolved.
func Create(callback func(...*transport.Response)) (*Promise, error) {
	promise, err := new()
	promise.callback = callback
	if err != nil {
		return promise, err
	}

	repository.mu.Lock()
	repository.promises[promise.Id] = promise
	repository.mu.Unlock()

	// Start a timer to remove the promise if it takes too long
	go func(promise *Promise) {
		<-time.After(time.Second * config.DEFAULT_TIMEOUT)
		if !promise.Closed {
			// Delete the promise after timeout
			promise.Close()
			logger.Error(fmt.Errorf("promise %s timed out", promise.Id))
		}
	}(promise)

	return promise, nil
}

// Find locates a promise in the repository by its ID.
// promiseId: The ID of the promise to locate.
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
