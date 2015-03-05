package fswatcher

import (
	"gopkg.in/fsnotify.v1"
	"gosync/dbsync"
	"gosync/utils"
	"gosync/storage"
)

func SysPathWatcher(path string) {
    utils.WriteF("Starting new watcher for:", path)
	watcher, err := fsnotify.NewWatcher()
	utils.Check("Cannot create new watcher for: "+path, 400, err)
	defer watcher.Close()
	done := make(chan bool)
	go func() {

		for {
			select {
			case event := <-watcher.Events:
				//logs.WriteLn("event:", event)
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
                    utils.WriteF("Chmod occurred on:", event.Name)
					runFileUpdate(path, event.Name, "chmod")
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
                    utils.WriteF("Rename occurred on:", event.Name)
					runFileUpdate(path, event.Name, "rename")
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
                    utils.WriteF("New File:", event.Name)
					runFileUpdate(path, event.Name, "create")
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
                    utils.WriteF("modified file:", event.Name)
					runFileUpdate(path, event.Name, "write")
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
                    utils.WriteF("Removed File: ", event.Name)
					runFileUpdate(path, event.Name, "remove")
				}
			case err := <-watcher.Errors:
                utils.WriteF("error:", err)
			}
		}

	}()
	err = watcher.Add(path)
    utils.Check("Cannot add watcher to: "+path, 400, err)
	<-done
}

func runFileUpdate(base_path, path, operation string) bool {
	listener := utils.GetListenerFromDir(base_path)
	rel_path := utils.GetRelativePath(listener, path)
	fsItem, err := utils.GetFileInfo(path)

	if err != nil {
        utils.WriteF("Error getting file details for %s: %+v", path, err)
	}

	dbItem, err := dbsync.GetOne(base_path, rel_path)
	if err != nil {
        utils.WriteF("Error occurred getting file row (%s): %+v", err.Error(), err)
	}

	switch operation {
	/*case "chmod:":
	  if dbItem.Perms != 0664{
	      logs.WriteLn("Perms don't match")
	  }*/
	case "create":
		if dbItem.Checksum != fsItem.Checksum {
			utils.WriteF("Creating:->")
			dbsync.Insert(listener, fsItem)
            utils.WriteF("Putting in storage:->")
			storage.PutFile(path, listener)
		}
	case "write":
		if dbItem.Checksum != fsItem.Checksum {
            utils.WriteF("Writing:->")
			dbsync.Insert(listener, fsItem)
            utils.WriteF("Putting in storage:->")
			storage.PutFile(path, listener)
		}
    case "remove":
        // Item was removed so we just wipe it from DB and storage
        dbsync.Insert(listener, fsItem)
	}
	return false
}
