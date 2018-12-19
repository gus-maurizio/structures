package duplexqueue

// minCapacity is the smallest capacity that duplexqueue may have.
// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
const minCapacity = 8

// Duplexqueue represents a single instance of the duplexqueue data structure.
type Duplexqueue struct {
	Buf   []interface{}		`json:"buffer"`
	Head  int 				`json:"qhead"`
	Tail  int 				`json:"qtail"`
	Count int 				`json:"qcount"`
}

// Len returns the number of elements currently stored in the queue.
func (q *Duplexqueue) Len() int {
	return q.Count
}

// PushBack appends an element to the back of the queue.  Implements FIFO when
// elements are removed with PopFront(), and LIFO when elements are removed
// with PopBack().
func (q *Duplexqueue) Init(qty int, elem interface{}) {
	q.Buf   = nil
	var qsz int
	if minCapacity > qty {
		qsz = minCapacity
	} else {
		qsz = (qty / 2) * 2
	}
	q.Buf = make([]interface{}, qsz)
	q.Head  = 0
	q.Tail  = 0
	q.Count = 0
	for i := 0; i < qty; i++ {
		q.Buf[i] = elem
	}
	q.Tail  += qty 
	q.Count  = qty
}

// PushBack appends an element to the back of the queue.  Implements FIFO when
// elements are removed with PopFront(), and LIFO when elements are removed
// with PopBack().
func (q *Duplexqueue) PushBack(elem interface{}) {
	q.growIfFull()

	q.Buf[q.Tail] = elem
	// Calculate new Tail position.
	q.Tail = q.next(q.Tail)
	q.Count++
}

// PushFront prepends an element to the front of the queue.
func (q *Duplexqueue) PushFront(elem interface{}) {
	q.growIfFull()

	// Calculate new Head position.
	q.Head = q.prev(q.Head)
	q.Buf[q.Head] = elem
	q.Count++
}


// PushPop pops back (returns that value) and pushes at front
func (q *Duplexqueue) PushPop(elem interface{}) interface{} {
	if q.Count <= 0 {
		panic("duplexqueue: PopBack() called on empty queue")
	}

	// Calculate new Tail position
	q.Tail = q.prev(q.Tail)
	// Remove value at Tail.
	ret := q.Buf[q.Tail]
	q.Buf[q.Tail] = nil
	// Calculate new Head position.
	q.Head = q.prev(q.Head)
	q.Buf[q.Head] = elem
	return ret
}



// PopFront removes and returns the element from the front of the queue.
// Implements FIFO when used with PushBack().  If the queue is empty, the call
// panics.
func (q *Duplexqueue) PopFront() interface{} {
	if q.Count <= 0 {
		panic("duplexqueue: PopFront() called on empty queue")
	}
	ret := q.Buf[q.Head]
	q.Buf[q.Head] = nil
	// Calculate new Head position.
	q.Head = q.next(q.Head)
	q.Count--

	q.shrinkIfExcess()
	return ret
}

// PopBack removes and returns the element from the back of the queue.
// Implements LIFO when used with PushBack().  If the queue is empty, the call
// panics.
func (q *Duplexqueue) PopBack() interface{} {
	if q.Count <= 0 {
		panic("duplexqueue: PopBack() called on empty queue")
	}

	// Calculate new Tail position
	q.Tail = q.prev(q.Tail)

	// Remove value at Tail.
	ret := q.Buf[q.Tail]
	q.Buf[q.Tail] = nil
	q.Count--

	q.shrinkIfExcess()
	return ret
}

// Front returns the element at the front of the queue.  This is the element
// that would be returned by PopFront().  This call panics if the queue is
// empty.
func (q *Duplexqueue) Front() interface{} {
	if q.Count <= 0 {
		panic("duplexqueue: Front() called when empty")
	}
	return q.Buf[q.Head]
}

// Back returns the element at the back of the queue.  This is the element
// that would be returned by PopBack().  This call panics if the queue is
// empty.
func (q *Duplexqueue) Back() interface{} {
	if q.Count <= 0 {
		panic("duplexqueue: Back() called when empty")
	}
	return q.Buf[q.prev(q.Tail)]
}

// At returns the element at index i in the queue without removing the element
// from the queue.  This method accepts only non-negative index values.  At(0)
// refers to the first element and is the same as Front().  At(Len()-1) refers
// to the last element and is the same as Back().  If the index is invalid, the
// call panics.
//
// The purpose of At is to allow Duplexqueue to serve as a more general purpose
// circular buffer, where items are only added to and removed from the the ends
// of the duplexqueue, but may be read from any place within the duplexqueue.  Consider the
// case of a fixed-size circular log buffer: A new entry is pushed onto one end
// and when full the oldest is popped from the other end.  All the log entries
// in the buffer must be readable without altering the buffer contents.
func (q *Duplexqueue) At(i int) interface{} {
	if i < 0 || i >= q.Count {
		panic("duplexqueue: At() called with index out of range")
	}
	// bitwise modulus
	return q.Buf[(q.Head+i)&(len(q.Buf)-1)]
}

