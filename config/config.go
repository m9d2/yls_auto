package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var (
	Conf = Config{}
)

type Config struct {
	Token string `yaml:"token"`
}

func InitConfig() {
	file, _ := os.ReadFile("./config.yaml")
	err := yaml.Unmarshal(file, &Conf)
	if err != nil {
		log.Fatal(err.Error())
	}
}
