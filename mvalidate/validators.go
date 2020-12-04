package mvalidate

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
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

// newly create validator with specified localeName
func NewValidator(localeName string) (Validator, error) {
	transInfo, ok := kTranslateRegFuncMap[localeName]
	if !ok {
		return nil, errors.Errorf("locale translator:[%s] not registed", localeName)
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
