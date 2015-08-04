// movingminmax.go, jpad 2015

/*
Package movingminmax provides an efficient O(1) moving minimum-maximum filter
that can be used in real-time contexts. It uses the algorithm from:

Daniel Lemire,
Streaming Maximum-Minimum Filter
Using No More than Three Comparisons per Element.
Nordic Journal of Computing, 13 (4), pages 328-339, 2006.

http://arxiv.org/abs/cs/0610046

This implementation uses a fixed amount of memory and makes no dynamic
allocations during updates.
*/
package movingminmax

import (
//	"fmt"
)

//// MovingMinMax //////////////////////////////////////////////////////////////

// MovingMinMax maintains moving minimum-maximum statistics.
type MovingMinMax struct {
	ww uint // window width
	n  uint // number of samples processed
	lo *deque_IV
	up *deque_IV
}

// NewMovingMinMax returns a new instance using a data window of size w.
func NewMovingMinMax(w uint) *MovingMinMax {
	return &MovingMinMax{
		ww: w,
		lo: newDeque_IV(w),
		up: newDeque_IV(w),
	}
}

// Update updates the statistics with the given sample value.
func (m *MovingMinMax) Update(value float32) {
	if m.up.Size() > 0 {
		if value > m.up.BackItem().v {
			m.up.PruneBack()
			for (m.up.Size() > 0) && (value >= m.up.BackItem().v) {
				m.up.PruneBack()
			}
		} else {
			m.lo.PruneBack()
			for (m.lo.Size() > 0) && (value <= m.lo.BackItem().v) {
				m.lo.PruneBack()
			}
		}
	}

	m.lo.PushBack(m.n, value)
	m.up.PushBack(m.n, value)

	if m.n == m.ww+m.lo.FrontItem().i {
		m.lo.PruneFront()
	}
	if m.n == m.ww+m.up.FrontItem().i {
		m.up.PruneFront()
	}
	m.n++
}

// Min returns the current moving minimum.
func (m *MovingMinMax) Min() float32 {
	return m.lo.FrontItem().v
}

// Max returns the current moving maximum.
func (m *MovingMinMax) Max() float32 {
	return m.up.FrontItem().v
}

//// MovingMin /////////////////////////////////////////////////////////////////

// MovingMin maintains moving minimum statistics.
type MovingMin struct {
	ww uint // sample data window width
	n  uint // number of samples processed
	lo *deque_IV
}

// NewMovingMin returns a new instance using a data window of size w.
func NewMovingMin(w uint) *MovingMin {
	return &MovingMin{
		ww: w,
		lo: newDeque_IV(w),
	}
}

// Update updates the statistics with the given sample value.
func (m *MovingMin) Update(value float32) {
	if (m.lo.Size() > 0) && (value < m.lo.BackItem().v) {
		m.lo.PruneBack()
		for (m.lo.Size() > 0) && (value <= m.lo.BackItem().v) {
			m.lo.PruneBack()
		}
	}

	m.lo.PushBack(m.n, value)

	if m.n == m.ww+m.lo.FrontItem().i {
		m.lo.PruneFront()
	}
	m.n++
}

// Min returns the current moving minimum.
func (m *MovingMin) Min() float32 {
	return m.lo.FrontItem().v
}

//// MovingMax /////////////////////////////////////////////////////////////////

// MovingMax maintains moving maximum statistics.
type MovingMax struct {
	ww uint // sample data window width
	n  uint // number of samples processed
	up *deque_IV
}

// NewMovingMax returns a new instance using a data window of size w.
func NewMovingMax(w uint) *MovingMax {
	return &MovingMax{
		ww: w,
		up: newDeque_IV(w),
	}
}

// Update updates the statistics with the given sample value.
func (m *MovingMax) Update(value float32) {
	if (m.up.Size() > 0) && (value > m.up.BackItem().v) {
		m.up.PruneBack()
		for (m.up.Size() > 0) && (value >= m.up.BackItem().v) {
			m.up.PruneBack()
		}
	}

	m.up.PushBack(m.n, value)

	if m.n == m.ww+m.up.FrontItem().i {
		m.up.PruneFront()
	}
	m.n++
}

// Max returns the current moving maximum.
func (m *MovingMax) Max() float32 {
	return m.up.FrontItem().v
}

//// MovingMean ////////////////////////////////////////////////////////////////

// MovingMean maintains moving mean statistics.
type MovingMean struct {
	ww   uint
	mean float32
	sum  float32
	invw float32
	dq   *deque_f32
}

// NewMovingMean returns a new instance using a data window of size w.
func NewMovingMean(w uint) *MovingMean {
	mmean := MovingMean{
		ww:   w,
		invw: 1.0 / float32(w),
		dq:   newDeque_f32(w),
	}
	// init deque
	mmean.dq.size = w
	mmean.dq.back = w
	return &mmean
}

func (m *MovingMean) Update(value float32) {
	m.dq.PushFront(value)
	m.sum -= m.dq.PopBack()
	m.sum += value
	m.mean = m.sum * m.invw
}

func (m *MovingMean) Mean() float32 {
	return m.mean
}

//// MovingMinMax0 /////////////////////////////////////////////////////////////

// MovingMinMax0 maintains moving minimum-maximum statistics.
type MovingMinMax0 struct {
	n  uint // how many values seen
	ww uint // sample data window width
	lo *intfloatqueue
	up *intfloatqueue
}

// NewMovingMinMax0 returns a new instance using a data window of size w.
func NewMovingMinMax0(w uint) *MovingMinMax0 {
	return &MovingMinMax0{
		ww: w,
		n:  0,
		lo: newintfloatqueue(w),
		up: newintfloatqueue(w),
	}
}

// Update updates the moving statistics with the given sample value.
func (m *MovingMinMax0) Update(value float32) {
	if m.up.nonempty() {
		if value > m.up.tailvalue() {
			m.up.prunetail()
			for (m.up.nonempty()) && (value >= m.up.tailvalue()) {
				m.up.prunetail()
			}
		} else {
			m.lo.prunetail()
			for (m.lo.nonempty()) && (value <= m.lo.tailvalue()) {
				m.lo.prunetail()
			}
		}
	}
	m.up.push(m.n, value)
	if m.n == m.ww+m.up.headindex() {
		m.up.prunehead()
	}
	m.lo.push(m.n, value)
	if m.n == m.ww+m.lo.headindex() {
		m.lo.prunehead()
	}
	m.n = m.n + 1
}

// Min returns the current moving minimum.
func (m *MovingMinMax0) Min() float32 {
	return m.lo.headvalue()
}

// Max returns the current moving maximum.
func (m *MovingMinMax0) Max() float32 {
	return m.up.headvalue()
}

//// _IV ///////////////////////////////////////////////////////////////////////

type _IV struct {
	i uint    // sample index
	v float32 // sample value
}
