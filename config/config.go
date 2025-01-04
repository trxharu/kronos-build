package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Source string `json:"source"`
	IncludeFiles []string `json:"includeFiles"`
	ExcludeDir []string `json:"excludeDir"`
	ServeDir string `json:"serveDir"`
	Listen string `json:"listen"`
}

func ReadConfigFromFile(filePath string) (Config, error) {
	config := Config {
		ServeDir: ".",
		Source: ".",
		IncludeFiles: []string{"**/*.*"},
		ExcludeDir: []string{".git"},
		Listen: "localhost:2345",
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
