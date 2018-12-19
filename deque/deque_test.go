package deque

import "testing"

func TestEmpty(t *testing.T) {
	var q Deque
	if q.Len() != 0 {
		t.Error("q.Len() =", q.Len(), "expect 0")
	}
}

func TestFrontBack(t *testing.T) {
	var q Deque
	q.PushLast("foo")
	q.PushLast("bar")
	q.PushLast("baz")
	if q.First() != "foo" {
		t.Error("wrong value at First of queue")
	}
	if q.Last() != "baz" {
		t.Error("wrong value at Last of queue")
	}

	if q.PopFirst() != "foo" {
		t.Error("wrong value removed from First of queue")
	}
	if q.First() != "bar" {
		t.Error("wrong value remaining at First of queue")
	}
	if q.Last() != "baz" {
		t.Error("wrong value remaining at Last of queue")
	}

	if q.PopLast() != "baz" {
		t.Error("wrong value removed from Last of queue")
	}
	if q.First() != "bar" {
		t.Error("wrong value remaining at First of queue")
	}
	if q.Last() != "bar" {
		t.Error("wrong value remaining at Last of queue")
	}
}

func TestGrowShrinkBack(t *testing.T) {
	var q Deque
	size := minSize * 2

	for i := 0; i < size; i++ {
		if q.Len() != i {
			t.Error("q.Len() =", q.Len(), "expected", i)
		}
		q.PushLast(i)
	}
	bufLen := len(q.buffer)

	// Remove from Last.
	for i := size; i > 0; i-- {
		if q.Len() != i {
			t.Error("q.Len() =", q.Len(), "expected", i)
		}
		x := q.PopLast()
		if x != i-1 {
			t.Error("q.PopLast() =", x, "expected", i-1)
		}
	}
	if q.Len() != 0 {
		t.Error("q.Len() =", q.Len(), "expected 0")
	}
	if len(q.buffer) == bufLen {
		t.Error("queue buffer did not shrink")
	}
}

func TestGrowShrinkFront(t *testing.T) {
	var q Deque
	size := minSize * 2

	for i := 0; i < size; i++ {
		if q.Len() != i {
			t.Error("q.Len() =", q.Len(), "expected", i)
		}
		q.PushLast(i)
	}
	bufLen := len(q.buffer)

	// Remove from First
	for i := 0; i < size; i++ {
		if q.Len() != size-i {
			t.Error("q.Len() =", q.Len(), "expected", minSize*2-i)
		}
		x := q.PopFirst()
		if x != i {
			t.Error("q.PopLast() =", x, "expected", i)
		}
	}
	if q.Len() != 0 {
		t.Error("q.Len() =", q.Len(), "expected 0")
	}
	if len(q.buffer) == bufLen {
		t.Error("queue buffer did not shrink")
	}
}

func TestSimple(t *testing.T) {
	var q Deque

	for i := 0; i < minSize; i++ {
		q.PushLast(i)
	}
	for i := 0; i < minSize; i++ {
		if q.First() != i {
			t.Error("peek", i, "had value", q.First())
		}
		x := q.PopFirst()
		if x != i {
			t.Error("remove", i, "had value", x)
		}
	}

	q.Clear()
	for i := 0; i < minSize; i++ {
		q.PushFirst(i)
	}
	for i := minSize - 1; i >= 0; i-- {
		x := q.PopFirst()
		if x != i {
			t.Error("remove", i, "had value", x)
		}
	}
}

func TestBufferWrap(t *testing.T) {
	var q Deque

	for i := 0; i < minSize; i++ {
		q.PushLast(i)
	}

	for i := 0; i < 3; i++ {
		q.PopFirst()
		q.PushLast(minSize + i)
	}

	for i := 0; i < minSize; i++ {
		if q.First().(int) != i+3 {
			t.Error("peek", i, "had value", q.First())
		}
		q.PopFirst()
	}
}

func TestBufferWrapReverse(t *testing.T) {
	var q Deque

	for i := 0; i < minSize; i++ {
		q.PushFirst(i)
	}
	for i := 0; i < 3; i++ {
		q.PopLast()
		q.PushFirst(minSize + i)
	}

	for i := 0; i < minSize; i++ {
		if q.Last().(int) != i+3 {
			t.Error("peek", i, "had value", q.First())
		}
		q.PopLast()
	}
}

func TestLen(t *testing.T) {
	var q Deque

	if q.Len() != 0 {
		t.Error("empty queue length not 0")
	}

	for i := 0; i < 1000; i++ {
		q.PushLast(i)
		if q.Len() != i+1 {
			t.Error("adding: queue with", i, "elements has length", q.Len())
		}
	}
	for i := 0; i < 1000; i++ {
		q.PopFirst()
		if q.Len() != 1000-i-1 {
			t.Error("removing: queue with", 1000-i-i, "elements has length", q.Len())
		}
	}
}

func TestBack(t *testing.T) {
	var q Deque

	for i := 0; i < minSize+5; i++ {
		q.PushLast(i)
		if q.Last() != i {
			t.Errorf("Last returned wrong value")
		}
	}
}

