package validation

import (
	"errors"
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

func (v *Validator) Email(email string) error {
	ok, _ := v.emailRegex.MatchString(email)
	if !ok {
		return errors.New("invalid email")
	}
	return nil
}

func (v *Validator) Username(username string) error {
	ok, _ := v.usernameRegex.MatchString(username)
	if !ok {
		return errors.New("invalid username")
	}
	return nil
}

func (v *Validator) Password(password string) error {
	ok, _ := v.passwordRegex.MatchString(password)
	if !ok {
		return errors.New("invalid password")
	}
	return nil
}

func (v *Validator) Login(login string) error {
	err := v.Email(login)
	if err == nil {
		return nil
	}
	err = v.Username(login)
	if err == nil {
		return nil
	}
	return errors.New("invalid login")
}
