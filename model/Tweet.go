package model

import (
	"github.com/jinzhu/gorm"
)

type Tweet struct {
	gorm.Model
	Content 	string 		`json:"content`
	UserID   	string 		`gorm:"primaryKey",json:"id"`
	Likes  		int32  		`gorm:"default:0",json:"id"`
}