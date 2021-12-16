package fs

import (
	match2 "github.com/whitesource/log4j-detect/fs/match"
	"os"
)

// Query represents a query for file system files.
// It matches files by their names with a match.FilenameMatcher,
// and can include an optional Matcher for arbitrary verifications
type Query struct {
	// Filename should match the file name (without parent directories)
	// if nil, all file names will match
	Filename match2.NameMatcher

	// Additional validations for the file
	// e.g.
	//  1. check file is not empty
	//  2. check if another file exists near this file match
	//     (take package-lock.json only if there is a package.json beside it)
	//  3. check that file is an executable
	Matcher match2.Matcher

	// list of paths matched
	Matched []string
}

// Match checks the following:
// * the Filename function matches the file name (if nil then this check passes)
// * the Matcher returns true when executed on name, path and mode (if it's nil then this check passes)
// If the above checks pass, then the full path is added to Matched, and true is returned.
// Otherwise, the path is discarded and false is returned.
func (q *Query) Match(name, path string, mode os.FileMode) bool {
	if q.Filename != nil && !q.Filename(name, mode) {
		return false
	}

	if q.Matcher != nil && !q.Matcher.Match(name, path, mode) {
		return false
	}

	q.Matched = append(q.Matched, path)
	return true
}
