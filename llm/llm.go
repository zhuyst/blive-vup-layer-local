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
		option.WithAPIKey(config.APIKey),
		option.WithBaseURL(config.BaseUrl),
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

type LLMResult struct {
	ReasoningContent string
	Content          string
}

func (llm *LLM) ChatWithLLM(ctx context.Context, messages []*ChatMessage) (*LLMResult, error) {
	if len(messages) == 0 {
		return nil, errors.New("no messages")
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
		Model:       openai.F(llm.cfg.Model),
	})
	if err != nil {
		log.Errorf("LLM err: %v", err)
		return nil, err
	}

	res := &LLMResult{}
	message := chatCompletion.Choices[0].Message
	if reasoningContentI, ok := message.JSON.ExtraFields["reasoning_content"]; ok {
		res.ReasoningContent = reasoningContentI.Raw()
	}
	resContent := chatCompletion.Choices[0].Message.Content
	resContent = strings.ReplaceAll(resContent, "喔~", "喵 ")
	resContent = strings.ReplaceAll(resContent, "~", " ")
	res.Content = strings.TrimSpace(resContent)
	log.Infof("LLM result reasoning_content: %s, content: %s", res.ReasoningContent, res.Content)
	return res, nil
}
