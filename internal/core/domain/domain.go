/*
Package name : domain
File name : domain.go
Author : Antony Injila
Description :
	- Host Portfolio entiry strunctures such as a User and a Project
	- User types have the GenerateHashPassord and CheckPasswordHarsh methods
*/
package domain

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        string              `json:"id"`
	FirstName string              `json:"firstname"`
	LastName  string              `json:"lastname"`
	Email     string              `json:"email"`
	Password  string              `json:"password"`
	Projects  map[string]*Project `json:"projects"`
}

type Project struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	UserID   string `json:"user_id"`
	Rate     int    `json:"rate"`
	CreateAt int64  `json:"created_at"`
}



func (u User) CheckPasswordHarsh(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}
