package main

import "fmt"

func main()  {
	str := "Golang is Good!"
	strPrt := &str

	fmt.Printf("str type is %T, value is %v, address is %p\n", str, str, &str)
	fmt.Printf("strPtr type is %T, and value is %v\n", strPrt, strPrt)

	//strPtrPtr := &strPrt
	//fmt.Printf("strPtrPtr type is %T, and value is %v\n", strPtrPtr, strPtrPtr)

	newStr := *strPrt
	fmt.Printf("newStr type is %T, value is %v, and address is %p\n", newStr, newStr, &newStr)


    *strPrt = "Java is Good too!"
	fmt.Printf("newStr type is %T, value is %v, and address is %p\n", newStr, newStr, &newStr)
	fmt.Printf("str type is %T, value is %v, address is %p\n", str, str, &str)


}