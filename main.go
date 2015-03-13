// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build !plan9,!solaris
package main

import (
	"flag"
	"gosync/fswatcher"
	"gosync/nodeinfo"
	"gosync/replicator"
	"gosync/utils"
	"log"
	"net/http"
	"os"
)

func StartWebFileServer(cfg *utils.Configuration) {
	nodeinfo.Initialize()
	log.Println("Starting Web File server and setting node as active")
	nodeinfo.SetAlive()

	var listenPort = ":" + cfg.ServerConfig.ListenPort
	for name, item := range cfg.Listeners {
		var section = "/" + name + "/"
		log.Printf("Adding section listener: %s, to serve directory: %s", section, item.Directory)
		http.Handle(section, http.StripPrefix(section, http.FileServer(http.Dir(item.Directory))))
	}
	log.Fatal(http.ListenAndServe(listenPort, nil))
}

func init() {
	var ConfigFile string
	flag.StringVar(&ConfigFile, "config", "/etc/gosync/config.cfg",
		"Please provide the path to the config file, defaults to: /etc/gosync/config.cfg")
	flag.Parse()
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		log.Fatalf("Configuration file does not exist or cannot be loaded:\n (%s)", ConfigFile)
	} else {
		utils.ReadConfigFromFile(ConfigFile)
	}
    cfg := utils.GetConfig()
    flag.Set("log_dir", cfg.ServerConfig.LogLocation)
    flag.Parse()
}

func main() {

	cfg := utils.GetConfig()
	replicator.InitialSync()
	for _, item := range cfg.Listeners {
		utils.WriteLn("Working with: " + item.Directory)
		go replicator.CheckIn(item.Directory)
		go fswatcher.SysPathWatcher(item.Directory)
	}
	StartWebFileServer(cfg)

}
