package gradle

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func Test_extractProjectDirs(t *testing.T) {
	type args struct {
		paths []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			args: args{
				paths: []string{
					"/a/settings.gradle",
					"/a/build.gradle",
					"/a/list/build.gradle",
					"/a/projects/build.gradle",
					"/a/utilities/build.gradle",
					"/b/build.gradle",
					"/c/settings.gradle",
					"/c/d/build.gradle",
				},
			},
			want: []string{"/a", "/b", "/c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractProjectDirs(tt.args.paths)
			sort.Strings(got)
			sort.Strings(tt.want)
			assert.EqualValues(t, tt.want, got)
		})
	}
}
