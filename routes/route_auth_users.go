package routes

import (
	controllerUsers "docker_test/src/controllers/users"
	"docker_test/src/helpers"
	structsDB "docker_test/structs/db"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

type roleSwagger struct {
	Role string `json:"role" swagger:"enum(user|admin)" query:"role" form:"role"`
}

func setupUsers(r echoswagger.ApiRoot, controllerUsers *controllerUsers.Users) {
	AuthUsers := r.Group("Auth Users", "/api/v1/me/")
	AuthUsers.SetDescription("Protected API on /api/v1")
	// decode token
	AuthUsers.EchoGroup().Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// get token from header
			token := c.Request().Header.Get("Authorization")

			claims, err := helpers.DecodeTokenHS256(token)
			if err != nil {
				return helpers.CustomHTTPError(c, err)
			}

			// set user id to context
			c.Set("userId", claims.UserID)
			c.Set("id", claims.Id)
			c.Set("user", claims)
			return next(c)
		}
	})

	AuthUsers.GET("info", controllerUsers.GetMeInfo).
		AddParamHeader("Authorization", "Authorization", "Authorization", true).
		AddResponse(http.StatusOK, "Success", []structsDB.Users{}, nil).
		AddResponse(http.StatusBadRequest, "Bad Request", nil, nil).
		// conflict
		AddResponse(http.StatusConflict, "Conflict Duplicate Data", nil, nil).
		AddResponse(http.StatusInternalServerError, "Internal Server Error", nil, nil)

	AuthUsers.PUT("info", controllerUsers.PutMeInfo).
		AddParamHeader("Authorization", "Authorization", "Authorization", true).
		AddParamFormNested(structsDB.POSTUserSwagger{}).
		AddResponse(http.StatusOK, "OK", structsDB.RESUser{}, nil).
		AddResponse(http.StatusBadRequest, "Bad Request", nil, nil).
		AddResponse(http.StatusInternalServerError, "Internal Server Error", nil, nil)

}
