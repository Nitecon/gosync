package storage

import (
	"fmt"
	"gosync/config"
)

var (
	// BackupConfig determines which storage system to use,
	// and we then assign that to "storage".
	BackupConfig = "GDrive"
	storage      Backupper
)

type Backupper interface {
	Upload(local_path, remote_path string) error
	Download(remote_path, local_path string) error
	CheckMD5(local_path, remote_path string) bool
}

func PutFile(local_path, remote_path string, cfg *config.Configuration) error {
	switch BackupConfig {
	case "GDrive":
		storage = &GDrive{}
	case "S3":
		storage = &S3{}
	}
	return storage.Upload(local_path, remote_path, cfg*config.Configuration)
}

func GetFile(local_path, remote_path string) error {
	switch BackupConfig {
	case "GDrive":
		storage = &GDrive{}
	case "S3":
		storage = &S3{}
	}
	return storage.Download(remote_path, local_path)
}

func CheckFileMD5(local_path, remote_path string, cfg *config.Configuration) bool {
	switch BackupConfig {
	case "GDrive":
		storage = &GDrive{}
	case "S3":
		storage = &S3{}
	}
	return storage.CheckMD5(local_path, remote_path)
}
