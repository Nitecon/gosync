package fswatcher

import (
	"gopkg.in/fsnotify.v1"
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
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("Rename occurred on:", event.Name)
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("New File:", event.Name)
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("Removed File: ", event.Name)
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
