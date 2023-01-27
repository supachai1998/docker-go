package structs

import (
	"docker_test/src/middlewares"
	"docker_test/structs"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Users struct {
	structs.Model
	Username    string `json:"username" gorm:"uniqueIndex" form:"username" query:"username" swaggertype:"string" example:"username"`
	Email       string `json:"email" form:"email" query:"email" swaggertype:"string" example:"ex@hotmail.com"`
	Password    string `json:"-" swaggertype:"string" form:"password" query:"password" example:"password"`
	FirstName   string `json:"first_name"  form:"first_name" swaggertype:"string" example:"supachai"`
	LastName    string `json:"last_name"  form:"last_name" swaggertype:"string" example:"last name"`
	PhoneNumber string `json:"phone_number" validate:"max=10" form:"phone_number" swaggertype:"string" example:"0899999999"`
	Address     string `json:"address" gorm:"type:text"  form:"address" swaggertype:"string" example:"address"`
	Province    string `json:"province"  form:"province" swaggertype:"string" example:"bankok"`
	Postcode    string `json:"postcode"  form:"postcode" swaggertype:"string" example:"10260"`
	UID         string `json:"user_id" gorm:"uniqueIndex"  form:"user_id" query:"user_id" swaggertype:"string" example:"U1234567890"`
	AcceptTerm  bool   `json:"accept_term"  form:"accept_term" swaggertype:"boolean" example:"true"`
}

// validate user
func (u *Users) Validate() error {
	if err := middlewares.Validate.Struct(u); err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

type GetUsers struct {
	structs.Model
	Username    string `json:"username"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Province    string `json:"province"`
	Postcode    string `json:"postcode"`
	UID         string `json:"user_id"`
	AcceptTerm  bool   `json:"accept_term"`
	RolesID     uint   `json:"roles_id"`
}

type UsersPut struct {
	structs.Model
	Username    string `json:"username" gorm:"uniqueIndex" form:"username" query:"username" swaggertype:"string" example:"username"`
	Email       string `json:"email" form:"email" query:"email" swaggertype:"string" example:"email@email.com"`
	Password    string `json:"-" swaggertype:"string" form:"password" example:"password"`
	FirstName   string `json:"first_name"  form:"first_name" swaggertype:"string" example:"supachai"`
	LastName    string `json:"last_name"  form:"last_name" swaggertype:"string" example:"last name"`
	PhoneNumber string `json:"phone_number" validate:"max=10" form:"phone_number" swaggertype:"string" example:"0899999999"`
	Address     string `json:"address" gorm:"type:text"  form:"address" swaggertype:"string" example:"address"`
	Province    string `json:"province"  form:"province" swaggertype:"string" example:"bankok"`
	Postcode    string `json:"postcode"  form:"postcode" swaggertype:"string" example:"10260"`
	UID         string `json:"user_id" gorm:"uniqueIndex"  form:"user_id" query:"user_id" swaggertype:"string" example:"U1234567890"`
	AcceptTerm  bool   `json:"accept_term"  form:"accept_term" swaggertype:"boolean" example:"true"`
	RolesID     uint   `json:"-"  form:"roles_id" swaggertype:"integer" example:"2"`
}

type UserIDOnly struct {
	UserID []uint `json:"user_id" form:"user_id" swaggertype:"integer" example:"1"`
}

// override table name
func (GetUsers) TableName() string {
	return "users"
}
func (UsersPut) TableName() string {
	return "users"
}
func (UserIDOnly) TableName() string {
	return "users"
}

// before save user generate password and encrypt by bcrypt
func (user *Users) BeforeSave(tx *gorm.DB) (err error) {
	var uidEmail uidEmail
	// find data if id or uid or username is exist
	if user.ID != nil {
		tx.Where("id = ?", user.ID).First(&uidEmail)
	} else if user.UID != "" {
		tx.Where("uid = ?", user.UID).First(&uidEmail)
	} else if user.Username != "" {
		tx.Where("username = ?", user.Username).First(&uidEmail)
	} else {
		return fmt.Errorf("id or uid or username is required")
	}
	// random username
	if len(user.Username) == 0 && uidEmail.Username == "" {
		user.Username = uuid.Must(uuid.NewV4()).String()
	}
	if len(user.Password) == 0 && uidEmail.Password == "" {
		// gen password if not set by uuid v4
		user.Password = uuid.Must(uuid.NewV4()).String()
	}
	if len(user.UID) == 0 && uidEmail.UID == "" {
		user.UID = user.Username
	}

	// encrypt password by bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return nil
}
func (user *UsersPut) BeforeSave(tx *gorm.DB) (err error) {
	var uidEmail uidEmail
	tx.Where("id = ? or uid = ? or username = ?", user.ID, user.UID, user.Username).First(&uidEmail)
	// random username
	if len(user.Username) == 0 && uidEmail.Username == "" {
		user.Username = uuid.Must(uuid.NewV4()).String()
	}
	if len(user.Password) == 0 && uidEmail.Password == "" {
		// gen password if not set by uuid v4
		user.Password = uuid.Must(uuid.NewV4()).String()
	}
	if len(user.UID) == 0 && uidEmail.UID == "" {
		user.UID = user.Username
	}
	// encrypt password by bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	return nil
}

type uidEmail struct {
	UID      string
	Email    string
	Username string
	Password string
}

func (*uidEmail) TableName() string {
	return "users"
}

type POSTUser struct {
	ID          uint   `json:"-" gorm:"primarykey" swaggertype:"integer" example:"1" body:"id"`
	UID         string `json:"-" gorm:"uniqueIndex"   swaggertype:"string" example:"U1234567890" body:"user_id" query:"user_id" form:"user_id"`
	Username    string `json:"username" form:"username" query:"username" swaggertype:"string" example:"username" body:"username"`
	Email       string `json:"email" form:"email" query:"email" swaggertype:"string" example:"ex@ex.com" body:"email"`
	Password    string `json:"passowrd"  swaggertype:"string" form:"password" example:"password" body:"password"`
	FirstName   string `json:"first_name"  form:"first_name" swaggertype:"string" example:"first name" body:"first_name"`
	LastName    string `json:"last_name"  form:"last_name" swaggertype:"string" example:"last name" body:"last_name"`
	PhoneNumber string `json:"phone_number" validate:"max=10" form:"phone_number" swaggertype:"string" example:"0899999999" body:"phone_number"`
	Address     string `json:"address" gorm:"type:text"  form:"address" swaggertype:"string" example:"address" body:"address"`
	Province    string `json:"province"  form:"province" swaggertype:"string" example:"bankok" body:"province"`
	Postcode    string `json:"postcode"  form:"postcode" swaggertype:"string" example:"10260" body:"postcode"`
}

func (POSTUser) TableName() string {
	return "users"
}

type POSTUserSwagger struct {
	UID         string `json:"user_id" form:"user_id" query:"user_id" swaggertype:"string" example:"user_id" body:"user_id"`
	Username    string `json:"username" form:"username" query:"username" swaggertype:"string" example:"username" body:"username"`
	Email       string `json:"email" form:"email" query:"email" swaggertype:"string" example:"ex@ex.com" body:"email"`
	Password    string `json:"passowrd"  swaggertype:"string" form:"password" example:"password" body:"password"`
	FirstName   string `json:"first_name"  form:"first_name" swaggertype:"string" example:"first name" body:"first_name"`
	LastName    string `json:"last_name"  form:"last_name" swaggertype:"string" example:"last name" body:"last_name"`
	PhoneNumber string `json:"phone_number" validate:"max=10" form:"phone_number" swaggertype:"string" example:"0899999999" body:"phone_number"`
	Address     string `json:"address" gorm:"type:text"  form:"address" swaggertype:"string" example:"address" body:"address"`
	Province    string `json:"province"  form:"province" swaggertype:"string" example:"bankok" body:"province"`
	Postcode    string `json:"postcode"  form:"postcode" swaggertype:"string" example:"10260" body:"postcode"`
}

type RESUser struct {
	ID          uint   `json:"-" gorm:"primarykey" swaggertype:"integer" example:"1" body:"id"`
	UID         string `json:"-" gorm:"uniqueIndex"   swaggertype:"string" example:"U1234567890" body:"user_id"`
	Password    string `json:"-" swaggertype:"string" form:"password" example:"password" body:"password"`
	Username    string `json:"username" form:"username" query:"username" swaggertype:"string" example:"username" body:"username"`
	Email       string `json:"email" form:"email" query:"email" swaggertype:"string" example:"ex@ex.com" body:"email"`
	FirstName   string `json:"first_name"  form:"first_name" swaggertype:"string" example:"first name" body:"first_name"`
	LastName    string `json:"last_name"  form:"last_name" swaggertype:"string" example:"last name" body:"last_name"`
	PhoneNumber string `json:"phone_number" validate:"max=10" form:"phone_number" swaggertype:"string" example:"0899999999" body:"phone_number"`
	Address     string `json:"address" gorm:"type:text"  form:"address" swaggertype:"string" example:"address" body:"address"`
	Province    string `json:"province"  form:"province" swaggertype:"string" example:"bankok" body:"province"`
	Postcode    string `json:"postcode"  form:"postcode" swaggertype:"string" example:"10260" body:"postcode"`
}

// // Bind the fields of ... into RESUser
func (user *RESUser) Bind(u *POSTUser) error {
	return mapstructure.Decode(u, user)
}

func (RESUser) TableName() string {
	return "users"
}

func (user *POSTUser) BeforeSave(tx *gorm.DB) (err error) {
	var uidEmail uidEmail
	tx.Where("id = ? or uid = ? or username = ?", user.ID, user.UID, user.Username).First(&uidEmail)
	// random username
	if len(user.Username) == 0 && uidEmail.Username == "" {
		user.Username = uuid.Must(uuid.NewV4()).String()
	}
	if len(user.Password) == 0 && uidEmail.Password == "" {
		// gen password if not set by uuid v4
		user.Password = uuid.Must(uuid.NewV4()).String()
	}
	if len(user.UID) == 0 && uidEmail.UID == "" {
		user.UID = user.Username
	}

	// encrypt password by bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return nil
}
