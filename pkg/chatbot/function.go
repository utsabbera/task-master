package chatbot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)

// Function wraps a callable function with OpenAI-compatible metadata and invocation logic.
//
// Example usage:
//
//	fn := NewFunction("my_func", "Does something", func(ctx context.Context, p MyParams) (MyResult, error) {
//	    // ...
//	})
//	resp := fn.Call(ctx, `{"field": "value"}`)
type Function struct {
	def  openai.FunctionDefinitionParam
	Func func(ctx context.Context, params string) (any, error)
}

// NewFunction creates a new Function with the given name, description, and implementation.
//
// P is the parameter type, R is the result type.
//
// The struct fields of P can be annotated with `jsonschema` tags to control the generated schema for OpenAI function parameters.
// For example, use `jsonschema:"minLength=1,required"` to specify validation constraints.
//
// See: https://pkg.go.dev/github.com/invopop/jsonschema for supported tags and options.
//
// Example:
//
//	type AddParams struct {
//		A int `json:"a" jsonschema:"minimum=0,required"`
//		B int `json:"b" jsonschema:"minimum=0,required"`
//	}
//	type AddResult struct { Sum int }
//	f := NewFunction[AddParams, AddResult]("add", "Adds two numbers", func(ctx context.Context, p AddParams) (AddResult, error) {
//		return AddResult{Sum: p.A + p.B}, nil
//	})
func NewFunction[P, R any](
	name, description string,
	fn func(context.Context, P) (R, error),
) Function {
	return Function{
		def: definition(name, description, fn),
		Func: func(ctx context.Context, args string) (any, error) {
			var param P
			if err := json.Unmarshal([]byte(args), &param); err != nil {
				// TODO: Handle this error properly
				return nil, fmt.Errorf("failed to decode function arguments: %w", err)
			}

			return fn(ctx, param)
		},
	}
}

// definition generates an OpenAI function definition for the given function signature.
func definition[P, R any](
	name, description string,
	_ func(context.Context, P) (R, error),
) openai.FunctionDefinitionParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}

	var param P
	schema := reflector.Reflect(param)

	funcParams := shared.FunctionParameters{
		"type":       schema.Type,
		"properties": schema.Properties,
		"required":   schema.Required,
	}

	return openai.FunctionDefinitionParam{
		Name:        name,
		Parameters:  funcParams,
		Description: openai.String(description),
	}
}

// FunctionResponse represents the result of a Function call.
//
// Either Data or Error will be set.
type FunctionResponse struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

// Error creates a FunctionResponse with the given error.
func Error(err error) FunctionResponse {
	return FunctionResponse{Error: err.Error()}
}

// Data creates a FunctionResponse with the given data.
func Data(data any) FunctionResponse {
	return FunctionResponse{Data: data}
}

// Call invokes the Function with the given arguments and returns a FunctionResponse.
//
// Example:
//
//	resp := fn.Call(ctx, `{"field": "value"}`)
func (f Function) Call(ctx context.Context, args string) FunctionResponse {
	result, err := f.Func(ctx, args)
	if err != nil {
		return Error(fmt.Errorf("function execution failed: %w", err))
	}

	return Data(result)
}
