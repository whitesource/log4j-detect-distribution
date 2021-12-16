package settings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSettings_GlobalExcludes(t *testing.T) {
	s := Settings{}
	assert.Equal(t, s.GlobalExcludes(), defaultExcludes)

	s = Settings{
		Excludes: []string{},
	}
	assert.NotEqual(t, s.GlobalExcludes(), defaultExcludes)

	s = Settings{
		Excludes: []string{"^\\.git$", "^node_modules$"},
	}
	exclude := s.GlobalExcludes()
	assert.True(t, exclude.Match(".git", "", 0))
	assert.True(t, exclude.Match("node_modules", "", 0))
}
