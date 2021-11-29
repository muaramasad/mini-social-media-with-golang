package model

import (
	"github.com/jinzhu/gorm"
	
)

type Relationship struct {
	gorm.Model
	FollowerID 	string 		`json:"follower_id`
	FollowedID 	string 		`json:"followed_id`
}