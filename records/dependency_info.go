package records

type DependencyInfo struct {
	GroupID        string             `json:"groupId"`
	ArtifactID     string             `json:"artifactId"`
	Version        string             `json:"version"`
	Type           string             `json:"type"`
	Classifier     string             `json:"classifier,omitempty"`
	Scope          string             `json:"scope,omitempty"`
	Sha1           string             `json:"sha1,omitempty"`
	AdditionalSha1 string             `json:"additionalSha1,omitempty"`
	SystemPath     string             `json:"systemPath,omitempty"`
	Filename       string             `json:"filename,omitempty"`
	DependencyType string             `json:"dependencyType,omitempty"`
	DependencyFile string             `json:"dependencyFile,omitempty"`
	Deduped        bool               `json:"deduped"`
	Commit         string             `json:"commit,omitempty"`
	Checksums      *map[string]string `json:"checksums,omitempty"`
	Children       *[]DependencyInfo  `json:"children,omitempty"`
}
