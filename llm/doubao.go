package llm

import (
	"blive-vup-layer/config"
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model/responses"
)

type doubaoProvider struct {
	client *arkruntime.Client
}

func newDoubaoProvider(cfg *config.LLMModelDoubaoConfig) *doubaoProvider {
	client := arkruntime.NewClientWithApiKey(
		cfg.APIKey,
		arkruntime.WithBaseUrl(cfg.BaseUrl),
	)
	return &doubaoProvider{client: client}
}

func (p *doubaoProvider) chatWithLLM(ctx context.Context, params *chatParams) (*chatResult, error) {
	maxToolCalls := int64(1)
	thinkType := responses.ThinkingType_enabled
	store := false
	temperature := 0.5
	topP := 0.5

	//schemaJson, err := json.Marshal(resultJsonSchema.Schema)
	//if err != nil {
	//	return nil, err
	//}

	systemMessage := &responses.ItemInputMessage{
		Role: responses.MessageRole_system,
		Content: []*responses.ContentItem{
			{
				Union: &responses.ContentItem_Text{
					Text: &responses.ContentItemText{
						Type: responses.ContentItemType_input_text,
						Text: params.Prompt,
					},
				},
			},
		},
	}
	contentMessage := &responses.ItemInputMessage{
		Role: responses.MessageRole_user,
		Content: []*responses.ContentItem{
			{
				Union: &responses.ContentItem_Text{
					Text: &responses.ContentItemText{
						Type: responses.ContentItemType_input_text,
						Text: params.Content,
					},
				},
			},
		},
	}
	createResponsesReq := &responses.ResponsesRequest{
		Model: params.Model,
		Thinking: &responses.ResponsesThinking{
			Type: &thinkType,
		},
		Reasoning: &responses.ResponsesReasoning{
			Effort: responses.ReasoningEffort_medium,
		},
		Store:       &store,
		Temperature: &temperature,
		TopP:        &topP,
		Input: &responses.ResponsesInput{
			Union: &responses.ResponsesInput_ListValue{
				ListValue: &responses.InputItemList{ListValue: []*responses.InputItem{
					{
						Union: &responses.InputItem_InputMessage{
							InputMessage: systemMessage,
						},
					},
					{
						Union: &responses.InputItem_InputMessage{
							InputMessage: contentMessage,
						},
					},
				}},
			},
		},
		//Text: &responses.ResponsesText{
		//	Format: &responses.TextFormat{
		//		Type: responses.TextType_json_object,
		//		//Name:        resultJsonSchema.Name,
		//		//Description: &resultJsonSchema.Description,
		//		//Schema:      &responses.Bytes{Value: schemaJson},
		//		//Strict:      &resultJsonSchema.Strict,
		//	},
		//},
		Tools: []*responses.ResponsesTool{
			{
				Union: &responses.ResponsesTool_ToolWebSearch{
					ToolWebSearch: &responses.ToolWebSearch{
						Type: responses.ToolType_web_search,
					},
				},
			},
		},
		MaxToolCalls: &maxToolCalls,
	}

	resp, err := p.client.CreateResponses(ctx, createResponsesReq)
	if err != nil {
		log.Errorf("LLM err: %v", err)
		return nil, err
	}

	res := &chatResult{}
	for _, output := range resp.Output {
		if reasoning := output.GetReasoning(); reasoning != nil {
			summary := reasoning.GetSummary()
			if len(summary) == 0 {
				continue
			}
			res.ReasoningContent = summary[0].GetText()
		}
		if message := output.GetOutputMessage(); message != nil {
			content := message.GetContent()
			if len(content) == 0 {
				continue
			}
			res.Content = content[0].GetText().GetText()
		}
	}

	return res, nil
}
