package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"net"
	"sync"
)

var (
	config     *Configuration
	configLock = new(sync.RWMutex)
)

type Configuration struct {
	ServerConfig BaseConfig
	S3Config     StorageS3 `toml:"StorageS3"`
	GDConfig     StorageGDrive
	Database     Database
	Listeners    map[string]Listener
}

type Database struct {
	Type string
	Dsn  string
}

type StorageS3 struct {
	Key    string `toml:"key"`
	Secret string `toml:"secret"`
	Region string `toml:"region"`
}

type StorageGDrive struct {
	Key    string
	Secret string
	Region string
}

type BaseConfig struct {
	ListenPort  string `toml:"listen_port"`
	RescanTime  int    `toml:"rescan"`
	StorageType string `toml:"storagetype"`
}

type Listener struct {
	Directory   string
	Key         string
	Secret      string
	Uid         int
	Gid         int
	StorageType string `toml:"storagetype"`
	Bucket      string `toml:"bucket"`
	BasePath    string `toml:"basepath"`
}


func ReadConfigFromFile(configfile string) {
	log.Println("Loading configuration from disk...")
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
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf("Error accessing network interfaces: %v", err.Error())
	}
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Printf("Error occurred getting local IP address: %v", err.Error())
		}
		// handle err
		for _, addr := range addrs {
			log.Println(addr)
			/*switch v := addr.(type) {
			case *net.IPAddr:
				// process IP address
			}*/

		}
	}
}
