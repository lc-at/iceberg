package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	ClientName      string `yaml:"client_name"`
	DbFilename      string `yaml:"db_filename"`
	SessionFilename string `yaml:"session_filename"`
}

func loadConfig(c *config) {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&c)
	if err != nil {
		log.Fatalln(err)
	}
}
