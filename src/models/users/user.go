package users

import (
	"docker_test/src/db"
	"docker_test/src/helpers"
	"docker_test/src/models"
	structs "docker_test/structs/db"
	structsDB "docker_test/structs/db"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
)

var Users = new(users)

func (*users) CreateUser(user *structsDB.Users) error {
	return db.Con.Create(&user).Error
}
func (*users) HardDeleteUser(user *structsDB.UsersPut) error {
	return db.Con.Unscoped().Delete(&user).Error
}
func (*users) SoftDeleteUser(user *structsDB.UsersPut) error {
	return db.Con.Delete(&user).Error
}

func (*users) UnSoftDeleteUser(user *structsDB.UsersPut) error {
	return db.Con.Model(&user).Unscoped().Clauses(clause.Returning{}).Where("uid = ? OR id = ?", user.UID, user.ID).Update("deleted_at", nil).Error
}

func (*users) UpdateUser(user *structsDB.RESUser) error {
	// encrypt password bcrypt
	if user.Password != "" {
		hashedPassword, err := helpers.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	return db.Con.Model(&user).Clauses(clause.Returning{}).Where("id = ?", user.ID).Updates(&user).Error
}
func (*users) GetUserByID(id int) (*structsDB.Users, error) {
	var user structsDB.Users
	err := db.Con.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (*users) GetAllUsers() ([]structsDB.Users, error) {
	var users []structsDB.Users
	err := db.Con.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (*users) BindParamIDOrUIDToUser(c echo.Context, user *structs.UsersPut) error {
	// get id or uid from context
	id, err := models.GetParamFormContextInt(c, "id")
	if err != nil {
		uid, err := models.GetParamFormContextStr(c, "id")
		if err != nil {
			return helpers.CustomHTTPError(c, echo.ErrNotFound)
		}
		user.UID = uid
		return nil
	}
	user.ID = helpers.PointerUint(uint(id))
	return nil
}
func (*users) GetUserByIDOrUserID(c echo.Context, userQuery *structsDB.UsersPut) (user *structsDB.Users, err error) {
	if err := Users.BindParamIDOrUIDToUser(c, userQuery); err != nil {
		return user, err
	}
	// get user from db
	if err := db.Con.Preload(clause.Associations).First(&user, "id = ? or uid = ?", userQuery.ID, userQuery.UID).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (u *users) GetUserByUserID(userID string) (*structsDB.Users, error) {
	var user structsDB.Users
	err := db.Con.Preload(clause.Associations).First(&user, "uid = ?", userID).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
