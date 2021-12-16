package clioptions

import (
	"io"
	"os"
)

// IOStreams provides the standard names for iostreams.
type IOStreams struct {
	// In think, os.Stdin
	In io.Reader
	// Out think, os.Stdout
	Out io.Writer
	// ErrOut think, os.Stderr
	ErrOut io.Writer
}

// StandardIOStreams returns an IOStreams from os.Stdin, os.Stdout
func StandardIOStreams() IOStreams {
	return IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
}
