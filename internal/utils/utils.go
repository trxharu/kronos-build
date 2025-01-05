package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else if errors.Is(err, fs.ErrNotExist) {
		return false 
	} else {
		return false
	}
}

func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if info.IsDir() {
		return true, nil
	} else {
		return false, nil
	}
}

func GetWatchableDirs(root string) ([]string, error) {
	var dirs []string
	absPath, _ := filepath.Abs(root)
	fileSystem := os.DirFS(absPath)
	
	err := fs.WalkDir(fileSystem, ".", func(src string, d fs.DirEntry, err error) error {
		if err != nil {
			return err	
		}
		if d.IsDir() {
			dirs = append(dirs, path.Join(absPath, src))
		}
		return nil
	})

	return dirs, err
}

func DirsExcludePatterns(dirs []string, patterns []string) []string {
	var newDirs []string
	regexPatterns := buildRegex(patterns)
	for _, dir := range dirs {
		absPath, _ := filepath.Abs(dir)
		if !matchDirPatterns(absPath, regexPatterns) {
			newDirs = append(newDirs, absPath)
		}
	}
	return newDirs 
}

func buildRegex(patterns []string) []regexp.Regexp {
	var regexPatterns []regexp.Regexp

	for _, pattern := range patterns {
		prefix := ".*"
		if filepath.IsAbs(pattern) {
			prefix = ""
		}
		strPattern := fmt.Sprintf("%s%s[/]{0,1}", prefix, pattern)
		regexPatterns = append(regexPatterns, *regexp.MustCompile(strPattern))
	}
	return regexPatterns
}

func matchDirPatterns(path string, patterns []regexp.Regexp) bool {
	for _, pattern := range patterns {
		if pattern.MatchString(path) {
			return true
		}
	}
	return false
}
