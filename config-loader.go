package main

import (
	"encoding/json"
	"github.com/imdario/mergo"
	"github.com/wen-bing/go-enterprise-web-toolkit/server"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func LoadConfig(env string, dir string) server.ServerConfig {
	var configName string
	if env == "production" {
		configName = "production.json"
	} else if env == "test" {
		configName = "test.json"
	} else {
		configName = "development.json"
	}

	localConfigFile := path.Join(dir, "local.json")
	configFile := path.Join(dir, configName)

	var localObj server.ServerConfig
	if _, err := os.Stat(localConfigFile); !os.IsNotExist(err) {
		raw, err2 := ioutil.ReadFile(localConfigFile)
		if err2 != nil {
			log.Printf("Read local.json error: %v", err)
			os.Exit(-1)
		}
		err2 = json.Unmarshal(raw, &localObj)
		if err2 != nil {
			log.Printf("local.json format error: %v", err)
			os.Exit(-1)
		}
	}

	var configObj server.ServerConfig
	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		raw, err2 := ioutil.ReadFile(configFile)
		if err2 != nil {
			log.Printf("Rread %s error: %v", configFile, err2)
			os.Exit(-1)
		}
		err2 = json.Unmarshal(raw, &configObj)
		if err2 != nil {
			log.Printf("%s format error: %v", configFile, err)
			os.Exit(-1)
		}
	}

	err := mergo.MergeWithOverwrite(&configObj, localObj)
	if err != nil {
		log.Printf("Merge config faileD: %v", err)
	}
	return configObj
}
