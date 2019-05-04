package main

import "fmt"

type Person struct {
	Name string	// 姓名
	Birth string	// 生日
	ID int64	// 身份证号
}

func main()  {
	// 声明实例化
	var p1 Person
	p1.Name =  "王小二"
	p1.Birth = "1990-12-11"


	// new函数实例化
	p2 := new(Person)
	p2.Name = "王二小"
	p2.Birth = "1990-12-22"


	// 取址实例化
	p3 := &Person{}
	p3.Name = "王三小"
	p3.Birth = "1990-12-23"

	// 初始化
	p4 := Person{
		Name:"王小四",
		Birth: "1990-12-23",
	}

	// 初始化
	p5 := &Person{
		"王五",
		"1990-12-23",
		5,
	}





}
