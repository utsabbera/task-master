package assistant

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
	def      FunctionDefinition
	callback func(ctx context.Context, params string) (any, error)
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
		def: DefinitionOf(name, description, fn),
		callback: func(ctx context.Context, args string) (any, error) {
			var param P
			if err := json.Unmarshal([]byte(args), &param); err != nil {
				// TODO: Handle this error properly
				return nil, fmt.Errorf("failed to decode function arguments: %w", err)
			}

			return fn(ctx, param)
		},
	}
}

func (f Function) Definition() FunctionDefinition {
	return f.def
}

// Call invokes the Function with the given arguments and returns a FunctionResponse.
//
// Example:
//
//	resp := fn.Call(ctx, `{"field": "value"}`)
func (f Function) Call(ctx context.Context, args string) FunctionResponse {
	result, err := f.callback(ctx, args)
	if err != nil {
		return Error(fmt.Errorf("function execution failed: %w", err))
	}

	return Data(result)
}

func toOpenAIFunctionDefinitionParam(def FunctionDefinition) openai.FunctionDefinitionParam {
	required := make([]string, 0, len(def.Parameters.Required))

	params := shared.FunctionParameters{
		"type":       def.Parameters.Type,
		"properties": def.Parameters.Properties,
		"required":   append(required, def.Parameters.Required...), // TODO: Handle this properly
	}

	return openai.FunctionDefinitionParam{
		Name:        def.Name,
		Parameters:  params,
		Description: openai.String(def.Description),
	}
}

type FunctionDefinition struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Parameters  *jsonschema.Schema `json:"parameters"`
}

// DefinitionOf generates an function definition for the given function signature.
func DefinitionOf[P, R any](
	name, description string,
	_ func(context.Context, P) (R, error),
) FunctionDefinition {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}

	var param P
	paramSchema := reflector.Reflect(param)

	return FunctionDefinition{
		Name:        name,
		Description: description,
		Parameters:  paramSchema,
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
