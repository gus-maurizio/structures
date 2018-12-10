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

package circularbuffer
 

type circularbuffer struct {
	head, curr	int
	Len		int
	buffer		[]interface{}
}

// We allocate a circularbuffer structure that basically is a fixed size
// array of pointers (interface{}). We give the option of specifying the
// initial size of the array and using make with cap.
// The non-exposed head and curr are actually pointers as indexes.
func New(size int) *circularbuffer {
	c := &circularbuffer{head: 0, 
		curr: 0, 
		win1: 0,
		win2: 0,
		Len: size,
	}
	c.buffer = make([]interface{},size,size)
	return c
}

// Length of buffer. Not really needed since we export Len.
func (c *circularbuffer) Length() int {
	return c.Len
}

// This sets a particular value of the buffer array. The function 
// returns the previous value stored in the idx position.
// idx is a zero based index. This functions should be used with
// extreme caution. The preferred way is to use Push.
// The idx argument can be positive or negative, always considered
// from the start of the buffer pointed to by the value of head.
func (c *circularbuffer) Set(idx int, value interface{}) interface{} {
	if idx <= -c.Len {idx = idx + c.Len * (-idx/c.Len)}
	where := (c.head + idx + c.Len) % c.Len
	oldvalue := c.buffer[where]
	c.buffer[where] = value
	return oldvalue
}

// This function retrieves a particular value of the buffer
// at idx relative to head.
func (c *circularbuffer) Get(idx int) interface{} {
        if idx <= -c.Len {idx = idx + c.Len * (-idx/c.Len)}
        where := (c.head + idx + c.Len) % c.Len
        return c.buffer[where]
}
i
// We use Push to add elements to the circular buffer.
// Upon creation the buffer is empty. Once the whole 
// capacity of the buffer (defined by Len elements) is
// exhausted, the first value is dropped and replaced
// with the pushed new value. The head of buffer is 
// incremented and wrapped around if necessary.
func (c *circularbuffer) Push(value interface{}) interface{} {
	var oldvalue interface{}
	if c.curr < c.Len {
		// buffer is still filling up
		oldvalue = c.buffer[c.curr]
		c.buffer[c.curr] = value
		c.curr = c.curr + 1
	} else {
		// buffer is full, so head should take the new value
		oldvalue = c.buffer[c.head]
		c.buffer[c.head] = value
		c.head = (c.head + 1) % c.Len
	}
        return oldvalue
}

// Get the ordered list of values based on the current state.
func (c *circularbuffer) GetValues() []interface{} {
	return append(c.buffer[c.head:c.curr], c.buffer[0:c.head]...)
}

