package user

import (
	"minitwitter/database"
	"minitwitter/model"
	"net/http"
	"github.com/labstack/echo/v4"
	"minitwitter/helper"
	"github.com/golang-jwt/jwt"
	"fmt"
	"strings"
	"golang.org/x/crypto/bcrypt"
	// "github.com/google/uuid"
	"io"
	"os"
)


func ViewProfile(c echo.Context) error {
	var user model.User
	var tweet []model.Tweet

	username := c.Param("username")
	db := database.Connect()

	db.Where("username = ?", username).First(&user)
	db.Where("user_id = ?", user.UUID).First(&tweet)

	var result helper.ResponseData
	result.Status = http.StatusOK
	result.Message = "success"
	result.Data = tweet

	return c.JSON(http.StatusOK, result)
}

func PostTweet(c echo.Context) error {
	authorizationHeader := c.Request().Header.Get("Authorization")
	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
	payload, _ := extractClaims(tokenString)

	content := c.FormValue("content")

	tweet := model.Tweet{
		Content : content,
		UserID: payload["user_id"].(string),
	}

	db := database.Connect()
	db.Create(&tweet)

	var result helper.ResponseData
	result.Status = http.StatusOK
	result.Message = "success"
	result.Data = tweet

	return c.JSON(http.StatusOK, result)
}

func UpdateProfile(c echo.Context) error {
	authorizationHeader := c.Request().Header.Get("Authorization")
	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
	payload, _ := extractClaims(tokenString)

	var user model.User

	fullname := c.FormValue("fullname")
	password := c.FormValue("password")
	photo, err := c.FormFile("photo")

	if err != nil {
		fmt.Printf("error 1 %s", err)
		return err
	}

	src, err := photo.Open()

	if err != nil {
		fmt.Printf("error 2 %s", err)
		return err
	}

	defer src.Close()

	dst, err := os.Create("assets/" + photo.Filename)

	if err != nil {
		fmt.Printf("error 3 %s", err)
		return err
	}

	defer dst.Close()

	_, err_copy := io.Copy(dst, src)

	if err_copy != nil {
		fmt.Printf("error 4 %s", err_copy)
		return err_copy
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	db := database.Connect()

	db.Where("uuid = ?", payload["user_id"].(string)).First(&user)

	user.Fullname = fullname
	user.Password = string(hashedPassword)
	user.PhotoPath = photo.Filename
	db.Save(&user)

	var result helper.ResponseData
	result.Status = http.StatusOK
	result.Message = "success"
	result.Data = nil

	return c.JSON(http.StatusOK, result)
}

func FollowUser(c echo.Context) error {
	authorizationHeader := c.Request().Header.Get("Authorization")
	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
	payload, _ := extractClaims(tokenString)

	var userFollowed model.User
	var relationship []model.Relationship

	username := c.Param("username")
	db := database.Connect()

	db.Where("username = ?", username).First(&userFollowed)
	db.Where("follower_id = ? AND followed_id = ?", payload["user_id"].(string), userFollowed.UUID).First(&relationship)

	if len(relationship) != 0 {
		return c.JSON(http.StatusOK, "anda sudah memfollow akun ini")
	}

	relationshipCreate := model.Relationship{
		FollowerID : payload["user_id"].(string),
		FollowedID: userFollowed.UUID,
	}

	db.Create(&relationshipCreate)

	var result helper.ResponseData
	result.Status = http.StatusOK
	result.Message = "success"
	result.Data = relationshipCreate

	return c.JSON(http.StatusOK, result)
}

func UnfollowUser(c echo.Context) error {
	authorizationHeader := c.Request().Header.Get("Authorization")
	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
	payload, _ := extractClaims(tokenString)

	var userFollowed model.User
	var relationship []model.Relationship

	username := c.Param("username")
	db := database.Connect()

	db.Where("username = ?", username).First(&userFollowed)
	check := db.Where("follower_id = ? AND followed_id = ?", payload["user_id"].(string), userFollowed.UUID).Delete(&relationship)

	if check.RowsAffected == 0 {
		return c.JSON(http.StatusOK, "anda belum follow akun ini")
	}

	var result helper.ResponseData
	result.Status = http.StatusOK
	result.Message = "success"
	result.Data = nil

	return c.JSON(http.StatusOK, result)
}

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := "AllYourBase"
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		 // check token signing method etc
		 return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		fmt.Printf("Invalid JWT Token")
		return nil, false
	}
}