package helpers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type OptionalCustomHTTPError struct {
	Code    *int    `json:"code"`
	Message *string `json:"message"`
}

/*
CustomErr is a custom error handler for echo

	@params
		err: error by panic error
		code: http status code
		message: error message from custom or using err.Error()
	return JSON format
*/
func CustomHTTPError(c echo.Context, err error, optionals ...OptionalCustomHTTPError) error {
	var optional OptionalCustomHTTPError
	if len(optionals) > 0 {
		optional = optionals[0]
	}

	var Errors echo.Map
	errs := json.Unmarshal([]byte(err.Error()), &Errors)
	if errs == nil {
		log.Println("[error] : ", Errors)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"status":  http.StatusBadRequest,
			"message": Errors,
		})
	}

	if err != nil {
		log.Println("[error] : ", err.Error())
		// if msg is empty using err.Error() by case
		var msg string
		var code int
		if optional.Message == nil {
			switch err {
			case gorm.ErrRecordNotFound:
				msg = "record not found"
				code = http.StatusBadRequest
			case gorm.ErrInvalidTransaction:
				msg = "invalid transaction"
				code = http.StatusBadRequest
			case echo.ErrForbidden:
				msg = "forbidden"
				code = http.StatusForbidden
			case echo.ErrUnauthorized:
				msg = "unauthorized"
				code = http.StatusUnauthorized
			case echo.ErrNotFound:
				msg = "not found"
				code = http.StatusNotFound
			case echo.ErrMethodNotAllowed:
				msg = "method not allowed"
				code = http.StatusMethodNotAllowed
			case echo.ErrBadRequest:
				msg = "bad request"
				code = http.StatusBadRequest
			case gorm.ErrEmptySlice:
				msg = "no data"
				code = http.StatusBadRequest
			default:
				if strings.Contains(err.Error(), "token contains an invalid number of segments") {
					msg = "token contains an invalid number of segments"
					code = http.StatusUnauthorized
				} else if strings.Contains(err.Error(), "is required") {
					msg = err.Error()
					code = http.StatusBadRequest
				} else if strings.Contains(err.Error(), "unique") {
					msg = "duplicate data"
					code = http.StatusConflict
				} else if strings.Contains(err.Error(), "เพิ่มเพื่อน") {
					msg = err.Error()
					code = http.StatusBadRequest
				} else if strings.Contains(err.Error(), "linebot:") {
					msg = err.Error()
					trim := strings.Trim(err.Error(), " ")
					getStatus := strings.Split(trim, "APIError")[1][1:4]
					codeStatus, _ := strconv.Atoi(getStatus)
					code = codeStatus
				} else if strings.Contains(err.Error(), "required") {
					var fieldValiation []string
					splitFor := strings.Split(err.Error(), "for")
					for _, val := range splitFor {
						if strings.Contains(val, "failed") {
							field := strings.Trim(strings.Split(val, "failed")[0], " ")
							// convert to sneak case
							field = SneakCase(field)
							fieldValiation = append(fieldValiation, field)
						}
					}
					msg = "field " + strings.Join(fieldValiation, ",") + " is required"
					code = http.StatusBadRequest
				} else if err.Error() == "user not lucky draw" {
					msg = err.Error()
					code = http.StatusBadRequest
				} else if strings.Contains(err.Error(), "SQLSTATE 42703") {
					msg = "field not found"
					code = http.StatusBadRequest
				} else {
					msg = "Internal Server Error"
					code = http.StatusInternalServerError
				}
			}
		} else {
			msg = *optional.Message
			code = *optional.Code
		}

		// custom return error here
		return c.JSON(code, echo.Map{
			"code":    code,
			"message": msg,
			"status":  code,
		})
	}
	return nil
}
