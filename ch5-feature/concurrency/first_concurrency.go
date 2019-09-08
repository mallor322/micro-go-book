package main

import (
	"fmt"
	"time"
)

func setVTo1(v *int)  {

	*v = 1
}

func setVTo2(v *int)  {
	*v = 2
}

func main()  {

	go func(name string) {
		fmt.Println("Hello " + name )
	}("xuan")

	time.Sleep(time.Second)
}




