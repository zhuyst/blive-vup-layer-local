package llm

import (
	"blive-vup-layer/config"
	"context"
	"errors"
	"fmt"
	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	log "github.com/sirupsen/logrus"
	"strings"
)

type LLM struct {
	cfg            *config.LLMConfig
	chatCompletion *qianfan.ChatCompletionV2
	client         *openai.Client
}

func NewLLM(config *config.LLMConfig) *LLM {
	client := openai.NewClient(
		// 替换下列示例中参数，将your_APIKey替换为真实值，如何获取API Key请查看https://cloud.baidu.com/doc/WENXINWORKSHOP/s/Um2wxbaps#步骤二-获取api-key
		option.WithAPIKey(config.APIKey),
		option.WithBaseURL("https://qianfan.baidubce.com/v2/"), // 千帆ModelBuilder平台地址
	)
	return &LLM{
		cfg:    config,
		client: client,
	}
}

type ChatMessage struct {
	User    string
	Message string
}

func (msg *ChatMessage) String() string {
	return fmt.Sprintf("用户【%s】说：%s", msg.User, msg.Message)
}

func (llm *LLM) ChatWithLLM(ctx context.Context, messages []*ChatMessage) (string, error) {
	if len(messages) == 0 {
		return "", errors.New("no messages")
	}
	contentSb := strings.Builder{}
	if len(messages) > 1 {
		contentSb.WriteString("以下是历史用户发言：\n")
		for _, msg := range messages[:len(messages)-1] {
			contentSb.WriteString(msg.String() + "\n")
		}
	}
	currentMsg := messages[len(messages)-1]
	contentSb.WriteString("以下是当前用户发言：\n")
	contentSb.WriteString(currentMsg.String())

	content := contentSb.String()
	log.Infof("LLM content: %s", content)

	chatCompletion, err := llm.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.ChatCompletionMessage{Role: "system", Content: llm.cfg.Prompt},
			openai.ChatCompletionMessage{Role: "user", Content: content},
		}),
		Temperature: openai.Float(0.5),
		TopP:        openai.Float(0.5),
		Model:       openai.F(llm.cfg.Model), //模型对应的model值，请查看支持的模型列表：https://cloud.baidu.com/doc/WENXINWORKSHOP/s/wm7ltcvgc
	})
	if err != nil {
		log.Errorf("LLM err: %v", err)
		return "", err
	}

	result := chatCompletion.Choices[0].Message.Content
	result = strings.ReplaceAll(result, "喔~", "喵 ")
	result = strings.ReplaceAll(result, "~", " ")
	result = strings.TrimSpace(result)
	log.Infof("LLM result: %s", result)
	return result, nil
}
