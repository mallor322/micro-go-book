package main

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"time"
)

type TokenService struct {

	tokenStore TokenStore
	tokenEnhancer TokenEnhancer

}

func (tokenService *TokenService) CreateAccessToken(oauth2Details *OAuth2Details) (*OAuth2Token, error) {

	existToken, err := tokenService.tokenStore.GetAccessToken(oauth2Details)
	var refreshToken *OAuth2Token
	if err == nil{
		// 存在未失效访问令牌，直接返回
		if !existToken.IsExpired(){
			tokenService.tokenStore.StoreAccessToken(existToken, oauth2Details)
			return existToken, nil

		}
		// 访问令牌已失效，移除
		tokenService.tokenStore.RemoveAccessToken(existToken.TokenValue)
		if existToken.RefreshToken != nil {
			refreshToken = existToken.RefreshToken
			tokenService.tokenStore.RemoveRefreshToken(refreshToken.TokenType)
		}
	}

	if refreshToken == nil || refreshToken.IsExpired(){
		refreshToken, err = tokenService.createRefreshToken(oauth2Details)
		if err != nil{
			return nil, err
		}
	}

	accessToken, err := tokenService.createAccessToken(refreshToken, oauth2Details)
	if err == nil{
		tokenService.tokenStore.StoreAccessToken(accessToken, oauth2Details)
		tokenService.tokenStore.StoreRefreshToken(refreshToken, oauth2Details)
	}
	return accessToken, err

}

func (tokenService *TokenService) createAccessToken(refreshToken *OAuth2Token, oauth2Details *OAuth2Details) (*OAuth2Token, error) {

	validitySeconds := oauth2Details.Client.AccessTokenValiditySeconds
	s, _ := time.ParseDuration(strconv.Itoa(validitySeconds) + "s")
	expiredTime := time.Now().Add(s)
	fmt.Println(expiredTime.Unix())
	fmt.Println(time.Now().Unix())
	accessToken := &OAuth2Token{
		RefreshToken:refreshToken,
		ExpiresTime:&expiredTime,
		TokenValue:uuid.NewV4().String(),
	}

	if tokenService.tokenEnhancer != nil{
		return tokenService.tokenEnhancer.Enhance(accessToken, oauth2Details)
	}
	return accessToken, nil
}

func (tokenService *TokenService) createRefreshToken(oauth2Details *OAuth2Details) (*OAuth2Token, error) {
	validitySeconds := oauth2Details.Client.RefreshTokenValiditySeconds
	s, _ := time.ParseDuration(strconv.Itoa(validitySeconds) + "s")
	expiredTime := time.Now().Add(s)
	refreshToken := &OAuth2Token{
		ExpiresTime:&expiredTime,
		TokenValue:uuid.NewV4().String(),
	}

	if tokenService.tokenEnhancer != nil{
		return tokenService.tokenEnhancer.Enhance(refreshToken, oauth2Details)
	}
	return refreshToken, nil
}

func (tokenService *TokenService) RefreshAccessToken(refreshTokenValue string) (*OAuth2Token, error){

	refreshToken, err := tokenService.tokenStore.ReadRefreshToken(refreshTokenValue)

	if err == nil{
		if refreshToken.IsExpired(){
			return nil, errors.New("Refresh Token " + refreshTokenValue + " is Expired")
		}
		oauth2Details, err := tokenService.tokenStore.ReadOAuth2DetailsForRefreshToken(refreshTokenValue)
		if err == nil{
			oauth2Token, err := tokenService.tokenStore.GetAccessToken(oauth2Details)
			// 移除原有的访问令牌
			if err == nil{
				tokenService.tokenStore.RemoveAccessToken(oauth2Token.TokenValue)
			}

			// 移除已使用的刷新令牌
			tokenService.tokenStore.RemoveRefreshToken(refreshTokenValue)
			refreshToken, err = tokenService.createRefreshToken(oauth2Details)
			if err == nil{
				accessToken, err := tokenService.createAccessToken(refreshToken, oauth2Details)
				if err == nil{
					tokenService.tokenStore.StoreAccessToken(accessToken, oauth2Details)
					tokenService.tokenStore.StoreRefreshToken(refreshToken, oauth2Details)
				}
				return accessToken, err;
			}
		}
	}
	return nil, err

}

func (tokenService *TokenService) GetAccessToken(details *OAuth2Details) (*OAuth2Token, error)  {
	return tokenService.tokenStore.GetAccessToken(details)
}

func (tokenService *TokenService) ReadAccessToken(tokenValue string) (*OAuth2Token, error){
	return tokenService.tokenStore.ReadAccessToken(tokenValue)
}

func (tokenService *TokenService) GetOAuth2DetailsByAccessToken(tokenValue string) (*OAuth2Details, error) {

	accessToken, err := tokenService.tokenStore.ReadAccessToken(tokenValue)
	if err == nil{
		if accessToken.IsExpired(){
			return nil, errors.New("Access Token " + tokenValue + " is Expired")
		}
		return tokenService.tokenStore.ReadOAuth2Details(tokenValue)
	}
	return nil, err
}




