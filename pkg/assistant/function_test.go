package assistant

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type AddParams struct {
	A int `json:"a" jsonschema:"description=First number to add,example=2"`
	B int `json:"b,omitempty" jsonschema:"description=Second number to add,example=3"`
}

type AddResult struct {
	Sum int `json:"sum"`
}

func addFunc(_ context.Context, p AddParams) (AddResult, error) {
	return AddResult{Sum: p.A + p.B}, nil
}

func errorFunc(_ context.Context, _ AddParams) (AddResult, error) {
	return AddResult{}, errors.New("fail")
}

func TestFunction_NewFunction(t *testing.T) {
	t.Run("should succeed with valid input", func(t *testing.T) {
		ctx := context.Background()
		fn := NewFunction("add", "Adds two numbers", addFunc)
		args := `{"a":2,"b":3}`

		res, err := fn.callback(ctx, args)
		result, ok := res.(AddResult)

		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, 5, result.Sum)
	})

	t.Run("should return error from function", func(t *testing.T) {
		ctx := context.Background()
		fn := NewFunction("add", "Adds two numbers", errorFunc)
		args := `{"a":1,"b":2}`

		_, err := fn.callback(ctx, args)

		assert.Error(t, err)
	})

	t.Run("should return error for bad json", func(t *testing.T) {
		ctx := context.Background()
		fn := NewFunction("add", "Adds two numbers", addFunc)

		_, err := fn.callback(ctx, `not-json`)

		assert.Error(t, err)
	})

	t.Run("should generate correct DefinitionOf", func(t *testing.T) {
		fn := NewFunction("add", "Adds two numbers", addFunc)

		assert.Equal(t, "add", fn.def.Name)
		assert.Equal(t, "Adds two numbers", fn.def.Description.Value)
		assert.NotEmpty(t, fn.def.Parameters)

		props, ok := fn.def.Parameters["properties"]
		assert.True(t, ok)

		expectedPropsJSON := `{
			"a": {
				"type": "integer",
				"description": "First number to add",
				"examples": [2]
			},
			"b": {
				"type": "integer",
				"description": "Second number to add",
				"examples": [3]
			}
		}`
		propsJSON, err := json.Marshal(props)
		assert.NoError(t, err)
		assert.JSONEq(t, expectedPropsJSON, string(propsJSON))

		required, hasRequired := fn.def.Parameters["required"].([]string)
		assert.True(t, hasRequired)
		assert.ElementsMatch(t, []string{"a"}, required)

		assert.Equal(t, "object", fn.def.Parameters["type"])
	})
}

func TestFunction_FunctionResponse(t *testing.T) {
	t.Run("should create data response", func(t *testing.T) {
		resp := Data("foo")

		assert.Equal(t, "foo", resp.Data)
		assert.Empty(t, resp.Error)
	})

	t.Run("should create error response", func(t *testing.T) {
		err := errors.New("fail")

		resp := Error(err)

		assert.Equal(t, err.Error(), resp.Error)
		assert.Empty(t, resp.Data)
	})
}

func TestFunction_Call(t *testing.T) {
	t.Run("should succeed with valid input", func(t *testing.T) {
		ctx := context.Background()
		fn := NewFunction("add", "Adds two numbers", addFunc)
		args := `{"a":4,"b":6}`

		resp := fn.Call(ctx, args)
		result, ok := resp.Data.(AddResult)

		assert.True(t, ok)
		assert.Empty(t, resp.Error)
		assert.Equal(t, 10, result.Sum)
	})

	t.Run("should return error from function", func(t *testing.T) {
		ctx := context.Background()
		fn := NewFunction("add", "Adds two numbers", errorFunc)
		args := `{"a":1,"b":2}`

		resp := fn.Call(ctx, args)

		assert.Equal(t, "function execution failed: fail", resp.Error)
	})

	t.Run("should return error for bad json", func(t *testing.T) {
		ctx := context.Background()
		fn := NewFunction("add", "Adds two numbers", addFunc)

		resp := fn.Call(ctx, `bad-json`)

		assert.NotEmpty(t, resp.Error)
		assert.Contains(t, resp.Error, "failed to decode function arguments")
	})
}
