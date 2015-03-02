package firstrun

import (
	"gosync/config"
	"gosync/dbsync"
	"gosync/fstools"
	"gosync/prototypes"
	"gosync/storage"
	"log"
	"os"
)

func getFileInDatabase(dbItem prototypes.DataTable, fsItems []fstools.FsItem) (bool, string) {
	for _, fsitem := range fsItems {
		if fsitem.Filename == dbItem.Path {
			return true, dbItem.Path
		}

	}
	return false, dbItem.Path
}

func InitialSync() {
	cfg := config.GetConfig()

	log.Println("Verifying DB Tables")
	dbsync.DBInit()
	log.Println("Initial sync starting...")

	for key, listener := range cfg.Listeners {
		// First check to see if the table is empty and do a full import false == not empty
		if dbsync.CheckEmpty(key) == false {
			// Database is not empty so pull the updates and match locally
			items := dbsync.FetchAll(key)
			// Walking the directory to get files.
			fsItems := fstools.ListFilesInDir(listener.Directory)
			for _, item := range items {
				itemMatch, pathMatch := getFileInDatabase(item, fsItems)
				if itemMatch {
					// Check to see if the checksum matches, if not check update times and upload / download
					if !item.IsDirectory {

						fileMD5 := fstools.GetMd5Checksum(pathMatch)
						if fileMD5 == item.Checksum {
							//log.Printf("Found %s in db and fs, matching md5...", pathMatch)
							//@TODO: download the file and set corrected params for file.
							hostname, _ := os.Hostname()
							if item.HostUpdated != hostname {
								if !storage.GetNodeCopy(item, key) {
									// The server must be down so lets get it from S3
									storage.GetFile(pathMatch, key)
								}
							}

						} else {
							//log.Printf("Found %s in db and fs, MD5 mismatch... DOWNLOADING", pathMatch)

							// Last resort download from S3
							storage.GetFile(pathMatch, key)
						}
					}
					// Now we check to make sure the files match correct users etc

				} else {
					// Item doesn't exist locally but exists in DB so restore it
					log.Printf("Item Deleted Locally: %s restoring from DB marker", pathMatch)
				}
			}

		} else {
			// Database is empty so lets import
			fsItems := fstools.ListFilesInDir(listener.Directory)
			for _, item := range fsItems {
				success := dbsync.Insert(key, item)
				if success != true {
					log.Printf("An error occurred inserting %x to database", item)
				}
			}
		}

	}

	log.Println("Initial sync completed...")
}
