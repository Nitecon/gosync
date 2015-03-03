package replicator

import (
	"gosync/config"
	"gosync/dbsync"
	"gosync/fstools"
	"gosync/storage"
	"gosync/utils"
	"log"
	"os"
    "time"
    "gosync/prototypes"
)

func getFileInDatabase(dbPath string, fsItems []fstools.FsItem) bool {
	for _, fsitem := range fsItems {
		if fsitem.Filename == dbPath {
			return true
		}
	}
	return false
}

func InitialSync() {
	cfg := config.GetConfig()

	log.Println("Verifying DB Tables")
	dbsync.CreateDB()
	log.Println("Initial sync starting...")

	for key, listener := range cfg.Listeners {
		// First check to see if the table is empty and do a full import false == not empty
		if dbsync.CheckEmpty(key) == false {
			// Database is not empty so pull the updates and match locally
			items := dbsync.FetchAll(key)
            handleDataChanges(items, listener, key)
		} else {
			// Database is empty so lets import
			fsItems := fstools.ListFilesInDir(listener.Directory)
			for _, item := range fsItems {
				success := dbsync.Insert(key, item)
				if success != true {
					log.Printf("An error occurred inserting %x to database", item)
				}
				if !item.IsDir {
					storage.PutFile(item.Filename, key)
				}

			}
		}

	}

	log.Println("Initial sync completed...")
}


func CheckIn(path string){
    log.Println("Starting db checking background script: " + path)
    ticker := time.NewTicker(10 * time.Second)
    quit := make(chan struct{})
    go func() {
        for {
            select {
            case <-ticker.C:
                log.Println("Checking all changed stuff in db for: " + path)
                listener := utils.GetListenerFromDir(path)
                items, err := dbsync.CheckIn(listener)
                if err != nil{
                    log.Printf("Error occurred getting data for %s (%s): %+v",listener, err.Error(), err)
                }
                cfg := config.GetConfig()
                handleDataChanges(items, cfg.Listeners[listener], listener)
            // @TODO: check that db knows Im alive.
            case <-quit:
                ticker.Stop()
                return
            }
        }
    }()
}

func handleDataChanges(items []prototypes.DataTable, listener config.Listener, listenerName string){
    // Walking the directory to get files.
    fsItems := fstools.ListFilesInDir(listener.Directory)
    for _, item := range items {
        absPath := utils.GetAbsPath(listenerName, item.Path)
        itemMatch := getFileInDatabase(absPath, fsItems)
        if itemMatch {
            // Check to make sure it's not a directory as directories don't need to be uploaded
            if !item.IsDirectory {
                fileMD5 := fstools.GetMd5Checksum(absPath)
                if fileMD5 != item.Checksum {
                    log.Printf("Found %s in db(%s) and fs(%s), NOT matching md5...", absPath, item.Checksum, fileMD5)
                    //@TODO: download the file and set corrected params for file.
                    hostname, _ := os.Hostname()
                    if item.HostUpdated != hostname {
                        if !storage.GetNodeCopy(item, listenerName) {
                            // The server must be down so lets get it from S3
                            storage.GetFile(absPath, listenerName)
                        }
                    }
                }
            }
            // Now we check to make sure the files match correct users etc

        } else {
            // Item doesn't exist locally but exists in DB so restore it
            log.Printf("Item Deleted Locally: %s restoring from DB marker", absPath)

            if item.IsDirectory{
                dirExists,_ := utils.ItemExists(absPath)
                if !dirExists{
                    os.MkdirAll(absPath, 0775)
                }
            }else{
                if !storage.GetNodeCopy(item, listenerName) {
                    log.Printf("Server is down for %s going to backup storage", listenerName)
                    // The server must be down so lets get it from S3
                    storage.GetFile(absPath, listenerName)
                }
            }


        }
    }
}