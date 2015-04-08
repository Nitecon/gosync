package main


import (
    "flag"
    log "github.com/cihub/seelog"
    "github.com/Nitecon/gosync/webserver"
    "os"
    "path/filepath"
    "github.com/Nitecon/gosync/config"
    "github.com/Nitecon/gosync/setup"
)

func getLoggerConfig() string {
    cfg := config.GetConfig()
    var loggerConfig = ""
    if cfg.ServerConfig.LogLocation != "stdout" {
        if _, err := os.Stat(filepath.Dir(cfg.ServerConfig.LogLocation)); os.IsNotExist(err) {
            os.Mkdir(filepath.Dir(cfg.ServerConfig.LogLocation), 0775)
        }
        loggerConfig = `<seelog type="asynctimer" asyncinterval="1000">
    <outputs formatid="main">
        <filter levels="`+cfg.ServerConfig.LogLevel+`">
          <file path="` + cfg.ServerConfig.LogLocation + `" />
        </filter>
    </outputs>
    <formats>
        <format id="main" format="%Date %Time [%LEVEL] %Msg%n"/>
    </formats>
    </seelog>`
    } else {
        loggerConfig = `<seelog type="asynctimer" asyncinterval="1000">
    <outputs formatid="main">
        <console/>
    </outputs>
    <formats>
        <format id="main" format="%Date %Time [%LEVEL] %Msg (%RelFile:%Func)%n"/>
    </formats>
    </seelog>`
    }

    return loggerConfig
}



func init() {
    var ConfigFile string
    flag.StringVar(&ConfigFile, "config", "config.cfg.example",
    "Please provide the path to the config file, defaults to: /etc/gosync/config.cfg")
    flag.Parse()
    if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
        log.Criticalf("Configuration file does not exist or cannot be loaded: (%s)", ConfigFile)
        os.Exit(1)
    } else {
        config.ReadConfigFromFile(ConfigFile)
    }

    logger, err := log.LoggerFromConfigAsString(getLoggerConfig())

    if err == nil {
        log.ReplaceLogger(logger)
    }

}

func main() {

    cfg := config.GetConfig()
    setup.FsVerify()
    setup.DbVerify()
    for _, item := range cfg.Listeners {
        log.Info("Working with: " + item.Directory)
        //go replicator.CheckIn(item.Directory)
        //fswatcher.SysPathWatcher(item.Directory)
    }
    webserver.StartWebFileServer()
}
