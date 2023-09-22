package buffer

import "sync"

var (
	// bufferPools holds the instances of buffer pools, keyed by buffer size
	bufferPools = make(map[int32]*Pool)
	// mutex to protect bufferPools map
	mu = sync.Mutex{}
)

type Pool struct {
	sync *sync.Pool
}

// Get retrieves a pointer to a byte slice from the pool.
// Input: None.
// Output: *[]byte - A pointer to a byte slice from the pool.
// Objective: To retrieve a pre-allocated byte slice pointer from the pool.
func (p *Pool) Get() *[]byte {
	return p.sync.Get().(*[]byte)
}

// Put returns a pointer to a byte slice to the pool.
// Input: x *[]byte - The pointer to a byte slice to be returned to the pool.
// Output: None.
// Objective: To return a used byte slice pointer back to the pool for reuse.
func (p *Pool) Put(x *[]byte) {
	p.sync.Put(x)
}

func ExistPool(size int32) *Pool {
	return bufferPools[size]
}

func GetPool(size int32) *Pool {
	mu.Lock()
	defer mu.Unlock()

	if pool := ExistPool(size); pool != nil {
		return pool
	}

	p := &Pool{
		sync: &sync.Pool{
			New: func() interface{} {
				b := make([]byte, size)
				return &b
			},
		},
	}

	bufferPools[size] = p

	return p
}
