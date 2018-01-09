package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"thaThrowdown/common/database/dgraph"
	"thaThrowdown/common/infrastructure"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// UserType is the type of the user that uses the system. At the moment
// it ca be either :
//    1. Client - which comes from client portal
//    2. Admin  - which comes from admin portal
type UserType int

// AuthHeader contains Authorization header
const AuthHeader = "Authorization"

// Client is a type of a user which comes from client portal
const Client UserType = 1

// Admin   is a type of a user which comes from admin portal
const Admin UserType = 2

type json map[string]interface{}

// User contains logged user details
type User struct {
	ID    dgraph.UID
	Email string
	Name  string
}

// UserValidatorFunc is a function that is passed to the Authenticate middleware.
// the function gets a set of MapClaims and then extracts the relevant information
// then validates the user accordingly (by either reading the data from the database
// and matching the data with the claims or some other form)
type UserValidatorFunc func(jwt.MapClaims) (User, error)

// Authenticate is a middleware that is used to authenticate the incoming request
func Authenticate(validator UserValidatorFunc, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var token string
		header := c.Request().Header[AuthHeader]
		if len(header) > 0 {
			token = header[0]
		} else {
			cookie, err := c.Cookie(infrastructure.COOKIE_NAME)
			if err != nil {
				log.Println("middleware.AuthenticateApi - error getting the cookie : ", err)
				return c.JSON(http.StatusUnauthorized, json{"message": "User is not authenticated !"})
			}
			token = cookie.Value
		}

		claims, tokenErr := infrastructure.DecodeAuthToken(token)
		if tokenErr == nil {
			usr, err := validator(claims)
			if err != nil {
				log.Println("middleware.AuthenticateApi - Error validating user : ", err)
				return c.JSON(http.StatusUnauthorized, json{"message": "Could not Validate user"})
			}
			c.Set("user", usr)
			return next(c)
		}

		log.Println("middleware.AuthenticateApi - Error validating token : ", tokenErr)
		return c.JSON(http.StatusUnauthorized, json{"message": "Invalid token"})
	}
}

// NewUser return a User populated from the claims
func NewUser(claims jwt.MapClaims) (User, error) {

	id, success := claims["userId"].(float64)
	if !success {
		return User{}, errors.New("invalid userId claims")
	}

	name, success := claims["name"].(string)
	if !success {
		return User{}, errors.New("invalid name claims")
	}

	email, success := claims["email"].(string)
	if !success {
		return User{}, errors.New("invalid email claims")
	}

	usr := User{
		ID:    dgraph.UID(id),
		Name:  name,
		Email: email,
	}

	fmt.Println("User : ", usr)
	return usr, nil
}

// UserFromToken return a user from token
func UserFromToken(token string) (User, error) {
	claims, err := infrastructure.DecodeAuthToken(token)
	if err == nil {
		usr, err := NewUser(claims)
		return usr, err
	}
	return User{}, err
}

// ToClaims converts the User to jwt Claims
func (u User) ToClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"userId": u.ID,
		"email":  u.Email,
		"name":   u.Name,
	}
}

// ToToken converts  user to token
func (u User) ToToken() (infrastructure.JwtToken, error) {
	claims := u.ToClaims()
	token, err := infrastructure.GetAuthToken(claims)
	if err != nil {
		return infrastructure.JwtToken(""), err
	}

	return infrastructure.JwtToken(token), nil
}

// ToS3Key build the amazon S3 bucket key
func (u User) ToS3Key(mediaType infrastructure.MediaType, filename string) string {

	path := fmt.Sprintf("%s/%s__%s", mediaType, strings.Replace(u.Email, "@", "_at_", -1), filename)
	log.Println("Amazon S3 Path : ", path)
	return path
}
