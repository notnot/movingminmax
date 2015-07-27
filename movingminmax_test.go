// movingminmax_test.go, jpad 2015

package movingminmax

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/notnot/container/deque_int"
)

//// tests /////////////////////////////////////////////////////////////////////

func TestReference(t *testing.T) {
	fmt.Printf("moving minmax, reference code (offline):\n")
	MovingMinMax_offline()

	// output results
	for i := 0; i < N; i++ {
		fmt.Printf("value[%02d] %.3f : min %.3f, max %.3f\n",
			i, values[i], mmins[i], mmaxs[i])
	}
}

//// benchmarks ////////////////////////////////////////////////////////////////

func BenchmarkReference(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MovingMinMax_offline()
	}
}

//// examples //////////////////////////////////////////////////////////////////

func ExampleEmpty() {

}

//// reference code ////////////////////////////////////////////////////////////

const (
	N = 25
	W = 5
)

var (
	values       []float32
	mmins, mmaxs []float32 // moving min and max results
)

func init() {
	fmt.Printf("initialising sample values...\n")

	values = make([]float32, N)
	mmins = make([]float32, N)
	mmaxs = make([]float32, N)

	rand.Seed(123)
	for i := 0; i < N; i++ {
		values[i] = rand.Float32()
	}
}

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
