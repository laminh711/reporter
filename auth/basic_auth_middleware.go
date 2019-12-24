package auth

import "github.com/labstack/echo"

// Middleware temporary middleware for authorization for go-echo
func Middleware(username, password string, c echo.Context) (bool, error) {
	if username == "admin" && password == "admin" {
		return true, nil
	}
	return false, nil
}
