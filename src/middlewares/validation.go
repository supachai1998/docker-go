package middlewares

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var Validate *validator.Validate

type CustomValidator struct {
	Validator *validator.Validate
}

func RegisterValidation() {
	Validate = validator.New()
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
