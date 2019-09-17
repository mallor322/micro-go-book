package main

import "errors"

type UserDetails struct {

	UserId int
	Username string
	Password string

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
