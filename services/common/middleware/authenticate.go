package middleware

import (
	"thaThrowdown/common/middleware"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//AuthenticateUser authenticates the user
func AuthenticateUser(next echo.HandlerFunc) echo.HandlerFunc {
	return middleware.Authenticate(validateUser, next)
}

func validateUser(claims jwt.MapClaims) (middleware.User, error) {
	return middleware.NewUser(claims)
}
