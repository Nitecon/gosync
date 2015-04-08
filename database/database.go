package database
import (
    "github.com/Nitecon/gosync/config"
    "github.com/Nitecon/gosync/filesystem"
    log "github.com/cihub/seelog"
)

var (
    dbstore Datastore
)

type Datastore interface {
    Insert(table string, item filesystem.FileData) bool
    Update(table string, item filesystem.FileData) bool
    Remove(table, path string) bool
    FetchAll(table string) ([]filesystem.FileData, error)
    FetchOne(listener, path string) (filesystem.FileData, error)
    Exists(listener string, path filesystem.FileData) bool
    Setup()
}

func setdbstoreEngine() {
    cfg := config.GetConfig()
    log.Infof("Database engine: %s", cfg.Database.Type)
    var engine = cfg.Database.Type
    switch engine {
        case "mysql":
        dbstore = &MySQLDB{config: cfg}
    }
}

func Add(table string, item filesystem.FileData) (updated bool) {
    setdbstoreEngine()
    if dbstore.Exists(table, item) {
        updated = dbstore.Update(table, item)
    }else{
        updated = dbstore.Insert(table, item)
    }
    return
}

func Exists(table string, fi filesystem.FileData) bool {
    setdbstoreEngine()
    return dbstore.Exists(table, fi)
}

func SetupDB(){
    setdbstoreEngine()
    dbstore.Setup()
}