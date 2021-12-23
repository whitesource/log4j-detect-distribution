package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileExists checks if a regular file exists
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// EnsurePath ensures the base directory of a given path exists.
// if the directory does not exist, it will be created with the given permissions.
func EnsurePath(path string, mod os.FileMode) {
	dir := filepath.Dir(path)

	if FileExists(dir) {
		return
	}

	if err := os.MkdirAll(dir, mod); err != nil {
		fmt.Printf("Unable to create dir %q %v\n", path, err)
		os.Exit(2)
	}
}

func ReadFileAsSlice(fileToRead string) ([]string, error) {
	lines := make([]string, 0)
	file, err := os.Open(fileToRead)
	if err != nil {
		return lines, err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return lines, err
	}

	return lines, err
}

// CreateTempFile creates a temporary file and returns the path
func CreateTempFile(contents string, pattern string) (string, error) {
	file, err := ioutil.TempFile("", pattern)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	_, err = file.WriteString(contents)
	return file.Name(), err
}
