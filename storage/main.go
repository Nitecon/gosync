package storage

import (
	"gosync/config"

	"strings"
    "log"
)

var (
	// BackupConfig determines which storage system to use,
	// and we then assign that to "storage".
	storage Backupper
)

type Backupper interface {
	Upload(local_path, remote_path string) error
	Download(remote_path, local_path string) error
	CheckMD5(local_path, remote_path string) bool
}

func setStorageEngine(listener string) {
	cfg := config.GetConfig()
	var engine = cfg.Listeners[listener].StorageType
	switch engine {
	case "gdrive":
		storage = &GDrive{config:cfg, listener:listener}
	case "s3":
		storage = &S3{config:cfg,listener:listener}
	}
}

func PutFile(local_path, listener string) error {
	setStorageEngine(listener)
	return storage.Upload(local_path, getRemotePath(listener, local_path))
}

func GetFile(local_path, listener string) error {
	setStorageEngine(listener)
    cfg := config.GetConfig()
    basePath := cfg.Listeners[listener].BasePath
    remotePath := basePath + strings.TrimPrefix(local_path, cfg.Listeners[listener].Directory)

	err := storage.Download(remotePath, local_path)
    if err != nil{
        log.Printf("Error downloading file from S3 (%s) : %+v", err.Error(), err)
    }
    return err
}

func CheckFileMD5(local_path, listener string) bool {
	setStorageEngine(listener)
	return storage.CheckMD5(local_path, getRemotePath(listener, local_path))
}

func getRemoteBasePath(listener string) string {
    cfg := config.GetConfig()
    return cfg.Listeners[listener].BasePath
}

func getBaseDir(listener string) string {
    cfg := config.GetConfig()
    return cfg.Listeners[listener].Directory
}

func getRemotePath(listener, local_path string) string {
    lPath := strings.TrimPrefix(local_path, getBaseDir(listener))
    log.Println("==>REMOTEPATH: "+ lPath)
    return getRemoteBasePath(listener) + lPath
}