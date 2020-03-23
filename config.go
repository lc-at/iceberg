package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	ClientName         string `yaml:"client_name"`
	DbConnectionString string `yaml:"db_connection_string"`
	SessionFilename    string `yaml:"session_filename"`
}

func loadConfig(c *config) {
	f, err := os.Open("config.yml")
	checkError(err)
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&c)
	checkError(err)
}
