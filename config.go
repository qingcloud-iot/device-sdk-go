package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var conf *Config

type Config struct {
	Device struct {
		Token string `yaml:"token"`
	}
	Mqttbroker struct {
		Addr string `yaml:"addr"`
	}
	Taskinfo struct {
		MiddleCredentials []string `yaml:"middle_credentials"`
	}
}

func InitConfig() *Config {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		panic(err)
	}
	return conf
}
