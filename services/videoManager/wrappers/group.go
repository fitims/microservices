package wrappers

import "github.com/labstack/echo"

// GroupWrapper is a wraper for Group struct in Echo
type GroupWrapper interface {
	Use(middleware ...echo.MiddlewareFunc)
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}
