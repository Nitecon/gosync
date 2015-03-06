package nodeinfo

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gosync/utils"
	"time"
    "pkg/os"
)

var (
	table    = "gosync_nodes"
)

type MySQLNodeDB struct {
	config *utils.Configuration
	db     *sqlx.DB
}

func (my *MySQLNodeDB) Insert() bool {
	tx := my.db.MustBegin()
	tx.MustExec("INSERT INTO "+table+" (hostname, host_ips, connected_on, last_update) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE hostname=?, host_ips=?, connected_on=?, last_update=?",
    utils.GetSystemHostname(),
    utils.GetLocalIp(),
		time.Now().UTC(),
		time.Now().UTC(),
    utils.GetSystemHostname(),
    utils.GetLocalIp(),
    time.Now().UTC(),
    time.Now().UTC()    )
	err := tx.Commit()
	utils.Check(err, 500, "Error setting node alive...")
	return true
}

func (my *MySQLNodeDB) GetAll() (utils.NodeTable, error) {
	nodeData := utils.NodeTable{}
    query := "select id, hostname, host_ips, connected_on, last_update from "+table+" where hostname != '"+utils.GetSystemHostname()+"' ORDER BY last_update DESC;"
    err := my.db.Select(&nodeData, query)
	return nodeData, err
}

func (my *MySQLNodeDB) Update() bool {
    tx := my.db.MustBegin()
    tx.MustExec("UPDATE "+table+" SET hostname=?, host_ips=?, last_update=? where hostname=?" ,
    utils.GetSystemHostname(),
    utils.GetLocalIp(),
    time.Now().UTC(),
    utils.GetSystemHostname())
    err := tx.Commit()
    utils.Check(err, 500, "Error updating node alive...")
    return true
}

// call this method when you want to close the connection
func (my *MySQLNodeDB) Close(deadHost string) {
    tx := my.db.MustBegin()
    tx.MustExec("DELETE FROM "+table+" where hostname='?'" ,
    deadHost)
    err := tx.Commit()
    utils.Check(err, 500, "Error updating node alive...")
}

func (my *MySQLNodeDB) initDB() {
	tempdb, err := sqlx.Connect("mysql", my.config.Database.Dsn+"&parseTime=True")
	if err != nil {
		utils.WriteLn(err.Error())
	}
	my.db = tempdb
}

func (my *MySQLNodeDB) CreateDB() {
	my.db.MustExec(createTableQuery(table))
}

func createTableQuery(table string) string {
	var createStmt = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
  hostname varchar(255) COLLATE utf8_unicode_ci NOT NULL UNIQUE,
  host_ips varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  connected_on timestamp NOT NULL,
  last_update timestamp NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
`, table)
	return createStmt
}

type NodeTable struct {
	Id         int       `db:"id"`
	HostName   string    `db:"hostname"`
	NodeIPs    string    `db:"host_ips"`
	Connected  time.Time `db:"connected_on"`
	LastUpdate time.Time `db:"last_update"`
}
