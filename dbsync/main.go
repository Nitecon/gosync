package dbsync

import (
	"gosync/config"
	"gosync/dbsync/dbadapter"
	"gosync/fstools"
	"gosync/prototypes"
	"log"
	"time"
)

func DBInit() {
    cfg:= config.GetConfig()
	if cfg.Database.Type == "mysql" {
		//log.Println("MySQL database in use... checking tables")
		dbadapter.MySQLSetupTables(cfg)
	}
}

func InsertItem(table string, item fstools.FsItem) bool {
    cfg := config.GetConfig()
	var updateSuccess = true
	if cfg.Database.Type == "mysql" {
		//log.Println("MySQL database adapter selected")
		updateSuccess = dbadapter.MySQLInsertItem(cfg, table, item)
	}
	return updateSuccess
}

func DBCheckEmpty(table string) bool {
    cfg := config.GetConfig()
	var isEmpty = true
	if cfg.Database.Type == "mysql" {
		//log.Println("MySQL database adapter selected")
		isEmpty = dbadapter.MySQLCheckEmpty(cfg, table)
	}
	if isEmpty {
		log.Println("Database is EMPTY, starting creation")
	} else {
		log.Println("Using existing table: " + table)
	}
	return isEmpty
}

func DBFetchAll(table string) []prototypes.DataTable {
    cfg := config.GetConfig()
	if cfg.Database.Type == "mysql" {
		items := dbadapter.MySQLFetchAll(cfg, table)
		return items
	}
	return nil
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
				// @TODO: check that db knows Im alive.
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
