package main

import (
	"ch6-config/conf"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/profiles", func(w http.ResponseWriter, req *http.Request) {
		//q := req.URL.Query().Get("q")
		_, _ = fmt.Fprintf(w, "个人信息：\n")
		_, _ = fmt.Fprintf(w, "姓名：%s，\n性别：%s，\n年龄 %d!", conf.Profile.Name, conf.Profile.Sex, conf.Profile.Age) //这个写入到w的是输出到客户端的
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
}
