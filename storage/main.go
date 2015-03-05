package storage

import (

    "gosync/utils"


)

var (
	// BackupConfig determines which storage system to use,
	// and we then assign that to "storage".
	storage Backupper
)

type Backupper interface {
	Upload(local_path, remote_path string) error
	Download(remote_path, local_path string, uid, gid int, perms string) error
	CheckMD5(local_path, remote_path string) bool
}

func setStorageEngine(listener string) {
	cfg := utils.GetConfig()
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
	return storage.Upload(local_path, utils.GetRelativeBasePath(listener, local_path))
}

func GetFile(local_path, listener string, uid, gid int, perms string) error {
	setStorageEngine(listener)
	err := storage.Download(utils.GetRelativeBasePath(listener, local_path), local_path, uid, gid, perms)
    if err != nil{
        utils.WriteF("Error downloading file from S3 (%s) : %+v", err.Error(), err)
    }
    return err
}

func CheckFileMD5(local_path, listener string) bool {
	setStorageEngine(listener)
	return storage.CheckMD5(local_path, utils.GetRelativeBasePath(listener, local_path))
}
