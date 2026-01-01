package config

const (
	ResultFilePath = "./result/"
)

type Config struct {
	DbPath            string          `toml:"db_path"`
	LLM               *LLMConfig      `toml:"llm"`
	AliyunTTS         *AliyunConfig   `toml:"aliyun_tts"`
	SpeechRecognition *AliyunConfig   `toml:"speech_recognition"`
	BiliBili          *BiliBiliConfig `toml:"biliBili"`
}

type LLMConfig struct {
	Model  *LLMModelConfig `toml:"model"`
	Prompt string          `toml:"prompt"`
}

type LLMModelConfig struct {
	Baidu  *LLMModelBaiduConfig  `toml:"baidu"`
	GLM    *LLMModelGLMConfig    `toml:"glm"`
	Doubao *LLMModelDoubaoConfig `toml:"doubao"`
	Qwen   *LLMModelQwenConfig   `toml:"qwen"`
}

type LLMModelBaiduConfig struct {
	BaseUrl       string `toml:"base_url"`
	APIKey        string `toml:"api_key"`
	ErnieModel    string `json:"ernie_model"`
	DeepSeekModel string `json:"deepseek_model"`
}

type LLMModelGLMConfig struct {
	BaseUrl  string `toml:"base_url"`
	APIKey   string `toml:"api_key"`
	GlmModel string `toml:"glm_model"`
}

type LLMModelDoubaoConfig struct {
	BaseUrl     string `toml:"base_url"`
	APIKey      string `toml:"api_key"`
	DoubaoModel string `toml:"doubao_model"`
}

type LLMModelQwenConfig struct {
	BaseUrl   string `toml:"base_url"`
	APIKey    string `toml:"api_key"`
	QwenModel string `toml:"qwen_model"`
}

type AliyunConfig struct {
	AccessKey string `toml:"access_key"`
	SecretKey string `toml:"secret_key"`
	AppKey    string `toml:"app_key"`
}

type BiliBiliConfig struct {
	AccessKey string `toml:"access_key"`
	SecretKey string `toml:"secret_key"`
	AppId     int64  `toml:"app_id"`
}
