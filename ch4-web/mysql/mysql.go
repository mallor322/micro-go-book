package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql",
		"root:a123456@tcp(47.96.140.41:3366)/user?charset=utf8")
	checkErr(err)

	//插入数据
	stmt, err := db.Prepare("INSERT INTO user SET name=?,habits=?,created_time=?")
	checkErr(err)

	res, err := stmt.Exec("aoho", "balls", "2019-4-09")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Printf("last insert id is: %d\n",id)
	//更新数据
	stmt, err = db.Prepare("update user set  habits=? where id=?")
	checkErr(err)

	res, err = stmt.Exec("running,hiking", id)
	checkErr(err)

	//查询数据
	fmt.Println("\nafter inserting records: ")
	queryTableRecords(db)

	//删除数据
	stmt, err = db.Prepare("delete from user where id=?")
	checkErr(err)

	res, err = stmt.Exec(id)
	checkErr(err)

	//查询删除之后的数据
	fmt.Println("\nafter deleting records: ")
	queryTableRecords(db)
	_ = db.Close()

}

func queryTableRecords(db *sql.DB)  {
	rows, err := db.Query("SELECT * FROM user")
	checkErr(err)

	for rows.Next() {
		var id int
		var name string
		var habits string
		var createdTime string
		err = rows.Scan(&id, &name, &habits, &createdTime)
		checkErr(err)
		fmt.Printf("[%d, %s, %s, %s]\n", id, name, habits, createdTime)
	}
}


func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
