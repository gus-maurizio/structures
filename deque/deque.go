// Copyright 2018 Gustavo Maurizio
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
// OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.
//

package deque

// Power of 2 for bitwise modulus: x % n == x & (n - 1).
const minSize = 64

// Deque represents a single instance of the data structure.
type Deque struct {
	buffer	[]interface{}
	first  	int
	last  	int
	size 	int
}

// Len returns the number of elements in the deque
func (q *Deque) Len() int {
	return q.size
}

//------------------------------------------------------------------------
// FIFO deque: 	add with PushLast()
// 				remove with PopFirst()
// LIFO deque: 	add with PushFirst()
// 				remove with PopLast()
// Pop from an empty deque returns nil
//------------------------------------------------------------------------

// PushLast appends an element to the Last of the deque 
func (q *Deque) PushLast(dequeitem interface{}) {
	q.resizeIfNeeded()

	q.buffer[q.last] = dequeitem
	// Calculate new last position.
	q.last = q.next(q.last)
	q.size++
}

// PushFirst adds an element to the First of the deque
func (q *Deque) PushFirst(dequeitem interface{}) {
	q.resizeIfNeeded()

	// Calculate new first position.
	q.first = q.prev(q.first)
	q.buffer[q.first] = dequeitem
	q.size++
}

// PopFirst removes and returns the first of the deque
func (q *Deque) PopFirst() interface{} {
	if q.size <= 0 { return nil }
	ret := q.buffer[q.first]
	q.buffer[q.first] = nil
	// Calculate new first position.
	q.first = q.next(q.first)
	q.size--

	q.compactIfNeeded()
	return ret
}

// PopLast removes and returns the element from the Last of the deque
func (q *Deque) PopLast() interface{} {
	if q.size <= 0 { return nil }

	// Calculate new last position
	q.last = q.prev(q.last)

	// Remove value at last.
	ret := q.buffer[q.last]
	q.buffer[q.last] = nil
	q.size--

	q.compactIfNeeded()
	return ret
}

// First returns (browse) the element at the First of the deque,
// that would be returned by PopFirst()
func (q *Deque) First() interface{} {
	if q.size <= 0 { return nil }
	return q.buffer[q.first]
}

// Last returns the element at the Last of the deque,
// that would be returned by PopLast()

func (q *Deque) Last() interface{} {
	if q.size <= 0 { return nil }
	return q.buffer[q.prev(q.last)]
}

// At returns (browse) the element at index i in the deque
// without removing the element. Index i is non negative.
// Index 0        is the first element and same as First()
// Index Len()-1  is the last  element and same as Last()
func (q *Deque) At(i int) interface{} {
	if i < 0 || i >= q.size { return nil }
	// bitwise modulus
	return q.buffer[(q.first+i)&(len(q.buffer)-1)]
}

// Clear removes all elements from the deque
func (q *Deque) Clear() {
	// bitwise modulus
	modBits := len(q.buffer) - 1
	for h := q.first; h != q.last; h = (h + 1) & modBits {
		q.buffer[h] = nil
	}
	q.first = 0
	q.last = 0
	q.size = 0
}

// Rotate rotates the deque +n steps First-to-Last
//                          -n steps Last-to-First
func (q *Deque) Rotate(n int) {
	if q.size <= 1 { return }
	// Rotating a multiple of q.size is same as no rotation.
	n %= q.size
	if n == 0 {	return }

	modBits := len(q.buffer) - 1
	// If no empty space in buffer, only move first and last indexes.
	if q.first == q.last {
		// Calculate new first and last using bitwise modulus.
		q.first = (q.first + n) & modBits
		q.last = (q.last + n) & modBits
		return
	}

	if n < 0 {
		// Rotate Last to First.
		for ; n < 0; n++ {
			// Calculate new first and last using bitwise modulus.
			q.first = (q.first - 1) & modBits
			q.last = (q.last - 1) & modBits
			// Put last value at first and remove value at last.
			q.buffer[q.first] = q.buffer[q.last]
			q.buffer[q.last] = nil
		}
		return
	}

	// Rotate First to Last.
	for ; n > 0; n-- {
		// Put first value at last and remove value at first.
		q.buffer[q.last] = q.buffer[q.first]
		q.buffer[q.first] = nil
		// Calculate new first and last using bitwise modulus.
		q.first = (q.first + 1) & modBits
		q.last = (q.last + 1) & modBits
	}
}

// prev returns the previous buffer position wrapping around buffer.
func (q *Deque) prev(i int) int {
	return (i - 1) & (len(q.buffer) - 1) // bitwise modulus
}

// next returns the next buffer position wrapping around buffer.
func (q *Deque) next(i int) int {
	return (i + 1) & (len(q.buffer) - 1) // bitwise modulus
}

// resizeIfNeeded resizes up if the buffer is full.
func (q *Deque) resizeIfNeeded() {
	if len(q.buffer) == 0 {
		q.buffer = make([]interface{}, minSize)
		return
	}
	if q.size == len(q.buffer) {
		q.resize()
	}
}

// compactIfNeeded resize down if the buffer 1/4 full.
func (q *Deque) compactIfNeeded() {
	if len(q.buffer) > minSize && (q.size<<2) == len(q.buffer) {
		q.resize()
	}
}

// resizes the deque to fit exactly twice its current contents
func (q *Deque) resize() {
	newBuf := make([]interface{}, q.size<<1)
	if q.last > q.first {
		copy(newBuf, q.buffer[q.first:q.last])
	} else {
		n := copy(newBuf, q.buffer[q.first:])
		copy(newBuf[n:], q.buffer[:q.last])
	}

	q.first = 0
	q.last = q.size
	q.buffer = newBuf
}
