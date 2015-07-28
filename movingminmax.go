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
Package movingminmax provides a moving minimum-maximum filter that can be used
in real-time contexts.
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
	ww  uint // sample data window width
	n   uint // number of samples processed
	lo  *deque.Deque
	up  *deque.Deque
}

// NewMovingMinMax returns a new instance using a data window of size w.
func NewMovingMinMax(w uint) *MovingMinMax {
	return &MovingMinMax{
		ww: w,
		lo: deque.New(),
		up: deque.New(),
	}
}

// Update updates the moving statistics with the given sample value.
func (m *MovingMinMax) Update(value float32) {
	switch {
	case m.n == 0: // initial minimum and maximum
		m.min = value
		m.max = value
	case m.n < m.ww: // absolute minimum and maximum
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

//// MovingMin /////////////////////////////////////////////////////////////////

type MovingMin struct {
	min float32
	ww  uint         // sample data window width
	n   uint         // number of samples processed
	iv  *deque.Deque // indices & values
}

func NewMovingMin(w uint) *MovingMin {
	return &MovingMin{
		ww: w,
		iv: deque.New(),
	}
}

func (m *MovingMin) Update(value float32) {
	// delete front item if it is too old
	if m.iv.Size() > 0 && m.iv.FrontItem().(_IV).i <= m.n-m.ww {
		m.iv.PopFront()
	}
	// delete items that can't become a minimum
	for m.iv.Size() > 0 && m.iv.BackItem().(_IV).v > value {
		m.iv.PopBack()
	}
	m.iv.PushBack(_IV{m.n, value})

	switch {
	case m.n == 0: // initial minimum
		m.min = value
	case m.n < m.ww: // absolute minimum
		if value < m.min {
			m.min = value
		}
	default: // moving minimum
		m.min = m.iv.FrontItem().(_IV).v
	}
	m.n++
}

func (m *MovingMin) Min() float32 {
	return m.min
}

//// MovingMax /////////////////////////////////////////////////////////////////

type MovingMax struct {
	max float32
	ww  uint         // sample data window width
	n   uint         // number of samples processed
	iv  *deque.Deque // indices & values
}

func NewMovingMax(w uint) *MovingMax {
	return &MovingMax{
		ww: w,
		iv: deque.New(),
	}
}

func (m *MovingMax) Update(value float32) {
	// delete front item if it is too old
	if m.iv.Size() > 0 && m.iv.FrontItem().(_IV).i <= m.n-m.ww {
		m.iv.PopFront()
	}
	// delete items that can't become a maximum
	for m.iv.Size() > 0 && m.iv.BackItem().(_IV).v < value {
		m.iv.PopBack()
	}
	m.iv.PushBack(_IV{m.n, value})

	switch {
	case m.n == 0: // initial minimum
		m.max = value
	case m.n < m.ww: // absolute minimum
		if value > m.max {
			m.max = value
		}
	default: // moving minimum
		m.max = m.iv.FrontItem().(_IV).v
	}
	m.n++
}

func (m *MovingMax) Max() float32 {
	return m.max
}

//// _IV ///////////////////////////////////////////////////////////////////////

type _IV struct {
	i uint    // sample index
	v float32 // sample value
}
