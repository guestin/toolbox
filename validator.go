package mod

import (
	"bytes"
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/pkg/errors"
)

/**
 * Created by xuchao on 2019-03-06 .
 * Update by suzhen on 2020-08-21 .
 */

// use a single instance , it caches struct info
var (
	uni          *ut.UniversalTranslator
	gValidate    *validator.Validate
	gValidatorZH = NewValidator(LangZH)
	gValidatorEN = NewValidator(LangEN)
)

func GlobalValidator() *validator.Validate {
	return gValidate
}

func DefaultValidatorZH() Validator {
	return gValidatorZH
}

func DefaultValidatorEN() Validator {
	return gValidatorEN
}

type LANG int

const (
	LangZH LANG = iota
	LangEN
)

const (
	StrLangZH = "zh"
	StrLangEN = "en"
)

func init() {
	zhTranslator := zh.New()
	enTranslator := en.New()
	uni = ut.New(enTranslator, zhTranslator, enTranslator)
	gValidate = validator.New()
	trans, found := uni.GetTranslator(StrLangZH)
	if found {
		if err := zhTranslations.RegisterDefaultTranslations(gValidate, trans); err != nil {
			panic(err)
		}
	}
}

var ValidateUnknownLANG = errors.New("unknown lang")

func (this *LANG) MarshalText() (text []byte, err error) {
	switch *this {
	case LangZH:
		return []byte(StrLangZH), nil
	case LangEN:
		return []byte(StrLangEN), nil
	}
	return nil, ValidateUnknownLANG
}

func (this *LANG) UnmarshalText(text []byte) error {
	switch string(text) {
	case StrLangZH:
		*this = LangZH
	case StrLangEN:
		*this = LangEN
	default:
		return ValidateUnknownLANG
	}
	return nil
}

type Validator interface {
	Validate(i interface{}) error
}

func ValidateStructWithLanguage(lang LANG, i interface{}) error {
	e := gValidate.Struct(i)
	if e != nil {
		if _, ok := e.(*validator.InvalidValidationError); ok {
			return e
		}
		// translate all error at once
		var buffer bytes.Buffer
		rawErrs := e.(validator.ValidationErrors)
		locale, err := lang.MarshalText()
		if err != nil {
			return err
		}
		trans, found := uni.GetTranslator(string(locale))
		if found {
			tansErrs := rawErrs.Translate(trans)
			for _, err := range tansErrs {
				buffer.WriteString(fmt.Sprintf("%s;", err))
			}
		} else {
			for _, err := range rawErrs {
				eStr := fmt.Sprintf("param:'%s' type:'%s' miss match with check:'%s';",
					err.Field(),
					err.Kind(),
					err.Tag())
				buffer.WriteString(eStr)
			}
		}
		return errors.New(buffer.String())
	}
	return nil
}

// @see ValidateStructWithLanguage LangEN
func ValidateStruct(i interface{}) error {
	return ValidateStructWithLanguage(LangEN, i)
}

// @see ValidateStructWithLanguage LangEN
func Verify(i interface{}) error {
	return ValidateStructWithLanguage(LangEN, i)
}

type __validator struct {
	lang LANG
}

func NewValidator(lang LANG) Validator {
	return &__validator{
		lang: lang,
	}
}

func (this *__validator) Validate(i interface{}) error {
	return ValidateStructWithLanguage(this.lang, i)
}
