package datastore

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gosync/utils"
	"os"
	"path"
	"time"
)

type MySQLDB struct {
	config *utils.Configuration
	db     *sqlx.DB
}

func (my *MySQLDB) Insert(table string, item utils.FsItem) bool {
	var isDirectory = 0
	if item.IsDir {
		isDirectory = 1
	}
	hostname, _ := os.Hostname()
	var keyExists = []utils.DataTable{}
    relFileName := string(utils.GetRelativePath(table, item.Filename))
	var query = fmt.Sprintf("SELECT id, path FROM %s WHERE path='%s' LIMIT 1", table, relFileName)

	err := my.db.Select(&keyExists, query)
	if err != nil {
        utils.LogWriteF("Error checking for existence of key: %s, in table %s\n %+v", utils.GetRelativePath(table, item.Filename), table, err)
	}
	tx := my.db.MustBegin()
	if len(keyExists) > 0 {
		rowId := fmt.Sprintf("%d", keyExists[0].Id)
		tx.MustExec("UPDATE "+table+" SET path=?, is_dir=?, filename=?, directory=?, checksum=?, atime=?, mtime=?, perms=?, host_updated=?, host_ips=?, last_update=? WHERE id='"+rowId+"'",
			utils.GetRelativePath(table, item.Filename),
			isDirectory,
			path.Base(item.Filename),
			path.Dir(item.Filename),
			utils.GetMd5Checksum(item.Filename),
			time.Now().UTC(),
			item.Mtime,
			item.Perms,
			hostname,
			utils.GetLocalIp(),
			time.Now().UTC(),
		)
		err = tx.Commit()
        utils.Check(err, 500, "Error updating data...")
	} else {
		tx.MustExec("INSERT INTO "+table+" (path, is_dir, filename, directory, checksum, atime, mtime, perms, host_updated, host_ips, last_update) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
			utils.GetRelativePath(table, item.Filename),
			isDirectory,
			path.Base(item.Filename),
			path.Dir(item.Filename),
			utils.GetMd5Checksum(item.Filename),
			time.Now().UTC(),
			item.Mtime,
			item.Perms,
			hostname,
            utils.GetLocalIp(),
			time.Now().UTC())
		err = tx.Commit()
		utils.Check(err, 500, "Error inserting data...")
	}
	return true
}

func (my *MySQLDB) CheckEmpty(table string) bool {
	var count int
	err := my.db.Get(&count, "SELECT count(*) FROM "+table+";")
	utils.ErrorCheckF(err, 500, "Error counting items in table (%s)", table)

	var isEmpty = true
	if count > 0 {
		isEmpty = false
	}
	return isEmpty
}

func (my *MySQLDB) FetchAll(table string) []utils.DataTable {
	dTable := []utils.DataTable{}
	query := "SELECT path, is_dir, checksum, mtime, perms, host_updated, host_ips FROM " + table + " ORDER BY is_dir DESC"
	err := my.db.Select(&dTable, query)
    utils.ErrorCheckF(err, 500, "Could not fetch items from database for %s", table)
	return dTable
}

func (my *MySQLDB) CheckIn(table string) ([]utils.DataTable, error) {
	hostname, _ := os.Hostname()
	dTable := []utils.DataTable{}
	query := fmt.Sprintf("SELECT path, is_dir, checksum, mtime, perms, host_updated, host_ips FROM %s where host_updated != '%s' ORDER BY is_dir DESC, last_update DESC", table, hostname)
	err := my.db.Select(&dTable, query)
	return dTable, err

}

func (my *MySQLDB) GetOne(listener, path string) (utils.DataTable, error){
    dItem := utils.DataTable{}
    query := fmt.Sprintf("SELECT path, is_dir, checksum, mtime, perms, host_ips FROM %s where path='%s'", listener, path)
    err := my.db.Get(&dItem, query)
    return dItem, err
}

func (my *MySQLDB) Remove(table string, item utils.FsItem) bool{
    return true
}

func (my *MySQLDB) CreateDB() {
    utils.WriteLn("Database initialized")
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
        utils.WriteLn(err.Error())
	}
	my.db = tempdb
}


func createTableQuery(table string) string {
	var createStmt = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
  path varchar(4096) COLLATE utf8_unicode_ci NOT NULL,
  is_dir int NOT NULL,
  filename varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  directory varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  checksum varchar(2048) COLLATE utf8_unicode_ci NOT NULL,
  atime timestamp NOT NULL,
  mtime timestamp NOT NULL,
  perms varchar(12) COLLATE utf8_unicode_ci NOT NULL,
  host_updated varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  host_ips varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  last_update timestamp default now() NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
`, table)

	return createStmt
}
