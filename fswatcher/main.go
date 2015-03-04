package fswatcher

import (
	"gopkg.in/fsnotify.v1"
	"gosync/dbsync"
	"gosync/fstools"
	"gosync/storage"
	"gosync/utils"
	"log"
)

func SysPathWatcher(path string) {
	log.Println("Starting new watcher for:", path)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {

		for {
			select {
			case event := <-watcher.Events:
				//log.Println("event:", event)
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Println("Chmod occurred on:", event.Name)
					runFileUpdate(path, event.Name, "chmod")
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("Rename occurred on:", event.Name)
					runFileUpdate(path, event.Name, "rename")
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("New File:", event.Name)
					runFileUpdate(path, event.Name, "create")
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					runFileUpdate(path, event.Name, "write")
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("Removed File: ", event.Name)
					runFileUpdate(path, event.Name, "remove")
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}

	}()
	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func runFileUpdate(base_path, path, operation string) bool {
	listener := utils.GetListenerFromDir(base_path)
	rel_path := utils.GetRelativePath(listener, path)
	fsItem, err := fstools.GetFileInfo(path)

	if err != nil {
		log.Fatalf("Error getting file details for %s: %+v", path, err)
	}

	dbItem, err := dbsync.GetOne(base_path, rel_path)
	if err != nil {
		log.Fatal("Error occurred getting file row (%s): %+v", err.Error(), err)
	}

	switch operation {
	/*case "chmod:":
	  if dbItem.Perms != 0664{
	      log.Println("Perms don't match")
	  }*/
	case "create":
		if dbItem.Checksum != fsItem.Checksum {
			log.Printf("Creating:->")
			dbsync.Insert(listener, fsItem)
			log.Printf("Putting in storage:->")
			storage.PutFile(path, listener)
		}
	case "write":
		if dbItem.Checksum != fsItem.Checksum {
			log.Printf("Writing:->")
			dbsync.Insert(listener, fsItem)
			log.Printf("Putting in storage:->")
			storage.PutFile(path, listener)
		}
    case "remove":
        // Item was removed so we just wipe it from DB and storage
        dbsync.Insert(listener, fsItem)
	}
	return false
}
