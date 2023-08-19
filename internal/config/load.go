package config

import (
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

const (
	delimeter = "."
	tagName   = "koanf"

	upTemplate     = "================ Loaded Configuration ================"
	bottomTemplate = "======================================================"
)

func Load(print bool) *Config {
	k := koanf.New(delimeter)

	// load default configuration from struct
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default: %s", err)
	}

	// load config from configmap
	if err := loadConfigmap(k); err != nil {
		log.Fatalf("Error loading from configmap: \n%v", err)
	}

	config := Config{}
	var tag = koanf.UnmarshalConf{Tag: tagName}
	if err := k.UnmarshalWithConf("", &config, tag); err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}

	if print {
		// pretty print loaded configuration using provided template
		log.Printf("%s\n%v\n%s\n", upTemplate, spew.Sdump(config), bottomTemplate)
	}

	return &config
}

func loadConfigmap(k *koanf.Koanf) error {
	if os.Getenv("RUNNING_INSIDE_POD") == "" {
		return nil
	}

	cm, err := os.ReadFile("/etc/phone-book/config.yaml")
	if err != nil {
		return fmt.Errorf("Error reading currnet namespace: %v", err)
	}

	if err := k.Load(rawbytes.Provider(cm), yaml.Parser()); err != nil {
		return fmt.Errorf("Error loading values: %s", err)
	}

	return nil
}
