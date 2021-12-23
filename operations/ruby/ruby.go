package ruby

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/whitesource/log4j-detect/records"
	"github.com/whitesource/log4j-detect/utils"
	"github.com/whitesource/log4j-detect/utils/exec"
	"os"
	"path/filepath"
	"strings"
)

//go:embed assets/gem_dependencies.rb
var gemDependenciesScript string

type Surgeon struct {
	commander exec.Commander
	logger    logr.Logger
}

func NewSurgeon(logger logr.Logger, commander exec.Commander) *Surgeon {
	return &Surgeon{
		commander: commander,
		logger:    logger.WithValues("surgeon", "ruby"),
	}
}

func (s Surgeon) Validate(_ []string) error {
	if !utils.IsInstalled("ruby") {
		return errors.New("ruby is not installed")
	}
	if !utils.IsInstalled("gem") {
		return errors.New("gem is not installed")
	}
	return nil
}

func (s Surgeon) Operate(paths []string) ([]records.OperationResult, error) {
	var results []records.OperationResult

	gemCacheDirs, err := s.discoverGemCacheDirs()
	if err != nil {
		return nil, fmt.Errorf("failed to get path to ruby gems: %w", err)
	}

	for _, p := range paths {
		result := s.singleProjectOperation(p, gemCacheDirs)
		s.warnIfSystemPathMissing(result)
		results = append(results, result)
	}

	return results, nil
}

func (s Surgeon) singleProjectOperation(lockFilePath string, gemCacheDirs []string) (result records.OperationResult) {
	deps, err := s.parseGemLock(lockFilePath)
	if err != nil {
		result.Err = fmt.Errorf("failed to parse lock file %s: %w", lockFilePath, err)
	}

	return records.OperationResult{
		ManifestFile:      lockFilePath,
		Direct:            s.extractDirects(deps),
		LibraryToChildren: s.extractLibrary2Children(deps),
		Libraries:         s.extractLibraries(deps, gemCacheDirs),
		LType:             records.LTRuby,
		Err:               nil,
		Organ:             records.ORuby,
	}
}

func (s Surgeon) discoverGemCacheDirs() ([]string, error) {
	output, err := s.commander.
		Command(s.logger, "gem", "environment", "gempath").
		Output()

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve ruby gem dir - %w", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("failed to retrieve ruby gem dir, unexpected output %s", output)
	}

	gemPath := lines[0]

	var cacheDirs []string
	for _, dir := range strings.Split(gemPath, ":") {
		cacheDir := filepath.Join(dir, "cache")
		if utils.DirExists(dir) {
			cacheDirs = append(cacheDirs, cacheDir)
		}
	}

	if len(cacheDirs) == 0 {
		return nil, fmt.Errorf("failed to find ruby gem cache directory from gempath %s", gemPath)
	}

	return cacheDirs, nil
}

func (s Surgeon) parseGemLock(path string) (*GemDependencies, error) {
	scriptPath, err := persistDepScript()
	if err != nil {
		return nil, fmt.Errorf("failed to create gem dependencies script file: %w", err)
	}

	cmd := s.commander.Command(s.logger, "ruby", scriptPath, path)
	cmd.SetDir(filepath.Dir(path))

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute gem dependencies script on %s: %w", path, err)
	}

	var deps GemDependencies
	err = json.Unmarshal(output, &deps)
	if err != nil {
		return nil, fmt.Errorf("failed to parse output of gem dependencies script %w", err)
	}

	return &deps, nil
}

func (s Surgeon) extractDirects(lock *GemDependencies) *[]records.Id {
	var directs []records.Id
	for _, d := range lock.Directs {
		directs = append(directs, records.Id(d))
	}
	return &directs
}

func (s Surgeon) extractLibraries(lock *GemDependencies, gemCacheDirs []string) *map[records.Id]records.Library {
	libs := map[records.Id]records.Library{}
	for id, dep := range lock.Dependencies {
		libs[records.Id(id)] = records.Library{
			Artifact:   dep.Name,
			Version:    dep.Version,
			LType:      records.LTRuby,
			SystemPath: s.systemPath(dep, gemCacheDirs),
		}
	}
	return &libs
}

func (s Surgeon) systemPath(dep Dependency, gemCacheDirs []string) string {
	filePattern := fmt.Sprintf("%s-%s*.gem", dep.Name, dep.Version)
	for _, dir := range gemCacheDirs {
		pattern := filepath.Join(dir, filePattern)
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) > 0 {
			return matches[0]
		}
	}
	return ""
}

func (s Surgeon) extractLibrary2Children(lock *GemDependencies) *map[records.Id][]records.Id {
	libToChildren := map[records.Id][]records.Id{}
	for id, children := range lock.DepsToChildren {
		var childrenIds []records.Id
		for _, c := range children {
			childrenIds = append(childrenIds, records.Id(c))
		}
		libToChildren[records.Id(id)] = childrenIds
	}
	return &libToChildren
}

func (s Surgeon) warnIfSystemPathMissing(result records.OperationResult) {
	for _, d := range *result.Libraries {
		if d.Artifact == "bundler" {
			continue
		}

		if !utils.FileExists(d.SystemPath) {
			fmt.Fprintf(os.Stderr, "warn: %s has missing gems. Please run `bundle install`.", result.ManifestFile)
			break
		}
	}
}

// persistDepScript saves the gem_dependencies.rb script to a temp file and returns the path
func persistDepScript() (string, error) {
	return utils.CreateTempFile(gemDependenciesScript, "ruby")
}
