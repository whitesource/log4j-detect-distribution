package gradle

import (
	"github.com/whitesource/log4j-detect/fs"
	"github.com/whitesource/log4j-detect/fs/match"
)

func Query() *fs.Query {
	return &fs.Query{
		Filename: match.NewExtensionMatcher("jar", "gem"),
	}
}
