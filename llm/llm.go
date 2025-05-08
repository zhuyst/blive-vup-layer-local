package llm

import (
	"blive-vup-layer/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"unicode/utf8"
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

type ChatWithLLMParams struct {
	ExtraInfos []string
	Messages   []*ChatMessage
}

type result struct {
	Reply string `json:"reply"`
}

const maxReplyLength = 30

func (llm *LLM) ChatWithLLM(ctx context.Context, params *ChatWithLLMParams) (*LLMResult, error) {
	requestId := uuid.NewV4().String()
	l := log.WithField("request_id", requestId)

	messages := params.Messages
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
	l.Infof("LLM content: %s", content)

	prompt := strings.TrimSpace(llm.cfg.Prompt)
	if len(params.ExtraInfos) > 0 {
		var extraInfoSb strings.Builder
		for _, extraInfo := range params.ExtraInfos {
			extraInfoSb.WriteString(fmt.Sprintf("- %s\n", extraInfo))
		}
		prompt = strings.Replace(prompt, "{{extra}}", extraInfoSb.String(), 1)
	} else {
		prompt = strings.Replace(prompt, "{{extra}}", "无", 1)
	}
	l.Infof("LLM prompt: %s", prompt)

	opts := []option.RequestOption{
		option.WithJSONSet("search_source", "baidu_search_v2"),
		option.WithJSONSet("prompt_template", prompt),
		option.WithJSONSet("enable_reasoning", true),
		option.WithJSONSet("response_format", "text"),
		option.WithJSONSet("enable_corner_markers", false),
	}

	chatCompletionParams := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			//openai.ChatCompletionMessage{Role: "system", Content: llm.cfg.Prompt},
			openai.ChatCompletionMessage{Role: "user", Content: content},
		}),
		Temperature: openai.Float(0.5),
		TopP:        openai.Float(0.5),
		Model:       openai.F(llm.cfg.Model),
	}
	chatCompletion, err := llm.client.Chat.Completions.New(ctx, chatCompletionParams, opts...)
	if err != nil {
		log.Errorf("LLM err: %v", err)
		return nil, err
	}

	res := &LLMResult{}
	message := chatCompletion.Choices[0].Message
	if reasoningContentI, ok := message.JSON.ExtraFields["reasoning_content"]; ok {
		res.ReasoningContent = convertUnicode(reasoningContentI.Raw())
	}
	l.Infof("LLM result reasoning_content: %s", res.ReasoningContent)

	resContent := chatCompletion.Choices[0].Message.Content
	l.Infof("LLM result content: %s", resContent)

	var llmResult result
	if err := json.Unmarshal([]byte(resContent), &llmResult); err != nil {
		return nil, err
	}

	resContent = llmResult.Reply
	resContent = convertUnicode(resContent)
	resContent = strings.ReplaceAll(resContent, "喔~", "喵 ")
	resContent = strings.ReplaceAll(resContent, "~", " ")
	res.Content = strings.TrimSpace(resContent)

	contentLength := utf8.RuneCountInString(res.Content)
	if contentLength > maxReplyLength {
		return nil, fmt.Errorf("LLM result content too long: %d", contentLength)
	}
	return res, nil
}

func convertUnicode(s string) string {
	res, err := strconv.Unquote(strings.Replace(strconv.Quote(s), `\\u`, `\u`, -1))
	if err != nil {
		return s
	}
	return res
}
