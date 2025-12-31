package models

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	Model
	Name     string `json:"name"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Avator   string `json:"avator"`
	RoleName string `json:"role_name" sql:"-"`
	RoleId   string `json:"role_id" sql:"-"`
	Role     string `json:"role"`
}

func CreateUser(name string, password string, avator string, nickname string, role string) uint {
	user := &User{
		Name:     name,
		Password: password,
		Avator:   avator,
		Nickname: nickname,
		Role:     role,
	}
	user.UpdatedAt = time.Now()
	DB.Create(user)
	return user.ID
}
func UpdateUser(id string, name string, password string, avator string, nickname string, role string) {
	user := &User{
		Avator:   avator,
		Nickname: nickname,
		Role:     role,
		Name:     name,
	}
	user.UpdatedAt = time.Now()
	if password != "" {
		user.Password = password
	}
	DB.Model(&User{}).Where("id = ?", id).Update(user)
}
func UpdateUserPass(name string, pass string) {
	user := &User{
		Password: pass,
	}
	user.UpdatedAt = time.Now()
	DB.Model(user).Where("name = ?", name).Update("Password", pass)
}
func UpdateUserAvator(name string, avator string) {
	user := &User{
		Avator: avator,
	}
	user.UpdatedAt = time.Now()
	DB.Model(user).Where("name = ?", name).Update("Avator", avator)
}
func FindUser(username string) User {
	var user User
	DB.Where("name = ?", username).First(&user)
	return user
}
func FindUserById(id interface{}) User {
	var user User
	DB.Select("user.*").Where("user.id = ?", id).First(&user)
	return user
}
func DeleteUserById(id string) {
	DB.Where("id = ?", id).Delete(User{})
}
func FindUsers() []User {
	var users []User
	DB.Select("user.*").Limit(50).Order("user.id desc").Find(&users)
	return users
}
func FindUserRole(query interface{}, id interface{}) User {
	var user User
	DB.Select(query).Where("user.id = ?", id).Joins("join user_role on user.id=user_role.user_id").Joins("join role on user_role.role_id=role.id").First(&user)
	return user
}
