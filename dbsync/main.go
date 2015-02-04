package dbsync

import (
	"gosync/config"
	"gosync/dbsync/dbadapter"
	"gosync/fstools"
	"log"
	"time"
)

func DBInit(cfg *config.Configuration) {
	if cfg.Database.Type == "mysql" {
		log.Println("MySQL database in use... checking tables")
		dbadapter.MySQLSetupTables(cfg)
	}
}

func InsertItem(cfg *config.Configuration, table string, item fstools.FsItem) bool {
	var updateSuccess = true
	if cfg.Database.Type == "mysql" {
		log.Println("MySQL database adapter selected")
		updateSuccess = dbadapter.MySQLInsertItem(cfg, table, item)
	}
	return updateSuccess
}

func DBCheckEmpty(cfg *config.Configuration, table string) bool {
	var isEmpty = true
	if cfg.Database.Type == "mysql" {
		log.Println("MySQL database adapter selected")
		isEmpty = dbadapter.MySQLCheckEmpty(cfg, table)
	}
	log.Printf("Database is empty?: %x", isEmpty)
	return isEmpty
}

func DBCheckin(path string, cfg *config.Configuration) {
	log.Println("Starting db checking background script: " + path)
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Checking all changed stuff in db for: " + path)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
