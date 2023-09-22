package socket

import (
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
)

// Promise is used to manage the async response mechanism.
type promise struct {
	wait chan *transport.Response // wait is a channel used to await the response.
}

// Wait waits for a response and returns it.
func (p *promise) Wait() *transport.Response {
	if p.wait == nil {
		return nil
	}

	res := <-p.wait

	return res
}

func (p *promise) Init() {
	p.wait = make(chan *transport.Response)
}

func Promise() *promise {
	return &promise{}
}

/*
package socket

import (
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
)

// Promise is used to manage the async response mechanism.
type promise struct {
	wait   chan *transport.Response // wait is a channel used to await the response.
	answer *transport.Response      // answer stores the received response.
}

// Wait waits for a response and returns it.
func (p *promise) Wait() *transport.Response {
	// If answer is already set, return it.
	if p.answer != nil {
		return p.answer
	}

	// If it's not set, the goroutine from Init will set it when a response is received.
	// We can simply wait for the channel to close to ensure the response has been set.
	for range p.wait {
	}

	return p.answer
}

func (p *promise) Init() {
	p.wait = make(chan *transport.Response, 1)

	go func() {
		p.answer = <-p.wait
		close(p.wait)
	}()
}

func Promise() *promise {
	return &promise{}
}
*/
