package validator

import (
	"errors"
	"strings"

	"github.com/admpub/log"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/webx-top/echo"
)

var DefaultLocale = `zh`

func NormalizeLocale(locale string) string {
	return strings.Replace(locale, `-`, `_`, -1)
}

func New(locales ...string) *Validate {
	var locale string
	if len(locales) > 0 {
		locale = locales[0]
	}
	if len(locale) == 0 {
		locale = DefaultLocale
	}
	translator, _ := UniversalTranslator().GetTranslator(locale)
	validate := validator.New()
	transtation, ok := Translations[locale]
	if !ok {
		args := strings.SplitN(locale, `_`, 2)
		if len(args) == 2 {
			transtation, ok = Translations[args[0]]
		} else {
			log.Warnf(`[validator] not found translation: %s`, locale)
		}
	}
	if ok {
		transtation(validate, translator)
	}
	return &Validate{
		validator:  validate,
		translator: translator,
	}
}

type Validate struct {
	validator  *validator.Validate
	translator ut.Translator
}

func (v *Validate) Object() *validator.Validate {
	return v.validator
}

// ValidateMap validates map data form a map of tags
func (v *Validate) ValidateMap(data map[string]interface{}, rules map[string]interface{}) map[string]interface{} {
	return v.validator.ValidateMap(data, rules)
}

// Struct 接收的参数为一个struct
func (v *Validate) Struct(i interface{}) error {
	return v.Error(v.validator.Struct(i))
}

// StructExcept 校验struct中的选项，不过除了fields里所给的字段
func (v *Validate) StructExcept(s interface{}, fields ...string) error {
	return v.Error(v.validator.StructExcept(s))
}

// StructFiltered 接收一个struct和一个函数，这个函数的返回值为bool，决定是否跳过该选项
func (v *Validate) StructFiltered(s interface{}, fn validator.FilterFunc) error {
	return v.Error(v.validator.StructFiltered(s, fn))
}

// StructPartial 接收一个struct和fields，仅校验在fields里的值
func (v *Validate) StructPartial(s interface{}, fields ...string) error {
	return v.Error(v.validator.StructPartial(s, fields...))
}

// Var 接收一个变量和一个tag的值，比如 validate.Var(i, "gt=1,lt=10")
func (v *Validate) Var(field interface{}, tag string) error {
	return v.Error(v.validator.Var(field, tag))
}

// VarWithValue 将两个变量进行对比，比如 validate.VarWithValue(s1, s2, "eqcsfield")
func (v *Validate) VarWithValue(field interface{}, other interface{}, tag string) error {
	return v.Error(v.validator.VarWithValue(field, other, tag))
}

// Validate 此处支持两种用法：
// 1. Validate(表单字段名, 表单值, 验证规则名)
// 2. Validate(结构体实例, 要验证的结构体字段1，要验证的结构体字段2)
// Validate(结构体实例) 代表验证所有带“valid”标签的字段
func (v *Validate) Validate(i interface{}, args ...string) echo.ValidateResult {
	e := echo.NewValidateResult()
	var err error
	switch m := i.(type) {
	case string:
		field := m
		var value, rule string
		switch len(args) {
		case 2:
			rule = args[1]
			fallthrough
		case 1:
			value = args[0]
		}
		if len(rule) == 0 {
			return e
		}
		err = v.validator.Var(value, rule)
		if err != nil {
			e.SetField(field)
			e.SetRaw(err)
			return e.SetError(v.Error(err))
		}
	default:
		if len(args) > 0 {
			err = v.validator.StructPartial(i, args...)
		} else {
			err = v.validator.Struct(i)
		}
		if err != nil {
			vErrors := err.(validator.ValidationErrors)
			e.SetField(vErrors[0].Field())
			e.SetRaw(vErrors[0])
			return e.SetError(v.Error(vErrors[0]))
		}
	}
	return e
}

func (v *Validate) Error(err error) error {
	if err == nil {
		return nil
	}
	switch rErr := err.(type) {
	case validator.FieldError:
		return errors.New(rErr.Translate(v.translator))
	case validator.ValidationErrors:
		return errors.New(rErr[0].Translate(v.translator))
	default:
		return err
	}
}