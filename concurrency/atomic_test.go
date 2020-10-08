package concurrency


import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestAtomicInt32(t *testing.T) {
	i := NewAtomicInt32(1)
	assert.Equal(t, int32(1), i.Get())

	i.Set(2)
	assert.Equal(t, int32(2), i.Get())

	i.Add(1)
	assert.Equal(t, int32(3), i.Get())

	i.CompareAndSwap(3, 4)
	assert.Equal(t, int32(4), i.Get())

	i.CompareAndSwap(3, 5)
	assert.Equal(t, int32(4), i.Get())
}

func TestAtomicInt64(t *testing.T) {
	i := NewAtomicInt64(1)
	assert.Equal(t, int64(1), i.Get())

	i.Set(2)
	assert.Equal(t, int64(2), i.Get())

	i.Add(1)
	assert.Equal(t, int64(3), i.Get())

	i.CompareAndSwap(3, 4)
	assert.Equal(t, int64(4), i.Get())

	i.CompareAndSwap(3, 5)
	assert.Equal(t, int64(4), i.Get())
}

func TestAtomicDuration(t *testing.T) {
	d := NewAtomicDuration(time.Second)
	assert.Equal(t, time.Second, d.Get())

	d.Set(time.Second * 2)
	assert.Equal(t, time.Second*2, d.Get())

	d.Add(time.Second)
	assert.Equal(t, time.Second*3, d.Get())

	d.CompareAndSwap(time.Second*3, time.Second*4)
	assert.Equal(t, time.Second*4, d.Get())

	d.CompareAndSwap(time.Second*3, time.Second*5)
	assert.Equal(t, time.Second*4, d.Get())
}

func TestAtomicString(t *testing.T) {
	var s AtomicString
	assert.Equal(t, "", s.Get())

	s.Set("a")
	assert.Equal(t, "a", s.Get())

	assert.Equal(t, false, s.CompareAndSwap("b", "c"))
	assert.Equal(t, "a", s.Get())

	assert.Equal(t, true, s.CompareAndSwap("a", "c"))
	assert.Equal(t, "c", s.Get())
}

func TestAtomicBool(t *testing.T) {
	b := NewAtomicBool(true)
	assert.Equal(t, true, b.Get())

	b.Set(false)
	assert.Equal(t, false, b.Get())

	b.Set(true)
	assert.Equal(t, true, b.Get())

	assert.Equal(t, false, b.CompareAndSwap(false, true))

	assert.Equal(t, true, b.CompareAndSwap(true, false))

	assert.Equal(t, true, b.CompareAndSwap(false, false))

	assert.Equal(t, true, b.CompareAndSwap(false, true))

	assert.Equal(t, true, b.CompareAndSwap(true, true))
}
