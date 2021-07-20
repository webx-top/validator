package validator

import (
	"context"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

type testBody struct {
	Name string `validate:"required" json:"name"`
	Age  int64  `validate:"required,min=18,max=200" json:"age"`
}

func TestValidator(t *testing.T) {
	ctx := defaults.NewMockContext()
	a := &testBody{
		Age: 6,
	}
	v := New(ctx)
	err := v.Struct(a)
	if err != nil {
		t.Log(err)
	}
	assert.Equal(t, `Name为必填字段`, err.Error())

	a.Name = "test"
	result := v.Validate(a)
	assert.Equal(t, false, result.Ok())
	assert.Equal(t, `Age最小只能为18`, result.Error().Error())

	a.Age = 20
	result = v.Validate(a)
	assert.Equal(t, true, result.Ok())
	assert.Equal(t, nil, result.Error())

	v2 := New(ctx)
	a.Name = ""
	err = v2.Struct(a)
	if err != nil {
		t.Log(err)
	}
	assert.Equal(t, `Name为必填字段`, err.Error())
}

type testBody2 struct {
	Name string `validate:"required,custom" json:"name"`
	Age  int64  `validate:"required,min=18,max=200" json:"age"`
}

func TestCustomValidator(t *testing.T) {
	RegisterCustomValidation(`custom`, func(ctx context.Context, f validator.FieldLevel) bool {
		eCtx := ctx.(echo.Context)
		assert.Equal(t, `json`, eCtx.Format())
		return f.Field().String() == `test`
	}, OptTranslations(map[string]*Translation{
		`zh`: {Text: `输入的名称无效`},
		`en`: {Text: `invalid name`},
	}))
	ctx := defaults.NewMockContext()
	ctx.SetFormat(`json`)
	a := &testBody2{
		Age: 6,
	}
	v := New(ctx)
	err := v.Struct(a)
	if err != nil {
		t.Log(err)
	}
	assert.Equal(t, `Name为必填字段`, err.Error())

	a.Name = "test2"
	result := v.Validate(a)
	assert.Equal(t, false, result.Ok())
	assert.Equal(t, `输入的名称无效`, result.Error().Error())

	a.Name = "test"
	result = v.Validate(a)
	assert.Equal(t, false, result.Ok())
	assert.Equal(t, `Age最小只能为18`, result.Error().Error())
}
