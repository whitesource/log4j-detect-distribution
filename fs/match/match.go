package match

import (
	"github.com/whitesource/log4j-detect/utils"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

// Matcher represents an arbitrary test for a file that returns either true or false.
// It can perform some check on the file's contents, path, permissions etc.
type Matcher interface {
	Match(name, path string, mode fs.FileMode) bool
}

// And is a Matcher that returns true only if all Matchers match
func And(matchers ...Matcher) Matcher {
	return and{Matchers: matchers}
}

type and struct {
	Matchers []Matcher
}

func (a and) Match(name, path string, mode fs.FileMode) bool {
	if a.Matchers == nil {
		return true
	}

	for _, m := range a.Matchers {
		if !m.Match(name, path, mode) {
			return false
		}
	}

	return true
}

// Or is a Matcher that returns true only if at least one of the Matchers matches
func Or(matchers ...Matcher) Matcher {
	return or{Matchers: matchers}
}

type or struct {
	Matchers []Matcher
}

func (o or) Match(name, path string, mode fs.FileMode) bool {
	if o.Matchers == nil {
		return true
	}

	for _, m := range o.Matchers {
		if m.Match(name, path, mode) {
			return true
		}
	}

	return false
}

// Func is a Matcher that performs an arbitrary match (i.e. the match is delegated to F)
func Func(f func(name, path string, mode fs.FileMode) bool) Matcher {
	return _func{F: f}
}

type _func struct {
	F func(name, path string, mode fs.FileMode) bool
}

func (f _func) Match(name, path string, mode os.FileMode) bool {
	return f.F(name, path, mode)
}

// Sibling is a Matcher that returns true only if the file (name) exists
// in the provided path's directory (path is the argument passed to Match)
func Sibling(name string) Matcher {
	return sibling{Name: name}
}

type sibling struct {
	Name string
}

func (s sibling) Match(_, path string, _ os.FileMode) bool {
	return utils.FileExists(filepath.Join(filepath.Dir(path), s.Name))
}

// FilenameRegex is a Matcher that returns true if the file name (without leading directories) matches the regex
func FilenameRegex(regex *regexp.Regexp) Matcher {
	return filenameRegex{regex: regex}
}

type filenameRegex struct {
	regex *regexp.Regexp
}

func (fr filenameRegex) Match(name, _ string, _ os.FileMode) bool {
	return fr.regex == nil || fr.regex.MatchString(name)
}
