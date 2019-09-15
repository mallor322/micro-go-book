package main

import (
	"ch5-feature/compute"
	"fmt"
)

func main()  {

	params := &compute.IntParams{
		P1:1,
		P2:2,
	}
	fmt.Println(params.Add())

}