package webserver

import (
    log "github.com/cihub/seelog"
    "net/http"
    "github.com/Nitecon/gosync/config"
)


func StartWebFileServer() {
    //nodeinfo.Initialize()
    //log.Info("Starting Web File server and setting node as active")
    //nodeinfo.SetAlive()
    cfg := config.GetConfig()
    var listenPort = ":" + cfg.ServerConfig.ListenPort
    for name, item := range cfg.Listeners {
        var section = "/" + name + "/"
        log.Infof("Adding section listener: %s, to serve directory: %s", section, item.Directory)
        http.Handle(section, http.StripPrefix(section, http.FileServer(http.Dir(item.Directory))))
    }
    log.Debug(http.ListenAndServe(listenPort, nil))
}