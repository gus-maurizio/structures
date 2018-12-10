package main

import (
	"fmt"
	"github.com/gus-maurizio/structures/circularbuffer"
	"time"
)

type measure struct {
	timestamp	int64
	message		string
}
	
func main() {  
	c := circularbuffer.New(5)
	for i := 1; i <= 21 ; i++ {
		msr := measure{ timestamp: time.Now().UnixNano(),
				message:   fmt.Sprintf("this is time %v",time.Now()),
				}
		old := c.Push(msr)
		fmt.Printf("iter %03d %v %+v \n", i, old, c.GetValues())
	}
	fmt.Printf("%#v \n", c.GetValues()[0:4])
}
