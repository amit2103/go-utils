package concurrency

import (
	"context"
	"time"
)

// Implements a basic counting semaphore
type Semaphore struct {
	slots   chan struct{}
	timeout time.Duration
}

// create a new semaphore with count and timeout
func NewSemaphore(count int, timeout time.Duration) *Semaphore {
	sem := &Semaphore{
		slots:   make(chan struct{}, count),
		timeout: timeout,
	}
	for i := 0; i < count; i++ {
		sem.slots <- struct{}{}
	}
	return sem
}

//
func (sem *Semaphore) Lock() bool {
	if sem.timeout == 0 {
		<-sem.slots
		return true
	}
	tm := time.NewTimer(sem.timeout)
	defer tm.Stop()
	select {
	case <-sem.slots:
		return true
	case <-tm.C:
		return false
	}
}

func (sem *Semaphore) AcquireContext(ctx context.Context) bool {
	select {
	case <-sem.slots:
		return true
	case <-ctx.Done():
		return false
	}
}

func (sem *Semaphore) TryLock() bool {
	select {
	case <-sem.slots:
		return true
	default:
		return false
	}
}

func (sem *Semaphore) ReleaseLock() {
	sem.slots <- struct{}{}
}

func (sem *Semaphore) Size() int {
	return len(sem.slots)
}
