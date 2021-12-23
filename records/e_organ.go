package records

// Organ Enum type
type Organ int8

const (
	OMaven Organ = iota
	OGradle
	OFS
	ORuby
)

func (organ Organ) String() string {
	// Order of strings should match the 'const' order
	return []string{
		"maven",
		"gradle",
		"file_system",
		"ruby",
	}[organ]
}
