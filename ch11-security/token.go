package main

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type OAuth2Details struct {
	Client *ClientDetails
	User *UserDetails

}


type OAuth2Token struct {

	// 刷新令牌
	RefreshToken *OAuth2Token
	// 令牌类型
	TokenType string
	// 令牌
	TokenValue string
	// 过期时间
	ExpiresTime *time.Time

}

func (oauth2Token *OAuth2Token) IsExpired() bool  {
	return oauth2Token.ExpiresTime != nil &&
		oauth2Token.ExpiresTime.Before(time.Now())
}

type TokenEnhancer interface {

	Enhance(oauth2Token *OAuth2Token, oauth2Details *OAuth2Details) (*OAuth2Token, error)

	Extract(tokenValue string) (*OAuth2Token, *OAuth2Details, error)

}

type OAuth2TokenCustomClaims struct {
	UserDetails UserDetails
	ClientDetails ClientDetails
	RefreshToken OAuth2Token
	jwt.StandardClaims
}

type JwtTokenEnhancer struct {
	secretKey []byte
}

func (enhancer *JwtTokenEnhancer) Enhance(oauth2Token *OAuth2Token, oauth2Details *OAuth2Details) (*OAuth2Token, error) {
	return enhancer.sign(oauth2Token, oauth2Details)
}

func (enhancer *JwtTokenEnhancer) Extract(tokenValue string) (*OAuth2Token, *OAuth2Details, error)  {

	token, err := jwt.ParseWithClaims(tokenValue, &OAuth2TokenCustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return enhancer.secretKey, nil
	})

	if err == nil{

		claims := token.Claims.(*OAuth2TokenCustomClaims)
		expiresTime := time.Unix(claims.ExpiresAt, 0)

		return &OAuth2Token{
			RefreshToken:&claims.RefreshToken,
			TokenValue:tokenValue,
			ExpiresTime: &expiresTime,
		}, &OAuth2Details{
			User:&claims.UserDetails,
			Client:&claims.ClientDetails,
		}, nil

	}
	return nil, nil, err

}

func (enhancer *JwtTokenEnhancer) sign(oauth2Token *OAuth2Token, oauth2Details *OAuth2Details)  (*OAuth2Token, error) {

	expireTime := oauth2Token.ExpiresTime
	clientDetails := *oauth2Details.Client
	userDetails := *oauth2Details.User
	clientDetails.ClientSecret = ""
	userDetails.Password = ""

	claims := OAuth2TokenCustomClaims{
		UserDetails:userDetails,
		ClientDetails:clientDetails,
		StandardClaims:jwt.StandardClaims{
			ExpiresAt:expireTime.Unix(),
			Issuer:"System",
		},
	}

	if oauth2Token.RefreshToken != nil{
		claims.RefreshToken = *oauth2Token.RefreshToken
	}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenValue, err := token.SignedString(enhancer.secretKey)

	if err == nil{
		oauth2Token.TokenValue = tokenValue
		return oauth2Token, nil;

	}
	return nil, err
}


