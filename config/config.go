package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Source string `json:"source"`
	WatchFileTypes []string `json:"watchFileTypes"`
	ExcludeDir []string `json:"excludeDir"`
	ServeDir string `json:"serveDir"`
	Listen string `json:"listen"`
	RunCmd []string `json:"runCmd"`
}

func ReadConfigFromFile(filePath string) (Config, error) {
	config := Config {
		ServeDir: ".",
		Source: ".",
		ExcludeDir: []string{},
		WatchFileTypes: []string{},
		Listen: "localhost:2345",
		RunCmd: []string{},
	}

	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	err = file.Close()

	return config, err
}
