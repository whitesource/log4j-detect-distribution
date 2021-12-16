package fs

import (
	"fmt"
	"github.com/karrick/godirwalk"
	"github.com/whitesource/log4j-detect/fs/match"
)

// ScanDirectory walks the tree under the provided directory recursively,
// and applies the provided queries to each file.
// If exclude matches a directory, it is skipped (and all files within it are ignored)
// If exclude matches a file, it will not be matched by any of the queries
func ScanDirectory(dir string, queries []*Query, exclude match.Matcher) error {
	if queries == nil {
		return nil
	}

	return godirwalk.Walk(
		dir,
		&godirwalk.Options{
			Unsorted:      true,
			ErrorCallback: walkErrorFunc,
			Callback: func(osPathname string, de *godirwalk.Dirent) error {
				return walkFunc(osPathname, de, queries, exclude)
			},
		},
	)
}

func walkFunc(path string, de *godirwalk.Dirent, queries []*Query, exclude match.Matcher) error {
	name := de.Name()
	mode := de.ModeType()

	if exclude != nil && exclude.Match(name, path, mode) {
		if de.IsDir() {
			return godirwalk.SkipThis
		} else {
			return nil
		}
	}

	for _, q := range queries {
		q.Match(name, path, mode)
	}

	return nil
}

func walkErrorFunc(_ string, err error) godirwalk.ErrorAction {
	// TODO: change to debug log file
	fmt.Println("fsWalkErrorFunc error: ", err)
	return godirwalk.SkipNode
}
