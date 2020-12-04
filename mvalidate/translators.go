package mvalidate

import (
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// provider by eg: zhTranslator
type RegisterTranslatorFunc func(v *validator.Validate, trans ut.Translator) (err error)
type LocaleTranslatorBuildFunc func() locales.Translator

type translatorInfo struct {
	RegisterTranslatorFunc    RegisterTranslatorFunc
	LocaleTranslatorBuildFunc LocaleTranslatorBuildFunc
}

func (this translatorInfo) BuildTranslator() ut.Translator {
	localeTranslator := this.LocaleTranslatorBuildFunc()
	uni := ut.New(localeTranslator)
	uniTranslator := uni.GetFallback()
	return uniTranslator
}

const DefaultTranslator = "zh"

var kTranslateRegFuncMap = map[string]translatorInfo{
	"zh": {
		RegisterTranslatorFunc:    zhTranslations.RegisterDefaultTranslations,
		LocaleTranslatorBuildFunc: zh.New},
	"en": {
		RegisterTranslatorFunc:    enTranslations.RegisterDefaultTranslations,
		LocaleTranslatorBuildFunc: en.New},
}

// register
func RegisterValidateTranslation(
	localeName string,
	regFunc RegisterTranslatorFunc,
	localeFunc LocaleTranslatorBuildFunc) {
	kTranslateRegFuncMap[localeName] = translatorInfo{
		RegisterTranslatorFunc:    regFunc,
		LocaleTranslatorBuildFunc: localeFunc,
	}
}
