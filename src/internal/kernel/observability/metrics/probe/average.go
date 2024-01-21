package probe

import (
	"time"
)

type Average struct {
	ticker     *time.Ticker
	counter    *Counter
	data       []uint64
	sliceTimes []int
	capacity   uint64

	current uint64
	size    uint64
}

func NewAverage(duration time.Duration, maxtimes ...time.Duration) *Average {
	var times uint64
	var defaultValue = time.Minute
	var highestDuration = time.Duration(0)

	if len(maxtimes) == 0 {
		maxtimes = append(maxtimes, defaultValue)
	}

	var sliceTimes = make([]int, 0, len(maxtimes))
	for _, maxtime := range maxtimes {
		sliceTimes = append(sliceTimes, int(maxtime/duration))
		if highestDuration < maxtime {
			highestDuration = maxtime
		}
	}

	times = uint64(highestDuration / duration)
	if times < 1 {
		times = 1
	}

	a := &Average{
		ticker:     time.NewTicker(duration),
		counter:    NewCounter(),
		data:       make([]uint64, times),
		capacity:   times,
		sliceTimes: sliceTimes,
	}

	go a.Start()

	return a
}

func (a *Average) Value() []float64 {
	var result = make([]float64, len(a.sliceTimes))

	for i, sliceTime := range a.sliceTimes {
		var sum uint64
		var count int
		for j := 0; j < sliceTime; j++ {
			index := (int(a.current) - j - 1 + int(a.capacity)) % int(a.capacity)
			sum += a.data[index]
			count++
		}

		if count > 0 {
			result[i] = float64(sum) / float64(count) // Calculez la moyenne
		}
	}

	return result
}

func (a *Average) Start() {
	for range a.ticker.C {
		a.data[a.current] = a.counter.Value()
		a.counter.Reset()

		a.current = (a.current + 1) % a.capacity
	}
}

func (a *Average) Hit() {
	a.counter.Increment()
}

func (a *Average) Reset() {
	a.counter.Reset()
	a.data = make([]uint64, a.capacity)
	a.current = 0
	a.size = 0
}
