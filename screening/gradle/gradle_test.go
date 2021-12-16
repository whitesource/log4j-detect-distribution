package gradle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGradleManifestNameGoodMatch(t *testing.T) {
	validNames := []string{
		"build.gradle",
		"build.gradle.kts",
		"settings.gradle",
		"settings.gradle.kts",
	}

	q := Query()
	for _, name := range validNames {
		assert.True(t, q.Match(name, "", 0))
	}
}

func TestGradleManifestNameNotMatch(t *testing.T) {
	invalidNames := []string{
		"bbuild.gradle",
		".build.gradle",
		"uild.gradle",
		"build.gradl",
		"build.gradle.",
		"build.gradle.k",
		"build.gradle.kt",
		"build.gradle.ktss",
		"build.gradle.ks",
		"build.gradle.ts",
		"buildbgradle",
		".settings.gradle.kts",
		".settings.gradle",
		"gradle.settings",
		"setting.gradle",
	}

	q := Query()
	for _, name := range invalidNames {
		assert.False(t, q.Match(name, "", 0))
	}
}
