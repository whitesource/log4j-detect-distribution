package ruby

// types for parsing output of gem dependencies script
type (
	GemDependencies struct {
		Directs        []string              `json:"directDependencies"`
		DepsToChildren map[string][]string   `json:"depsToChildren"`
		Dependencies   map[string]Dependency `json:"dependencies"`
	}

	Dependency struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
)
