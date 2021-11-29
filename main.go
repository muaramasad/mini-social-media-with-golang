package main

import (
	"minitwitter/model"
	"minitwitter/controller/auth"
	"minitwitter/controller/user"
	"minitwitter/controller/tweet"
	"minitwitter/database"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	// "github.com/golang-jwt/jwt"
)

func main() {

	db := database.Connect()
	db.AutoMigrate(&model.User{}, &model.Tweet{}, &model.Relationship{})

	e := echo.New()

	e.POST("/api/login", auth.PostLogin)
	e.POST("/api/register", auth.PostRegisterUser)

	api := e.Group("/api")
	api.Use(echoMiddleware.JWT([]byte("AllYourBase")))
	api.GET("/user/:username", user.ViewProfile)
	api.POST("/user/tweet", user.PostTweet)
	api.POST("/user/update", user.UpdateProfile)
	api.GET("/user/:username/follow", user.FollowUser)
	api.GET("/user/:username/unfollow", user.UnfollowUser)
	api.GET("/tweet/:id/like", tweet.LikeTweet)

	e.Logger.Fatal(e.Start(":3003"))
}