package promise

import (
	"sync"

	"github.com/kodmain/kitsune/src/internal/core/server/transport"
)

// Promise is used to manage the async response mechanism.
type Promise struct {
	Id                string
	responsesRequired int
	responses         []*transport.Response
	callback          func(...*transport.Response)
	mu                sync.Mutex
}

func (p *Promise) Add(req *transport.Request) {
	req.Pid = p.Id
	p.responsesRequired++
}

func (p *Promise) Resolve(res *transport.Response) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.responses = append(p.responses, res)
	if len(p.responses) == p.responsesRequired {
		p.callback(p.responses...)
		repository.mu.Lock()
		delete(repository.promises, p.Id)
		repository.mu.Unlock()
	}
}
