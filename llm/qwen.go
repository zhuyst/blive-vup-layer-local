package llm

import (
	"blive-vup-layer/config"
	"context"
	"errors"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type qwenProvider struct {
	client openai.Client
}

func newQwenProvider(cfg *config.LLMModelQwenConfig) *qwenProvider {
	client := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseUrl),
	)
	return &qwenProvider{client: client}
}

type QwenResponseFormat struct {
	Type       string      `json:"type"`
	JsonSchema *JsonSchema `json:"json_schema"`
}

type JsonSchema struct {
	Name        string            `json:"name"`
	Strict      bool              `json:"strict"`
	Description string            `json:"description"`
	Schema      *JsonSchemaSchema `json:"schema"`
}

type JsonSchemaSchema struct {
	Type                 string                         `json:"type"`
	Properties           map[string]*JsonSchemaProperty `json:"properties"`
	Required             []string                       `json:"required"`
	AdditionalProperties bool                           `json:"additionalProperties"`
}

type JsonSchemaProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

var resultJsonSchema = JsonSchema{
	Name:        "reply_result",
	Description: "回复的响应结构体",
	Strict:      true,
	Schema: &JsonSchemaSchema{
		Type: "object",
		Properties: map[string]*JsonSchemaProperty{
			"reply": {Type: "string", Description: "回复的内容"},
		},
		Required: []string{"reply"},
	},
}

func (p *qwenProvider) chatWithLLM(ctx context.Context, params *chatParams) (*chatResult, error) {
	opts := []option.RequestOption{
		option.WithJSONSet("enable_thinking", true),
		option.WithJSONSet("thinking_budget", "200"),
		option.WithJSONSet("enable_search", true),
		//option.WithJSONSet("response_format", QwenResponseFormat{
		//	Type:       "json_schema",
		//	JsonSchema: &resultJsonSchema,
		//}),
	}

	chatCompletionParams := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(params.Prompt),
			openai.UserMessage(params.Content),
		},
		Temperature: openai.Float(0.5),
		TopP:        openai.Float(0.5),
		Model:       params.Model,
	}
	stream := p.client.Chat.Completions.NewStreaming(ctx, chatCompletionParams, opts...)
	chatCompletion := openai.ChatCompletionAccumulator{}
	reasoningBuilder := strings.Builder{}
	for stream.Next() {
		chunk := stream.Current()
		chatCompletion.AddChunk(chunk)

		if reasoningContentI, ok := chunk.JSON.ExtraFields["reasoning_content"]; ok {
			reasoningBuilder.WriteString(reasoningContentI.Raw())
		}
	}
	if err := stream.Err(); err != nil {
		return nil, err
	}

	if len(chatCompletion.Choices) == 0 {
		return nil, errors.New("empty choices")
	}

	res := &chatResult{}
	res.ReasoningContent = reasoningBuilder.String()
	res.Content = chatCompletion.Choices[0].Message.Content
	return res, nil
}
