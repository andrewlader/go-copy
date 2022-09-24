package copylib

import (
	"strings"

	"github.com/spf13/viper"
)

type replaceMode int8

const (
	replaceNever replaceMode = iota
	replaceSkipIfSame
	replaceAlways
)

type configuration struct {
	name         string
	source       string
	destinations []string
	replace      replaceMode
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
	}

	switch strings.ToLower(config["replace"].(string)) {
	case "always":
		configObj.replace = replaceAlways
		break
	case "never":
		configObj.replace = replaceNever
		break
	case "skip":
		configObj.replace = replaceSkipIfSame
		break
	}

	return configObj
}
