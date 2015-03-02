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
	tx := my.db.MustBegin()
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
	err := tx.Commit()
	checkErr(err, "Error inserting data")
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
	//log.Println("Executing scan all items...")
	err := my.db.Select(&dTable, query)
	checkErr(err, "Error occurred getting file details for: "+table)
	//log.Println("Executing scan all items... COMPLETE")
	return dTable
}

func (my *MySQLDB) CheckIn(table string) {

}

func (my *MySQLDB) DBInit() {

	log.Println("Database initialized")
	for key, _ := range my.config.Listeners {
		my.db.MustExec(createTableQuery(key))
	}

}

func (my *MySQLDB) Close() error {
	return my.db.Close()
}

func (my *MySQLDB) initDB() {

	//cfg := my.config.config.GetConfig()
	log.Printf("%+v", my.config)
	log.Println("Starting DB Initialization")
	log.Println("Getting Config")

	log.Println("Getting DB Connection")
	tempdb, err := sqlx.Connect("mysql", my.config.Database.Dsn+"&parseTime=True")
	if err != nil {
		log.Println(err.Error())
	}
	my.db = tempdb
	//return db
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func createTableQuery(table string) string {
	var createStmt = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
  path text COLLATE utf8_unicode_ci NOT NULL,
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
