package llm

import (
	"net/http"
	"time"

	"github.com/sashabaranov/go-openai"
)

func NewClient(APIKey, URL, timeout string) *openai.Client {
	// Set up OpenAI client
	if APIKey == "" {
		//log.Fatal("OPENAI_API_KEY environment variable not set")
		APIKey = "sk-xxx"
	}
	config := openai.DefaultConfig(APIKey)
	config.BaseURL = URL

	dur, err := time.ParseDuration(timeout)
	if err != nil {
		dur = 150 * time.Second
	}

	transport := &headerTransport{
		headers: map[string]string{
			"HTTP-Referer": "https://iodevs.net",
			"X-Title":      "LocalAGI",
		},
		base: http.DefaultTransport,
	}

	config.HTTPClient = &http.Client{
		Timeout:   dur,
		Transport: transport,
	}
	return openai.NewClientWithConfig(config)
}

type headerTransport struct {
	headers map[string]string
	base    http.RoundTripper
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Set(k, v)
	}
	return t.base.RoundTrip(req)
}
