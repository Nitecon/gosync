package prototypes

import (
	"time"
)

type DataTable struct {
	Id          int       `db:"id"`
	Path        string    `db:"path"`
	IsDirectory bool      `db:"is_dir"`
	Filename    string    `db:"filename"`
	Directory   string    `db:"directory"`
	Checksum    string    `db:"checksum"`
	Atime       time.Time `db:"atime"`
	Mtime       time.Time `db:"mtime"`
	Perms       string    `db:"perms"`
	HostUpdated string    `db:"host_updated"`
	LastUpdate  time.Time `db:"last_update"`
}
