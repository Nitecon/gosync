package dbsync

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gosync/config"
	"gosync/fstools"
	"gosync/prototypes"
	"log"
	//"net/url"
	"os"
	"path"
	"time"
)

type MySQLDB struct {
	config *config.Configuration
	db     *sqlx.DB
}

func (my *MySQLDB) Insert(table string, item fstools.FsItem) bool {
	var isDirectory = 0
	if item.IsDir {
		isDirectory = 1
	}
	hostname, _ := os.Hostname()
	var keyExists = []prototypes.DataTable{}
	var query = fmt.Sprintf("SELECT id, path FROM %s WHERE %s='%s' LIMIT 1", table, "path", item.Filename)
	err := my.db.Select(&keyExists, query)
	if err != nil {
		log.Fatalf("Error checking for existence of key: %s, in table %s\n %+v", item.Filename, table, err)
	}
	tx := my.db.MustBegin()
	if len(keyExists) > 0 {
		rowId := fmt.Sprintf("%d",keyExists[0].Id)
		tx.MustExec("UPDATE "+table+" SET path=?, is_dir=?, filename=?, directory=?, checksum=?, atime=?, mtime=?, perms=?, host_updated=?, last_update=? WHERE id='"+rowId+"'",
        item.Filename,
        isDirectory,
        path.Base(item.Filename),
        path.Dir(item.Filename),
        item.Checksum,
        time.Now().UTC(),
        item.Mtime,
        item.Perms,
        hostname,
        time.Now().UTC(),
        )
		err = tx.Commit()
		checkErr(err, "Error inserting data")
	} else {
		tx.MustExec("INSERT INTO "+table+" (path, is_dir, filename, directory, checksum, atime, mtime, perms, host_updated, last_update) VALUES (?,?,?,?,?,?,?,?,?,?)",
			item.Filename,
			isDirectory,
			path.Base(item.Filename),
			path.Dir(item.Filename),
			item.Checksum,
			time.Now().UTC(),
			item.Mtime,
			item.Perms,
			hostname,
			time.Now().UTC())
		err = tx.Commit()
		checkErr(err, "Error inserting data")
	}

	return true
}

func (my *MySQLDB) CheckEmpty(table string) bool {
	var count int
	err := my.db.Get(&count, "SELECT count(*) FROM "+table+";")
	if err != nil {
		checkErr(err, "Error counting items in table: "+table)
	}
	var isEmpty = true
	if count > 0 {
		isEmpty = false
	}
	return isEmpty
}

func (my *MySQLDB) FetchAll(table string) []prototypes.DataTable {
	dTable := []prototypes.DataTable{}
	query := "SELECT path, is_dir, checksum, mtime, perms, host_updated FROM " + table + " ORDER BY last_update ASC"
	err := my.db.Select(&dTable, query)
	checkErr(err, "Error occurred getting file details for: "+table)
	return dTable
}

func (my *MySQLDB) CheckIn(table string) []prototypes.DataTable{
    hostname, _ := os.Hostname()
    dTable := []prototypes.DataTable{}
    query := fmt.Sprintf("SELECT path, is_dir, checksum, mtime, perms, host_updated FROM %s where host_updated != '%s' ORDER BY last_update ASC", table, hostname)
    //log.Printf("Executing: %s", query)
    err := my.db.Select(&dTable, query)
    checkErr(err, "Error occurred getting file details for: "+table)
    //log.Printf("Changed Items: %+v", dTable)
    return dTable

}

func (my *MySQLDB) CreateDB() {

	log.Println("Database initialized")
	for key, _ := range my.config.Listeners {
		my.db.MustExec(createTableQuery(key))
	}

}

func (my *MySQLDB) Close() error {
	return my.db.Close()
}

func (my *MySQLDB) initDB() {
	tempdb, err := sqlx.Connect("mysql", my.config.Database.Dsn+"&parseTime=True")
	if err != nil {
		log.Println(err.Error())
	}
	my.db = tempdb
}

func checkExists(db *sqlx.DB, table, key, val string) bool {
	var keyExists = 0
	var query = fmt.Sprintf("SELECT 1 FROM ? WHERE ?='?' LIMIT 1", table, key, val)
	err := db.Get(&keyExists, query)
	if err != nil {
		checkErr(err, "Error checking existence of ("+key+") in table: "+table)
	}
	if keyExists > 0 {
		return true
	} else {
		return false
	}
	return false

}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func createTableQuery(table string) string {
	var createStmt = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
  path varchar(4096) COLLATE utf8_unicode_ci NOT NULL,
  is_dir int NOT NULL,
  filename varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  directory varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  checksum varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  atime timestamp NOT NULL,
  mtime timestamp NOT NULL,
  perms varchar(12) COLLATE utf8_unicode_ci NOT NULL,
  host_updated varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  last_update timestamp default now() NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
`, table)

	return createStmt
}
