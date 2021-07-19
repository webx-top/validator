package validator

import "github.com/go-playground/validator"

type CustomValidation struct {
	Tag        string
	Func       validator.Func
	CallIfNull bool
}

var CustomValidations = map[string]*CustomValidation{}

func RegisterCustomValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) {
	var callIfNull bool
	if len(callValidationEvenIfNull) > 0 {
		callIfNull = callValidationEvenIfNull[0]
	}
	CustomValidations[tag] = &CustomValidation{Tag: tag, Func: fn, CallIfNull: callIfNull}
}
