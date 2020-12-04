package mvalidate

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator0(t *testing.T) {
	v := struct {
		Age  int     `validate:"gt=1,lt=3"`
		Name *string `validate:"required"`
	}{}

	arr := []struct {
		Age  int     `validate:"gt=1,lt=3"`
		Name *string `validate:"required"`
	}{v, v}

	validator, err := NewValidatorWithTranslator(zhTranslations.RegisterDefaultTranslations, zh.New())
	assert.NoError(t, err)
	v.Age = 4
	err = validator.Struct(v)
	fmt.Println("e1", err)
	err = validator.Var(arr, `required,dive`)
	fmt.Println("e2", err)
}
