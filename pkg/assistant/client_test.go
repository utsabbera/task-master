package assistant

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Chat(t *testing.T) {
	ctx := context.Background()

	ts := NewTestServer(t)
	defer ts.Close()

	t.Run("should return message from OpenAI API", func(t *testing.T) {
		config := Config{
			BaseURL: ts.URL,
			Model:   "echo",
		}
		cli := NewClient(config)
		cli.Init()

		msg, err := cli.Chat(ctx, "Hello, world!")

		assert.NoError(t, err)
		assert.Equal(t, "Hello, world!", msg)
	})

	t.Run("should handle tool call and marshal response", func(t *testing.T) {
		ctx := context.Background()
		called := false

		config := Config{
			BaseURL: ts.URL,
			Model:   "tool-call",
		}
		cli := NewClient(config)
		cli.RegisterFunction(NewFunction("test", "desc", func(context.Context, struct{}) (any, error) {
			called = true
			return "Done!", nil
		}))
		cli.Init()

		msg, err := cli.Chat(ctx, "test")
		assert.NoError(t, err)
		assert.Equal(t, "```json\n"+`{"data":"Done!"}`+"\n```", msg)
		assert.True(t, called)
	})

	t.Run("should handle tool error and send it in response", func(t *testing.T) {
		ctx := context.Background()
		called := false

		config := Config{
			BaseURL: ts.URL,
			Model:   "tool-call",
		}
		cli := NewClient(config)
		cli.RegisterFunction(NewFunction("test", "desc", func(context.Context, struct{}) (any, error) {
			called = true
			return nil, errors.New("failed to process request")
		}))
		cli.Init()

		msg, err := cli.Chat(ctx, "test")
		assert.NoError(t, err)
		assert.Equal(t, "```json\n"+`{"error":"function execution failed: failed to process request"}`+"\n```", msg)
		assert.True(t, called)
	})

	t.Run("should return error if tool function not found", func(t *testing.T) {
		config := Config{
			BaseURL: ts.URL,
			Model:   "tool-call",
		}
		cli := NewClient(config)
		cli.Init()

		msg, err := cli.Chat(ctx, "notfound")

		assert.Empty(t, msg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "function notfound not found")
	})

	t.Run("should return error if function response cannot be marshalled", func(t *testing.T) {
		ctx := context.Background()
		config := Config{
			BaseURL: ts.URL,
			Model:   "tool-call",
		}
		cli := NewClient(config)
		cli.RegisterFunction(NewFunction("bad", "desc", func(context.Context, struct{}) (any, error) {
			return make(chan int), nil
		}))
		cli.Init()

		msg, err := cli.Chat(ctx, "bad")
		assert.Empty(t, msg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error marshalling function bad response")
	})
}

func TestClient_RegisterFunction(t *testing.T) {
	t.Run("should register a single function", func(t *testing.T) {
		config := Config{
			APIKey:         "test-key",
			BaseURL:        "http://localhost",
			Model:          "gpt-3.5-turbo",
			AppName:        "TestApp",
			AppDescription: "Test app desc",
		}

		cli := NewClient(config)

		fn := NewFunction("fn1", "desc1", func(context.Context, struct{}) (struct{}, error) { return struct{}{}, nil })
		cli.RegisterFunction(fn)

		c := cli.(*client)
		assert.Contains(t, c.funcs, "fn1")
		assert.Equal(t, fn.def, c.funcs["fn1"].def)
	})
}

func TestClient_RegisterFunctions(t *testing.T) {
	t.Run("should register multiple functions", func(t *testing.T) {
		config := Config{
			APIKey:         "test-key",
			BaseURL:        "http://localhost",
			Model:          "gpt-3.5-turbo",
			AppName:        "TestApp",
			AppDescription: "Test app desc",
		}
		cli := NewClient(config)
		fn1 := NewFunction("fn1", "desc1", func(context.Context, struct{}) (struct{}, error) { return struct{}{}, nil })
		fn2 := NewFunction("fn2", "desc2", func(context.Context, struct{}) (struct{}, error) { return struct{}{}, nil })

		cli.RegisterFunctions(fn1, fn2)

		c := cli.(*client)
		assert.Contains(t, c.funcs, "fn1")
		assert.Equal(t, fn1.def, c.funcs["fn1"].def)
		assert.Contains(t, c.funcs, "fn2")
		assert.Equal(t, fn2.def, c.funcs["fn2"].def)
	})
}

func TestClient_Init(t *testing.T) {
	t.Run("should initialize client with registered functions", func(t *testing.T) {
		config := Config{
			APIKey:         "test-key",
			BaseURL:        "http://localhost",
			Model:          "gpt-3.5-turbo",
			AppName:        "TestApp",
			AppDescription: "Test app desc",
		}
		cli := NewClient(config)
		cli.RegisterFunction(NewFunction("fn", "desc", func(context.Context, struct{}) (struct{}, error) { return struct{}{}, nil }))

		cli.Init()

		c := cli.(*client)
		assert.NotNil(t, c.openai)
		assert.NotEmpty(t, c.params.Tools)

		// TODO:  assert system prompt is set
	})
}
