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
 

import	(
	"fmt"
	"testing"
	)

func dump(c *circularbuffer) {
	if c == nil {
		fmt.Println("empty buffer")
		return
	}
	fmt.Printf("%#v\n",*c)
}

func TestNew(t * testing.T) {
	for i := 0; i < 10; i++ {
		cbuf := New(i,0)
		if cbuf.Len == i {fmt.Printf("iteration %d ok\n", i)}
	}
}

func TestInit(t * testing.T) {
	cbuf := New(30,5)
	cbuf.Init(6)
	total := 0
	for _, value := range cbuf.GetValues() {
		total += value.(int)
	}
	if total == 6 * 30 { 
		fmt.Printf("Match init\n")
	} else {
		fmt.Printf("MISMATCH 180 != %d\n",total)
		t.Fail()
	}
}

func TestMain(t * testing.T) {
        cbuf := New(5,9)
        for i := 1; i <= 8 ; i++ {
                old := cbuf.Push(i)
		dump(cbuf)
                fmt.Printf("iter %d last value: %d (old value: %v)  %+v \n", i, cbuf.Get(-1).(int), old, cbuf.GetValues())
        }
        fmt.Printf("%#v \n", cbuf.GetValues())
}
