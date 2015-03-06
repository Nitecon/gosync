package fswatcher

import (
	"gopkg.in/fsnotify.v1"
	"gosync/dbsync"
	"gosync/storage"
	"gosync/utils"
)

func SysPathWatcher(path string) {
	utils.LogWriteF("Starting new watcher for %s:", path)
	watcher, err := fsnotify.NewWatcher()
	if !utils.CheckF(err, 400, "Cannot create new watcher for %s ", path) {
		defer watcher.Close()
		done := make(chan bool)
		go func() {

			for {
				select {
				case event := <-watcher.Events:
					//logs.WriteLn("event:", event)
					if event.Op&fsnotify.Chmod == fsnotify.Chmod {
						utils.LogWriteF("Chmod occurred on:", event.Name)
						runFileUpdate(path, event.Name, "chmod")
					}
					if event.Op&fsnotify.Rename == fsnotify.Rename {
						utils.LogWriteF("Rename occurred on:", event.Name)
						runFileUpdate(path, event.Name, "rename")
					}
					if event.Op&fsnotify.Create == fsnotify.Create {
						utils.LogWriteF("New File:", event.Name)
						runFileUpdate(path, event.Name, "create")
					}
					if event.Op&fsnotify.Write == fsnotify.Write {
						utils.LogWriteF("modified file:", event.Name)
						runFileUpdate(path, event.Name, "write")
					}
					if event.Op&fsnotify.Remove == fsnotify.Remove {
						utils.LogWriteF("Removed File: ", event.Name)
						runFileUpdate(path, event.Name, "remove")
					}
				case err := <-watcher.Errors:
					utils.LogWriteF("error:", err)
				}
			}

		}()
		err = watcher.Add(path)
		utils.CheckF(err, 500, "Cannot add watcher to %s ", path)
		<-done
	}

}

func runFileUpdate(base_path, path, operation string) bool {
	listener := utils.GetListenerFromDir(base_path)
	rel_path := utils.GetRelativePath(listener, path)
	fsItem, err := utils.GetFileInfo(path)

	if err != nil {
		utils.LogWriteF("Error getting file details for %s: %+v", path, err)
	}

	dbItem, err := dbsync.GetOne(base_path, rel_path)

	if !utils.CheckF(err, 400, "Error getting file row (%s)", rel_path) {
		switch operation {
		/*case "chmod:":
		  if dbItem.Perms != 0664{
		      logs.WriteLn("Perms don't match")
		  }*/
		case "create":
			if dbItem.Checksum != fsItem.Checksum {
				utils.LogWriteF("Creating:->")
				dbsync.Insert(listener, fsItem)
				utils.LogWriteF("Putting in storage:->")
				storage.PutFile(path, listener)
			}
		case "write":
			if dbItem.Checksum != fsItem.Checksum {
				utils.LogWriteF("Writing:->")
				dbsync.Insert(listener, fsItem)
				utils.LogWriteF("Putting in storage:->")
				storage.PutFile(path, listener)
			}
		case "remove":
			// Item was removed so we just wipe it from DB and storage
			dbsync.Remove(listener, fsItem)
		}
	}

	return false
}
