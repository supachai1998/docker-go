package Users

import (
	"docker_test/src/db"
	"docker_test/src/helpers"
	ModelUsers "docker_test/src/models/users"
	"docker_test/structs"
	structsDB "docker_test/structs/db"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var modelUsers = ModelUsers.NewUsers()

func (*Users) CreateUser(c echo.Context) error {
	// bind & validate
	postUser := new(structs.Login)
	c.Bind(postUser)
	if err := c.Validate(postUser); err != nil {
		return helpers.CustomHTTPError(c, err)
	}

	user := structsDB.Users{
		Username: postUser.Username,
		Password: postUser.Password,
	}

	// create user in db
	if err := modelUsers.CreateUser(&user); err != nil {
		return helpers.CustomHTTPError(c, err)
	}
	return c.NoContent(http.StatusCreated)
}

func (*Users) GetAllUsers(c echo.Context) error {
	var pagination structs.Pagination
	db := db.Con
	db, _ = pagination.InitPanigation(c, db, []string{
		"username", "email", "phone_number", "first_name", "last_name", "province", "postcode", "address", "uid",
	})
	// get all users
	var users []structsDB.Users
	if err := db.Debug().Preload("Role").Find(&users).Limit(-1).Count(&pagination.Total).Error; err != nil {
		return helpers.CustomHTTPError(c, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":       users,
		"pagination": pagination,
	})
}
func (*Users) PutUsers(c echo.Context) error {
	claims := c.Get("user").(*structs.CustomClaims)
	user := new(structsDB.POSTUser)
	user.UID = claims.UserID
	c.Bind(user)
	if err := c.Validate(user); err != nil {
		return helpers.CustomHTTPError(c, err)
	}
	RESUser := new(structsDB.RESUser)
	RESUser.Bind(user)
	// update user in db
	if err := modelUsers.UpdateUser(RESUser); err != nil {
		return helpers.CustomHTTPError(c, err)
	}
	return c.NoContent(http.StatusOK)
}

func (*Users) SoftDeleteUsers(c echo.Context) error {
	user := new(structsDB.UsersPut)
	modelUsers.BindParamIDOrUIDToUser(c, user)
	c.Bind(user)
	if err := c.Validate(user); err != nil {
		return helpers.CustomHTTPError(c, err)
	}
	if err := modelUsers.SoftDeleteUser(user); err != nil {
		return helpers.CustomHTTPError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (*Users) UnSoftDeleteUsers(c echo.Context) error {
	user := new(structsDB.UsersPut)
	modelUsers.BindParamIDOrUIDToUser(c, user)
	c.Bind(user)
	if err := c.Validate(user); err != nil {
		return helpers.CustomHTTPError(c, err)
	}
	if err := modelUsers.UnSoftDeleteUser(user); err != nil {
		return helpers.CustomHTTPError(c, err)
	}
	return c.JSON(http.StatusOK, user)
}

func (*Users) HardDeleteUsers(c echo.Context) error {
	user := new(structsDB.UsersPut)
	modelUsers.BindParamIDOrUIDToUser(c, user)
	c.Bind(user)
	if err := c.Validate(user); err != nil {
		return helpers.CustomHTTPError(c, err)
	}

	if err := modelUsers.HardDeleteUser(user); err != nil {
		return helpers.CustomHTTPError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (*Users) PutUsersRole(c echo.Context) (err error) {
	user := new(structsDB.UsersPut)
	modelUsers.BindParamIDOrUIDToUser(c, user)
	c.Bind(user)
	if err := c.Validate(user); err != nil {
		return helpers.CustomHTTPError(c, err)
	}
	roleName := c.FormValue("role")
	if len(roleName) == 0 {
		return helpers.CustomHTTPError(c, gorm.ErrRecordNotFound)
	}

	return c.NoContent(http.StatusOK)
}

func (*Users) GetUserByIDOrUID(c echo.Context) error {
	userQuery := new(structsDB.UsersPut)
	user, err := modelUsers.GetUserByIDOrUserID(c, userQuery)
	if err != nil {
		return helpers.CustomHTTPError(c, err)
	}
	return c.JSON(http.StatusOK, user)
}

func (*Users) GetMeInfo(c echo.Context) error {
	claims := c.Get("user").(*structs.CustomClaims)
	user, err := modelUsers.GetUserByUserID(claims.UserID)
	if err != nil {
		return helpers.CustomHTTPError(c, err)
	}
	return c.JSON(http.StatusOK, user)
}

func (*Users) PutMeInfo(c echo.Context) error {
	claims := c.Get("user").(*structs.CustomClaims)
	id := c.Get("id").(*uint)
	user := new(structsDB.POSTUser)
	c.Bind(user)
	if user.UID == "" {
		user.UID = claims.UserID
	}
	if err := c.Validate(user); err != nil {
		return helpers.CustomHTTPError(c, err)
	}

	RESUser := new(structsDB.RESUser)
	RESUser.Bind(user)
	if id != nil {
		RESUser.ID = *id
	}
	// update user in db
	if err := modelUsers.UpdateUser(RESUser); err != nil {
		return helpers.CustomHTTPError(c, err)
	}

	if claims.UserID != RESUser.UID {

		c, err := helpers.ResetTokenByUserIdUsername(helpers.PointerUint(RESUser.ID), RESUser.UID, helpers.ConcatUsername(RESUser.FirstName, RESUser.LastName), c)
		if err != nil {
			return helpers.CustomHTTPError(c, err)
		}
		return c.JSON(http.StatusOK, echo.Map{
			"message": "user updated",
			"data":    RESUser,
			"token":   c.Response().Header().Get("Authorization"),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "user updated",
		"data":    RESUser,
	})
}
