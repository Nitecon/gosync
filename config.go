package gosync
import (
    "github.com/BurntSushi/toml"
    "io/ioutil"
    "sync"
    "log"
)

var (
    config     *Configuration
    configLock = new(sync.RWMutex)
)

func ReadConfigFromFile(configfile string) {
    config_file, err := ioutil.ReadFile(configfile)
    if err != nil {
        log.Fatalf("Could not read configuration file: %s\n\n%s", configfile, err.Error())
    }
    tempConf := new(Configuration)
    _, err = toml.Decode(string(config_file), &tempConf)
    if err != nil {
        log.Fatalf("Invalid items in configuration file:\n%s", err.Error())
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


