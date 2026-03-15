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

// print displaysa text representation of the configuration.
func (config *configuration) print() {
	if config == nil {
		return
	}

	// dests := strings.Join(config.destinations, ", ")
	replaceStr := "unknown"
	switch config.replace {
	case replaceNever:
		replaceStr = "never"
	case replaceSkipIfSame:
		replaceStr = "skip"
	case replaceAlways:
		replaceStr = "always"
	}

	PrintKeyValue("Name: ", config.name)
	PrintKeyValue("  Source: ", config.source)
	PrintKeyValueArray("  Destinations: ", config.destinations)
	PrintKeyValue("  Replace: ", replaceStr)
}

// ListConfigurations displays the string representations for all configurations.
func ListConfigurations() {
	keys := viper.AllSettings()
	for key := range keys {
		cfg := getConfiguration(key)
		if cfg != nil {
			cfg.print()
		}
	}

	Print("End of List")
}

func getConfiguration(key string) *configuration {
	var destinations []string

	config := viper.GetStringMap(key)
	if config == nil {
		return nil
	}

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

	case "never":
		configObj.replace = replaceNever

	case "skip":
		configObj.replace = replaceSkipIfSame
	}

	return configObj
}
