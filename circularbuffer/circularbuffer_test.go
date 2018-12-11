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

func TestNew(t *testing.T) {
	for i := 0; i < 10; i++ {
		cbuf := New(i,0)
		if cbuf.Len == i {fmt.Printf("iteration %d ok\n", i)}
	}
}

func TestInit(t *testing.T) {
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

func TestMain(t *testing.T) {
        cbuf := New(5,9)
        for i := 1; i <= 8 ; i++ {
                old := cbuf.Push(i)
		dump(cbuf)
                fmt.Printf("iter %d last value: %d (old value: %v)  %+v \n", i, cbuf.Get(-1).(int), old, cbuf.GetValues())
        }
        fmt.Printf("%#v \n", cbuf.GetValues())
}

func TestGetSet(t *testing.T) {
	cbuf := New(10,0)
	for i := -10; i < 10; i++ {
		old := cbuf.Set(i, i*i)
		fmt.Printf("Setting %d with %d (old %d) for %v\n",i,i*i, old.(int), cbuf.GetValues())
	}
}


func TestRollingWindow(t *testing.T) {
	cbuf := New(20,0)
	for i := 1; i < 20; i++ { cbuf.Set(i, i*i)}
	fmt.Printf("Init values: %v\n",cbuf.GetValues())
	
	fmt.Printf("Roll Window 05 : %v\n",cbuf.GetValues()[cbuf.Len-5:cbuf.Len])
}


func TestDoFunction(t *testing.T) {
	cbuf := New(20,0)
	for i := 1; i < 20; i++ { cbuf.Set(i, i*i)}
	n := 0
	s := 0
	cbuf.Do( func(p interface{}) {
		n++
		if p != nil { s += p.(int) }
		})
	fmt.Printf("N and S: %d %d %v\n",n,s,cbuf.GetValues())
}

