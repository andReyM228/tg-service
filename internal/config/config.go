package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type (
	Config struct {
		TgBot   TgBot   `yaml:"tg-bot" validate:"required"`
		ChatGPT ChatGPT `yaml:"chat-gpt" validate:"required"`
	}

	TgBot struct {
		Token string `yaml:"token" validate:"required"`
	}

	ChatGPT struct {
		Key   string `yaml:"key" validate:"required"`
		Model string `yaml:"model" validate:"required"`
	}
)

func ParseConfig() (Config, error) {
	file, err := os.ReadFile("./cmd/config.yaml")
	if err != nil {
		fmt.Errorf("parseConfig: %s", err)
	}

	var cfg Config

	if err := yaml.Unmarshal(file, &cfg); err != nil {
		fmt.Errorf("parseConfig: %s", err)
	}

	return cfg, nil
}
