package main

import (
	"fmt"
	"github.com/gus-maurizio/structures/circularbuffer"
	"time"
)

type measure struct {
	msgNum		int
	timestamp	int64
	isAlert		bool
	message		string
}
	
func main() {  
	c := circularbuffer.New(5)
	c.Init(measure{})
	for i := 1; i <= 21 ; i++ {
		msr := measure{ msgNum:	   i,
				timestamp: time.Now().UnixNano(),
				message:   "message", 
				isAlert:   true,
				}
		old := c.Push(msr)
		fmt.Printf("iter %03d %v %+v \n", i, old, c.GetValues())
	}
	fmt.Printf("%#v \n", c.GetValues()[0:4])
}