func checkRotate(t *testing.T, size int) {
	var q Deque
	for i := 0; i < size; i++ {
		q.PushLast(i)
	}

	for i := 0; i < q.Len(); i++ {
		x := i
		for n := 0; n < q.Len(); n++ {
			if q.At(n) != x {
				t.Fatalf("a[%d] != %d after rotate and copy", n, x)
			}
			x++
			if x == q.Len() {
				x = 0
			}
		}
		q.Rotate(1)
		if q.Last().(int) != i {
			t.Fatal("wrong value during rotation")
		}
	}
	for i := q.Len() - 1; i >= 0; i-- {
		q.Rotate(-1)
		if q.First().(int) != i {
			t.Fatal("wrong value during reverse rotation")
		}
	}
}

func TestRotate(t *testing.T) {
	checkRotate(t, 10)
	checkRotate(t, minSize)
	checkRotate(t, minSize+minSize/2)

	var q Deque
	for i := 0; i < 10; i++ {
		q.PushLast(i)
	}
	q.Rotate(11)
	if q.First() != 1 {
		t.Error("rotating 11 places should have been same as one")
	}
	q.Rotate(-21)
	if q.First() != 0 {
		t.Error("rotating -21 places should have been same as one -1")
	}
	q.Rotate(q.Len())
	if q.First() != 0 {
		t.Error("should not have rotated")
	}
	q.Clear()
	q.PushLast(0)
	q.Rotate(13)
	if q.First() != 0 {
		t.Error("should not have rotated")
	}
}

func TestAt(t *testing.T) {
	var q Deque

	for i := 0; i < 1000; i++ {
		q.PushLast(i)
	}

	// First to Last.
	for j := 0; j < q.Len(); j++ {
		if q.At(j).(int) != j {
			t.Errorf("index %d doesn't contain %d", j, j)
		}
	}

	// Last to First
	for j := 1; j <= q.Len(); j++ {
		if q.At(q.Len()-j).(int) != q.Len()-j {
			t.Errorf("index %d doesn't contain %d", q.Len()-j, q.Len()-j)
		}
	}
}

func TestClear(t *testing.T) {
	var q Deque

	for i := 0; i < 100; i++ {
		q.PushLast(i)
	}
	if q.Len() != 100 {
		t.Error("push: queue with 100 elements has length", q.Len())
	}
	cap := len(q.buffer)
	q.Clear()
	if q.Len() != 0 {
		t.Error("empty queue length not 0 after clear")
	}
	if len(q.buffer) != cap {
		t.Error("queue capacity changed after clear")
	}

	// Check that there are no remaining references after Clear()
	for i := 0; i < len(q.buffer); i++ {
		if q.buffer[i] != nil {
			t.Error("queue has non-nil deleted elements after Clear()")
			break
		}
	}
}

func TestInsert(t *testing.T) {
	q := new(Deque)
	for _, x := range "ABCDEFG" {
		q.PushLast(x)
	}
	insert(q, 4, 'x') // ABCDxEFG
	if q.At(4) != 'x' {
		t.Error("expected x at position 4")
	}

	insert(q, 2, 'y') // AByCDxEFG
	if q.At(2) != 'y' {
		t.Error("expected y at position 2")
	}
	if q.At(5) != 'x' {
		t.Error("expected x at position 5")
	}

	insert(q, 0, 'b') // bAByCDxEFG
	if q.First() != 'b' {
		t.Error("expected b inserted at First")
	}

	insert(q, q.Len(), 'e') // bAByCDxEFGe

	for i, x := range "bAByCDxEFGe" {
		if q.PopFirst() != x {
			t.Error("expected", x, "at position", i)
		}
	}
}

func TestRemove(t *testing.T) {
	q := new(Deque)
	for _, x := range "ABCDEFG" {
		q.PushLast(x)
	}

	if remove(q, 4) != 'E' { // ABCDFG
		t.Error("expected E from position 4")
	}

	if remove(q, 2) != 'C' { // ABDFG
		t.Error("expected C at position 2")
	}
	if q.Last() != 'G' {
		t.Error("expected G at Last")
	}

	if remove(q, 0) != 'A' { // BDFG
		t.Error("expected to remove A from First")
	}
	if q.First() != 'B' {
		t.Error("expected G at Last")
	}

	if remove(q, q.Len()-1) != 'G' { // BDF
		t.Error("expected to remove G from Last")
	}
	if q.Last() != 'F' {
		t.Error("expected F at Last")
	}

	if q.Len() != 3 {
		t.Error("wrong length")
	}
}

