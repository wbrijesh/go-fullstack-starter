package utils

import (
	"fmt"
	"go-fullstack-starter/schema"
	"os"

	"gopkg.in/yaml.v2"
)

func ReadConfig() schema.ConfigType {
	var Config schema.ConfigType

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("Error reading config file" + err.Error())
	}

	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		fmt.Println(err)
	}

	return Config
}
