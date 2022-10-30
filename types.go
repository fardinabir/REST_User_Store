package main

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}

type UserTag struct {
	gorm.Model
	UserId uint
	User   User
	Tag    string `json:"tag"`
}

type UserTot struct {
	gorm.Model
	UserId     uint
	User       User
	UserTagID  uint
	UserTag    UserTag
	Name       string `json:"fullName"`
	ExpiryTime int64  `json:"expiryTime"`
}

type GetResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"fullName"`
	Phone string `json:"phone"`
}

type PostResponse struct {
	ID uint `json:"id"`
}

type TagPostReq struct {
	Tags   []string `json:"tags"`
	Expiry int64    `json:"expiry"`
}
