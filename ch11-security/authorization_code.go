package main

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

type AuthorizationCodeServices interface {
	// 生成对应的授权码
	CreateAuthorizationCode(details *OAuth2Details) string
	// 消费授权码，返回对应的用户和客户端信息
	ConsumeAuthorizationCode(code string) (*OAuth2Details, error)
}

type InMemoryAuthorizationCodeServices struct {
	codeDict  *sync.Map
}


func NewInMemoryAuthorizationCodeServices() *InMemoryAuthorizationCodeServices{
	return &InMemoryAuthorizationCodeServices{
		codeDict:&sync.Map{},
	}
}

func (service *InMemoryAuthorizationCodeServices)CreateAuthorizationCode(details *OAuth2Details) string {
	code := service.getRandomString(6)
	service.codeDict.Store(code, details)
	return code
}

func (service *InMemoryAuthorizationCodeServices)ConsumeAuthorizationCode(code string) (*OAuth2Details, error)  {
	oauth2Token, ok := service.codeDict.Load(code)
	if ok{
		return oauth2Token.(*OAuth2Details), nil
	}
	return nil, errors.New("Code " + code + " is not exist")
}


func  (service *InMemoryAuthorizationCodeServices) getRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}