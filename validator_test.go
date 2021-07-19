package validator

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

type testBody struct {
	Name string `validate:"required" json:"name"`
	Age  int64  `validate:"required,min=18,max=200" json:"age"`
}

func TestValidator(t *testing.T) {
	a := &testBody{
		Age: 6,
	}
	v := New()
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
}
