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
	APIKey string `toml:"api_key"`
	Model  string `json:"model"`
	Prompt string `toml:"prompt"`
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
