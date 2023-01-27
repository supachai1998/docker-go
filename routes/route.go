package routes

import (
	"docker_test/src/config"
	"docker_test/src/controllers"
	controllerUsers "docker_test/src/controllers/users"
	"docker_test/src/helpers"
	"docker_test/structs"
	structsDB "docker_test/structs/db"
	"net/http"

	_ "docker_test/docs"

	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

func RoutePath(r echoswagger.ApiRoot) {
	// instantiate controller
	controller := controllers.Controller{}
	controllerUsers := controllerUsers.NewUsers()

	r.GET("/", func(c echo.Context) error {
		return helpers.CustomHTTPError(c, echo.ErrForbidden)
	}).
		AddResponse(http.StatusForbidden, "Forbidden", nil, nil)

	r.GET("/rate-limited", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"max per minute": config.RATELIMITPERMINUTE,
		})
	}).AddResponse(http.StatusOK, "Success", echo.Map{}, nil)
	// group v1
	v1 := r.Group("base", "/api/v1/")
	v1.SetDescription("Base API on /api/v1")

	// ------- api user -------
	v1.POST("user/login", controller.GenerateToken).
		AddParamFormNested(structs.Login{}).
		AddResponse(http.StatusOK, "Success", echo.Map{}, nil).
		AddResponse(http.StatusBadRequest, "Bad Request", nil, nil).
		AddResponse(http.StatusInternalServerError, "Internal Server Error", nil, nil).
		AddResponse(http.StatusUnauthorized, "Unauthorized", nil, nil).
		AddResponse(http.StatusForbidden, "Forbidden", nil, nil)

	v1.POST("users", controllerUsers.CreateUser).
		AddParamFormNested(structs.Login{}).AddResponse(http.StatusOK, "Success", structsDB.Users{}, nil).
		AddResponse(http.StatusBadRequest, "Bad Request", nil, nil).
		AddResponse(http.StatusInternalServerError, "Internal Server Error", nil, nil)
	// ------- end api user -------

	setupUsers(r, controllerUsers)

}
