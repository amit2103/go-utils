package concurrency

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSemaNoTimeout(t *testing.T) {
	s := NewSemaphore(1, 0)
	s.Lock()
	released := false
	go func() {
		released = true
		s.ReleaseLock()
	}()
	s.Lock()
	assert.True(t, released)
}

func TestSemaTimeout(t *testing.T) {
	s := NewSemaphore(1, 1*time.Millisecond)
	s.Lock()
	release := make(chan struct{})
	released := make(chan struct{})
	go func() {
		<-release
		s.ReleaseLock()
		released <- struct{}{}
	}()
	assert.False(t, s.Lock())
	release <- struct{}{}
	<-released
	assert.True(t, s.Lock())
}

func TestSemaAcquireContext(t *testing.T) {
	s := NewSemaphore(1, 0)
	s.Lock()
	release := make(chan struct{})
	released := make(chan struct{})
	go func() {
		<-release
		s.ReleaseLock()
		released <- struct{}{}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	assert.False(t, s.AcquireContext(ctx))
	release <- struct{}{}
	<-released
	assert.True(t, s.AcquireContext(context.Background()))
}

func TestSemaTryAcquire(t *testing.T) {
	s := NewSemaphore(1, 0)
	assert.True(t, s.TryLock())
	assert.False(t, s.TryLock())
	s.ReleaseLock()
	assert.True(t, s.TryLock())
}
