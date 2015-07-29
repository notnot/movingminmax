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
	"github.com/notnot/container/deque"
	"fmt"
)

/////////////////////////////////
// MovingMinMax is an implementation of the algorithm described in the following paper :
//
// Daniel Lemire, Streaming Maximum-Minimum Filter Using No More than
// Three Comparisons per Element. Nordic Journal of Computing, 13 (4), pages 328-339, 2006.
//
// http://arxiv.org/abs/cs/0610046
//
//
/////////////////////////////////

////////
// To implement the algorithm from Lemire 2006, we can use a fixed amount of memory
// and minimal memory allocation.
// Following custom queue serves this purpose.
//////
type intfloatnode struct {
	index uint
	value float32
}

func newintfloatnode(i uint, v float32) *intfloatnode {
	return &intfloatnode{
		index: i,
		value: v,
	}
}

type intfloatqueue struct {
	nodes []*intfloatnode
	head  uint
	tail  uint
	count uint
	size  uint
}

// generate a power of two that is larger or equal than x 
func newpoweroftwo(x uint) uint {
	answer := uint(1)
	for answer < x { // (is there are faster way to do that?)
		answer *= 2 // let us prey that this does not overflow
	}
	return answer
}

func newintfloatqueue(size uint) *intfloatqueue {
	size = newpoweroftwo(size)
	n := make([]*intfloatnode, size)
	for i := uint(0); i < size; i++ {
		n[i] = &intfloatnode{}
	}
	return &intfloatqueue{
		nodes: n,
		size:  size,
	}
}
func (q *intfloatqueue) push(index uint, value float32) {
	q.nodes[q.tail].index = index
	q.nodes[q.tail].value = value
	q.tail = (q.tail + 1) & (q.size - 1)
	q.count++
}

func (q *intfloatqueue) tailnode() *intfloatnode {
	return q.nodes[(q.tail+q.size-1)&(q.size-1)]
}

func (q *intfloatqueue) poptail() *intfloatnode {
	q.tail = (q.tail + q.size - 1) & (q.size - 1)
	node := q.nodes[q.tail]
	q.count--
	return node
}

func (q *intfloatqueue) prunetail() {
	q.tail = (q.tail + q.size - 1) & (q.size - 1)
	q.count--
}

func (q *intfloatqueue) pop() *intfloatnode {
	node := q.nodes[q.head]
	q.head = (q.head + 1) & (q.size - 1)
	q.count--
	return node
}

func (q *intfloatqueue) headnode() *intfloatnode {
	return q.nodes[q.head]
}

//// MovingMinMax //////////////////////////////////////////////////////////////
//
// MovingMinMax maintains moving minimum-maximum statistics in O(1) time per
// value seen. It uses the algorithm from
//
// Daniel Lemire, Streaming Maximum-Minimum Filter Using No More than
// Three Comparisons per Element. Nordic Journal of Computing, 13 (4), pages 328-339, 2006.
//
// http://arxiv.org/abs/cs/0610046
//
// This implementation uses a fixed amount of memory and makes no dynamic allocation
// after the first call to Update.
//
type MovingMinMax struct {
	min float32
	max float32
	n   uint // how many values seen
	ww  uint // sample data window width
	lo  *intfloatqueue
	up  *intfloatqueue
}

// NewMovingMinMax returns a new instance using a data window of size w.
func NewMovingMinMax(w uint) *MovingMinMax {
	return &MovingMinMax{
		ww: w,
		n:  0,
		lo: newintfloatqueue(w),
		up: newintfloatqueue(w),
	}
}

// Update updates the moving statistics with the given sample value.
func (m *MovingMinMax) Update(value float32) {
	if m.up.count > 0 {
		if value > m.up.tailnode().value {
			m.up.prunetail()
			for (m.up.count > 0) && (value >= m.up.tailnode().value) {
				m.up.prunetail()
			}
		} else {
			m.lo.prunetail()
			for (m.lo.count > 0) && (value <= m.lo.tailnode().value) {
				m.lo.prunetail()
			}
		}
	}
	m.up.push(m.n, value)
	m.lo.push(m.n, value)
	if m.n == m.ww+m.lo.headnode().index {
		m.lo.pop()
	}
	if m.n == m.ww+m.up.headnode().index {
		m.up.pop()
	}
	m.max = m.up.headnode().value
	m.min = m.lo.headnode().value
	fmt.Println("val = ", value, " max = ", m.max, " min = ", m.min)
	m.n = m.n + 1
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
