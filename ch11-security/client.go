package main

import "github.com/pkg/errors"

type ClientDetails struct {
	// client 的标识
	ClientId string
	// client 的密钥
	ClientSecret string
	// 访问令牌有效时间，秒
	AccessTokenValiditySeconds int
	// 刷新令牌有效时间，秒
	RefreshTokenValiditySeconds int
	// 重定向地址，授权码类型中使用
	RegisteredRedirectUri string
	// 可以使用的授权类型
	AuthorizedGrantTypes []string
}

func (clientDetails *ClientDetails) IsMatch(clientId string, clientSecret string) bool {
	return clientId == clientDetails.ClientId && clientSecret == clientDetails.ClientSecret
}

type ClientDetailService interface {
	GetClientDetailByClientId(clientId string) (*ClientDetails, error)
}

type InMemoryClientDetailService struct {
	clientDetailsDict map[string]*ClientDetails

}

func (clientDetailsService *InMemoryClientDetailService)GetClientDetailByClientId(clientId string) (*ClientDetails, error) {

	clientDetails := clientDetailsService.clientDetailsDict[clientId]

	if clientDetails == nil{
		return nil, errors.New("ClientId " + clientId + " is not exist")
	}
	return clientDetails, nil
}

func NewInMemoryClientDetailService(clientDetailsList []*ClientDetails ) *InMemoryClientDetailService{
	clientDetailsDict := make(map[string]*ClientDetails)

	if clientDetailsList != nil {
		for _, value := range clientDetailsList {
			clientDetailsDict[value.ClientId] = value
		}
	}

	return &InMemoryClientDetailService{
		clientDetailsDict:clientDetailsDict,
	}
}


