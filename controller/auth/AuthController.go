package auth

import (
	"time"
	"minitwitter/database"
	"minitwitter/model"
	"minitwitter/helper"
	"fmt"
	"net/http"

	//"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func PostLogin(c echo.Context) error {
	var user model.User
	db := database.Connect()
	email := c.FormValue("email")
	password := c.FormValue("password")

	err := db.Where("email = ?", email).First(&user)

	if err != nil {
		fmt.Println(err.Error)
	}

	isSuccess := CheckPasswordHash(password, user.Password)

	var result helper.ResponseData

	if isSuccess {
		// session, _ := session.Get("session-name", c)
		// session.Values["email"] = user.Email
		// session.Values["full_name"] = user.Fullname
		// session.Save(c.Request(), c.Response())
		// return c.Redirect(http.StatusSeeOther, "/")
		// result.Status = http.StatusOK
		// result.Message = "User berhasil login"
		// result.Data = user

		mySigningKey := []byte("AllYourBase")

		claims := model.JWTToken{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
				Issuer: user.Username,
			},
			UUID: user.UUID,

		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// fmt.Println(token)
		SignedToken, err := token.SignedString(mySigningKey)
		if err != nil {
			fmt.Println(SignedToken)
		}
		result.Status = http.StatusOK
		result.Message = "Success"
		result.Data = string(SignedToken)

		//return c.Redirect(http.StatusSeeOther, "/")		

	} else {
		result.Status = http.StatusOK
		result.Message = "email/password salah"
		result.Data = nil
	}

	return c.JSON(http.StatusOK, result)
}

func PostRegisterUser(c echo.Context) error {
	email := c.FormValue("email")
	fullName := c.FormValue("fullname")
	password := c.FormValue("password")
	username := c.FormValue("username")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	user := model.User{
		Fullname: fullName,
		Email: email,
		Password: string(hashedPassword),
		Username: username,
		UUID: GenerateUUID(),
	}

	db := database.Connect()
	errCreate := db.Create(&user)
	// fmt.Println(errCreate)
	if errCreate.Error != nil {
		return c.JSON(http.StatusOK, "Username sudah terpakai")
	}
	var result helper.ResponseData
	result.Status = http.StatusOK
	result.Message = "User berhasil ditambahkan"
	result.Data = nil

	return c.JSON(http.StatusOK, result)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	return err == nil
}

func GenerateUUID() string {
	uuidWithHyphen := uuid.New().String()
    //uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return uuidWithHyphen
}