package settings

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"github.com/whitesource/log4j-detect/fs"
	"github.com/whitesource/log4j-detect/fs/match"
	"github.com/whitesource/log4j-detect/operations"
	"github.com/whitesource/log4j-detect/records"
	"github.com/whitesource/log4j-detect/utils/exec"
	"regexp"
)

type Flags struct {
	mavenOnly  bool
	gradleOnly bool
}

func (f *Flags) ToSettings(logger logr.Logger) (*Settings, error) {
	if f.mavenOnly && f.gradleOnly {
		return nil, fmt.Errorf("bad")
	}

	s := &Settings{
		Resolvers: Resolvers{
			Gradle: GradleResolver{
				Disabled: f.mavenOnly,
			},
			Maven: MavenResolver{
				Disabled: f.gradleOnly,
			},
			Fs: FilesystemResolver{
				Disabled: false,
			},
		},
		logger: logger,
	}
	return s, nil
}

func AddFlags(cmd *cobra.Command, f *Flags) {
	cmd.Flags().BoolVarP(&f.mavenOnly, "maven-only", "m", false, "only scan for maven projects")
	cmd.Flags().BoolVarP(&f.gradleOnly, "gradle-only", "g", false, "only scan for gradle projects")
}

// Settings represents all settings.
// this includes logging parameters and other parameters that may modify the
// behavior of resolvers for different package managers.
// It does not contain authentication information - this should be stored in a profile or passed as arguments.
type Settings struct {
	Resolvers Resolvers
	Excludes  []string
	logger    logr.Logger
}

type Resolvers struct {
	Gradle GradleResolver
	Maven  MavenResolver
	Fs     FilesystemResolver
}

type Resolver interface {
	Queries() map[records.Organ]*fs.Query
	Surgeons(logger logr.Logger, commander exec.Commander) map[records.Organ]operations.Surgeon
}

var defaultExcludes = match.Or(
	match.FilenameRegex(regexp.MustCompile("^target$")),
)

func (s *Settings) GlobalExcludes() match.Matcher {
	if s.Excludes == nil {
		return defaultExcludes
	}

	if len(s.Excludes) == 0 {
		return nil
	}

	var excludeMatchers []match.Matcher
	for _, e := range s.Excludes {
		if r, err := regexp.Compile(e); err == nil {
			excludeMatchers = append(excludeMatchers, match.FilenameRegex(r))
		} else {
			s.logger.Info("invalid regex syntax found in excludes", "regex", e)
		}
	}
	return match.Or(excludeMatchers...)
}

func (r *Resolvers) ManifestQueries() map[records.Organ]*fs.Query {
	return mergeQueries(r.Maven, r.Gradle, r.Fs)
}

func (r *Resolvers) Surgeons(logger logr.Logger, commander exec.Commander) map[records.Organ]operations.Surgeon {
	return mergeSurgeons(logger, commander, r.Maven, r.Gradle, r.Fs)
}
