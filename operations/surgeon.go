package operations

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/whitesource/log4j-detect/records"
	"os"
)

type Surgeon interface {
	Validate(paths []string) error

	// Operate confirms the availability of all requirements
	// it processes the manifest files to generate the []utils.OperationResult
	Operate(paths []string) ([]records.OperationResult, error)
}

// Perform performs the dependency resolution for all matching manifest files, with the corresponding Surgeons
// detected is a map from utils.Organ to a list of matching manifest files
// return: a list of scan result per manifest file
//  (might be less for multi-module projects)
//  A scan result will be returned only in a case a Surgeon matching the utils.Organ is found
func Perform(logger logr.Logger, detected map[records.Organ][]string, surgeons map[records.Organ]Surgeon) (results []records.OperationResult) {
	for o, s := range surgeons {
		paths := detected[o]
		if len(paths) == 0 {
			continue
		}

		if err := s.Validate(paths); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error: %s project detected (manifest: %s), but an error occurred: %v\n", o, paths[0], err)
			continue
		}

		r, err := s.Operate(paths)
		if err != nil {
			logger.Error(err, "failed to scan projects", "projectType", o)
			continue
		}

		results = append(results, r...)
	}
	return results
}
