package mvalidate

import (
	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"strings"
)

// provider by eg: zhTranslator
type RegisterTranslatorFunc func(v *validator.Validate, trans ut.Translator) (err error)

type Validator interface {
	//
	Struct(v interface{}) error
	//
	Var(v interface{}, tag string) error
	// warning! only can be used on Raw() validated values
	Translate(error) error
	// raw validator
	Raw() *validator.Validate
}

type _validator struct {
	raw          *validator.Validate
	utTranslator ut.Translator
}

func (this *_validator) Translate(err error) error {
	if err == nil {
		return nil
	}
	if this.utTranslator == nil {
		return err
	}
	if validateErr, ok := err.(validator.ValidationErrors); ok {
		var sb strings.Builder
		for k, v := range validateErr.Translate(this.utTranslator) {
			sb.WriteString(k)
			sb.WriteString(":")
			sb.WriteString(v)
			sb.WriteString(";")
		}
		return errors.New(sb.String())
	}
	return err
}

func (this *_validator) Struct(v interface{}) error {
	err := this.raw.Struct(v)
	return this.Translate(err)
}

func (this *_validator) Var(v interface{}, tag string) error {
	err := this.raw.Var(v, tag)
	return this.Translate(err)
}

func (this *_validator) Raw() *validator.Validate {
	return this.raw
}

// newly create validator with translator
func NewValidatorWithTranslator(regFunc RegisterTranslatorFunc, primary locales.Translator) (Validator, error) {
	uniTranslator := ut.New(primary)
	primaryUniTranslator, _ := uniTranslator.GetTranslator(primary.Locale())
	validatorIns := validator.New()
	err := regFunc(validatorIns, primaryUniTranslator)
	if err != nil {
		return nil, err
	}
	return &_validator{
		validatorIns,
		primaryUniTranslator,
	}, nil
}
