package mvalidate

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/guestin/mob/merrors"
	"strings"
)

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

type ValidateError []string

func (this ValidateError) Error() string {
	var sb strings.Builder
	for idx, it := range this {
		sb.WriteString(it)
		if idx != len(this)-1 {
			sb.WriteString(";\n")
		}
	}
	return sb.String()
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
		ve := make(ValidateError, 0, len(validateErr))
		for k, v := range validateErr.Translate(this.utTranslator) {
			ve = append(ve, fmt.Sprintf("%s:%s", k, v))
		}
		return ve
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

// newly create validator with specified localeName
func NewValidator(localeName string) (Validator, error) {
	transInfo, ok := kTranslateRegFuncMap[localeName]
	if !ok {
		return nil, merrors.ErrorWrapf(nil, "locale translator:[%s] not registed", localeName)
	}
	validatorIns := validator.New()
	uniTranslator := transInfo.BuildTranslator()
	err := transInfo.RegisterTranslatorFunc(validatorIns, uniTranslator)
	if err != nil {
		return nil, err
	}
	return &_validator{
		validatorIns,
		uniTranslator,
	}, nil
}
