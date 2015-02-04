package dbadapter

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"gosync/config"
	"gosync/fstools"
	"log"
	"os"
	"time"
)

/*
CREATE TABLE IF NOT EXISTS `backups` (
`id` int(10) unsigned NOT NULL,
  `path` text COLLATE utf8_unicode_ci NOT NULL,
  `filename` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `checksum` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `atime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `mtime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `uid` int(5) NOT NULL,
  `gid` int(5) NOT NULL,
  `perms` int(4) NOT NULL,
  `host_updated` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `last_update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
*/

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
	dbmap := initDb(cfg)
	defer dbmap.Db.Close()
	log.Println("Database initialized")
	for key, _ := range cfg.Listeners {
		//table := dbmap.AddTableWithName(FsTable{}, key).SetKeys(true, "Id")
		dbmap.AddTableWithName(FsTable{}, key).SetKeys(true, "Id")
		err := dbmap.CreateTablesIfNotExists()
		checkErr(err, "Create tables failed")
		count, err := dbmap.SelectInt("select count(*) from " + key)
		checkErr(err, "select count(*) failed")
		if count < 1 {
			log.Println("New table build starting for: " + key)
		}

	}

}

func MySQLInsertItem(cfg *config.Configuration, table string, item fstools.FsItem) bool {
	dbmap := initDb(cfg)
	defer dbmap.Db.Close()
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
	err := dbmap.Insert(row)
	if err != nil {
		checkErr(err, "Error occurred adding item to table: "+table)
		return false
	} else {
		return true
	}

}

func MySQLCheckEmpty(cfg *config.Configuration, table string) bool {
	dbmap := initDb(cfg)
	defer dbmap.Db.Close()
	count, err := dbmap.SelectInt("select count(*) from " + table)
	checkErr(err, "select count(*) failed on "+table)
	var isEmpty = true
	if count > 0 {
		isEmpty = false
	}
	return isEmpty
}

func initDb(cfg *config.Configuration) *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	//root:pw@unix(/tmp/mysql.sock)/myDatabase?loc=Local
	//user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true

	//log.Println("CONNECT_STRING:" + cfg.Database.Dsn)
	db, err := sql.Open("mysql", cfg.Database.Dsn)
	if err != nil {
		panic(err)
	}
	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func main() {
	db, err := sql.Open("mysql", "user:password@/database")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO squareNum VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT squareNumber FROM squarenum WHERE number = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	// Insert square numbers for 0-24 in the database
	for i := 0; i < 25; i++ {
		_, err = stmtIns.Exec(i, (i * i)) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	var squareNum int // we "scan" the result in here

	// Query the square-number of 13
	err = stmtOut.QueryRow(13).Scan(&squareNum) // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of 13 is: %d", squareNum)

	// Query another number.. 1 maybe?
	err = stmtOut.QueryRow(1).Scan(&squareNum) // WHERE number = 1
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of 1 is: %d", squareNum)
}
