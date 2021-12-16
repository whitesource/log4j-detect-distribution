package hash

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
	"testing"
)

func Test_fromReader(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		offset   int
		expected string
	}{
		{
			name:     "basic",
			str:      "this is my string",
			offset:   5,
			expected: "is my string",
		},
		{
			name: "complex",
			str: `
this is a more complicated string

with multiple lines
`,
			offset: 19,
			expected: `plicated string

with multiple lines
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &offsetReader{reader: strings.NewReader(tt.str), offset: tt.offset}
			bs, _ := ioutil.ReadAll(r)
			assert.Equal(t, tt.expected, string(bs))
		})
	}
}

func Test_countingReader(t *testing.T) {
	tests := []string{
		"this is my string",
		`
this is a more complicated string

with multiple lines
`,
		"and one last string\r\r\n",
	}

	for _, str := range tests {
		r := &countingReader{reader: strings.NewReader(str)}
		_, _ = ioutil.ReadAll(r)
		assert.Equal(t, len(str), r.count)
	}
}

func Test_platformTransformer(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected string
	}{
		{
			name:     "crlf_basic",
			str:      "this is\r\nmy string with\r\n crlf line \rendings\r\n\r\n",
			expected: "this is\nmy string with\n crlf line \rendings\n\n",
		},
		{
			name:     "lf_basic",
			str:      "this is\nmy string with\n lf line \nendings\n \n",
			expected: "this is\r\nmy string with\r\n lf line \r\nendings\r\n \r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := transform.NewReader(strings.NewReader(tt.str), &lineEndingTransformer{})
			got, err := ioutil.ReadAll(r)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, string(got))
		})
	}
}

func Test_platformTransformerFiles(t *testing.T) {
	tests := []struct {
		name         string
		sourcePath   string
		expectedPath string
	}{
		{
			name:         "go1_crlf.txt",
			sourcePath:   "testdata/other_platform/go1_crlf.txt",
			expectedPath: "testdata/other_platform/go1_lf.txt",
		},
		{
			name:         "go1_lf.txt",
			sourcePath:   "testdata/other_platform/go1_lf.txt",
			expectedPath: "testdata/other_platform/go1_crlf.txt",
		},
		{
			name:         "README_lf.md",
			sourcePath:   "testdata/other_platform/README_lf.md",
			expectedPath: "testdata/other_platform/README_crlf.md",
		},
		{
			name:         "README_crlf.md",
			sourcePath:   "testdata/other_platform/README_crlf.md",
			expectedPath: "testdata/other_platform/README_lf.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, err := ioutil.ReadFile(tt.sourcePath)
			assert.Nil(t, err)
			expected, err := ioutil.ReadFile(tt.expectedPath)
			assert.Nil(t, err)
			r := transform.NewReader(bytes.NewReader(source), &lineEndingTransformer{})
			got, err := ioutil.ReadAll(r)
			assert.Nil(t, err)
			assert.Equal(t, expected, got, "contents are not equal")
		})
	}
}
