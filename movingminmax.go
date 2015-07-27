// movingminmax.go, jpad 2015
/*
Implementation based on the paper:

STREAMING MAXIMUM-MINIMUM FILTER
USING NO MORE THAN THREE COMPARISONS PER ELEMENT
Daniel Lemire
University of Quebec at Montreal (UQAM),
UER ST 100 Sherbrooke West, Montreal (Quebec), H2X 3P2 Canada
lemire@acm.org
*/

/*
Package movingminmax offers a moving minimum-maximum filter that can be used in
real-time contexts.
*/
package movingminmax

import (
	"fmt"

	"github.com/notnot/container/deque"
)

//// MovingMinMax //////////////////////////////////////////////////////////////

// MovingMinMax maintains moving minimum-maximum statistics.
type MovingMinMax struct {
	min float32
	max float32
	w   uint // sample data window width
	n   uint // number of samples processed
	lo  *deque.Deque
	up  *deque.Deque
}

// New returns a new MovingMinMax instance using a data window of size w.
func New(w uint) *MovingMinMax {
	return &MovingMinMax{
		w:  w,
		lo: deque.New(),
		up: deque.New(),
	}
}

// String returns a human readable description of the instance.
func (m *MovingMinMax) String() string {
	return fmt.Sprintf("moving min %.3f, max %.3f (%d samples)")
}

// Update updates the moving statistics with the given sample value.
func (m *MovingMinMax) Update(value float32) {
	switch {
	case m.n == 0: // initial minimum and maximum
		m.min = value
		m.max = value
	case m.n < m.w: // absolute minimum and maximum
		if value < m.min {
			m.min = value
		} else if value > m.max {
			m.max = value
		}
	default: // moving minimum and maximum
		//m.min =
		//m.max =
	}
	m.n++
}

// Min returns the current moving minimum.
func (m *MovingMinMax) Min() float32 {
	return m.min
}

// Max returns the current moving maximum.
func (m *MovingMinMax) Max() float32 {
	return m.max
}
