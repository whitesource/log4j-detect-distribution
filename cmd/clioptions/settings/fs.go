package settings

import (
	"github.com/go-logr/logr"
	"github.com/whitesource/log4j-detect/fs"
	"github.com/whitesource/log4j-detect/operations"
	fsop "github.com/whitesource/log4j-detect/operations/fs"
	rc "github.com/whitesource/log4j-detect/records"
	fsscreen "github.com/whitesource/log4j-detect/screening/fs"
	"github.com/whitesource/log4j-detect/utils/exec"
)

type FilesystemResolver struct {
	Disabled bool
}

func (r FilesystemResolver) Queries() map[rc.Organ]*fs.Query {
	if r.Disabled {
		return nil
	}

	return map[rc.Organ]*fs.Query{rc.OFS: fsscreen.Query()}
}

func (r FilesystemResolver) Surgeons(logger logr.Logger, commander exec.Commander) map[rc.Organ]operations.Surgeon {
	if r.Disabled {
		return nil
	}

	return map[rc.Organ]operations.Surgeon{
		rc.OFS: fsop.NewSurgeon(logger, commander),
	}
}
