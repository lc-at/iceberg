package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

// config for Iceberg
type config struct {
	ClientName         string            `yaml:"client_name"`
	DbConnectionString string            `yaml:"db_connection_string"`
	SessionFilename    string            `yaml:"session_filename"`
	Days               map[int]string    `yaml:"days"`
	MessageTemplates   map[string]string `yaml:"message_templates"`
}

func loadConfig(c *config) {
	f, err := os.Open("config.yml")
	checkError(err)
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&c)
	checkError(err)
}

func (c config) getMessageTemplate(key string) string {
	template, ok := c.MessageTemplates[key]
	if !ok {
		return ""
	}
	return template
}

func (c config) getDayByName(name string) (int, bool) {
	for k, v := range c.Days {
		if v == name {
			return k, true
		}
	}
	return 0, false
}

func (c config) getNameByDay(num int) (string, bool) {
	name, ok := c.Days[num]
	if !ok {
		return "", false
	}
	return name, true
}
