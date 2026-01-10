package llm

import (
	"blive-vup-layer/config"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	log "github.com/sirupsen/logrus"
)

type glmProvider struct {
	client openai.Client
}

func newGLMProvider(cfg *config.LLMModelGLMConfig) *glmProvider {
	client := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseUrl),
	)
	return &glmProvider{client: client}
}

type GlmThinking struct {
	Type          string `json:"type"`
	ClearThinking bool   `json:"clear_thinking"`
}

type GlmTool struct {
	Type      string        `json:"type"`
	WebSearch *GlmWebSearch `json:"web_search"`
}

type GlmWebSearch struct {
	Enable       string `json:"enable"`
	SearchEngine string `json:"search_engine"`
	SearchIntent string `json:"search_intent"`
	SearchPrompt string `json:"search_prompt"`
	Count        int    `json:"count"`
	ContentSize  string `json:"content_size"`
}

type GlmResponseFormat struct {
	Type string `json:"type"`
}

func (p *glmProvider) chatWithLLM(ctx context.Context, params *chatParams) (*chatResult, error) {
	searchPrompt := fmt.Sprintf("请用简洁的语言总结网络搜索{search_result}中的关键信息，按重要性排序并引用来源日期。现在的时间是%s。", time.Now().Format("2006年01月02日 15:04:05"))
	opts := []option.RequestOption{
		option.WithJSONSet("tools", []GlmTool{
			{
				Type: "web_search",
				WebSearch: &GlmWebSearch{
					Enable:       "true",
					SearchEngine: "search_pro",
					SearchIntent: "false",
					SearchPrompt: searchPrompt,
					Count:        5,
					ContentSize:  "medium",
				},
			},
		}),
		option.WithJSONSet("thinking", GlmThinking{
			Type:          "disabled",
			ClearThinking: true,
		}),
		option.WithJSONSet("response_format", GlmResponseFormat{
			Type: "json_object",
		}),
	}

	chatCompletionParams := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(params.Prompt),
			openai.UserMessage(params.Content),
		},
		Temperature: openai.Float(0.5),
		TopP:        openai.Float(0.5),
		Model:       params.ModelName,
	}
	chatCompletion, err := p.client.Chat.Completions.New(ctx, chatCompletionParams, opts...)
	if err != nil {
		log.Errorf("LLM err: %v", err)
		return nil, err
	}
	if len(chatCompletion.Choices) == 0 {
		return nil, errors.New("empty choices")
	}

	res := &chatResult{}

	message := chatCompletion.Choices[0].Message
	if reasoningContentI, ok := message.JSON.ExtraFields["reasoning_content"]; ok {
		res.ReasoningContent = reasoningContentI.Raw()
	}

	res.Content = chatCompletion.Choices[0].Message.Content
	return res, nil
}
