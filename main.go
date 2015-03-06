// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build !plan9,!solaris
package main

import (
	"flag"
    "gosync/nodeinfo"
    "gosync/replicator"
	"gosync/fswatcher"

	"net/http"
    "gosync/utils"
	"os"
)


func StartWebFileServer(cfg *utils.Configuration) {
    nodeinfo.Initialize()
    utils.WriteLn("Starting web listener")
    nodeinfo.SetAlive()

	var listenPort = ":" + cfg.ServerConfig.ListenPort
	for name, item := range cfg.Listeners {
		var section = "/" + name + "/"
        utils.WriteLn("Adding section listener: " + section + "| Serving directory: " + item.Directory)
		http.Handle(section, http.StripPrefix(section, http.FileServer(http.Dir(item.Directory))))
	}
    utils.LogWriteF("%v", http.ListenAndServe(listenPort, nil))
}

func main() {
	var ConfigFile string
	flag.StringVar(&ConfigFile, "config", "/etc/gosync/config.cfg",
    "Please provide the path to the config file, defaults to: /etc/gosync/config.cfg")
	flag.Parse()
	if _, err := os.Stat(ConfigFile); !utils.Check(err, 404, "No config file specified") {
        utils.ReadConfigFromFile(ConfigFile)
		cfg := utils.GetConfig()
        replicator.InitialSync()
		for _, item := range cfg.Listeners {
            utils.WriteLn("Working with: " + item.Directory)
			go replicator.CheckIn(item.Directory)
			go fswatcher.SysPathWatcher(item.Directory)
		}
		StartWebFileServer(cfg)
	} else {
        utils.LogWriteF("Config file specified does not exist (%s)", ConfigFile)
	}

}