/*
func TestFrontBackOutOfRangePanics(t *testing.T) {
	const msg = "should panic when peeking empty queue"
	var q Deque
	assertPanics(t, msg, func() {
		q.First()
	})
	assertPanics(t, msg, func() {
		q.Last()
	})

	q.PushLast(1)
	q.PopFirst()

	assertPanics(t, msg, func() {
		q.First()
	})
	assertPanics(t, msg, func() {
		q.Last()
	})
}

func TestPopFrontOutOfRangePanics(t *testing.T) {
	var q Deque

	assertPanics(t, "should panic when removing empty queue", func() {
		q.PopFirst()
	})

	q.PushLast(1)
	q.PopFirst()

	assertPanics(t, "should panic when removing emptied queue", func() {
		q.PopFirst()
	})
}

func TestPopBackOutOfRangePanics(t *testing.T) {
	var q Deque

	assertPanics(t, "should panic when removing empty queue", func() {
		q.PopLast()
	})

	q.PushLast(1)
	q.PopLast()

	assertPanics(t, "should panic when removing emptied queue", func() {
		q.PopLast()
	})
}

func TestAtOutOfRangePanics(t *testing.T) {
	var q Deque

	q.PushLast(1)
	q.PushLast(2)
	q.PushLast(3)

	assertPanics(t, "should panic when negative index", func() {
		q.At(-4)
	})

	assertPanics(t, "should panic when index greater than length", func() {
		q.At(4)
	})
}

func TestInsertOutOfRangePanics(t *testing.T) {
	q := new(Deque)

	assertPanics(t, "should panic when inserting out of range", func() {
		insert(q, 1, "X")
	})

	q.PushLast("A")

	assertPanics(t, "should panic when inserting at negative index", func() {
		insert(q, -1, "Y")
	})

	assertPanics(t, "should panic when inserting out of range", func() {
		insert(q, 2, "B")
	})
}

func TestRemoveOutOfRangePanics(t *testing.T) {
	q := new(Deque)

	assertPanics(t, "should panic when removing from empty queue", func() {
		remove(q, 0)
	})

	q.PushLast("A")

	assertPanics(t, "should panic when removing at negative index", func() {
		remove(q, -1)
	})

	assertPanics(t, "should panic when removing out of range", func() {
		remove(q, 1)
	})
}

func assertPanics(t *testing.T, name string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s: didn't panic as expected", name)
		}
	}()

	f()
}

*/

func BenchmarkPushFront(b *testing.B) {
	var q Deque
	for i := 0; i < b.N; i++ {
		q.PushFirst(i)
	}
}

func BenchmarkPushBack(b *testing.B) {
	var q Deque
	for i := 0; i < b.N; i++ {
		q.PushLast(i)
	}
}

func BenchmarkSerial(b *testing.B) {
	var q Deque
	for i := 0; i < b.N; i++ {
		q.PushLast(i)
	}
	for i := 0; i < b.N; i++ {
		q.PopFirst()
	}
}

func BenchmarkSerialReverse(b *testing.B) {
	var q Deque
	for i := 0; i < b.N; i++ {
		q.PushFirst(i)
	}
	for i := 0; i < b.N; i++ {
		q.PopLast()
	}
}

func BenchmarkRotate(b *testing.B) {
	q := new(Deque)
	for i := 0; i < b.N; i++ {
		q.PushLast(i)
	}
	b.ResetTimer()
	// N complete rotations on length N - 1.
	for i := 0; i < b.N; i++ {
		q.Rotate(b.N - 1)
	}
}

func BenchmarkInsert(b *testing.B) {
	q := new(Deque)
	for i := 0; i < b.N; i++ {
		q.PushLast(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		insert(q, q.Len()/2, -i)
	}
}

func BenchmarkRemove(b *testing.B) {
	q := new(Deque)
	for i := 0; i < b.N; i++ {
		q.PushLast(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		remove(q, q.Len()/2)
	}
}

// insert is used to insert an element into the middle of the queue, before the
// element at the specified index.  insert(0,e) is the same as PushFirst(e) and
// insert(Len(),e) is the same as PushLast(e).  Complexity is constant plus
// linear in the lesser of the distances between i and either of the ends of
// the queue.  Accepts only non-negative index values, and panics if index is
// out of range.
func insert(q *Deque, i int, elem interface{}) {
	if i < 0 || i > q.Len() {
		panic("deque: Insert() called with index out of range")
	}
	if i == 0 {
		q.PushFirst(elem)
		return
	}
	if i == q.Len() {
		q.PushLast(elem)
		return
	}
	if i <= q.Len()/2 {
		q.Rotate(i)
		q.PushFirst(elem)
		q.Rotate(-i)
	} else {
		rots := q.Len() - i
		q.Rotate(-rots)
		q.PushLast(elem)
		q.Rotate(rots)
	}
}

// remove removes and returns an element from the middle of the queue, at the
// specified index.  remove(0) is the same as PopFirst() and remove(Len()-1) is
// the same as PopLast().  Complexity is constant plus linear in the lesser of
// the distances between i and either of the ends of the queue.  Accepts only
// non-negative index values, and panics if index is out of range.
func remove(q *Deque, i int) interface{} {
	if i < 0 || i >= q.Len() {
		panic("deque: Remove() called with index out of range")
	}
	if i == 0 {
		return q.PopFirst()
	}
	if i == q.Len()-1 {
		return q.PopLast()
	}
	if i <= q.Len()/2 {
		q.Rotate(i)
		elem := q.PopFirst()
		q.Rotate(-i)
		return elem
	}
	rots := q.Len() - 1 - i
	q.Rotate(-rots)
	elem := q.PopLast()
	q.Rotate(rots)
	return elem
}
