package screening

import (
	"fmt"
	"github.com/go-logr/logr"
	fs2 "github.com/whitesource/log4j-detect/fs"
	"github.com/whitesource/log4j-detect/fs/match"
	rc "github.com/whitesource/log4j-detect/records"
)

// ScreenDirectory scans the given path for the supported manifest files
// and returns a map from Organ to detected manifest files
func ScreenDirectory(logger logr.Logger, dir string, queryMap map[rc.Organ]*fs2.Query, exclude match.Matcher) (map[rc.Organ][]string, error) {
	logger = logger.WithName("screen").WithValues("directory", dir)

	var queries []*fs2.Query
	for _, q := range queryMap {
		queries = append(queries, q)
	}

	err := fs2.ScanDirectory(dir, queries, exclude)
	if err != nil {
		return nil, fmt.Errorf("failed to screen directory %s: %w", dir, err)
	}

	result := map[rc.Organ][]string{}
	for o, q := range queryMap {
		if len(q.Matched) > 0 {
			logger.Info("found manifest files for organ", "organ", o, "matched", q.Matched)
			result[o] = q.Matched
		}
	}
	return result, nil
}
