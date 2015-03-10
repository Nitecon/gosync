package storage

import (
    "gosync/utils"
    "log"
)

type GDrive struct {
	config *utils.Configuration
    listener string
}

func (g *GDrive) Upload(local_path, remote_path string) error {
    utils.WriteLn("Doing an upload to GDrive")
	return nil
}

func (g *GDrive) Download(remote_path, local_path string, uid, gid int, perms, listener string) error {
    utils.WriteLn("Doing a download from GDrive")
	return nil
}

func (g *GDrive) Remove(remote_path string) bool {
    log.Printf("Removing %s from Google Drive Storage", remote_path)
	return true
}
