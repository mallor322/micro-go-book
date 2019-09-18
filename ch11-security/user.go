package main

import "errors"

type UserDetails struct {
	// 用户标识
	UserId int
	// 用户名 唯一
	Username string
	// 用户密码
	Password string
	// 用户具有的权限
	Authorities []string // 具备的权限
}

func (userDetails *UserDetails)IsMatch(username string, password string) bool {
	return userDetails.Password == password && userDetails.Username == username
}

type UserDetailsService interface {
	GetUserDetailByUsername(username string)(*UserDetails, error)
}


type InMemoryUserDetailService struct {
	userDetailsDict map[string]*UserDetails

}

func (userDetailsService *InMemoryUserDetailService)GetUserDetailByUsername(username string) (*UserDetails, error) {

	userDetails := userDetailsService.userDetailsDict[username]

	if userDetails == nil{
		return nil, errors.New("Username " + username + " is not exist")
	}
	return userDetails, nil
}

func NewInMemoryUserDetailService(userDetailsList []*UserDetails ) *InMemoryUserDetailService{
	userDetailsDict := make(map[string]*UserDetails)

	if userDetailsList != nil {
		for _, value := range userDetailsList {
			userDetailsDict[value.Username] = value
		}
	}

	return &InMemoryUserDetailService{
		userDetailsDict:userDetailsDict,
	}
}
