package fswatcher

import (
	"gopkg.in/fsnotify.v1"
	"gosync/datastore"
	"gosync/storage"
	"gosync/utils"
	"log"
	"os"
	"strconv"
)

func SysPathWatcher(path string) {
	log.Printf("Starting new watcher for %s:", path)
	watcher, err := fsnotify.NewWatcher()
	if !utils.ErrorCheckF(err, 400, "Cannot create new watcher for %s ", path) {
		defer watcher.Close()
		done := make(chan bool)
		go func() {
			listener := utils.GetListenerFromDir(path)
			rel_path := utils.GetRelativePath(listener, path)
			for {
				select {
				case event := <-watcher.Events:
					//logs.WriteLn("event:", event)
					if event.Op&fsnotify.Chmod == fsnotify.Chmod {
						runFileChmod(path, event.Name)
					}
					if event.Op&fsnotify.Rename == fsnotify.Rename {
						log.Printf("Rename occurred on:", event.Name)
						if checksumItem(path, rel_path, event.Name) {
							runFileRename(path, event.Name)
						}
					}
					if (event.Op&fsnotify.Create == fsnotify.Create) || (event.Op&fsnotify.Write == fsnotify.Write) {
						log.Printf("New / Modified File: %s", event.Name)
						if checksumItem(path, rel_path, event.Name) {
							runFileCreateUpdate(path, event.Name, "create")
						}
					}

					if event.Op&fsnotify.Remove == fsnotify.Remove {
						runFileRemove(path, event.Name)
						log.Printf("Removed File: ", event.Name)
					}
				case err := <-watcher.Errors:
					log.Printf("error:", err)
				}
			}

		}()
		err = watcher.Add(path)
		utils.ErrorCheckF(err, 500, "Cannot add watcher to %s ", path)
		<-done
	}

}

func runFileRename(base_path, path string) bool {
	log.Printf("File removed %s: ", path)

	return true
}

func runFileRemove(base_path, path string) bool {
	listener := utils.GetListenerFromDir(base_path)
	rel_path := utils.GetRelativePath(listener, path)
	// Remove the item from the database (it's already removed from filesystem
	datastore.Remove(listener, rel_path)
	// Remove the item from the backup storage location
	storage.RemoveFile(rel_path, listener)
	return true
}

func runFileChmod(base_path, path string) bool {
	listener := utils.GetListenerFromDir(base_path)
	rel_path := utils.GetRelativePath(listener, path)
	dbItem, err := datastore.GetOne(base_path, rel_path)
	if err != nil {
		log.Printf("Error occurred trying to get %s from DB\nError: %s", rel_path, err.Error())
		return false
	}
	fsItem, err := utils.GetFileInfo(path)
	if err != nil {
		log.Printf("Could not find item on filesystem: %s\nError:%s", path, err.Error())
	}
	if dbItem.Perms != fsItem.Perms {
		iPerm, _ := strconv.Atoi(dbItem.Perms)
		mode := int(iPerm)
        if _, err := os.Stat(path); os.IsNotExist(err) {
            log.Printf("File no longer exists returning")
            return true
        }else{
            err := os.Chmod(path, os.FileMode(mode))
            if err != nil {
                log.Printf("Error occurred changing file modes: %s", err.Error())
                return false
            }
            return true
        }
	} else {
		log.Printf("File modes are correct changing nothing")
		return false
	}
	return true
}

func runFileCreateUpdate(base_path, path, operation string) bool {
	listener := utils.GetListenerFromDir(base_path)
	rel_path := utils.GetRelativePath(listener, path)
	fsItem, err := utils.GetFileInfo(path)

	if err != nil {
		log.Printf("Error getting file details for %s: %+v", path, err)
	}

	dbItem, err := datastore.GetOne(base_path, rel_path)

	if !utils.ErrorCheckF(err, 400, "Error getting file row (%s)", rel_path) {
		switch operation {
		case "create":
			if dbItem.Checksum != fsItem.Checksum {
				log.Printf("Creating:-> %s", rel_path)
				datastore.Insert(listener, fsItem)
				log.Printf("Putting in storage:-> %s", rel_path)
				storage.PutFile(path, listener)
			}
		case "write":
			if dbItem.Checksum != fsItem.Checksum {
				log.Printf("Writing:->")
				datastore.Insert(listener, fsItem)
				log.Printf("Putting in storage:->")
				storage.PutFile(path, listener)
			}
		}
	}

	return false
}

func checksumItem(base_path, rel_path, abspath string) bool {
	dbItem, err := datastore.GetOne(base_path, rel_path)
	if err != nil {
		log.Printf("Error occurred getting item from DB: %s", err.Error())
		return false
	} else {
		if dbItem.Checksum != utils.GetMd5Checksum(abspath) {
			log.Printf("%s != DB item checksum", abspath)
			return true
		} else {
			log.Printf("%s === DB item checksum", abspath)
			return false
		}
	}
	return false
}
