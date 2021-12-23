package gradle

import (
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/go-logr/logr"
	"github.com/whitesource/log4j-detect/records"
	"github.com/whitesource/log4j-detect/utils"
	"github.com/whitesource/log4j-detect/utils/exec"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
)

type Surgeon struct {
	additionalArgs []string
	configIncludes []*regexp.Regexp
	configExcludes []*regexp.Regexp

	commander exec.Commander
	logger    logr.Logger
}

//go:embed assets/init.ws.gradle
var initScriptSource string

func NewSurgeon(logger logr.Logger, commander exec.Commander, additionalArgs, configIncludes, configExcludes []string) *Surgeon {
	return &Surgeon{
		additionalArgs: additionalArgs,
		configIncludes: utils.ConvertToRegexSlice(configIncludes),
		configExcludes: utils.ConvertToRegexSlice(configExcludes),
		commander:      commander,
		logger:         logger.WithValues("surgeon", "gradle"),
	}
}

func (s Surgeon) Validate(paths []string) error {
	dirs := extractProjectDirs(paths)
	for _, d := range dirs {
		binary := s.determineGradlePath(d)
		if binary == "gradle" && !utils.IsInstalled("gradle") {
			return errors.New("gradle is not installed")
		} else {
			return nil
		}
	}
	return nil
}

func (s Surgeon) Operate(paths []string) (results []records.OperationResult, err error) {
	initScriptPath, err := persistInitScript()
	if err != nil {
		return nil, err
	}

	projectDirs := extractProjectDirs(paths)
	manifest2Result := make(map[string]bool)
	for _, dir := range projectDirs {
		if fullDir, er := filepath.Abs(dir); er == nil {
			dir = fullDir
		}

		dirOpResults, opErr := s.singleDirOperation(dir, initScriptPath)
		if opErr != nil {
			s.logger.Error(opErr, "operation failed", "dir", dir)
			results = append(results, records.OperationResult{
				ManifestFile: dir,
				Err:          opErr,
			})
			continue
		}

		for _, or := range dirOpResults {
			if _, exists := manifest2Result[or.ManifestFile]; !exists {
				manifest2Result[or.ManifestFile] = true
				results = append(results, or)
			}
		}
	}
	return
}

func (s Surgeon) singleDirOperation(dir, initScriptPath string) (results []records.OperationResult, err error) {
	// TODO handle temp files
	outputFile, err := ioutil.TempFile("", "gradle")
	if err != nil {
		return
	}

	err = s.executeInitScript(initScriptPath, dir, outputFile.Name())
	if err != nil {
		return
	}

	var projects []Project
	err = json.NewDecoder(outputFile).Decode(&projects)
	if err != nil {
		return
	}

	for _, p := range projects {
		p.ConfigGraphs = s.cleanConfigurations(p.ConfigGraphs)
		r := buildOpResult(p)
		if r.Direct != nil {
			results = append(results, r)
		}
	}

	return
}

func (s Surgeon) executeInitScript(initScriptPath, projectDir, outputPath string) error {
	var args []string
	args = append(args, s.additionalArgs...)
	args = append(args, "--init-script", initScriptPath,
		"-DWS_INIT_GRADLE_OUTPUT_PATH="+outputPath,
		"-Dorg.gradle.parallel=",
		"-DWS_INIT_GRADLE_INCLUDE_DEPENDENCIES=true",
		"whitesourceDependenciesTask",
	)

	binary := s.determineGradlePath(projectDir)

	cmd := s.commander.Command(s.logger, binary, args...)
	cmd.SetDir(projectDir)
	return cmd.Run()
}

// determineGradlePath checks if the gradle wrapper exists in the project directory
// if it does, the wrapper path is returned - otherwise gradle is returned
func (s Surgeon) determineGradlePath(projectDir string) (gradlePath string) {
	defer s.logger.Info("determined gradle path", "projectDir", projectDir, "gradlePath", gradlePath)

	gradlew := filepath.Join(projectDir, "gradlew")
	if runtime.GOOS == "windows" {
		gradlew += ".bat"
	}

	if utils.FileExists(gradlew) {
		gradlePath, _ = filepath.Abs(gradlew)
	} else {
		gradlePath = "gradle"
	}

	return
}

