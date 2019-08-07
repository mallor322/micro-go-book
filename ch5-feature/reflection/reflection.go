package main

import (
	"fmt"
	"reflect"
)

// 定义一个人的接口
type Person interface {

	// 和人说hello
	SayHello(name string)
	// 跑步
	Run() string
}

type Hero struct {
	Name string
	Age int
	Speed int
}
func (hero *Hero) SayHello(name string)  {
	fmt.Println("Hello " + name, ", I am " + hero.Name)
}

func (hero *Hero) Run() string {
	fmt.Println("I am running at speed " + string(hero.Speed))
	return ""
}

func main()  {


	//typeOfHero := reflect.TypeOf(Hero{})
	////fmt.Printf("Hero's type is %s, kind is %s\n", typeOfHero, typeOfHero.Kind())
	////fmt.Printf("*Hero's type is %s, kind is %s",reflect.TypeOf(&Hero{}).Elem().Field(0), reflect.TypeOf(&Hero{}).Kind())
	//
	//// 通过 #NumField 获取结构体字段的数量
	//for i := 0 ; i < typeOfHero.NumField(); i++{
	//	fmt.Printf("field'name is %s, type is %s, kind is %s\n", typeOfHero.Field(i).Name, typeOfHero.Field(i).Type, typeOfHero.Field(i).Type.Kind())
	//}
	//
	//
	//for i := 0 ; i < typeOfHero.NumMethod(); i++{
	//	fmt.Printf("method is %s, type is %s, kind is %s\n", typeOfHero.Method(i).Name, typeOfHero.Method(i).Type, typeOfHero.Method(i).Type.Kind())
	//}

	var person Person = &Hero{}
	// 获取接口Person的类型对象
	typeOfPerson := reflect.TypeOf(person)
	// 打印Person的方法类型和名称
	for i := 0 ; i < typeOfPerson.NumMethod(); i++{
		fmt.Printf("method is %s, type is %s, kind is %s\n", typeOfPerson.Method(i).Name, typeOfPerson.Method(i).Type, typeOfPerson.Method(i).Type.Kind())
	}


}



