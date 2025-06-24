package assistant

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openai/openai-go"
	"github.com/stretchr/testify/require"
)

// NewTestServer returns an httptest.Server that mocks an LLM (Large Language Model) API endpoint for unit testing.
// It inspects incoming requests, decodes OpenAI chat completion parameters, and responds with predefined outputs.
//
// Supported models:
//   - "echo": Responds with the user's message content as the reply.
//   - "tool-call": If the last message is a tool message, replies with its content. Otherwise, replies with a tool call where the function name matches the user's message.
//
// The server uses testify/require for request validation and fails the test on unexpected input or model.
//
// Example:
//
//	ts := NewTestServer(t)
//	defer ts.Close()
//	config := Config{BaseURL: ts.URL, Model: "echo"}
//	cli := NewClient(config)
//	cli.Init()
//	msg, err := cli.Chat(context.TODO(), "Hello, world!")
//
// Use this helper to simulate LLM responses in tests without real API calls.
func NewTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params openai.ChatCompletionNewParams
		err := json.NewDecoder(r.Body).Decode(&params)
		require.NoError(t, err)

		message := params.Messages[len(params.Messages)-1]

		switch params.Model {
		case "echo":
			resp := openai.ChatCompletion{
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Content: message.OfUser.Content.OfString.String(),
						},
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			require.NoError(t, json.NewEncoder(w).Encode(resp))
			return

		case "tool-call":
			if toolMessage := message.OfTool; toolMessage != nil {
				resp := openai.ChatCompletion{
					Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{
						Content: toolMessage.Content.OfString.String(),
					}}},
				}

				w.Header().Set("Content-Type", "application/json")
				require.NoError(t, json.NewEncoder(w).Encode(resp))
				return
			}

			resp := openai.ChatCompletion{
				Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{
					Content: "",
					ToolCalls: []openai.ChatCompletionMessageToolCall{{
						ID: "1",
						Function: openai.ChatCompletionMessageToolCallFunction{
							Name:      message.OfUser.Content.OfString.String(),
							Arguments: "{}",
						},
					}},
				}}},
			}

			w.Header().Set("Content-Type", "application/json")
			require.NoError(t, json.NewEncoder(w).Encode(resp))
			return

		default:
			require.Fail(t, "Unexpected model: %s", params.Model)
		}
	}))
}
