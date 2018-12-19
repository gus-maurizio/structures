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
