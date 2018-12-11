package main

import (
	"fmt"
	"github.com/gus-maurizio/structures/circularbuffer"
	//"time"
)

type measure struct {
	msgNum		int
	isAlert		bool
}
	
func main() {  
	c := circularbuffer.New(5,measure{})
	c.Init(measure{})
	for i := 1; i <= 8 ; i++ {
		msr := measure{ msgNum:	   i,
				isAlert:   true,
				}
		old := c.Push(msr)
		fmt.Printf("iter %03d %v %+v \n", i, old, c.GetValues())
	}
	fmt.Printf("%#v \n", c.GetValues()[0:4])
}
