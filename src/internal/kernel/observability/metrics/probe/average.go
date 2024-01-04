package probe

import (
	"sync"
	"sync/atomic"
	"time"
)

type Average struct {
	counters map[string]*circularBuffer
	ticker   *time.Ticker
	mu       sync.Mutex
}

// circularBuffer représente un buffer circulaire pour les compteurs
type circularBuffer struct {
	buffer []uint64
	index  int
}

// newCircularBuffer crée un nouveau buffer circulaire avec une taille donnée
func newCircularBuffer(size int) *circularBuffer {
	return &circularBuffer{
		buffer: make([]uint64, size),
		index:  0,
	}
}

// next avance l'index du buffer circulaire
func (cb *circularBuffer) next() {
	cb.index = (cb.index + 1) % len(cb.buffer)
	cb.buffer[cb.index] = 0 // Réinitialiser le compteur à l'index actuel
}

// increment augmente le compteur à l'index actuel
func (cb *circularBuffer) increment() {
	atomic.AddUint64(&cb.buffer[cb.index], 1)
}

func NewAverage(duration time.Duration, historySize int) *Average {
	a := &Average{
		counters: make(map[string]*circularBuffer),
		ticker:   time.NewTicker(duration),
	}

	go func() {
		for range a.ticker.C {
			a.advanceCounters()
		}
	}()

	return a
}

func (a *Average) advanceCounters() {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, cb := range a.counters {
		cb.next()
	}
}

func (a *Average) Hit(zone string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.counters[zone]; !ok {
		a.counters[zone] = newCircularBuffer(0)
	}

	a.counters[zone].increment()
}
