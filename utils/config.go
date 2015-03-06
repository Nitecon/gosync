package utils

import (
    "github.com/BurntSushi/toml"
    "io/ioutil"
    "net"
    "sync"
    "os"
)

var (
    config     *Configuration
    configLock = new(sync.RWMutex)
)

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

func GetLocalIp() {
    LogWriteF("Ip address: %s", net.LookupIP(os.Hostname()))
    ifaces, err := net.Interfaces()
    if err != nil {
        LogWriteF("Error accessing network interfaces: %v", err.Error())
    }
    // handle err
    for _, i := range ifaces {
        addrs, err := i.Addrs()
        if err != nil {
            LogWriteF("Error occurred getting local IP address: %v", err.Error())
        }
        // handle err
        for _, addr := range addrs {
            LogWriteF("%s",addr.Network())
            LogWriteF("%s",addr.String())
            /*switch v := addr.(type) {
            case *net.IPAddr:
                // process IP address
            }*/

        }
    }
}
