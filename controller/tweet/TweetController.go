package tweet

import (
	"minitwitter/database"
	"minitwitter/model"
	"net/http"
	"github.com/labstack/echo/v4"
	"minitwitter/helper"
)

func LikeTweet(c echo.Context) error {
	var tweet model.Tweet

	tweetId := c.Param("id")

	db := database.Connect()
	db.Where("id = ?", tweetId).First(&tweet)

	tweet.Likes = int32(tweet.Likes) + 1
	db.Save(&tweet)

	var result helper.ResponseData
	result.Status = http.StatusOK
	result.Message = "success"
	result.Data = tweet

	return c.JSON(http.StatusOK, result)
}