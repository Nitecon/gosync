package dbadapter

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gosync/config"
	"gosync/fstools"
	"log"
	"net/url"
	"os"
	"path"
	"time"
)

func createTableQuery(table string) string {
	var createStmt = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
  path text COLLATE utf8_unicode_ci NOT NULL,
  is_dir int NOT NULL,
  filename varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  directory varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  checksum varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  atime int(10) NOT NULL,
  mtime int(10) NOT NULL,
  uid int(5) NOT NULL,
  gid int(5) NOT NULL,
  perms varchar(12) COLLATE utf8_unicode_ci NOT NULL,
  host_updated varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  last_update int(10) NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
	
`, table)

	return createStmt
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
	query := fmt.Sprintf("INSERT INTO %s (path, is_dir, filename, directory, checksum, atime, mtime, uid, gid, perms, host_updated, last_update) VALUES (\"%s\", %d, \"%s\", \"%s\", \"%s\",%d, %d, %d, %d, \"%s\", \"%s\", %d )",
		table,
		url.QueryEscape(item.Filename),
		isDirectory,
		url.QueryEscape(path.Base(item.Filename)),
		url.QueryEscape(path.Dir(item.Filename)),
		item.Checksum,
		time.Now().Unix(),
		item.Mtime,
		item.Uid,
		item.Gid,
		item.Perms,
		hostname,
		time.Now().Unix())
	log.Printf("Executing Query: %s", query)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Critical error, cannot insert into table %s : %v", table, err.Error())
		return false
	} else {
		return true
	}
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
