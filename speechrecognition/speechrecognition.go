package speechrecognition

import (
	"blive-vup-layer/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type SpeechRecognition struct {
	cfg    *config.AliyunConfig
	client *sdk.Client
}

func NewSpeechRecognition(cfg *config.AliyunConfig) (*SpeechRecognition, error) {
	aliyunConfig := sdk.NewConfig()
	credential := &credentials.AccessKeyCredential{
		AccessKeyId:     cfg.AccessKey,
		AccessKeySecret: cfg.SecretKey,
	}
	client, err := sdk.NewClientWithOptions("cn-beijing", aliyunConfig, credential)
	if err != nil {
		return nil, err
	}
	return &SpeechRecognition{
		cfg:    cfg,
		client: client,
	}, nil
}

type TokenResult struct {
	ErrMsg string `json:"ErrMsg"`
	Token  *Token `json:"Token"`
}

type Token struct {
	UserId     string `json:"UserId"`
	Id         string `json:"Id"`
	ExpireTime int64  `json:"ExpireTime"`
}

func (sr *SpeechRecognition) getToken() (*Token, error) {
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Domain = "nls-meta.cn-shanghai.aliyuncs.com"
	request.ApiName = "CreateToken"
	request.Version = "2019-02-28"
	response, err := sr.client.ProcessCommonRequest(request)
	if err != nil {
		return nil, err
	}
	httpStatus := response.GetHttpStatus()
	if httpStatus != http.StatusOK {
		return nil, fmt.Errorf("get token failed: status_code: %d, content: %s", httpStatus, response.GetHttpContentString())
	}

	var tr TokenResult
	if err := json.Unmarshal([]byte(response.GetHttpContentString()), &tr); err != nil {
		return nil, err
	}
	return tr.Token, nil
}

type FlashRecognizerResult struct {
	TaskId      string       `json:"task_id"`
	Result      string       `json:"result"`
	Status      int          `json:"status"`
	Message     string       `json:"message"`
	FlashResult *FlashResult `json:"flash_result"`
}

type FlashResult struct {
	Duration  int64       `json:"duration"`
	Completed bool        `json:"completed"`
	Latency   int64       `json:"latency"`
	Sentences []*Sentence `json:"sentences"`
}

type Sentence struct {
	Text      string `json:"text"`
	BeginTime int64  `json:"begin_time"`
	EndTime   int64  `json:"end_time"`
	ChannelId int64  `json:"channel_id"`
}

func (sr *SpeechRecognition) RecognitionWav(ctx context.Context, wavFileBytes []byte) (string, error) {
	u := "https://nls-gateway-cn-beijing.aliyuncs.com/stream/v1/FlashRecognizer"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBuffer(wavFileBytes))
	if err != nil {
		return "", err
	}
	token, err := sr.getToken()
	if err != nil {
		return "", err
	}
	query := url.Values{}
	query.Set("appkey", sr.cfg.AppKey)
	query.Set("token", token.Id)
	query.Set("format", "wav")
	query.Set("sample_rate", "16000")
	query.Set("enable_inverse_text_normalization", "true")
	query.Set("enable_timestamp_alignment", "true")
	req.URL.RawQuery = query.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("recognition failed: status_code: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var res FlashRecognizerResult
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}
	textSb := strings.Builder{}
	for _, s := range res.FlashResult.Sentences {
		textSb.WriteString(s.Text)
	}
	return textSb.String(), nil
}
