package chatbot

// Config holds the configuration for the chatbot client.
type Config struct {
	// Model is the name of the LLM model to use.
	Model string `json:"model" yaml:"model"`
	// BaseURL is the base URL of the LLM service.
	BaseURL string `json:"baseUrl" yaml:"baseUrl"`
	// APIKey is the API token for the LLM service.
	APIKey string `json:"apiKey" yaml:"apiKey"`
	// AppName is the name of the application using the LLM.
	AppName string `json:"appName" yaml:"appName"`
	// AppDescription is the description of the application using the LLM.
	AppDescription string `json:"appDescription" yaml:"appDescription"`

	// TODO: Allow additional instructions to be passed to the LLM
}
