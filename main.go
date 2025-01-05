package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"

	"trxharu.dev/build-tool/config"
	"trxharu.dev/build-tool/internal/server"
	"trxharu.dev/build-tool/internal/utils"
	"trxharu.dev/build-tool/internal/watcher"
	"trxharu.dev/build-tool/internal/websocket"
)


func main() {
	config, err := config.ReadConfigFromFile("settings.json")
	fatalErrorExit(err)

	dirs, err := utils.GetWatchableDirs(config.Source)
	fatalErrorExit(err)
	dirs = utils.DirsExcludePatterns(dirs, config.ExcludeDir)

	regexPatterns := compilePatterns(config.WatchFileTypes)

	ws := &websocket.WebSocket{
		WsChan: make(chan websocket.Message),
	}

	fswatcher := watcher.Watcher{}
	fswatcher.WatchOverDirs(dirs, func(event int, args string) {
		if event == watcher.DIR_CREATED {
			ws.WsChan <- websocket.Message{
				Type: "newdir",
				Data: args,
			}
		}
		if event == watcher.FILE_CREATED || event == watcher.FILE_MODIFIED {
			pass := filterFileTypes(args, regexPatterns)
			if !pass { return }
			ws.WsChan <- websocket.Message{ 
				Type: "update", 
				Data: args,
			}
		}

		if event == watcher.REMOVE_EVENT {
			ws.WsChan <- websocket.Message{ 
				Type: "remove", 
				Data: args,
			}
		}
	})
	

	intSignal := make(chan os.Signal, 1)
	signal.Notify(intSignal, os.Interrupt)

	s, err := server.StartServer(config.ServeDir, config.Listen, ws)
	fatalErrorExit(err)
	<- intSignal

	fswatcher.Close()
	ws.Close()
	s.Close()	
}

func compilePatterns(patterns []string) []regexp.Regexp {
	var regexps []regexp.Regexp
	for _, patt := range patterns {
		regexps = append(regexps, *regexp.MustCompile(".*" + patt + "$"))	
	}
	return regexps	
}

func filterFileTypes(filename string, patterns []regexp.Regexp) bool {
	for _, regex := range patterns {
		if regex.MatchString(filename) { return true }
	}
	return false
}

func fatalErrorExit(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}
