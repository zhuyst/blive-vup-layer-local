package llm

import (
	"blive-vup-layer/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type ModelType string

const (
	ModelTypeErnie    ModelType = "ernie"
	ModelTypeDeepSeek ModelType = "deepseek"
	ModelTypeGLM      ModelType = "glm"
	ModelTypeQwen     ModelType = "qwen"
	ModelTypeDoubao   ModelType = "doubao"
)

type LLM struct {
	cfg         *config.LLMConfig
	providerMap map[ModelType]*providerWithModel
}

type providerWithModel struct {
	Provider  provider
	ModelType ModelType
	ModelName string
}

func NewLLM(cfg *config.LLMConfig) *LLM {
	baiduP := newBaiduProvider(cfg.Model.Baidu)
	glmP := newGLMProvider(cfg.Model.GLM)
	doubaoP := newDoubaoProvider(cfg.Model.Doubao)
	qwenP := newQwenProvider(cfg.Model.Qwen)

	pm := make(map[ModelType]*providerWithModel)
	pm[ModelTypeErnie] = &providerWithModel{
		Provider:  baiduP,
		ModelType: ModelTypeErnie,
		ModelName: cfg.Model.Baidu.ErnieModel,
	}
	pm[ModelTypeDeepSeek] = &providerWithModel{
		Provider:  baiduP,
		ModelType: ModelTypeDeepSeek,
		ModelName: cfg.Model.Baidu.DeepSeekModel,
	}
	pm[ModelTypeGLM] = &providerWithModel{
		Provider:  glmP,
		ModelType: ModelTypeGLM,
		ModelName: cfg.Model.GLM.GlmModel,
	}
	pm[ModelTypeDoubao] = &providerWithModel{
		Provider:  doubaoP,
		ModelType: ModelTypeDoubao,
		ModelName: cfg.Model.Doubao.DoubaoModel,
	}
	pm[ModelTypeQwen] = &providerWithModel{
		Provider:  qwenP,
		ModelType: ModelTypeQwen,
		ModelName: cfg.Model.Qwen.QwenModel,
	}
	return &LLM{
		cfg:         cfg,
		providerMap: pm,
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
	Model      ModelType
	ExtraInfos []string
	Messages   []*ChatMessage
}

type result struct {
	Reply string `json:"reply"`
}

const maxReplyLength = 30

func (llm *LLM) ChatWithLLM(ctx context.Context, params *ChatWithLLMParams) (*LLMResult, error) {
	requestId := uuid.NewV4().String()
	l := log.WithFields(log.Fields{
		"request_id": requestId,
		"model":      params.Model,
	})

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

	p, ok := llm.providerMap[params.Model]
	if !ok {
		return nil, errors.New("invalid model")
	}

	chatPar := &chatParams{
		Prompt:    prompt,
		Content:   content,
		ModelType: p.ModelType,
		ModelName: p.ModelName,
	}
	chatRes, err := p.Provider.chatWithLLM(ctx, chatPar)
	if err != nil {
		return nil, err
	}

	res := &LLMResult{}
	resContent := convertUnicode(chatRes.Content)
	res.ReasoningContent = convertUnicode(chatRes.ReasoningContent)
	l.Infof("LLM result reasoning_content: %s", res.ReasoningContent)
	l.Infof("LLM result content: %s", resContent)

	resContent = strings.TrimPrefix(resContent, "```json")
	resContent = strings.TrimSuffix(resContent, "```")

	var llmResult result
	if err := json.Unmarshal([]byte(resContent), &llmResult); err != nil {
		return nil, err
	}

	resContent = llmResult.Reply
	resContent = convertUnicode(resContent)
	resContent = strings.ReplaceAll(resContent, "喔~", "喵 ")
	resContent = strings.ReplaceAll(resContent, "~", " ")
	res.Content = strings.TrimSpace(resContent)

	if res.Content == "" {
		return nil, errors.New("LLM return empty")
	}

	contentLength := utf8.RuneCountInString(res.Content)
	if contentLength > maxReplyLength {
		return nil, fmt.Errorf("LLM result content too long: %d", contentLength)
	}
	return res, nil
}

type chatParams struct {
	Prompt    string
	Content   string
	ModelType ModelType
	ModelName string
}

type chatResult struct {
	Content          string
	ReasoningContent string
}

type provider interface {
	chatWithLLM(ctx context.Context, params *chatParams) (*chatResult, error)
}

func convertUnicode(s string) string {
	res, err := strconv.Unquote(strings.Replace(strconv.Quote(s), `\\u`, `\u`, -1))
	if err != nil {
		return s
	}
	return res
}
