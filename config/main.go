package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type Configuration struct {
	ServerConfig BaseConfig
	Database     Database
	Listeners    map[string]listener
}

type Database struct {
	Type string
	Dsn  string
}

type BaseConfig struct {
	ListenPort string `toml:"listen_port"`
}

type listener struct {
	Directory string
	Key       string
	Secret    string
}

func ReadConfigFromFile(configfile string) *Configuration {
	config_file, err := ioutil.ReadFile(configfile)
	if err != nil {
		panic(err.Error())
	}
	var conf Configuration
	_, err = toml.Decode(string(config_file), &conf)
	if err != nil {
		panic(err.Error())
	}
	return &conf
}
