package firstrun

import (
	"gosync/config"
	"gosync/dbsync"
	"gosync/fstools"
	"gosync/prototypes"
	"log"
)

func getFileInDatabase(dbItem prototypes.DataTable, fsItems []fstools.FsItem) (bool, string) {
	for _, fsitem := range fsItems {
		if fsitem.Filename == dbItem.Path {
			return true, dbItem.Path
		}

	}
	return false, dbItem.Path
}

func InitialSync(cfg *config.Configuration) {
	log.Println("Verifying DB Tables")
	dbsync.DBInit(cfg)
	log.Println("Initial sync starting...")
	for key, listener := range cfg.Listeners {
		// First check to see if the table is empty and do a full import false == not empty
		if dbsync.DBCheckEmpty(cfg, key) == false {
			// Database is not empty so pull the updates and match locally
			items := dbsync.DBFetchAll(cfg, key)
			// Walking the directory to get files.
			fsItems := fstools.ListFilesInDir(listener.Directory)
			for _, item := range items {
				itemMatch, pathMatch := getFileInDatabase(item, fsItems)
				if itemMatch {
					// Check to see if the checksum matches, if not check update times and upload / download
					if !item.IsDirectory {
						fileMD5 := fstools.GetMd5Checksum(pathMatch)
						if fileMD5 == item.Checksum {
							log.Printf("Found %s in db and fs, matching md5...", pathMatch)
							//@TODO: download the file and set corrected params for file.
						} else {
							log.Printf("Found %s in db and fs, MD5 mismatch... DOWNLOADING", pathMatch)
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
				success := dbsync.InsertItem(cfg, key, item)
				if success != true {
					log.Printf("An error occurred inserting %x to database", item)
				}
			}
		}

	}

	log.Println("Initial sync completed...")
}
