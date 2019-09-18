package main

import "github.com/pkg/errors"

type TokenStore interface {

	// 存储访问令牌
	StoreAccessToken(oauth2Token *OAuth2Token, oauth2Details *OAuth2Details)
	// 根据令牌值获取访问令牌结构体
	ReadAccessToken(tokenValue string) (*OAuth2Token, error)
	// 根据令牌值获取令牌对应的客户端和用户信息
	ReadOAuth2Details(tokenValue string)(*OAuth2Details, error)
	// 根据客户端信息和用户信息获取访问令牌
	GetAccessToken(oauth2Details *OAuth2Details)(*OAuth2Token, error);
	// 移除存储的访问令牌
	RemoveAccessToken(tokenValue string)
	// 存储刷新令牌
	StoreRefreshToken(oauth2Token *OAuth2Token, oauth2Details *OAuth2Details)
	// 移除存储的刷新令牌
	RemoveRefreshToken(oauth2Token string)
	// 根据令牌值获取刷新令牌
	ReadRefreshToken(tokenValue string)(*OAuth2Token, error)
	// 根据令牌值获取刷新令牌对应的客户端和用户信息
	ReadOAuth2DetailsForRefreshToken(tokenValue string)(*OAuth2Details, error)

}


type JwtTokenStore struct {
	jwtTokenEnhancer *JwtTokenEnhancer
}

func (tokenStore *JwtTokenStore)StoreAccessToken(oauth2Token *OAuth2Token, oauth2Details *OAuth2Details){

}

func (tokenStore *JwtTokenStore)ReadAccessToken(tokenValue string) (*OAuth2Token, error){
	oauth2Token, _, err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return oauth2Token, err

}
// 根据令牌值获取令牌对应的客户端和用户信息
func  (tokenStore *JwtTokenStore) ReadOAuth2Details(tokenValue string)(*OAuth2Details, error){
	_, oauth2Details, err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return oauth2Details, err

}
// 根据客户端信息和用户信息获取访问令牌
func (tokenStore *JwtTokenStore) GetAccessToken(oauth2Details *OAuth2Details)(*OAuth2Token, error){
	return nil, errors.New("JwtTokenStore not support")
}
// 移除存储的访问令牌
func (tokenStore *JwtTokenStore) RemoveAccessToken(tokenValue string){

}
// 存储刷新令牌
func (tokenStore *JwtTokenStore) StoreRefreshToken(oauth2Token *OAuth2Token, oauth2Details *OAuth2Details){

}
// 移除存储的刷新令牌
func (tokenStore *JwtTokenStore)RemoveRefreshToken(oauth2Token string){

}
// 根据令牌值获取刷新令牌
func (tokenStore *JwtTokenStore) ReadRefreshToken(tokenValue string)(*OAuth2Token, error){
	oauth2Token, _, err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return oauth2Token, err
}
// 根据令牌值获取刷新令牌对应的客户端和用户信息
func (tokenStore *JwtTokenStore)ReadOAuth2DetailsForRefreshToken(tokenValue string)(*OAuth2Details, error){
	_, oauth2Details, err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return oauth2Details, err
}
