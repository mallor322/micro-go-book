package main

import "fmt"

type aliasInt = int // 定义一个类型别名
type myInt int // 定义一个新的类型

func main()  {

	var alias aliasInt
	fmt.Printf("alias value is %v, type is %T\n", alias, alias)

	var myint myInt
	fmt.Printf("myint value is %v, type is %T\n", myint, myint)
	

}

