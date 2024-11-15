package validator

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEchoCustomValidator_Validate(t *testing.T) {
	v := NewEchoCustomValidator()
	type args struct {
		Username string `validate:"username,min=3,max=52"`
		Password string `validate:"password"`
	}

	assert.Nil(t, v.Validate(args{Username: "test", Password: "Test1234"}))
	dd := v.Validate(args{Username: "test", Password: "test"})
	log.Println(dd)
	assert.Error(t, dd)
	assert.Error(t, v.Validate(args{Username: "te", Password: "1234sadasD"}))
	assert.Error(t, v.Validate(args{Username: "testasdadsadsasdadsasddasdsadasdadsadsadsaadsdasdasda", Password: "1234asdDD"}))
	assert.Error(t, v.Validate(args{Username: "", Password: ""}))
}
