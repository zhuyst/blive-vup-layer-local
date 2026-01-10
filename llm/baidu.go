// https://ai.baidu.com/ai-doc/AppBuilder/amaxd2det

package llm

import (
	"blive-vup-layer/config"
	"context"
	"errors"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	log "github.com/sirupsen/logrus"
)

type baiduProvider struct {
	client openai.Client
}

func newBaiduProvider(cfg *config.LLMModelBaiduConfig) *baiduProvider {
	client := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseUrl),
	)
	return &baiduProvider{client: client}
}

type BaiduSearchItemsPostprocess struct {
	WindowSize int `json:"window_size"`
	StrideSize int `json:"stride_size"`
	MaxSlice   int `json:"max_slice"`
}

type BaiduResponseFormat struct {
	Type       string      `json:"type"`
	JsonSchema *JsonSchema `json:"json_schema"`
}

func (p *baiduProvider) chatWithLLM(ctx context.Context, params *chatParams) (*chatResult, error) {
	opts := []option.RequestOption{
		option.WithJSONSet("search_source", "baidu_search_v2"),
		option.WithJSONSet("prompt_template", params.Prompt),
		option.WithJSONSet("enable_reasoning", true),
		option.WithJSONSet("response_format", "text"),
		option.WithJSONSet("enable_corner_markers", false),
		option.WithJSONSet("search_items_postprocess", BaiduSearchItemsPostprocess{
			WindowSize: 400,
			StrideSize: 300,
			MaxSlice:   4,
		}),
	}

	chatCompletionParams := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
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
