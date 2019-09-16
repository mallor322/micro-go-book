package main

import "net/http"

type TokenService struct {

	tokenStore TokenStore
}

func (tokenService *TokenService) CreateAccessToken(userDetails *UserDetails, clientDetail *ClientDetails) (*OAuth2Token, error) {
	
}


func (tokenService *TokenService) RefreshAccessToken(refreshTokenValue string) (*OAuth2Token, error){

}

func (tokenService *TokenService) GetAccessToken(userDetails *UserDetails)  {

}

func (tokenService *TokenService) ReadAccessToken(tokenValue string) (*OAuth2Token, error){
	
}

func GetUserDetailsByAccessToken(tokenValue string) (*UserDetails, error) {
	
}




