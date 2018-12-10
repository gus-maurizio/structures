// Package circular buffer implements operations on circular lists.
//package circbuffer
package main
import (
	"fmt"
	"time"
)

// A CircBuffer is an element of a circular list, or ring.
//
type CircBuffer struct {
	head, curr int
	Length     int
	Value      []interface{} // for use by client; untouched by this library
}

func (c *CircBuffer) init() *CircBuffer{
	c.head = 0
	c.curr = 0
	return c 
}

// New creates a circular buffer of n elements.
func New(n int) *CircBuffer {
	if n <= 0 {
		n = 1
	}
	c := new(CircBuffer) 
	c.Value = make([]interface{}, n, n)
	return c
}


func main() {
	fmt.Printf("Started %d\n",time.Now().UnixNano())
	buff := CircBuffer.New(40)
	fmt.Printf("Buffer  %#v\n",buff)
	
}
