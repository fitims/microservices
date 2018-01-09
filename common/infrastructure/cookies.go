package infrastructure

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"thaThrowdown/common/database/dgraph"
)

// COOKIE_NAME contains the name of the access token cookie
const COOKIE_NAME = "x-access-token"

// BuildAuthCookie build and returns an authorization cookie that contains userId and email
func BuildAuthCookie(claims jwt.MapClaims) (*http.Cookie, error) {
	token, err := GetAuthToken(claims)
	if err != nil {
		return nil, err
	}

	cookie := new(http.Cookie)
	cookie.Name = COOKIE_NAME
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * 7 * time.Hour) // expires in 1 days
	return cookie, nil
}

// BuildUserAuthCookie build and returns an authorization cookie that contains userId and email
func BuildUserAuthCookie(id dgraph.UID, email string, name string) (*http.Cookie, error) {
	token, err := GetUserAuthToken(id, email, name)
	if err != nil {
		return nil, err
	}

	cookie := new(http.Cookie)
	cookie.Name = COOKIE_NAME
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * 7 * time.Hour) // expires in 1 days
	return cookie, nil
}

// GetBlankAuthCookie build and returns a blank authorization cookie. It is used to nullfiy an existing cookie
func GetBlankAuthCookie() *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = COOKIE_NAME
	cookie.Value = "deleted"
	cookie.Expires = time.Unix(0, 0)
	return cookie
}
