package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Prompt struct {
	Title  string `yaml:"title"`
	Prefix string `yaml:"prefix"`
	Suffix string `yaml:"suffix"`
}

var prompTmp *[]Prompt

// LoadPrompt 加载Prompt
func LoadPrompt() *[]Prompt {
	data, err := os.ReadFile("prompt.yml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &prompTmp)
	if err != nil {
		log.Fatal(err)
	}
	return prompTmp
}