func (q *Duplexqueue) Index(i int) interface{} {
	if i == 0 { return q.Buf[q.Head] }
	if i >= q.Count  || i < -q.Count  { i = i % q.Count }
	if i < 0 { i += q.Count }
	// bitwise modulus
	return q.Buf[(q.Head+i)&(len(q.Buf)-1)]
}


// Clear removes all elements from the queue, but retains the current capacity.
// This is useful when repeatedly reusing the queue at high frequency to avoid
// GC during reuse.  The queue will not be resized smaller as long as items are
// only added.  Only when items are removed is the queue subject to getting
// resized smaller.
func (q *Duplexqueue) Clear() {
	// bitwise modulus
	modBits := len(q.Buf) - 1
	for h := q.Head; h != q.Tail; h = (h + 1) & modBits {
		q.Buf[h] = nil
	}
	q.Head = 0
	q.Tail = 0
	q.Count = 0
}

// Rotate rotates the duplexqueue n steps front-to-back.  If n is negative, rotates
// back-to-front.  Having Duplexqueue provide Rotate() avoids resizing that could
// happen if implementing rotation using only Pop and Push methods.
func (q *Duplexqueue) Rotate(n int) {
	if q.Count <= 1 {
		return
	}
	// Rotating a multiple of q.Count is same as no rotation.
	n %= q.Count
	if n == 0 {
		return
	}

	modBits := len(q.Buf) - 1
	// If no empty space in buffer, only move Head and Tail indexes.
	if q.Head == q.Tail {
		// Calculate new Head and Tail using bitwise modulus.
		q.Head = (q.Head + n) & modBits
		q.Tail = (q.Tail + n) & modBits
		return
	}

	if n < 0 {
		// Rotate back to front.
		for ; n < 0; n++ {
			// Calculate new Head and Tail using bitwise modulus.
			q.Head = (q.Head - 1) & modBits
			q.Tail = (q.Tail - 1) & modBits
			// Put Tail value at Head and remove value at Tail.
			q.Buf[q.Head] = q.Buf[q.Tail]
			q.Buf[q.Tail] = nil
		}
		return
	}

	// Rotate front to back.
	for ; n > 0; n-- {
		// Put Head value at Tail and remove value at Head.
		q.Buf[q.Tail] = q.Buf[q.Head]
		q.Buf[q.Head] = nil
		// Calculate new Head and Tail using bitwise modulus.
		q.Head = (q.Head + 1) & modBits
		q.Tail = (q.Tail + 1) & modBits
	}
}

// prev returns the previous buffer position wrapping around buffer.
func (q *Duplexqueue) prev(i int) int {
	return (i - 1) & (len(q.Buf) - 1) // bitwise modulus
}

// next returns the next buffer position wrapping around buffer.
func (q *Duplexqueue) next(i int) int {
	return (i + 1) & (len(q.Buf) - 1) // bitwise modulus
}

// growIfFull resizes up if the buffer is full.
func (q *Duplexqueue) growIfFull() {
	if len(q.Buf) == 0 {
		q.Buf = make([]interface{}, minCapacity)
		return
	}
	if q.Count == len(q.Buf) {
		q.resize()
	}
}

// shrinkIfExcess resize down if the buffer 1/4 full.
func (q *Duplexqueue) shrinkIfExcess() {
	if len(q.Buf) > minCapacity && (q.Count<<2) == len(q.Buf) {
		q.resize()
	}
}

// resize resizes the duplexqueue to fit exactly twice its current contents.  This is
// used to grow the queue when it is full, and also to shrink it when it is
// only a quarter full.
func (q *Duplexqueue) resize() {
	newBuf := make([]interface{}, q.Count<<1)
	if q.Tail > q.Head {
		copy(newBuf, q.Buf[q.Head:q.Tail])
	} else {
		n := copy(newBuf, q.Buf[q.Head:])
		copy(newBuf[n:], q.Buf[:q.Tail])
	}

	q.Head = 0
	q.Tail = q.Count
	q.Buf = newBuf
}

func (q *Duplexqueue) Do(f func(interface{})) {
	for i := 0; i < q.Count; i++ {
		f(q.Buf[(q.Head+i)&(len(q.Buf)-1)])
	}
}

func (q *Duplexqueue) DoIndex(idx int, f func(interface{})) {
	for i := 0; i < q.Count; i++ {
		f(q.Index(idx + i))
	}
}

func (q *Duplexqueue) DoFor(idx int, cnt int, f func(interface{})) {
	if cnt > q.Count { cnt %= q.Count }
	for i := 0; i < cnt; i++ {
		f(q.Index(idx + i))
	}
}
