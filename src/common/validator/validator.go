package appvalidator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	"sync/atomic"
)

type (
	translatorHolder struct {
		t ut.Translator
	}
)

var (
	globalTran = defaultValue()
)

func defaultValue() *atomic.Value {
	v := &atomic.Value{}
	v.Store(translatorHolder{t: nil})
	return v
}

func GetTranslator() ut.Translator {
	return globalTran.Load().(translatorHolder).t
}

func setTranslator(translator ut.Translator) {
	globalTran.Store(translatorHolder{t: translator})
}

func RegisterGinValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		trans := registerTranslate(v)
		if trans != nil {
			setTranslator(trans)
		}
	}
}

func registerTranslate(v *validator.Validate) ut.Translator {
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)

	trans, _ := uni.GetTranslator("en")
	err := entranslations.RegisterDefaultTranslations(v, trans)

	if err != nil {
		panic(err)
		return nil
	}
	return trans
}

func registerValidator(v *validator.Validate, trans ut.Translator, tag string, validator func(fl validator.FieldLevel) bool, errMessage string) {
	addValidation(v, tag, validator)
	addTranslation(v, trans, tag, errMessage)
}

func addValidation(v *validator.Validate, tag string, validatorFunc func(fl validator.FieldLevel) bool) {
	err := v.RegisterValidation(tag, validatorFunc)
	if err != nil {
		panic(err)
	}
}

func addTranslation(v *validator.Validate, trans ut.Translator, tag string, errMessage string) {
	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, errMessage, false)
	}

	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		tag := fe.Tag()

		t, err := ut.T(tag, fe.Field(), param)
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}

	_ = v.RegisterTranslation(tag, trans, registerFn, transFn)
}
