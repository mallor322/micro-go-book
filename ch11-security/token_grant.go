package main

import (
	"errors"
	"net/http"
)

type TokenGranter interface {
	grant(grantType string, client *ClientDetails, reader *http.Request) (*OAuth2Token, error)
}



type ComposeTokenGranter struct {
	tokenGrantDict map[string] TokenGranter
}

func NewComposeTokenGranter(tokenGrantDict map[string] TokenGranter) *ComposeTokenGranter {
	return &ComposeTokenGranter{
		tokenGrantDict:tokenGrantDict,
	}
}

func (tokenGranter *ComposeTokenGranter) grant(grantType string, client *ClientDetails, reader *http.Request) (*OAuth2Token, error) {

	dispatchGranter := tokenGranter.tokenGrantDict[grantType]

	if dispatchGranter == nil{
		return nil, errors.New("Grant Type " + grantType + " is not supported")
	}

	return dispatchGranter.grant(grantType, client, reader)
}

type UsernamePasswordTokenGranter struct {
	supportGrantType string
	userDetailsService UserDetailsService
	tokenService *TokenService

}

func (tokenGranter *UsernamePasswordTokenGranter) grant(grantType string, client *ClientDetails, reader *http.Request) (*OAuth2Token, error) {
	if grantType != tokenGranter.supportGrantType{
		return nil, errors.New("Target Grant Type is " + grantType + ", but current grant type is " + tokenGranter.supportGrantType)
	}

	username := reader.FormValue("username")
	password := reader.FormValue("password")

	if username == "" || password == ""{
		return nil, errors.New( "Please provide correct user information")
	}

	userDetails, err := userDetailsService.GetUserDetailByUsername(username)

	if err != nil{
		return nil, errors.New( "Username "+ username +" is not exist")}

	if !userDetails.IsMatch(username, password){
		return nil, errors.New( "Username or password is not corrent")
	}

	return tokenGranter.tokenService.CreateAccessToken(&OAuth2Details{
		Client:client,
		User:userDetails,

	})

}






