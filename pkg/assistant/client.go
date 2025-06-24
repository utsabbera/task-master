// Package assistant provides an interface for AI chat clients to send messages and receive responses.
//
// Example usage:
//
//	package main
//
//	import (
//		"context"
//		"github.com/yourorg/yourrepo/pkg/assistant"
//	)
//
//	func main() {
//		config := assistant.Config{ /* ... */ }
//		client := assistant.NewClient(config)
//		client.Init()
//		ctx := context.Background()
//		resp, err := client.Chat(ctx, "Hello!")
//		if err != nil {
//			// handle error
//		}
//		fmt.Println(resp)
//	}
package assistant

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
	"github.com/utsabbera/task-master/pkg/util"
)

//go:generate mockgen -destination=client_mock.go -package=assistant . Client

// Client defines the interface for an AI chat client that can register functions and process messages.
type Client interface {
	// Init initializes the client and its dependencies.
	Init()
	// RegisterFunction registers a single function for use by the chat client.
	RegisterFunction(fn Function)
	// RegisterFunctions registers multiple functions for use by the chat client.
	RegisterFunctions(funcs ...Function)
	// Chat sends a prompt to the chat client and returns the response.
	Chat(ctx context.Context, prompt string) (string, error)
}

type client struct {
	config Config
	params openai.ChatCompletionNewParams
	openai openai.Client
	funcs  map[string]Function
}

func NewClient(config Config) Client {
	return &client{
		config: config,
		funcs:  make(map[string]Function),
	}
}

func (c *client) RegisterFunction(fn Function) {
	c.funcs[fn.def.Name] = fn
}

func (c *client) RegisterFunctions(funcs ...Function) {
	for _, fn := range funcs {
		c.RegisterFunction(fn)
	}
}

func (c *client) Init() {
	c.initOpenAIClient()
	c.initParams()
}

func (c *client) initOpenAIClient() {
	options := make([]option.RequestOption, 0, 2)

	if c.config.BaseURL != "" {
		options = append(options, option.WithBaseURL(c.config.BaseURL))
	}

	if c.config.APIKey != "" {
		options = append(options, option.WithAPIKey(c.config.APIKey))
	}

	c.openai = openai.NewClient(options...)
}

func (c *client) initParams() {
	tools := make([]openai.ChatCompletionToolParam, 0, len(c.funcs))

	for _, fn := range c.funcs {
		tools = append(tools, openai.ChatCompletionToolParam{
			Function: fn.def,
		})
	}

	c.params = openai.ChatCompletionNewParams{
		Model: shared.ChatModel(c.config.Model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(c.systemPrompt()),
		},
		Tools:       tools,
		Seed:        openai.Int(0),
		Temperature: openai.Float(0.2),
	}
}

func (c *client) systemPrompt() string {
	return fmt.Sprintf(`
You are a helpful, concise, and context-aware chat assistant for a %s application - %s. 

Your capabilities are limited to the following tasks:
`+strings.Join(
		util.Map(
			util.Values(c.funcs),
			func(f Function) string {
				return fmt.Sprintf("- %s: %s", f.def.Name, f.def.Description)
			},
		),
		"\n",
	)+`
Always confirm actions with the user, provide clear feedback, and handle errors gracefully. If a requested task is not found or an operation fails, inform the user and suggest next steps. Use natural, friendly language and keep responses brief and actionable.
`, c.config.AppName, c.config.AppDescription)
}

func (c *client) Chat(ctx context.Context, prompt string) (string, error) {
	c.params.Messages = append(c.params.Messages, openai.UserMessage(prompt))
	return c.process(ctx)
}

func (c *client) process(ctx context.Context) (string, error) {
	completion, err := c.openai.Chat.Completions.New(ctx, c.params)
	if err != nil {
		return "", err
	}

	if len(completion.Choices) > 0 {
		response := completion.Choices[0].Message

		c.params.Messages = append(c.params.Messages, response.ToParam())

		// TODO: handle refusal

		if len(response.ToolCalls) > 0 {
			err := c.handleToolCalls(ctx, response.ToolCalls)
			if err != nil {
				return "", err
			}

			return c.process(ctx)
		}

		return response.Content, nil
	}

	return "", nil
}

func (c *client) handleToolCalls(ctx context.Context, calls []openai.ChatCompletionMessageToolCall) error {
	// TODO: Handle the errors properly - loop it through chat

	for _, call := range calls {
		// TODO: Pass these errors for the AI model to handle

		fn, exists := c.funcs[call.Function.Name]
		if !exists {
			return fmt.Errorf("function %s not found", call.Function.Name)
		}

		response := fn.Call(ctx, call.Function.Arguments)
		respBytes, err := json.Marshal(response)
		if err != nil {
			return fmt.Errorf("error marshalling function %s response: %w", call.Function.Name, err)
		}

		message := "```json\n" + string(respBytes) + "\n```"
		c.params.Messages = append(c.params.Messages, openai.ToolMessage(message, call.ID))
	}

	return nil
}
