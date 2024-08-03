package validation

import (
	"github.com/dlclark/regexp2"
)

type (
	Config struct {
		EmailRegex    string
		UsernameRegex string
		PasswordRegex string
	}

	Validator struct {
		emailRegex    *regexp2.Regexp
		usernameRegex *regexp2.Regexp
		passwordRegex *regexp2.Regexp
	}
)

func NewValidator(config Config) (*Validator, error) {
	emailRegex, err := regexp2.Compile(config.EmailRegex, 0)
	if err != nil {
		return nil, err
	}

	usernameRegex, err := regexp2.Compile(config.UsernameRegex, 0)
	if err != nil {
		return nil, err
	}

	passwordRegex, err := regexp2.Compile(config.PasswordRegex, 0)
	if err != nil {
		return nil, err
	}

	return &Validator{
		emailRegex:    emailRegex,
		usernameRegex: usernameRegex,
		passwordRegex: passwordRegex,
	}, nil
}

func (v *Validator) Email(email string) bool {
	ok, _ := v.emailRegex.MatchString(email)
	return ok
}

func (v *Validator) Username(username string) bool {
	ok, _ := v.usernameRegex.MatchString(username)
	return ok
}

func (v *Validator) Password(password string) bool {
	ok, _ := v.passwordRegex.MatchString(password)
	return ok
}
