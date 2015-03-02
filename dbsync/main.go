package dbsync

import (
	"gosync/config"
	"gosync/fstools"
	"gosync/prototypes"
	"log"
	"time"
)

var (
    storage Datastore
)

type Datastore interface {
    Insert(table string, item fstools.FsItem) bool
    CheckEmpty(table string) bool
    FetchAll(table string) []prototypes.DataTable
    CheckIn(path string) []prototypes.DataTable
    CreateDB()
    Close() error // call this method when you want to close the connection
    initDB()
}

func setStorageEngine() {
    cfg := config.GetConfig()
    var engine = cfg.Database.Type
    switch engine {
        case "mysql":
        storage = &MySQLDB{config:cfg}
        storage.initDB()
        //case "pgsql":
        //storage = &PgSQLDB{}
    }
}

func Insert(table string, item fstools.FsItem) bool {
    setStorageEngine()
    return storage.Insert(table, item)
}

func CheckEmpty(table string) bool{
    setStorageEngine()
    empty := storage.CheckEmpty(table)
    if empty {
        log.Println("Database is EMPTY, starting creation")
    }else{
        log.Println("Using existing table: " + table)
    }
    return empty
}

func FetchAll(table string) []prototypes.DataTable{
    setStorageEngine()
    return storage.FetchAll(table)
}

func CheckIn(path string){
    log.Println("Starting db checking background script: " + path)
    ticker := time.NewTicker(10 * time.Second)
    quit := make(chan struct{})
    go func() {
        for {
            select {
            case <-ticker.C:
                log.Println("Checking all changed stuff in db for: " + path)
                listener := getListener(path)
                storage.CheckIn(listener)
            // @TODO: check that db knows Im alive.
            case <-quit:
                ticker.Stop()
                return
            }
        }
    }()
}

func CreateDB(){
    setStorageEngine()
    storage.CreateDB()
}

func getListener(dir string) string {
    cfg := config.GetConfig()
    var listener = ""
    for lname, ldata := range cfg.Listeners{
        if ldata.Directory == dir{
            listener = lname
        }
    }
    return listener
}