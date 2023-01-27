package models

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetParamFormContextInt(c echo.Context, name string) (int, error) {
	val, err := strconv.Atoi(c.Param(name))
	if err == nil {
		return val, nil
	}
	return 0, echo.ErrNotFound
}
func GetParamFormContextStr(c echo.Context, name string) (string, error) {
	val := c.Param(name)
	if len(val) > 0 {
		return val, nil
	}
	return "", echo.ErrNotFound
}

func ComparePassBcrypt(pass string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
}
