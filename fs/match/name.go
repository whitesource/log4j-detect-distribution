package match

import (
	"io/fs"
	"regexp"
)

// NameMatcher represents an arbitrary test for a file name (without leading directories)
// that returns either true or false.
// the matcher should rely only on the file name, or the file type without performing complex checks,
// as they should be implemented in a Matcher
type NameMatcher func(name string, mode fs.FileMode) bool

// NewSimpleNameMatcher is a NameMatcher that checks if the file name is equal to one of the provided names
func NewSimpleNameMatcher(name ...string) func(name string, _ fs.FileMode) bool {
	names := name
	return func(name string, _ fs.FileMode) bool {
		for _, n := range names {
			if name == n {
				return true
			}
		}
		return false
	}
}

func NewNameRegexMatcher(rs ...*regexp.Regexp) func(name string, _ fs.FileMode) bool {
	return func(name string, _ fs.FileMode) bool {
		for _, r := range rs {
			if r.MatchString(name) {
				return true
			}
		}
		return false
	}
}
