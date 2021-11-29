package middlewarecustom

import(
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func (c echo.Context) error {
		session, _ := session.Get("session-name",c)
		if session.Values["full_name"] == nil {
			return c.Redirect(http.StatusSeeOther, "/")
		}
		return next(c)
	}
}

func JWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func (c echo.Context) error {
		mySigningKey := []byte("AllYourBase")
		authorizationHeader := c.Request().Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			c.JSON(http.StatusUnauthorized, "Invalid token")
		}
		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Signing method invalid")
			} else if method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("Signing method invalid")
			}
			return mySigningKey, nil
		})
		_, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			return c.JSON(http.StatusUnauthorized, "Invalid JWT Token")
		} else {
			return next(c)
		}
	}
}