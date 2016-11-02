package oots

import (
	"runtime"
	"sync"
)

// ThreadLimiter can be used to limit the amount
// of threads to runtime.NumCPU()
type ThreadLimiter struct {
	current int64
	mutex   *sync.Mutex
}

// NewThreadLimiter constructs a new ThreadLimiter
func NewThreadLimiter() *ThreadLimiter {
	return &ThreadLimiter{
		current: 0,
		mutex:   &sync.Mutex{},
	}
}

// Add increments the internal counter
func (tl *ThreadLimiter) Add(count int64) {
	tl.mutex.Lock()
	tl.current += count
	tl.mutex.Unlock()
}

// Done decrements the internal counter
func (tl *ThreadLimiter) Done() {
	tl.Add(-1)
}

// WaitTurn waits until internal counter is less than runtime.NumCPU()
func (tl *ThreadLimiter) WaitTurn() {
	for {
		tl.mutex.Lock()
		val := tl.current
		tl.mutex.Unlock()

		if val < int64(runtime.NumCPU()) {
			break
		}
	}
}

// Wait blocks until counter reaches zero
func (tl *ThreadLimiter) Wait() {
	for {
		tl.mutex.Lock()
		val := tl.current
		tl.mutex.Unlock()

		if val == 0 {
			break
		}

		runtime.Gosched()
	}
}
