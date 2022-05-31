package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	MySQL map[string]string
}

func NewConfig() Configuration {

	file, err := os.ReadFile("/app/config/conf.json")

	if err != nil {
		fmt.Println(err)
	}

	configuration := Configuration{}

	err = json.Unmarshal(file, &configuration)
	if err != nil {
		fmt.Println(err)
	}
	return configuration

}
