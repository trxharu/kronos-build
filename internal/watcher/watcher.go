package watcher

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"trxharu.dev/build-tool/internal/utils"
)

const (
	DIR_CREATED = iota
	FILE_CREATED
	FILE_MODIFIED
	REMOVE_EVENT
	RENAME_EVENT
)

type Watcher struct {
	fsWatcher *fsnotify.Watcher
}

type WatchCallback func(event int, args string)

func (w *Watcher) WatchOverDirs(dirs []string, callback WatchCallback) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	w.fsWatcher = watcher
	
	go func() {
		for {
			select {
			case event, ok := <- watcher.Events:
				if !ok { return }

				if event.Has(fsnotify.Write) {
					callback(FILE_MODIFIED, event.Name)
				}

				if event.Has(fsnotify.Create) {
					if ok, _ := utils.IsDir(event.Name); ok {
						w.AddDir(event.Name)
						callback(DIR_CREATED, event.Name)
					} else {
						callback(FILE_CREATED, event.Name)
					}
				}

				if event.Has(fsnotify.Rename) {
					callback(RENAME_EVENT, event.Name)
				}
				
				if event.Has(fsnotify.Remove) {
					w.RemoveDir(event.Name)
					callback(REMOVE_EVENT, event.Name)
				}
			case err, ok := <- watcher.Errors:
				if !ok { return }
				log.Println("error:", err)
			}
		}
	}()
	// Initial watch on given folders
	for _, dir := range dirs {
		err = watcher.Add(dir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (w *Watcher) AddDir(path string) {
	if !utils.IsPathExists(path) { return }
	err := w.fsWatcher.Add(path)	
	if err != nil {
		log.Fatal(err)
	}
}

func (w *Watcher) RemoveDir(path string) {
	watchDirs := w.fsWatcher.WatchList()
	for _, dir := range watchDirs {
		if dir == path {
			err := w.fsWatcher.Remove(path)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (w *Watcher) Close() {
	err := w.fsWatcher.Close()
	if err != nil {
		log.Fatal(err)
	}
}
