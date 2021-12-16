package maven

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMavenManifestNameGoodMatch(t *testing.T) {
	validNames := []string{
		"pom.xml",
	}

	q := Query()
	for _, name := range validNames {
		assert.True(t, q.Match(name, "", 0))
	}
}

func TestMavenManifestNameNotMatch(t *testing.T) {
	invalidNames := []string{
		"ppom.xml",
		".pom.xml",
		"pom.xm",
		"om.xml",
		"pom.xml.",
		"pom.xml.b",
		"test_pom.xml",
		"pom-xml",
	}

	q := Query()
	for _, name := range invalidNames {
		assert.False(t, q.Match(name, "", 0))
	}
}
