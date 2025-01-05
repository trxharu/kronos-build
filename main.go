package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/gorilla/websocket"
	"trxharu.dev/build-tool/config"
	"trxharu.dev/build-tool/internal/server"
	"trxharu.dev/build-tool/internal/utils"
	"trxharu.dev/build-tool/internal/watcher"
)


func main() {
	config, err := config.ReadConfigFromFile("settings.json")
	fatalErrorExit(err)

	dirs, err := utils.GetWatchableDirs(config.Source)
	fatalErrorExit(err)
	dirs = utils.DirsExcludePatterns(dirs, config.ExcludeDir)

	regexPatterns := compilePatterns(config.WatchFileTypes)

	fswatcher := watcher.Watcher{}

	fswatcher.WatchOverDirs(dirs, func(event int, args string) {
		if event == watcher.FILE_CREATED || event == watcher.FILE_MODIFIED {
			pass := filterFileTypes(args, regexPatterns)
			if !pass { return }
			fmt.Println(args)
		}
		if event == watcher.REMOVE_EVENT {
			fmt.Println("REMOVED", args)
		}
	})

	defer fswatcher.Close()

	wsHandler := func(ws *websocket.Conn) {
		ws.WriteJSON("text ws")
		ws.Close()
	}

	err = server.StartServer(config.ServeDir, config.Listen, wsHandler)
	fatalErrorExit(err)
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
