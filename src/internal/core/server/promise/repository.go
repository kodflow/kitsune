package promise

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
)

var repository = &Repository{
	promises: map[string]*Promise{},
}

type Repository struct {
	promises map[string]*Promise
	mu       sync.Mutex
}

func (r *Repository) promise() (*Promise, error) {
	p, err := new()

	if err != nil {
		return p, err
	}

	r.mu.Lock()
	r.promises[p.Id] = p
	r.mu.Unlock()

	return p, nil
}

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

func Create(callback func(...*transport.Response)) (*Promise, error) {
	promise, err := new()
	promise.callback = callback
	if err != nil {
		return promise, err
	}

	repository.mu.Lock()
	repository.promises[promise.Id] = promise
	repository.mu.Unlock()

	return promise, nil
}

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
