package model

import (
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UUID   		string 		`gorm:"primaryKey"`
	Fullname 	string 		`json:"fullname`
	Email 		string 		`json:"email"`
	Password 	string 		`json:"password"`
	Username 	string 		`gorm:"unique",json:"username"`
	PhotoPath   string		`json:"photo_path"`
}

type JWTToken struct {
	jwt.StandardClaims
	UUID   		string 		`json:"user_id"`
}