package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

var (
	config         Config
	once           sync.Once
	configFilePath = "./configs/default.json"
)

type Config struct {
	HTTP           HTTPConf           `json:"http"`
	CacheImagesLRU CacheImagesLRUConf `json:"cacheImagesLru"`
}

type HTTPConf struct {
	Port string `json:"port"`
}

type CacheImagesLRUConf struct {
	Capacity int    `json:"capacity"`
	Dir      string `json:"dir"`
}

func SetFilePath(filePath string) {
	configFilePath = filePath
}

func Get() *Config {
	once.Do(func() {
		jsonString, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			log.Fatal(err)
		}

		if err = json.Unmarshal(jsonString, &config); err != nil {
			log.Fatal(err)
		}
	})
	return &config
}
