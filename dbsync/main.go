package dbsync

import (
	"gosync/config"
	"gosync/fstools"
	"gosync/prototypes"
	"log"
)

var (
	dbstore Datastore
)

type Datastore interface {
	Insert(table string, item fstools.FsItem) bool
	CheckEmpty(table string) bool
	FetchAll(table string) []prototypes.DataTable
	CheckIn(listener string) ([]prototypes.DataTable, error)
	CreateDB()
	Close() error // call this method when you want to close the connection
	initDB()
}

func setdbstoreEngine() {
	cfg := config.GetConfig()
	var engine = cfg.Database.Type
	switch engine {
	case "mysql":
		dbstore = &MySQLDB{config: cfg}
		dbstore.initDB()
		//case "pgsql":
		//dbstore = &PgSQLDB{}
	}
}

func Insert(table string, item fstools.FsItem) bool {
	setdbstoreEngine()
	return dbstore.Insert(table, item)
}

func CheckEmpty(table string) bool {
	setdbstoreEngine()
	empty := dbstore.CheckEmpty(table)
	if empty {
		log.Println("Database is EMPTY, starting creation")
	} else {
		log.Println("Using existing table: " + table)
	}
	return empty
}

func FetchAll(table string) []prototypes.DataTable {
	setdbstoreEngine()
	return dbstore.FetchAll(table)
}

func CheckIn(listener string) ([]prototypes.DataTable, error) {
	log.Println("Starting db checking background script for: " + listener)
	data,err := dbstore.CheckIn(listener)
    return data, err

}

func CreateDB() {
	setdbstoreEngine()
	dbstore.CreateDB()
}
