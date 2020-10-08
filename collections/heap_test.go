package collections

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHeap(t *testing.T) {
	Convey("Given a heap", t, func() {
		h := newHeap()

		Convey("When Heap is empty", func() {

			Convey("size is o", func() {
				size := h.size()
				So(size, ShouldEqual, 0)
			})

			Convey("No maximum node present", func() {
				max := h.maximum()
				So(max.value, ShouldBeEmpty)
			})

		})

		Convey("When insert 50 nodes", func() {
			for i := 0; i < 50; i++ {
				h.insert([]float64{}, float64(i), i)
			}
			max1 := h.maximum()
			h.extractMax()
			h.extractMax()
			h.extractMax()
			max2 := h.maximum()

			Convey("The max1.length should be 49", func() {
				So(max1.length, ShouldEqual, 49)
			})
			Convey("The max2.length should be 46", func() {
				So(max2.length, ShouldEqual, 46)
			})
		})

	})
}
