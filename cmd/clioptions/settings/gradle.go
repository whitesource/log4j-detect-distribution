package settings

import (
	"github.com/go-logr/logr"
	"github.com/whitesource/log4j-detect/fs"
	"github.com/whitesource/log4j-detect/operations"
	gradleS "github.com/whitesource/log4j-detect/operations/gradle"
	rc "github.com/whitesource/log4j-detect/records"
	gradleQ "github.com/whitesource/log4j-detect/screening/gradle"
	"github.com/whitesource/log4j-detect/utils/exec"
)

type GradleResolver struct {
	Disabled       bool
	AdditionalArgs []string
	Configurations struct {
		Include []string
		Exclude []string
	}
}

func (r GradleResolver) Queries() map[rc.Organ]*fs.Query {
	if r.Disabled {
		return nil
	}

	return map[rc.Organ]*fs.Query{rc.OGradle: gradleQ.Query()}
}

func (r GradleResolver) Surgeons(logger logr.Logger, commander exec.Commander) map[rc.Organ]operations.Surgeon {
	if r.Disabled {
		return nil
	}

	return map[rc.Organ]operations.Surgeon{
		rc.OGradle: gradleS.NewSurgeon(logger, commander, r.AdditionalArgs, r.Configurations.Include, r.Configurations.Exclude),
	}
}
