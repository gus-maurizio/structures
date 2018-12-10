package main

import (
	"fmt"
	"github.com/gus-maurizio/structures/circularbuffer"
)

func main() {  
	c := circularbuffer.New(5)
	for i := 1; i <= 21 ; i++ {
		old := c.Push(i)
		fmt.Printf("iter %03d %v %+v \n", i, old, c.GetValues())
	}
	fmt.Printf("%+v \n", *c)
}
