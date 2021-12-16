package settings

import (
	"github.com/go-logr/logr"
	"github.com/whitesource/log4j-detect/fs"
	"github.com/whitesource/log4j-detect/operations"
	"github.com/whitesource/log4j-detect/records"
	"github.com/whitesource/log4j-detect/utils/exec"
)

func mergeQueries(resolvers ...Resolver) map[records.Organ]*fs.Query {
	m := map[records.Organ]*fs.Query{}
	for _, r := range resolvers {
		qs := r.Queries()
		if qs == nil {
			continue
		}

		for k, v := range qs {
			m[k] = v
		}
	}
	return m
}

func mergeSurgeons(logger logr.Logger, commander exec.Commander, resolvers ...Resolver) map[records.Organ]operations.Surgeon {
	m := map[records.Organ]operations.Surgeon{}
	for _, r := range resolvers {
		surgeons := r.Surgeons(logger, commander)
		if surgeons == nil {
			continue
		}

		for k, v := range surgeons {
			m[k] = v
		}
	}
	return m
}
