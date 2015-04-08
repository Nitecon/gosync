package database

import (
	"fmt"
	"github.com/Nitecon/gosync/config"
	"github.com/Nitecon/gosync/filesystem"
	log "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"path"
	"strings"
    "time"
)

type MySQLDB struct {
	config *config.Configuration
	db     *sqlx.DB
}

func (my *MySQLDB) Insert(table string, item filesystem.FileData)  bool {
	log.Debugf("MySQL Insert on table %s", table)
	my.db = my.GetConn()
	defer my.db.Close()
	tx := my.db.MustBegin()
	tx.MustExec("INSERT INTO "+table+" (is_dir, filename, directory, checksum, atime, mtime, perms, host_updated, host_ips, last_update) VALUES (?,?,?,?,?,?,?,?,?,?)",
		item.IsDir,
		item.Filename,
		item.Directory,
		item.Checksum,
		time.Now().UTC(),
		item.Mtime,
		item.Perms,
		item.HostUpdated,
		item.HostIPs,
		time.Now().UTC())
	err := tx.Commit()
	if err != nil {
		log.Infof("Error inserting data to DB table %s, for file %s\n%s\n%+v", table, item.Filename, err.Error())
		return false
	}
	return true
}

func (my *MySQLDB) Update(table string, item filesystem.FileData)  bool {
	log.Debugf("MySQL Update on table %s", table)
	my.db = my.GetConn()
	defer my.db.Close()
	tx := my.db.MustBegin()
	tx.MustExec("UPDATE "+table+" SET is_dir=?, checksum=?, atime=?, mtime=?, perms=?, host_updated=?, host_ips=?, last_update=? WHERE filename=? AND directory=?",
	item.IsDir,
	item.Checksum,
	time.Now().UTC(),
	item.Mtime,
	item.Perms,
	item.HostUpdated,
	item.HostIPs,
	time.Now().UTC(),
	item.Filename,
	item.Directory)
	err := tx.Commit()
	if err != nil {
		log.Infof("Error updating data in DB table %s, for file %s\n%s\n%+v", table, item.Filename, err.Error())
		return false
	}
	return true
}

func (my *MySQLDB) Remove(table, path string) (success bool) {
	log.Debugf("MySQL Remove on table %s", table)

	return
}

func (my *MySQLDB) FetchAll(table string) (fsitem []filesystem.FileData, err error) {
	log.Debugf("MySQL GetAll on table %s", table)

	return
}

func (my *MySQLDB) FetchOne(table, path string) (fi filesystem.FileData, err error) {
	log.Debugf("MySQL GetOne on table %s", table)

	return
}

func (my *MySQLDB) Exists(table string, fi filesystem.FileData) (exists bool) {
	log.Debugf("MySQL Exists on table %s", table)
	my.db = my.GetConn()
	defer my.db.Close()
	query := fmt.Sprintf("SELECT id FROM %s where filename='%s' AND directory='%s';", table, fi.Filename, fi.Directory)
	_, err := my.db.Query(query)
	if err == nil {
		log.Debugf("Item exists: %s", fi.Directory+"/"+fi.Filename)
		exists = true
		return
	}
	log.Debugf("Item does not exist: %s", fi.Directory+"/"+fi.Filename)
	return
}

func (my *MySQLDB) GetConn() (tempdb *sqlx.DB) {
	dbc := my.config.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?autocommit=true&parseTime=True", dbc.User, dbc.Pass, dbc.Host, dbc.Port, dbc.DBase)
	tempdb, err := sqlx.Connect(dbc.Type, dsn)
	if err != nil {
		log.Info(err.Error())
	}
	return
}

func createTableQuery(table string) string {
	var createStmt = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
  is_dir int NOT NULL,
  filename varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  directory varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  checksum varchar(2048) COLLATE utf8_unicode_ci NOT NULL,
  atime timestamp NOT NULL,
  mtime timestamp NOT NULL,
  perms varchar(12) COLLATE utf8_unicode_ci NOT NULL,
  host_updated varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  host_ips varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  last_update timestamp NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
`, table)

	return createStmt
}

func (my *MySQLDB) Setup() {
	my.db = my.GetConn()
	defer my.db.Close()
	for key, _ := range my.config.Listeners {
		log.Debugf("MySQL Setup on table %s", key)
		my.db.MustExec(createTableQuery(key))
	}
}

func (my *MySQLDB) fmtFilename(fpath string) (filename string) {
	filename = path.Base(fpath)
	return
}

func (my *MySQLDB) fmtDir(fpath string, table string) (dir string) {
	// We get the directory of the file first
	rdir := path.Dir(fpath)
	// now we strip out the directory base
	ldir := my.config.Listeners[table].Directory
	dir = strings.TrimPrefix(rdir, ldir)
	return
}
