package maven

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/creekorful/mvnparser"
	"github.com/go-logr/logr"
	"github.com/whitesource/log4j-detect/records"
	"github.com/whitesource/log4j-detect/utils"
	"github.com/whitesource/log4j-detect/utils/exec"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const (
	// regex to catch the dependencies lice of the command:
	// mvn dependency:list -DoutputAbsoluteArtifactFilename=true
	rgxStr = "\\[INFO\\]\\s+([^:]*):([^:]*):([\\w]*):([^:]*):([\\w]*):([\\S]*)[\\s]?"

	// indexes of part to capture group
	iGroup    = 1
	iArtifact = 2
	iType     = 3
	iVersion  = 4
	iScope    = 5
	iPath     = 6
)

var lineRegex = regexp.MustCompile(rgxStr)

type Surgeon struct {
	additionalArgs []string
	scopeIncludes  []*regexp.Regexp
	scopeExcludes  []*regexp.Regexp

	commander exec.Commander
	logger    logr.Logger
}

func NewSurgeon(logger logr.Logger, commander exec.Commander, additionalArgs, scopeIncludes, scopeExcludes []string) *Surgeon {
	return &Surgeon{
		additionalArgs: additionalArgs,
		scopeIncludes:  utils.ConvertToRegexSlice(scopeIncludes),
		scopeExcludes:  utils.ConvertToRegexSlice(scopeExcludes),
		commander:      commander,
		logger:         logger.WithValues("surgeon", "maven"),
	}
}

func (s Surgeon) Validate(_ []string) error {
	if !utils.IsInstalled("mvn") {
		return errors.New("mvn is not installed")
	}
	return nil
}

func (s Surgeon) Operate(paths []string) ([]records.OperationResult, error) {
	var results []records.OperationResult

	if err := prepOperationsRoom(); err != nil {
		return nil, err
	}

	for _, f := range paths {
		project, err := readPomXml(f)
		if err != nil {
			continue
		}
		if project.Packaging == "pom" {
			// skip aggregator pom
			continue
		}

		results = append(results, s.singleProjectOperation(f))
	}

	return results, nil
}

// TODO validate maven installed
func prepOperationsRoom() (err error) {
	return nil
}

func (s Surgeon) singleProjectOperation(pomPath string) (result records.OperationResult) {
	result.ManifestFile, _ = filepath.Abs(pomPath)
	result.LType = records.LTJava

	withTransitive, noTransitive, err := s.execAllMavenList(pomPath)
	if err != nil {
		result.Err = err
		return
	}

	result.Libraries = s.parseMavenDepsList(withTransitive)
	result.LibraryToChildren = s.createLibrary2Children(result.Libraries)
	result.Direct = s.findProjectDirects(noTransitive, result)
	return result
}

func (s Surgeon) execAllMavenList(pomPath string) (withTransitive, noTransitive []string, err error) {
	var err1, err2 error
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		withTransitive, err1 = s.execMavenList(pomPath, true, s.additionalArgs)
	}()
	go func() {
		defer wg.Done()
		noTransitive, err2 = s.execMavenList(pomPath, false, s.additionalArgs)
	}()
	wg.Wait()

	if err1 != nil {
		return nil, nil, err1
	}
	if err2 != nil {
		return nil, nil, err2
	}
	return
}

func (s Surgeon) execMavenList(path string, transitive bool, addArgs []string) (outLines []string, err error) {
	args := []string{"dependency:list",
		"-DoutputAbsoluteArtifactFilename=true",
		fmt.Sprintf("-DexcludeTransitive=%t", !transitive),
	}
	args = append(args, addArgs...)

	cmd := s.commander.Command(s.logger, "mvn", args...)
	cmd.SetDir(filepath.Dir(path))
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve maven dependency list: %w", err)
	}

	outLines = strings.Split(string(output), "\n")
	return
}

func (s Surgeon) parseMavenDepsList(lines []string) *map[records.Id]records.Library {
	result := map[records.Id]records.Library{}

	for _, line := range lines {
		l, ok := s.matchLineAndBuildLibrary(line)
		if !ok {
			continue
		}
		id := generateId(l.GroupId, l.Artifact)
		result[id] = l
	}

	return &result
}

func (s Surgeon) matchLineAndBuildLibrary(line string) (records.Library, bool) {
	groups := lineRegex.FindStringSubmatch(line)
	if len(groups) != 7 {
		return records.Library{}, false
	}

	scope := groups[iScope]
	if len(s.scopeExcludes) > 0 && utils.RegexMatch(scope, s.scopeExcludes) {
		return records.Library{}, false
	}
	if len(s.scopeIncludes) > 0 && !utils.RegexMatch(scope, s.scopeIncludes) {
		return records.Library{}, false
	}

	l := records.Library{
		Artifact:   groups[iArtifact],
		Version:    groups[iVersion],
		LType:      records.LTJava,
		LScope:     records.LibScopeByText(scope),
		SystemPath: groups[iPath],
		GroupId:    groups[iGroup],
	}
	return l, true
}

func (s Surgeon) createLibrary2Children(libraries *map[records.Id]records.Library) *map[records.Id][]records.Id {
	result := map[records.Id][]records.Id{}

	for id, lib := range *libraries {
		sp := lib.SystemPath
		pom := sp[0:len(sp)-len(filepath.Ext(sp))] + ".pom"
		if !utils.FileExists(pom) {
			s.logger.Info("pom file not found", "path", pom)
			continue
		}
		prj, err := readPomXml(pom)
		if err != nil {
			s.logger.Error(err, "failed to parse pom.xml file", "path", pom)
			continue
		}

		children := extractChildren(prj, libraries)
		if len(children) == 0 {
			continue
		}
		result[id] = children
	}

	return &result
}

func (s Surgeon) findProjectDirects(mavenDepsList []string, opr records.OperationResult) *[]records.Id {
	// find direct according to maven
	notTransitive := s.parseMavenDepsList(mavenDepsList)

	// build a map with all libraries that appears as child of any other library
	allChildren := map[records.Id]bool{}
	for _, children := range *opr.LibraryToChildren {
		for _, c := range children {
			allChildren[c] = true
		}
	}

	// add any library that does not have a parent as direct
	// so far important only for relocated libraries
	for id := range *opr.Libraries {
		if _, found := allChildren[id]; !found {
			// value not important, need the key only
			(*notTransitive)[id] = records.Library{}
		}
	}

	var direct []records.Id
	for key := range *notTransitive {
		direct = append(direct, key)
	}

	return &direct
}

// prj: the library pom to look for its transitive libraries
// libraries: all libraries available in the class path of the main project
// returns:
//  the list of children needed by the library prj
//  filtering out libraries that are not part of the class path
func extractChildren(prj mvnparser.MavenProject, libraries *map[records.Id]records.Library) (children []records.Id) {
	for _, dep := range prj.Dependencies {
		child := generateId(dep.GroupId, dep.ArtifactId)
		if _, found := (*libraries)[child]; found {
			children = append(children, child)
		}
	}
	return children
}

// generates the utils.Id key
func generateId(group, artifact string) records.Id {
	return records.Id(fmt.Sprintf("%s:%s", group, artifact))
}

func readPomXml(pomPath string) (result mvnparser.MavenProject, err error) {
	bs, err := ioutil.ReadFile(pomPath)
	if err != nil {
		return
	}

	reader := bytes.NewReader(bs)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	if err = decoder.Decode(&result); err != nil {
		return
	}
	return result, nil
}
