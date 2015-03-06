package replicator

import (
	"gosync/dbsync"
	"gosync/storage"
	"gosync/utils"
	"os"
    "time"
    "fmt"
)

func getFileInDatabase(dbPath string, fsItems []utils.FsItem) bool {
	for _, fsitem := range fsItems {
		if fsitem.Filename == dbPath {
			return true
		}
	}
	return false
}

func InitialSync() {
	cfg := utils.GetConfig()

    utils.WriteLn("Verifying DB Tables")
	dbsync.CreateDB()
    utils.WriteLn("Initial sync starting...")

	for key, listener := range cfg.Listeners {
		// First check to see if the table is empty and do a full import false == not empty
		if dbsync.CheckEmpty(key) == false {
			// Database is not empty so pull the updates and match locally
			items := dbsync.FetchAll(key)
            handleDataChanges(items, listener, key)
		} else {
			// Database is empty so lets import
			fsItems := utils.ListFilesInDir(listener.Directory)
			for _, item := range fsItems {
				success := dbsync.Insert(key, item)
				if success != true {
                    utils.LogWriteF("An error occurred inserting %x to database", item)
				}
				if !item.IsDir {
					storage.PutFile(item.Filename, key)
				}

			}
		}

	}

    utils.WriteLn("Initial sync completed...")
}


func CheckIn(path string){
    utils.WriteLn("Starting db checking background script: " + path)
    ticker := time.NewTicker(10 * time.Second)
    quit := make(chan struct{})
    go func() {
        for {
            select {
            case <-ticker.C:
                utils.WriteLn("Checking all changed stuff in db for: " + path)
                listener := utils.GetListenerFromDir(path)
                items, err := dbsync.CheckIn(listener)
                if err != nil{
                    utils.LogWriteF("Error occurred getting data for %s (%s): %+v",listener, err.Error(), err)
                }
                cfg := utils.GetConfig()
                handleDataChanges(items, cfg.Listeners[listener], listener)
            // @TODO: check that db knows Im alive.
            case <-quit:
                ticker.Stop()
                return
            }
        }
    }()
}

func handleDataChanges(items []utils.DataTable, listener utils.Listener, listenerName string){
    // Walking the directory to get files.
    fsItems := utils.ListFilesInDir(listener.Directory)
    for _, item := range items {
        absPath := utils.GetAbsPath(listenerName, item.Path)
        itemMatch := getFileInDatabase(absPath, fsItems)
        perms := fmt.Sprintf("%#o",item.Perms)
        if itemMatch {
            // Check to make sure it's not a directory as directories don't need to be uploaded
            if !item.IsDirectory {
                fileMD5 := utils.GetMd5Checksum(absPath)
                if fileMD5 != item.Checksum {
                    utils.LogWriteF("Found %s in db(%s) and fs(%s), NOT matching md5...", absPath, item.Checksum, fileMD5)
                    //@TODO: download the file and set corrected params for file.
                    hostname, _ := os.Hostname()
                    if item.HostUpdated != hostname {

                        if !storage.GetNodeCopy(item, listenerName, listener.Uid,listener.Gid, perms) {
                            // The server must be down so lets get it from S3
                            storage.GetFile(absPath, listenerName, listener.Uid,listener.Gid, perms)
                        }
                    }
                }
            }
            // Now we check to make sure the files match correct users etc

        } else {
            // Item doesn't exist locally but exists in DB so restore it
            utils.LogWriteF("Item Deleted Locally: %s restoring from DB marker", absPath)

            if item.IsDirectory{
                dirExists,_ := utils.ItemExists(absPath)
                if !dirExists{
                    os.MkdirAll(absPath, 0775)
                }
            }else{
                if !storage.GetNodeCopy(item, listenerName, listener.Uid,listener.Gid, perms) {
                    utils.LogWriteF("Server is down for %s going to backup storage", listenerName)
                    // The server must be down so lets get it from S3
                    storage.GetFile(absPath, listenerName, listener.Uid,listener.Gid, perms)
                }
            }


        }
    }
}