package watcher

import (
	"log"
	"github.com/fsnotify/fsnotify"
)

type WatcherCallback func(filename string)

func WatchOverDirs(dirs []string, callback WatcherCallback) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}
	
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <- watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					callback(event.Name)
				}
			case err, ok := <- watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	for _, dir := range dirs {
		err = watcher.Add(dir)
		if err != nil {
			log.Fatal(err)
		}
	}
	<-make(chan struct{})
}
