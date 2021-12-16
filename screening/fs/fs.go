package gradle

import (
	"github.com/whitesource/log4j-detect/fs"
	"github.com/whitesource/log4j-detect/fs/match"
	"regexp"
)

func Query() *fs.Query {
	return &fs.Query{
		Filename: match.NewNameRegexMatcher(
			regexp.MustCompile("^.*\\.jar$"),
		),
	}
}
