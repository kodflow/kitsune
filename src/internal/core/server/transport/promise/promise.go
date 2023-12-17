// Package promise provides functionalities for managing asynchronous responses.
package promise

import (
	"sync"

	"github.com/kodmain/kitsune/src/internal/core/server/transport/proto/generated"
)

// Promise struct represents an asynchronous operation that may complete at some point.
// It manages the lifecycle of an asynchronous request, collecting responses until it is resolved.
type Promise struct {
	Id        string                       // Id is a unique identifier for the promise.
	Length    int                          // Length is the number of responses required to resolve the promise.
	Closed    bool                         // Closed indicates whether the promise has been resolved.
	responses []*generated.Response        // Responses accumulates the responses received.
	callback  func(...*generated.Response) // Callback is a function to be called when the promise is resolved.
	mu        sync.Mutex                   // Mutex ensures thread safety when accessing the promise.
}

// Add adds a request to the promise and increments the number of responses required.
// It associates the promise ID with the request, indicating that the request is part of the promise.
//
// Parameters:
// - req: *generated.Request The request to be added to the promise.
func (p *Promise) Add(req *generated.Request) {
	req.Pid = p.Id // Associate the request with the promise ID.
	p.Length++     // Increment the number of expected responses.
}

// Resolve accumulates responses and resolves the promise if enough responses have been received.
// It adds the response to the promise and checks if the promise has received all expected responses.
// If all responses are received, it triggers the callback function and closes the promise.
//
// Parameters:
// - res: *generated.Response The response to be added to the promise.
func (p *Promise) Resolve(res *generated.Response) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.responses = append(p.responses, res) // Accumulate the response.
	// Check if the promise has received all expected responses.
	if len(p.responses) == p.Length {
		p.Close() // Resolve the promise.
	}
}

// Close marks the promise as closed and executes the callback if all responses are received.
// It also removes the promise from the repository of active promises.
func (p *Promise) Close() {
	p.Closed = true // Mark the promise as resolved.
	// Check if all responses are received and execute the callback.
	if len(p.responses) == p.Length {
		p.callback(p.responses...)
	}

	// Remove the promise from the active repository.
	repository.mu.Lock()
	delete(repository.promises, p.Id)
	repository.mu.Unlock()
}
