package gradle

type Dependency struct {
	Type          string `json:"type"`
	Version       string `json:"version"`
	GroupId       string `json:"groupId"`
	IsInnerModule bool   `json:"isInnerModule"`
	SystemPath    string `json:"systemPath"`
	ArtifactId    string `json:"artifactId"`
}

type Project struct {
	Dependencies map[string]Dependency  `json:"dependencies"`
	Manifest     string                 `json:"manifest"`
	ConfigGraphs map[string]ConfigGraph `json:"configGraphs"`
	Version      string                 `json:"version"`
	GroupId      string                 `json:"groupId"`
	ArtifactId   string                 `json:"artifactId"`
}

type ConfigGraph struct {
	Deps2Children map[string][]string `json:"deps2Children"`
	DirectDeps    []string            `json:"directDeps"`
}
