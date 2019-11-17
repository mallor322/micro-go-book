package model

import (
	"github.com/gohouse/gorose/v2"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/mysql"
	"log"
)

type User struct {
	UserId     int    `json:"user_id"`     //Id
	UserName   string `json:"user_name"`   //用户名称
	Password   string `json:"password"`    //密码
	Age        int    `json:"age"`         //年龄
	CreateTime string `json:"create_time"` // 创建时间
	ModifyTime string `json:"modify_time"` //修改时间
}

type UserModel struct {
}

func NewUserModel() *UserModel {
	return &UserModel{}
}

func (p *UserModel) getTableName() string {
	return "user"
}

func (p *UserModel) GetUserList() ([]gorose.Data, error) {
	conn := mysql.DB()
	list, err := conn.Table(p.getTableName()).Get()
	if err != nil {
		log.Printf("Error : %v", err)
		return nil, err
	}
	return list, nil
}

/*func (p *UserModel) GetUserByUsername(username string) (*User,  error)  {

	conn := mysql.DB()
	if result, err := conn.Table(p.getTableName()).Where(map[string]interface{}{"username": username}).First(); err == nil{

	}else {
		return nil, err
	}

}*/

func (p *UserModel) CheckUser(username string, password string) (bool, error) {
	conn := mysql.DB()
	num, err := conn.Table(p.getTableName()).Where(map[string]interface{}{"username": username, "password": password}).Count("*")
	if err != nil {
		log.Printf("Error : %v", err)
		return false, err
	}
	return num > 0, nil
}

func (p *UserModel) CreateUser(user *User) error {
	conn := mysql.DB()
	_, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"user_id":     user.UserId,
		"user_name":   user.UserName,
		"password":    user.Password,
		"age":         user.Age,
		"create_time": user.CreateTime,
		"modify_time": user.ModifyTime,
	}).Insert()
	if err != nil {
		log.Printf("Error : %v", err)
		return err
	}
	return nil
}
