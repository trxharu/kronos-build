package main

import (
	"fmt"
	"os"

	"trxharu.dev/build-tool/config"
	"trxharu.dev/build-tool/internal/utils"
	// "trxharu.dev/build-tool/internal/watcher"
	// "trxharu.dev/build-tool/server"
)


func main() {
	config, err := config.ReadConfigFromFile("settings.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	dirs, err := utils.GetWatchableDirs(config.Source)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	dirs = utils.DirsExcludePatterns(dirs, config.ExcludeDir)
	fmt.Println(dirs)
	// watcher.WatchOverDirs(dirs, func(filename string) {
	// 	fmt.Println("File Modified: ", filename)
	// })
	// err = server.StartServer(config.ServeDir, config.Listen)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}
