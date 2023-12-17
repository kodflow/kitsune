package buffer

import "sync"

var (
	// bufferPools holds the instances of buffer pools, keyed by buffer size.
	// It maps buffer sizes to their respective Pool instances.
	bufferPools = make(map[int32]*Pool)

	// mu is used to synchronize access to the bufferPools map.
	mu = sync.Mutex{}
)

// Pool represents a pool of byte slices.
// It provides functionality to get and put byte slices to avoid frequent allocation and deallocation,
// which can improve performance in scenarios where byte slices of a fixed size are repeatedly used.
type Pool struct {
	sync *sync.Pool // sync.Pool internally used for pooling byte slices.
}

// Get retrieves a pointer to a byte slice from the pool.
// It returns a pre-allocated byte slice pointer from the pool, reducing the need for memory allocation.
//
// Output:
// - *[]byte: A pointer to a byte slice from the pool.
func (p *Pool) Get() *[]byte {
	return p.sync.Get().(*[]byte)
}

// Put returns a pointer to a byte slice to the pool.
// This method is used to return a used byte slice pointer back to the pool for reuse,
// aiding in efficient memory management.
//
// Input:
// - x *[]byte: The pointer to the byte slice to be returned to the pool.
func (p *Pool) Put(x *[]byte) {
	p.sync.Put(x)
}

// ExistPool checks for the existence of a Pool for the given buffer size.
// It returns the Pool if it exists in bufferPools, otherwise nil.
//
// Parameters:
// - size: int32 The size of the buffer pool to look for.
//
// Returns:
// - *Pool: The existing buffer pool, or nil if not found.
func ExistPool(size int32) *Pool {
	return bufferPools[size]
}

// GetPool retrieves or creates a Pool for the given buffer size.
// It ensures thread-safe access to the bufferPools.
// If a Pool does not exist for the given size, it is created and added to bufferPools.
//
// Parameters:
// - size: int32 The size of the buffer pool to retrieve or create.
//
// Returns:
// - *Pool: A buffer pool for the given size.
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
