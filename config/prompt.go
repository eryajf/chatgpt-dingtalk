package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Prompt struct {
	Title  string `yaml:"title"`
	Prefix string `yaml:"prefix"`
	Suffix string `yaml:"suffix"`
}

var prompTmp *[]Prompt

// LoadPrompt 加载Prompt
func LoadPrompt() *[]Prompt {
	data, err := ioutil.ReadFile("prompt.yml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &prompTmp)
	if err != nil {
		log.Fatal(err)
	}
	return prompTmp
}
