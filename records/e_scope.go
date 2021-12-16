package records

// LibScope Enum type
type LibScope int8

const (
	LSDefault LibScope = iota
	LSEmpty
	LSCompile
	LSProvided
	LSRuntime
	LSSystem
	LSDev
	LSOptional
	LSOptionalDev
	LSTest
	LSImport
)

var libScopesText = []string{
	"default", "", "compile", "provided", "runtime", "system", "dev", "optional", "opt-dev", "test", "import",
}

func (ls LibScope) String() string {
	// Order of strings should match the 'const' order
	return libScopesText[ls]
}

func LibScopeByText(s string) LibScope {
	for idx, text := range libScopesText {
		if text == s {
			return LibScope(idx)
		}
	}
	return LSEmpty
}
