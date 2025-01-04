package utils

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
)

func GetWatchableDirs(root string) ([]string, error) {
	var dirs []string
	fileSystem := os.DirFS(root)
	
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err	
		}
		if d.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})

	return dirs, err
}

func DirsIncludePatterns(dirs []string, patterns []string) []string {

	return dirs
}

func DirsExcludePatterns(dirs []string, patterns []string) []string {
	var newDirs []string
	regexPatterns := buildRegex(patterns)
	for _, dir := range dirs {
		if !matchDirPatterns(dir, regexPatterns) {
			newDirs = append(newDirs, dir)
		}
	}
	return newDirs 
}

func buildRegex(patterns []string) []regexp.Regexp {
	var regexPatterns []regexp.Regexp
	for _, pattern := range patterns {
		strPattern := fmt.Sprintf("^%s[/]{0,1}", pattern)
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
