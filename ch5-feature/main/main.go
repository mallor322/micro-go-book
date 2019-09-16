package main

import (
	"fmt"
	"github.com/keets2012/Micro-Go-Pracrise/ch5-feature/compute"
)

func main()  {

	params := &compute.IntParams{
		P1:1,
		P2:2,
	}
	fmt.Println(params.Add())

}