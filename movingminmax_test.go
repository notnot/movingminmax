// movingminmax_test.go, jpad 2015

package movingminmax

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/notnot/container/deque_int"
)

const (
	N = 25
	W = 5
)

var (
	values       []float32 // input sample values
	mmins, mmaxs []float32 // reference results
)

// init initializes the test values and computes and displays the reference
// results.
func init() {
	values = make([]float32, N)
	mmins = make([]float32, N)
	mmaxs = make([]float32, N)

	// generate random sample values in [0.0, 1.0)
	rand.Seed(123)
	for i := 0; i < N; i++ {
		values[i] = rand.Float32()
	}

	fmt.Printf("moving minmax, reference code (offline):\n")
	MovingMinMax_offline()

	// display results
	for i := 0; i < N; i++ {
		fmt.Printf("value[%02d] %.3f : min %.3f, max %.3f\n",
			i, values[i], mmins[i], mmaxs[i])
	}
}

//// tests /////////////////////////////////////////////////////////////////////

func TestReference(t *testing.T) {
	minmax := New(W)
	for i := range values {
		minmax.Update(values[i])

		min := minmax.Min()
		if min != mmins[i] {
			t.Errorf("values[%d] Min() got: %.3f, want: %.3f",
				i, min, mmins[i])
		}

		max := minmax.Max()
		if max != mmaxs[i] {
			t.Errorf("values[%d] Max() got: %.3f, want: %.3f",
				i, max, mmaxs[i])
		}
	}
}

//// benchmarks ////////////////////////////////////////////////////////////////

func BenchmarkReference(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MovingMinMax_offline()
	}
}

func BenchmarkMovingMinMax(b *testing.B) {
	minmax := New(W)
	for j := 0; j < b.N; j++ {
		for i := range values {
			minmax.Update(values[i])
		}
	}
}

//// examples //////////////////////////////////////////////////////////////////

func ExampleEmpty() {

}

//// reference code ////////////////////////////////////////////////////////////

func MovingMinMax_offline() {
	U := deque_int.New() // upper -> indices of maxima
	L := deque_int.New() // lower -> indices of minima

	// initial minimum and maximum
	mmins[0] = values[0]
	mmaxs[0] = values[0]

	for i := 1; i < len(values); i++ {
		if i < W {
			// 'absolute' minimum and maximum
			if values[i] < mmins[i-1] {
				mmins[i] = values[i]
			} else {
				mmins[i] = mmins[i-1]
			}
			if values[i] > mmaxs[i-1] {
				mmaxs[i] = values[i]
			} else {
				mmaxs[i] = mmaxs[i-1]
			}
		} else {
			// 'moving' minimum and maximum
			mmins[i-1] = values[iExtreme(L, i)]
			mmaxs[i-1] = values[iExtreme(U, i)]
		}

		// update monotonic wedge
		if values[i] > values[i-1] {
			L.PushBack(i - 1)
			if i == W+L.FrontItem() {
				L.PopFront()
			}
			for U.Size() > 0 {
				if values[i] <= values[U.BackItem()] {
					if i == W+U.FrontItem() {
						U.PopFront()
					}
					break
				}
				U.PopBack()
			}
		} else {
			U.PushBack(i - 1)
			if i == W+U.FrontItem() {
				U.PopFront()
			}
			for L.Size() > 0 {
				if values[i] >= values[L.BackItem()] {
					if i == W+L.FrontItem() {
						L.PopFront()
					}
					break
				}
				L.PopBack()
			}
		}
	}

	// final minimum and maximum
	i := len(values) - 1
	mmins[i] = values[iExtreme(L, i)]
	mmaxs[i] = values[iExtreme(U, i)]
}

func iExtreme(d *deque_int.Deque, i int) int {
	if d.Size() > 0 {
		return d.FrontItem()
	} else {
		return i - 1
	}
}
