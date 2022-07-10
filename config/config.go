package config

import (
	"gopkg.in/ini.v1"
	"fmt"
	"os"
)


type Config struct {
	BucketName string `ini:"bucket"`
	AccessKeyID string `ini:"aws_access_key"`
	SecretAccessKey string `ini:"aws_secret_access_key"`
	Region string `ini:"region"`
	EndpointUrl string `ini:"endpoint_url"`
	PathToConfig string
}


func LoadConfigIni(pathToFile string) (config Config, err error) {
	cfg, err := ini.Load(pathToFile)
	if err != nil {
		fmt.Println("could not open the config file")
		os.Exit(1)
	}
	keys := []string{"bucket", "aws_access_key_id", "aws_secret_access_key", "region", "endpoint_url"}
	CheckIfAllKeysArePresent(keys, cfg)
	err = cfg.Section("default").MapTo(&config)
	return
}


func CheckIfAllKeysArePresent(keys []string, cfg *ini.File) {
	for _, key := range keys {
		if yes := cfg.Section("default").HasKey(key); yes != true {
			fmt.Println("Not all keys are present in the config file!")
			os.Exit(1)
		}
	}
}
