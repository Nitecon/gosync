package firstrun

import (
	"gosync/config"
	"gosync/dbsync"
	"gosync/fstools"
	"log"
)

func InitialSync(cfg *config.Configuration) {
	log.Println("Verifying DB Tables")
	dbsync.DBInit(cfg)
	log.Println("Initial sync starting...")
	for key, listener := range cfg.Listeners {
		// First check to see if the table is empty and do a full import false == not empty
		if dbsync.DBCheckEmpty(cfg, key) == false {
			// Database is not empty so pull the updates and match locally
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
