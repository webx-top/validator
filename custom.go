package validator

import (
	"log"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func NewCustomValidation(tag string, fn validator.Func, options ...CustomValidationOption) *CustomValidation {
	c := &CustomValidation{
		Tag:          tag,
		Func:         fn,
		CallIfNull:   false,
		Translations: map[string]*Translation{},
	}
	for _, option := range options {
		option(c)
	}
	return c
}

type CustomValidation struct {
	Tag          string
	Func         validator.Func
	CallIfNull   bool
	Translations map[string]*Translation
}

type CustomValidationOption func(*CustomValidation)

func OptCallIfNull(on bool) CustomValidationOption {
	return func(c *CustomValidation) {
		c.CallIfNull = on
	}
}

func OptTranslation(locale string, translation *Translation) CustomValidationOption {
	return func(c *CustomValidation) {
		c.Translations[locale] = translation
	}
}

func OptTranslations(translations map[string]*Translation) CustomValidationOption {
	return func(c *CustomValidation) {
		c.Translations = translations
	}
}

type Translation struct {
	Text     string //{0}必须是一个有效的ISBN编号
	Override bool
}

func (v *CustomValidation) Register(validate *validator.Validate, translator ut.Translator, locale string) {
	validate.RegisterValidation(v.Tag, v.Func, v.CallIfNull)
	translation, ok := v.Translations[locale]
	if ok {
		validate.RegisterTranslation(v.Tag, translator, func(translator ut.Translator) error {
			return translator.Add(v.Tag, translation.Text, translation.Override)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, err := ut.T(fe.Tag(), fe.Field())
			if err != nil {
				log.Printf("警告: 翻译字段错误: %#v", fe)
				return fe.(error).Error()
			}
			return t
		})
	}
}

var CustomValidations = map[string]*CustomValidation{}

func RegisterCustomValidation(tag string, fn validator.Func, options ...CustomValidationOption) {
	CustomValidations[tag] = NewCustomValidation(tag, fn, options...)
}
