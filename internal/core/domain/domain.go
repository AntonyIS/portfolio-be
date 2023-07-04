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
	Id             string           `json:"id"`
	FirstName      string           `json:"firstname"`
	LastName       string           `json:"lastname"`
	Email          string           `json:"email"`
	Title          string           `json:"title"`
	Password       string           `json:"password"`
	Projects       []*Project       `json:"projects"`
	Certifications []*Certification `json:"certification"`
}

type Certification struct {
	Id             string `json:"id"`
	UserID         string `json:"user_id"`
	Title          string `json:"title"`
	Institution    string `json:"institution"`
	State          string `json:"state"`
	IssuedDate     string `json:"issued_date"`
	CredentialLink string `json:"credential_link"`
	Decription     string `json:"decription"`
}
type Project struct {
	Id        string `json:"id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	UserName  string `json:"user_name"`
	UserTitle string `json:"user_title"`
	Rate      int    `json:"rate"`
	CreateAt  int64  `json:"created_at"`
}

func (u User) CheckPasswordHarsh(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}


