package copylib

import (
	"github.com/spf13/viper"
)

type configuration struct {
	name         string
	source       string
	destinations []string
	replace      bool
}

func getConfiguration(key string) *configuration {
	var destinations []string

	config := viper.GetStringMap(key)
	dests := config["destinations"].([]interface{})

	for _, destInst := range dests {
		destinations = append(destinations, destInst.(string))
	}

	configObj := &configuration{
		name:         config["name"].(string),
		source:       config["source"].(string),
		destinations: destinations,
		replace:      config["replace"].(bool),
	}

	return configObj
}
