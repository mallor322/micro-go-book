package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
)

// 基于图灵API一个简单的聊天机器人
type requestBody struct {

	Key string  `json:"key"`
	Info string  `json:"info"`
	UserId string  `json:"userid"`
}

type responseBody struct {

	Code int `json:"code"`
	Text string	`json:"text"`
	List []string	`json:"list"`
	Url string	`json:"url"`
	
}

func process(inputChan <-chan string, userid string)  {

	for{

		input := <- inputChan
		if input == "EOF"{
			break
		}

		reqData := &requestBody{
			Key: "792bcf45156d488c92e9d11da494b085",
			Info : input,
			UserId: userid,
		}

		byteData,_ := json.Marshal(&reqData)

		req, err := http.NewRequest("POST",
			"http://www.tuling123.com/openapi/api",
			bytes.NewReader(byteData))

		req.Header.Set("Content-Type", "application/json;charset=UTF-8")

		client := http.Client{}

		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("Network Error!")
		}else {
			body, _ := ioutil.ReadAll(resp.Body)
			var respData responseBody
			json.Unmarshal(body, &respData)
			fmt.Println("AI: " + respData.Text)

		}
		resp.Body.Close()

	}

}

func main()  {

	var input string

	fmt.Println("Enter 'EOF' to shut down: ")

	channel := make(chan string)

	go process(channel, string(rand.Int63()))

	 for {

	 	fmt.Scanf("%s", &input)
	 	channel <- input
		 if input == "EOF"{
			 fmt.Println("Bye!")
			 break
		 }

	 }



}