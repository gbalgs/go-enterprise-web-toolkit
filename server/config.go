package server

import (
	"encoding/json"
	"github.com/fatih/structs"
	"github.com/wen-bing/go-enterprise-web-toolkit/core"
	"io/ioutil"
	"log"
	"os"
)

type ServerConfig struct {
	Port int           `json:"port"`
	DB   core.DBConfig `json:"db"`
}

func initConfig(env string) ServerConfig {
	var configName string
	if env == "production" {
		configName = "production.json"
	} else if env == "test" {
		configName = "test.json"
	} else {
		configName = "development.json"
	}

	localConfigFile := "./configs/local.json"
	configFile := "./configs/" + configName

	var localObj ServerConfig
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

	var configObj ServerConfig
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

	//merge local to override config object
	//TODO
	//refacto to use reflect
	if !structs.IsZero(localObj) {
		if localObj.Port != 0 {
			configObj.Port = localObj.Port
		}

		if !structs.IsZero(localObj.DB) {
			if localObj.DB.Port != 0 {
				configObj.DB.Port = localObj.DB.Port
			}

			if localObj.DB.Host != "" {
				configObj.DB.Host = localObj.DB.Host
			}
			if localObj.DB.Type != "" {
				configObj.DB.Type = localObj.DB.Type
			}
			if localObj.DB.Name != "" {
				configObj.DB.Name = localObj.DB.Name
			}
			if localObj.DB.User != "" {
				configObj.DB.User = localObj.DB.User
			}
			if localObj.DB.Password != "" {
				configObj.DB.Password = localObj.DB.Password
			}
		}
	}

	return configObj

}
