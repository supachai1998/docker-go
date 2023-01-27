package controllers

import (
	"docker_test/src/db"
	"docker_test/src/helpers"
	"docker_test/src/models"
	"docker_test/structs"
	structsDB "docker_test/structs/db"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (*Controller) GenerateToken(c echo.Context) error {
	token := helpers.GetTokenFormHeader(c)

	user := new(structsDB.UsersPut)
	var login structs.Login
	c.Bind(&login)
	// setup content
	content := make(map[string]string)
	// loop in query params to content
	for key, _ := range c.QueryParams() {
		content[key] = c.QueryParam(key)
	}
	var count int64

	// find user in db
	if err := db.Con.Where("email = ? or username = ?", login.Username, login.Username).First(user).Limit(-1).Count(&count).Error; err != nil {
		err := fmt.Errorf("IP [" + c.RealIP() + "] try login >> " + login.Username + " & invalid username or email")
		return helpers.CustomHTTPError(c, err, helpers.OptionalCustomHTTPError{
			Code:    helpers.PointerInt(http.StatusUnauthorized),
			Message: helpers.PointerString("Invalid username or password"),
		})
	}
	if count > 1 {
		err := fmt.Errorf("IP [" + c.RealIP() + "] try login >> " + login.Username + " & duplicate username or email")
		return helpers.CustomHTTPError(c, err, helpers.OptionalCustomHTTPError{
			Code:    helpers.PointerInt(http.StatusUnauthorized),
			Message: helpers.PointerString("Invalid username or password"),
		})
	}

	// check password
	if !models.ComparePassBcrypt(login.Password, user.Password) {
		err := fmt.Errorf("IP [" + c.RealIP() + "] try login >> " + login.Username + " & invalid password")
		return helpers.CustomHTTPError(c, err, helpers.OptionalCustomHTTPError{
			Code:    helpers.PointerInt(http.StatusUnauthorized),
			Message: helpers.PointerString("Invalid username or password"),
		})
	}
	if user.UID == "" {
		err := errors.New("user not register")
		return helpers.CustomHTTPError(c, err, helpers.OptionalCustomHTTPError{
			Code:    helpers.PointerInt(http.StatusNotFound),
			Message: helpers.PointerString("กรุณาเพิ่มเพื่อน GHB Buddy ที่ <a href=\"https://line.me/R/ti/p/@ghbbuddy\" target=\"_blank\">https://line.me/R/ti/p/@ghbbuddy</a>"),
		})
	}
	// generate token
	token, _, err := helpers.GenTokenHS256(user.UID, user.Username)
	if err != nil {
		return helpers.CustomHTTPError(c, err)
	}
	// save token to header bearer
	c.Response().Header().Set("Authorization", "Bearer "+token)
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