func (s Surgeon) cleanConfigurations(graphs map[string]ConfigGraph) map[string]ConfigGraph {
	result := make(map[string]ConfigGraph)

	for key, value := range graphs {
		if len(s.configExcludes) > 0 && utils.RegexMatch(key, s.configExcludes) {
			continue
		}
		if len(s.configIncludes) == 0 || utils.RegexMatch(key, s.configIncludes) {
			result[key] = value
		}
	}

	return result
}

// extractProjectDirs extracts all root gradle project directories that need to be scanned,
// excluding build.gradle files with a parent directory containing a settings.gradle(.kts) file
func extractProjectDirs(manifests []string) []string {
	settingsDirs := extractSettingsDirs(manifests)

	dirs := make(map[string]bool)
	for _, sd := range settingsDirs {
		dirs[sd] = true
	}

	for _, m := range manifests {
		if dirs[filepath.Dir(m)] {
			continue
		}

		isSub := false
		for _, sd := range settingsDirs {
			if isSubDir(filepath.Dir(m), sd) {
				isSub = true
				break
			}
		}

		if !isSub {
			dirs[filepath.Dir(m)] = true
		}
	}
	return toStringList(dirs)
}

// isSubDir checks if child is a subdirectory of parent
func isSubDir(child string, parent string) bool {
	var prev string
	for prev != child {
		if parent == child {
			return true
		}
		prev = child
		child = path.Dir(child)
	}

	return false
}

// extractSettingsDirs locates all directories containing a settings.gradle(.kts) file
func extractSettingsDirs(manifests []string) []string {
	dirs := make(map[string]bool)
	for _, m := range manifests {
		name := filepath.Base(m)
		if name == "settings.gradle" || name == "settings.gradle.kts" {
			dirs[filepath.Dir(m)] = true
		}
	}
	return toStringList(dirs)
}

func toStringList(strMap map[string]bool) (values []string) {
	for k := range strMap {
		values = append(values, k)
	}
	return
}

// persistInitScript saves the gradle init script to a temp file and returns the path
func persistInitScript() (string, error) {
	return utils.CreateTempFile(initScriptSource, "gradle")
}

func buildOpResult(project Project) records.OperationResult {
	return records.OperationResult{
		ManifestFile:      project.Manifest,
		Direct:            extractDirectDeps(project),
		LibraryToChildren: extractLibraryToChildren(project),
		Libraries:         extractLibraries(project),
		LType:             records.LTJava,
		Err:               nil,
	}
}

func extractLibraries(project Project) *map[records.Id]records.Library {
	libraries := make(map[records.Id]records.Library)
	for id, dep := range project.Dependencies {
		if !isValidGradleDependency(dep) {
			continue
		}
		libraries[records.Id(id)] = extractLibrary(dep)
	}
	return &libraries
}

func isValidGradleDependency(dep Dependency) bool {
	if dep.IsInnerModule {
		return true
	}
	if dep.GroupId != "" && dep.ArtifactId != "" && dep.Version != "" {
		return true
	}
	return dep.SystemPath != "" && utils.FileExists(dep.SystemPath)
}

func extractLibrary(dep Dependency) records.Library {
	return records.Library{
		LType:      records.LTJava,
		GroupId:    dep.GroupId,
		Artifact:   dep.ArtifactId,
		Version:    dep.Version,
		SystemPath: dep.SystemPath,
		IsProject:  dep.IsInnerModule,
	}
}

func extractLibraryToChildren(project Project) *map[records.Id][]records.Id {
	set := make(map[string]map[string]bool)
	for _, graph := range project.ConfigGraphs {
		for k, v := range graph.Deps2Children {
			if _, ok := set[k]; !ok {
				set[k] = make(map[string]bool)
			}

			for _, d := range v {
				set[k][d] = true
			}
		}
	}

	libraryToChildren := make(map[records.Id][]records.Id)
	for k, v := range set {
		libraryToChildren[records.Id(k)] = idSet2List(v)
	}
	return &libraryToChildren
}

func extractDirectDeps(project Project) *[]records.Id {
	directSet := make(map[string]bool)
	for _, graph := range project.ConfigGraphs {
		for _, d := range graph.DirectDeps {
			directSet[d] = true
		}
	}

	directs := idSet2List(directSet)
	return &directs
}

func idSet2List(ids map[string]bool) []records.Id {
	var list []records.Id
	for id := range ids {
		list = append(list, records.Id(id))
	}
	return list
}
