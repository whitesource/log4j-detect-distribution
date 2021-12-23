package settings

import (
	"github.com/go-logr/logr"
	"github.com/whitesource/log4j-detect/fs"
	"github.com/whitesource/log4j-detect/operations"
	rubyS "github.com/whitesource/log4j-detect/operations/ruby"
	rc "github.com/whitesource/log4j-detect/records"
	rubyQ "github.com/whitesource/log4j-detect/screening/ruby"
	"github.com/whitesource/log4j-detect/utils/exec"
)

type RubyResolver struct {
	Disabled bool
}

func (r RubyResolver) Queries() map[rc.Organ]*fs.Query {
	if r.Disabled {
		return nil
	}

	return map[rc.Organ]*fs.Query{rc.ORuby: rubyQ.Query()}
}

func (r RubyResolver) Surgeons(logger logr.Logger, commander exec.Commander) map[rc.Organ]operations.Surgeon {
	if r.Disabled {
		return nil
	}

	return map[rc.Organ]operations.Surgeon{
		rc.ORuby: rubyS.NewSurgeon(logger, commander),
	}
}
