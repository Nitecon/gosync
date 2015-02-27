package storage

import "log"

type GDrive struct {
	config string
}

func (g *GDrive) Upload(local_path, remote_path string) error {
	log.Println("Doing an upload to GDrive")
	return nil
}

func (g *GDrive) Download(remote_path, local_path string) error {
	log.Println("Doing a download from GDrive")
	return nil
}

func (g *GDrive) CheckMD5(local_path, remote_path string) bool {
	log.Println("Doing md5 check on GDrive")
	return true
}
