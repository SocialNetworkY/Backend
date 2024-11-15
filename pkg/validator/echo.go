package validator

import (
	"net/http"

	"github.com/dlclark/regexp2"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type EchoCustomValidator struct {
	v *validator.Validate
}

var (
	usernameRegex = regexp2.MustCompile("^[a-zA-Z0-9]{3,52}$", 0)
	passwordRegex = regexp2.MustCompile("^(?=.*[A-Za-z])(?=.*\\d)[A-Za-z\\d]{8,255}$", 0)
)

func NewEchoCustomValidator() *EchoCustomValidator {
	v := validator.New()
	v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		b, _ := usernameRegex.MatchString(fl.Field().String())
		return b
	})
	v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		b, _ := passwordRegex.MatchString(fl.Field().String())
		return b
	})

	return &EchoCustomValidator{
		v: v,
	}
}

func (cv *EchoCustomValidator) Validate(i interface{}) error {
	if err := cv.v.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
