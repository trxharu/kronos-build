package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"

	"trxharu.dev/kronos-build/config"
	"trxharu.dev/kronos-build/internal/server"
	"trxharu.dev/kronos-build/internal/utils"
	"trxharu.dev/kronos-build/internal/watcher"
	"trxharu.dev/kronos-build/internal/websocket"
	"trxharu.dev/kronos-build/trigger"
)


func main() {
	args := os.Args
	var cfg config.Config
	var err error

	if len(args) > 1 {
		cfg, err = config.ReadConfigFromFile(args[1])
		fatalErrorExit(err)
	} else {
		cfg, err = config.ReadConfigFromFile("kronos.build.json")
		fatalErrorExit(err)
	}

	dirs, err := utils.GetWatchableDirs(cfg.Source)
	fatalErrorExit(err)
	
	trigger.SetupPipeline(cfg.ServeDir)			
	trigger.TriggerBuildPipeline(cfg.RunCmd)

	dirs = utils.DirsExcludePatterns(dirs, cfg.ExcludeDir)
	regexPatterns := compilePatterns(cfg.WatchFileTypes)

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
			trigger.TriggerBuildPipeline(cfg.RunCmd)
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

	s, err := server.StartServer(cfg.ServeDir, cfg.Listen, ws)
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
