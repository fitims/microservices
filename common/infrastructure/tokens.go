package infrastructure

import (
	"time"

	"github.com/satori/go.uuid"

	"errors"
	"thaThrowdown/common/database/dgraph"

	"github.com/dgrijalva/jwt-go"
)

const secret = "世thathrowdowneshtesenimireshume界"

// JwtToken encapsulated JWT Token
type JwtToken string

var (
	InvalidTokenError = errors.New("Invalid Token")
	TokenExpiredError = errors.New("Token has expired")
	ClaimsError       = errors.New("Could not decode claims")
	UnknownTokenError = errors.New("Uknown error")
)

// GetAuthToken returns a valid authentication token for the claims provided
func GetAuthToken(claims jwt.MapClaims) (string, error) {

	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // expires in 24 hours

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

// GetRegistrationToken returns registration token for the new user
func GetRegistrationToken(userID dgraph.UID, email string, name string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"email":  email,
		"name":   name,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Hour * 48).Unix(), // expires in 48 hours
	})

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

// GetPasswordToken returns a valid token for user that forgot his password
func GetPasswordToken(userID dgraph.UID, email string) (string, string, error) {
	u1 := uuid.NewV4()
	guid := u1.String()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":  userID,
		"email":   email,
		"request": guid,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour * 48).Unix(), // expires in 48 hours
	})

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, guid, err
}

// GetUserAuthToken returns a valid authentication token for user
func GetUserAuthToken(userID dgraph.UID, email string, name string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"email":  email,
		"name":   name,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // expires in 24 hours
	})

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

// DecodeAuthToken decodes the provided authentication token
func DecodeAuthToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err == nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			return claims, nil
		}
		return nil, ClaimsError
	}

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, InvalidTokenError
		} else if ve.Errors&(jwt.ValidationErrorNotValidYet|jwt.ValidationErrorExpired) != 0 {
			return nil, TokenExpiredError
		}
	}

	return nil, UnknownTokenError
}
