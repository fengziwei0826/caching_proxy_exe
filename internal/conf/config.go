package conf

import (
	"encoding/json"
	"io"
	"os"
)

type GlobalConfig struct {
	HTTPConfig  `json:"http"`
	RedisConfig `json:"redis"`
}

type HTTPConfig struct {
	Port      int
	ProxyAddr string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

var glabalConfig *GlobalConfig

func init() {
	glabalConfig = &GlobalConfig{
		HTTPConfig: HTTPConfig{
			Port:      3000,
			ProxyAddr: "http://localhost:8080",
		},
		RedisConfig: RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
	}
}

func GetGlobalConfig() *GlobalConfig {
	return glabalConfig
}

func LoadConfigFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	bs, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bs, glabalConfig); err != nil {
		return err
	}
	return nil
}
