package settings

import (
	"github.com/go-logr/logr"
	"github.com/whitesource/log4j-detect/fs"
	"github.com/whitesource/log4j-detect/operations"
	mavenS "github.com/whitesource/log4j-detect/operations/maven"
	rc "github.com/whitesource/log4j-detect/records"
	mavenQ "github.com/whitesource/log4j-detect/screening/maven"
	"github.com/whitesource/log4j-detect/utils/exec"
)

type MavenResolver struct {
	Disabled       bool
	AdditionalArgs []string
	Scopes         struct {
		Include []string
		Exclude []string
	}
}

func (r MavenResolver) Queries() map[rc.Organ]*fs.Query {
	if r.Disabled {
		return nil
	}

	return map[rc.Organ]*fs.Query{rc.OMaven: mavenQ.Query()}
}

func (r MavenResolver) Surgeons(logger logr.Logger, commander exec.Commander) map[rc.Organ]operations.Surgeon {
	if r.Disabled {
		return nil
	}

	return map[rc.Organ]operations.Surgeon{
		rc.OMaven: mavenS.NewSurgeon(logger, commander, r.AdditionalArgs, r.Scopes.Include, r.Scopes.Exclude),
	}
}
