package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type User struct {
	Id          int
	Name        string
	habits      string
	createdTime string
}

func main() {
	session, err := mgo.Dial("mongodb://root:example@47.96.140.41:27017/user")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")
	err = c.Insert(&User{9, "Ale", "running", "2019-4-09"},
		&User{10, "Cla", "hiking", "2019-4-09"})
	if err != nil {
		log.Fatal(err)
	}

	result := User{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Name:", result.Name)
}
