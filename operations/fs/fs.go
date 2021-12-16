package maven

import (
	"github.com/go-logr/logr"
	"github.com/whitesource/log4j-detect/records"
	"github.com/whitesource/log4j-detect/utils/exec"
)

type Surgeon struct {
	commander exec.Commander
	logger    logr.Logger
}

func NewSurgeon(logger logr.Logger, commander exec.Commander) *Surgeon {
	return &Surgeon{
		commander: commander,
		logger:    logger.WithValues("surgeon", "fs"),
	}
}

func (s Surgeon) Validate(_ []string) error {
	return nil
}

func (s Surgeon) Operate(paths []string) ([]records.OperationResult, error) {
	return []records.OperationResult{
		{
			ManifestFile: "",
			Direct:       s.asDirect(paths),
			Libraries:    s.asLibs(paths),
			LType:        records.LTFs,
			Organ:        records.OFS,
		},
	}, nil
}

func (s Surgeon) asDirect(paths []string) *[]records.Id {
	var ids []records.Id
	for _, p := range paths {
		ids = append(ids, records.Id(p))
	}
	return &ids
}

func (s Surgeon) asLibs(paths []string) *map[records.Id]records.Library {
	id2Lib := map[records.Id]records.Library{}
	for _, p := range paths {
		id2Lib[records.Id(p)] = records.Library{
			Artifact:   p,
			Version:    "",
			LType:      records.LTFs,
			SystemPath: p,
		}
	}
	return &id2Lib
}
