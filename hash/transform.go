package hash

import (
	"bytes"
	"golang.org/x/text/transform"
	"io"
	"regexp"
)

var (
	crlf = []byte("\r\n")
	lf   = []byte("\n")
)

var noWs = noWhitespaceTransformer{}

// spaceRegex is a regex for removing whitespace - ** based on WS agents api **
var spaceRegex = regexp.MustCompile("[\t\r\n ]")

// noWhitespaceTransformer implements the transform.Transformer interface
// it removes all whitespace from an io.Reader
type noWhitespaceTransformer struct{}

func (noWhitespaceTransformer) Reset() {}
func (noWhitespaceTransformer) Transform(dst, src []byte, _ bool) (nDst, nSrc int, err error) {
	trimmed := spaceRegex.ReplaceAll(src, []byte{})
	copy(dst, trimmed)
	return len(trimmed), len(src), nil
}

// countingReader implements the io.Reader interface
// it counts the amount of bytes in reader
type countingReader struct {
	reader io.Reader
	count  int
}

func (r *countingReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	r.count += n
	return n, err
}

// offsetReader is an io.Reader that only reads bytes starting from offset
type offsetReader struct {
	reader io.Reader

	// offset is the offset in reader to start reading from
	offset int

	// pos is the current position in reader
	pos int
}

func (c *offsetReader) Read(p []byte) (n int, err error) {
	src := make([]byte, len(p))

	n, err = c.reader.Read(src)
	if err != nil {
		return
	}

	idx := c.offset - c.pos
	c.pos += n
	if idx < 0 {
		copy(p, src[:n])
	} else if idx < n {
		copy(p, src[idx:n])
		n -= idx
	} else {
		n = 0
	}
	return
}

// lineEnding represents the line endings type
type lineEnding int

const (
	crlfLE lineEnding = iota
	lfLE
	unknownLE
)

// lineEndingTransformer implements the transform.Transformer interface
// it switches the line ending type for a reader:
// crlf -> lf ("\r\n" -> "\r")
// lf -> crlf ("\r" -> "\r\n")
type lineEndingTransformer struct {
	// target line ending type for the current reader
	target lineEnding

	// contains the previously read bit, for CRLF matching
	prev byte
}

func (t *lineEndingTransformer) Reset() {
	t.target = unknownLE
	t.prev = 0
}

// Transform determines the line ending type in the current buffer
// and sets the target line ending to the opposite type
// if no line endings are found, src is copied to dst as is
func (t *lineEndingTransformer) Transform(dst, src []byte, isEof bool) (nDst, nSrc int, err error) {
	switch {
	case t.target == lfLE || bytes.Contains(src, crlf):
		t.target = lfLE
		return t.transformToLF(dst, src, isEof)
	case t.target == crlfLE || bytes.Contains(src, lf):
		t.target = crlfLE
		return t.transformToCRLF(dst, src, isEof)
	default:
		nDst = copy(dst, src)
		nSrc = nDst
		return
	}
}

// transformToLF converts CRLF line endings to LF
func (t *lineEndingTransformer) transformToLF(dst, src []byte, _ bool) (nDst, nSrc int, err error) {
	for nDst < len(dst) && nSrc < len(src) {
		c := src[nSrc]
		if t.prev == '\r' && c != '\n' {
			if nDst+2 > len(dst) {
				break
			}
			dst[nDst] = '\r'
			dst[nDst+1] = c
			nSrc++
			nDst += 2
			t.prev = 0
		} else if c == '\r' {
			t.prev = '\r'
			nSrc++
		} else {
			dst[nDst] = c
			t.prev = c
			nDst++
			nSrc++
		}
	}
	if nSrc < len(src) {
		err = transform.ErrShortDst
	}
	return
}

// transformToLF converts LF line endings to CRLF
func (lineEndingTransformer) transformToCRLF(dst, src []byte, _ bool) (nDst, nSrc int, err error) {
	for nDst < len(dst) && nSrc < len(src) {
		if c := src[nSrc]; c == '\n' {
			if nDst+1 == len(dst) {
				break
			}
			dst[nDst] = '\r'
			dst[nDst+1] = '\n'
			nSrc++
			nDst += 2
		} else {
			dst[nDst] = c
			nSrc++
			nDst++
		}
	}
	if nSrc < len(src) {
		err = transform.ErrShortDst
	}
	return
}
