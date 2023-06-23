package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type (
	Config struct {
		TgBot   TgBot   `yaml:"tg-bot"`
		ChatGPT ChatGPT `yaml:"chat-gpt"`
	}

	TgBot struct {
		Token string `yaml:"token"`
	}

	ChatGPT struct {
		Key   string `yaml:"key"`
		Model string `yaml:"model"`
	}
)

func ParseConfig() (Config, error) {
	file, err := os.ReadFile("C:\\Users\\admin\\Desktop\\projects\\buycars\\tg-service\\cmd\\config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config

	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Fatal(err)
	}

	return cfg, nil
}
