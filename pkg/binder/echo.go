package binder

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mime/multipart"
	"reflect"
	"strings"
)

type EchoCustomBinder struct {
	binder echo.DefaultBinder
}

func NewEchoCustomBinder() *EchoCustomBinder {
	return &EchoCustomBinder{}
}

func (cb *EchoCustomBinder) Bind(i interface{}, c echo.Context) error {
	if err := cb.binder.Bind(i, c); err != nil {
		return err
	}

	contentType := c.Request().Header.Get(echo.HeaderContentType)

	if !strings.HasPrefix(contentType, echo.MIMEApplicationForm) && !strings.HasPrefix(contentType, echo.MIMEMultipartForm) {
		return nil
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	return bindFiles(i, form.File)
}

func bindFiles(i interface{}, files map[string][]*multipart.FileHeader) error {
	iValue := reflect.Indirect(reflect.ValueOf(i))
	if iValue.Kind() != reflect.Struct {
		err := fmt.Errorf("input is not struct pointer, indirect type is %s", iValue.Type().String())
		return err
	}

	iType := iValue.Type()
	for i := 0; i < iType.NumField(); i++ {
		fValue := iValue.Field(i)
		if !fValue.CanSet() {
			continue
		}

		fType := iType.Field(i)

		switch fType.Type {
		case reflect.TypeOf((*multipart.FileHeader)(nil)):
			file := getFiles(files, fType.Tag.Get("form"), fType.Name)
			if len(file) > 0 {
				fValue.Set(reflect.ValueOf(file[0]))
			}
		case reflect.TypeOf(([]*multipart.FileHeader)(nil)):
			file := getFiles(files, fType.Tag.Get("form"), fType.Name)
			if len(file) > 0 {
				fValue.Set(reflect.ValueOf(file))
			}
		}
	}
	return nil
}

func getFiles(files map[string][]*multipart.FileHeader, names ...string) []*multipart.FileHeader {
	for _, name := range names {
		file, ok := files[name]
		if ok {
			return file
		}
	}
	return nil
}
