package main

import (
	"errors"
	"net/http"
)

type TokenGranter interface {
	grant(grantType string, clientId string, reader *http.Request) (*OAuth2Token, error)
}


type ComposeTokenGranter struct {
	tokenGrantDict map[string] TokenGranter
}

func (tokenGranter *ComposeTokenGranter) grant(grantType string, clientId string, reader *http.Request) (*OAuth2Token, error) {

	dispatchGranter := tokenGranter.tokenGrantDict[grantType]

	if dispatchGranter == nil{
		return nil, errors.New("Grant Type " + grantType + " is not supported")
	}

	return dispatchGranter.grant(grantType, reader)
}

type UsernamePasswordTokenGranter struct {
	supportGrantType string
	userDetailsService UserDetailsService

}

func (tokenGranter *UsernamePasswordTokenGranter) grant(grantType string, clientId string, reader *http.Request) (*OAuth2Token, error) {
	if grantType != tokenGranter.supportGrantType{
		return nil, errors.New("Target Grant Type is " + grantType + ", but current grant type is " + tokenGranter.supportGrantType)
	}

	username := reader.Form.Get("username")
	password := reader.Form.Get("password")

	if username == "" || password == ""{
		return nil, errors.New( "Please provide correct user information")
	}

	userDetails, err := userDetailsService.GetUserDetailByUsername(username)

	if err != nil{
		return nil, errors.New( "Username "+ username +" is not exist")}

	if !userDetails.IsMatch(username, password){
		return nil, errors.New( "Username or password is not corrent")
	}




}






