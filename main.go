// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build !plan9,!solaris
package main

import (
	"flag"
	"gosync/config"
	"gosync/dbsync"
	"gosync/firstrun"
	"gosync/fswatcher"
	"log"
	"net/http"
	"os"
)


func StartWebFileServer(cfg *config.Configuration) {
	log.Println("Starting web listener")
	var listenPort = ":" + cfg.ServerConfig.ListenPort
	for name, item := range cfg.Listeners {
		var section = "/" + name + "/"
		log.Println("Adding section listener: " + section + "| Serving directory: " + item.Directory)
		http.Handle(section, http.StripPrefix(section, http.FileServer(http.Dir(item.Directory))))
	}
	log.Printf("%v", http.ListenAndServe(listenPort, nil))
}

func main() {
	var ConfigFile string
	flag.StringVar(&ConfigFile, "config", "/etc/gosync/config.cfg",
    "Please provide the path to the config file, defaults to: /etc/gosync/config.cfg")
	flag.Parse()
	if _, err := os.Stat(ConfigFile); err == nil {
		log.Printf("Using %s as config file", ConfigFile)
		config.ReadConfigFromFile(ConfigFile)
		cfg := config.GetConfig()
		firstrun.InitialSync()
		for _, item := range cfg.Listeners {
			log.Println("Working with: " + item.Directory)
			go dbsync.DBCheckin(item.Directory, cfg)
			go fswatcher.SysPathWatcher(item.Directory)
		}
		StartWebFileServer(cfg)
	} else {
		log.Fatalf("Config file specified does not exist (%s)", ConfigFile)

	}

}
