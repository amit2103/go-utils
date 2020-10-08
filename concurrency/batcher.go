package concurrency

import (
	"time"
)

//The purpose of this utility is to delay concurrent operations for a configured period of time
type Batcher struct {
	interval time.Duration
	queue    chan int
	waiters  AtomicInt32
	nextID   AtomicInt32
	after    func(time.Duration) <-chan time.Time
}

func NewBatcher(interval time.Duration) *Batcher {
	return &Batcher{
		interval: interval,
		queue:    make(chan int),
		waiters:  NewAtomicInt32(0),
		nextID:   NewAtomicInt32(0),
		after:    time.After,
	}
}

func newBatcherForTest(interval time.Duration, after func(time.Duration) <-chan time.Time) *Batcher {
	return &Batcher{
		interval: interval,
		queue:    make(chan int),
		waiters:  NewAtomicInt32(0),
		nextID:   NewAtomicInt32(0),
		after:    after,
	}
}

func (b *Batcher) Wait() int {
	numWaiters := b.waiters.Add(1)
	if numWaiters == 1 {
		b.newBatch()
	}
	return <-b.queue
}

func (b *Batcher) newBatch() {
	go func() {
		<-b.after(b.interval)

		id := b.nextID.Add(1)

		waiters := b.waiters.Get()
		for !b.waiters.CompareAndSwap(waiters, 0) {
			waiters = b.waiters.Get()
		}

		for i := int32(0); i < waiters; i++ {
			b.queue <- int(id)
		}
	}()
}
