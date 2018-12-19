# duplexqueue

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Fast ring-buffer duplex queue, or duplexqueue, or [double-ended queue](https://en.wikipedia.org/wiki/Double-ended_queue)) implementation.

[![GoDoc](https://godoc.org/github.com/gus-maurizio/structures/duplexqueue?status.svg)](https://godoc.org/github.com/gus-maurizio/structures/duplexqueue)

## Installation

```
$ go get github.com/gus-maurizio/structures/duplexqueue
```

## Duplexqueue data structure

Duplexqueue generalizes a queue and a stack, to efficiently add and remove items at either end with O(1) performance.  [Queue](https://en.wikipedia.org/wiki/Queue_(abstract_data_type)) (FIFO) operations are supported using `PushBack()` and `PopFront()`.  [Stack](https://en.wikipedia.org/wiki/Stack_(abstract_data_type)) (LIFO) operations are supported using `PushBack()` and `PopBack()`.

## Ring-buffer Performance

This duplexqueue implementation is optimized for CPU and GC performance.  The circular buffer automatically re-sizes by powers of two, growing when additional capacity is needed and shrinking when only a quarter of the capacity is used, and uses bitwise arithmetic for all calculations.  Since growth is by powers of two, adding elements will only cause O(log n) allocations.

The ring-buffer implementation significantly improves memory and time performance with fewer GC pauses, compared to implementations based on slices and linked lists.  By wrapping around the buffer, previously used space is reused, making allocation unnecessary until all buffer capacity is used.

For maximum speed, this duplexqueue implementation leaves concurrency safety up to the application to provide, however the application chooses, if needed at all.

## Reading Empty Duplexqueue

Since it is OK for the duplexqueue to contain a nil value, it is necessary to either panic or return a second boolean value to indicate the duplexqueue is empty, when reading or removing an element.  This duplexqueue panics when reading from an empty duplexqueue.  This is a run-time check to help catch programming errors, which may be missed if a second return value is ignored.  Simply check Duplexqueue.Len() before reading from the duplexqueue.

## Example

```go
package main

import (
    "fmt"
    "github.com/gus-maurizio/structures/duplexqueue"
)

func main() {
    var q duplexqueue.Duplexqueue
    q.PushBack("foo")
    q.PushBack("bar")
    q.PushBack("baz")

    fmt.Println(q.Len())   // Prints: 3
    fmt.Println(q.Front()) // Prints: foo
    fmt.Println(q.Back())  // Prints: baz

    q.PopFront() // remove "foo"
    q.PopBack()  // remove "baz"

    q.PushFront("hello")
    q.PushBack("world")

    // Consume duplexqueue and print elements.
    for q.Len() != 0 {
        fmt.Println(q.PopFront())
    }
}
```

## Uses

Duplexqueue can be used as both a:
- [Queue](https://en.wikipedia.org/wiki/Queue_(abstract_data_type)) using `PushBack` and `PopFront`
- [Stack](https://en.wikipedia.org/wiki/Stack_(abstract_data_type)) using `PushBack` and `PopBack`
