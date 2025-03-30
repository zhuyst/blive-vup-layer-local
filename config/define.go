package config

const (
	ResultFilePath = "./result/"
)

type Config struct {
	DbPath    string           `toml:"db_path"`
	LLM       *LLMConfig       `toml:"llm"`
	AliyunTTS *AliyunTTSConfig `toml:"aliyun_tts"`
	BiliBili  *BiliBiliConfig  `toml:"biliBili"`
}

type LLMConfig struct {
	APIKey string `toml:"api_key"`
	Model  string `json:"model"`
	Prompt string `toml:"prompt"`
}

type AliyunTTSConfig struct {
	AccessKey string `toml:"access_key"`
	SecretKey string `toml:"secret_key"`
	AppKey    string `toml:"app_key"`
}

type BiliBiliConfig struct {
	AccessKey string `toml:"access_key"`
	SecretKey string `toml:"secret_key"`
	AppId     int64  `toml:"app_id"`
}
