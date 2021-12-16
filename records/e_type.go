package records

// LibType Enum type
type LibType int8

const (
	LTEmpty LibType = iota
	LTGolang
	LTJava
	LTFs
)

func (lt LibType) String() string {
	// Order of strings should match the 'const' order
	return [...]string{
		"", "GO", "JAVA", "FS",
	}[lt]
}
