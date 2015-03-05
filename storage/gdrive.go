package storage

import (
    "gosync/utils"
)

type GDrive struct {
	config *utils.Configuration
    listener string
}

func (g *GDrive) Upload(local_path, remote_path string) error {
    utils.WriteLn("Doing an upload to GDrive")
	return nil
}

func (g *GDrive) Download(remote_path, local_path string, uid, gid int, perms string) error {
    utils.WriteLn("Doing a download from GDrive")
	return nil
}

func (g *GDrive) CheckMD5(local_path, remote_path string) bool {
    utils.WriteLn("Doing md5 check on GDrive")
	return true
}
