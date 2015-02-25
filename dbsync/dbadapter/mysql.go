package dbadapter

import (
	//_ "github.com/lib/pq"
	//"database/sql"
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

func MySQLSetupTables(cfg *config.Configuration) {
	var db *sqlx.DB
	db = initDb(cfg)
	defer db.Close()
	log.Println("Database initialized")

	for key, _ := range cfg.Listeners {
		db.MustExec(createTableQuery(key))
	}

}

func MySQLInsertItem(cfg *config.Configuration, table string, item fstools.FsItem) bool {
	var db *sqlx.DB
	db = initDb(cfg)
	defer db.Close()

	var isDirectory = 0
	if item.IsDir {
		isDirectory = 1
	}
	hostname, _ := os.Hostname()
	tx := db.MustBegin()

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

func MySQLFetchAll(cfg *config.Configuration, table string) []prototypes.DataTable {
	var db *sqlx.DB
	db = initDb(cfg)
	defer db.Close()

	dTable := []prototypes.DataTable{}
	query := "SELECT path, is_dir, checksum, mtime, perms, host_updated FROM " + table + " ORDER BY last_update ASC"
	//log.Println("Executing scan all items...")
	err := db.Select(&dTable, query)
	checkErr(err, "Error occurred getting file details for: "+table)
	//log.Println("Executing scan all items... COMPLETE")
	return dTable
}

func MySQLCheckEmpty(cfg *config.Configuration, table string) bool {
	var db *sqlx.DB
	db = initDb(cfg)
	defer db.Close()

	var count int

	err := db.Get(&count, "SELECT count(*) FROM "+table+";")
	if err != nil {
		checkErr(err, "Error counting items in table: "+table)
	}
	var isEmpty = true
	if count > 0 {
		isEmpty = false
	}
	return isEmpty
}

func initDb(cfg *config.Configuration) *sqlx.DB {
	/*db, err := sql.Open("mysql", cfg.Database.Dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err.Error())
	}
	*/
	db, err := sqlx.Connect("mysql", cfg.Database.Dsn+"&parseTime=True")
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
