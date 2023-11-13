// Package promise provides functionalities for managing asynchronous responses.
package promise

import (
	"sync"

	"github.com/kodmain/kitsune/src/internal/core/server/transport"
)

// Promise struct represents an asynchronous operation that may complete at some point.
type Promise struct {
	Id        string // Unique identifier for the promise
	Length    int    // Number of responses required to resolve the promise
	Closed    bool
	responses []*transport.Response        // Accumulates the responses received
	callback  func(...*transport.Response) // Function to be called when the promise is resolved
	mu        sync.Mutex                   // Mutex to ensure thread safety
}

// Add adds a request to the promise and increments the number of responses required.
// req: The request to be added to the promise
func (p *Promise) Add(req *transport.Request) {
	req.Pid = p.Id
	p.Length++
}

// Resolve accumulates responses and resolves the promise if enough responses have been received.
// res: The response to be added to the promise
func (p *Promise) Resolve(res *transport.Response) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.responses = append(p.responses, res)
	if len(p.responses) == p.Length {
		p.Close()
	}
}

func (p *Promise) Close() {
	p.Closed = true
	if len(p.responses) == p.Length {
		p.callback(p.responses...)
	}

	repository.mu.Lock()
	delete(repository.promises, p.Id)
	repository.mu.Unlock()
}
