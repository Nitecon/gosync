package prototypes

import (
	"time"
)

//path, is_dir, filename, directory, checksum, atime, mtime, uid, gid, perms, host_updated, last_update
type DataTable struct {
	Id          int       `db:"id"`
	Path        string    `db:"path"`
	IsDirectory bool      `db:"is_dir"`
	Filename    string    `db:"filename"`
	Directory   string    `db:"directory"`
	Checksum    string    `db:"checksum"`
	Atime       time.Time `db:"atime"`
	Mtime       time.Time `db:"mtime"`
	Uid         int       `db:"uid"`
	Gid         int       `db:"gid"`
	Perms       string    `db:"perms"`
	HostUpdated string    `db:"host_updated"`
	LastUpdate  time.Time `db:"last_update"`
}
