package dbadapter

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gosync/config"
	"gosync/fstools"
	"log"
	"os"
	"time"
)

func createTableQuery(table string) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id int(10) unsigned NOT NULL,
  path text COLLATE utf8_unicode_ci NOT NULL,
  filename varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  checksum varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  atime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  mtime timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  uid int(5) NOT NULL,
  gid int(5) NOT NULL,
  perms int(4) NOT NULL,
  host_updated varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  last_update timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
`, table)
}

type FsTable struct {
	Id         int    `id`
	Path       string `path`
	IsDir      int    `is_dir`
	Filename   string `filename`
	Checksum   string `checksum`
	Mtime      int    `mtime`
	Uid        int    `uid`
	Gid        int    `gid`
	Perms      string `perms`
	HostName   string `host_updated`
	LastUpdate int    `last_update`
}

func MySQLSetupTables(cfg *config.Configuration) {
	var db *sql.DB
	db = initDb(cfg)
	defer db.Close()
	log.Println("Database initialized")

	for key, _ := range cfg.Listeners {
		_, err := db.Query(createTableQuery(key))
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

	}

}

func MySQLInsertItem(cfg *config.Configuration, table string, item fstools.FsItem) bool {
	var db *sql.DB
	db = initDb(cfg)
	defer db.Close()

	var isDirectory = 0
	if item.IsDir {
		isDirectory = 1
	}

	hostname, _ := os.Hostname()
	row := &FsTable{
		Path:       item.Filename,
		IsDir:      isDirectory,
		Filename:   item.Filename,
		Checksum:   item.Checksum,
		Mtime:      item.Mtime,
		Uid:        item.Uid,
		Gid:        item.Gid,
		Perms:      item.Perms,
		HostName:   hostname,
		LastUpdate: int(time.Now().Unix()),
	}
	/*
		err := dbmap.Insert(row)
		if err != nil {
			checkErr(err, "Error occurred adding item to table: "+table)
			return false
		} else {
			return true
		}*/
	log.Printf("Stub in for adding data %v", row)
	return true
}

func MySQLCheckEmpty(cfg *config.Configuration, table string) bool {
	var db *sql.DB
	db = initDb(cfg)
	defer db.Close()

	var count int
	query := fmt.Sprintf("SELECT count(*) FROM %s", table)
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Fatalf("Critical error, cannot read from table %s : %v", table, err.Error())
	}
	var isEmpty = true
	if count > 0 {
		isEmpty = false
	}
	return isEmpty
}

func initDb(cfg *config.Configuration) *sql.DB {
	db, err := sql.Open("mysql", cfg.Database.Dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err.Error())
	}
	return db
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
