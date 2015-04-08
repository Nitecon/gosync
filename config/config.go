package config

import (
    "github.com/BurntSushi/toml"
    "io/ioutil"
    "sync"
   // log "github.com/cihub/seelog"
)

var (
    config     *Configuration
    configLock = new(sync.RWMutex)
)

type Configuration struct {
    ServerConfig ServerConf `toml:"ServerConfig"`
    StorageSetup StorageConf `toml:"RemoteStorage"`
    Database     Database
    Listeners    map[string]Listener
}

type Database struct {
    Type string `toml:"type"`
    User string `toml:"user"`
    Pass string `toml:"pass"`
    Host string `toml:"host"`
    Port string `toml:"port"`
    DBase string `toml:"dbase"`
}

type StorageConf struct {
    Key    string `toml:"key"`
    Secret string `toml:"secret"`
    Region string `toml:"region"`
}

type ServerConf struct {
    ListenPort  string `toml:"listen_port"`
    RescanTime  int    `toml:"rescan"`
    StorageType string `toml:"storagetype"`
    LogLocation string `toml:"log_location"`
    LogLevel    string `toml:"log_level"`
}

type Listener struct {
    Directory   string
    Uid         int
    Gid         int
    Bucket      string `toml:"bucket"`
    BasePath    string `toml:"basepath"`
}

func ReadConfigFromFile(configfile string) {
    config_file, err := ioutil.ReadFile(configfile)
    if err != nil {
        panic(err.Error())
    }
    tempConf := new(Configuration)
    _, err = toml.Decode(string(config_file), &tempConf)
    if err != nil {
        panic(err.Error())
    }
    configLock.Lock()
    config = tempConf
    configLock.Unlock()
}

func GetConfig() *Configuration {
    configLock.RLock()
    defer configLock.RUnlock()
    return config
}