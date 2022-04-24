package source

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type Source interface {
	Name() string
	// Peek returns up to len length Data without advancing the position
	Peek(size int) Data
	// Read returns up to len length Data and advances the current position
	Read(size int) Data
}

type source struct {
	data string
	name string
	// position is the 0 index current location within data for the next read
	position int
}

func (s *source) locationAtPosition(pos int) (l Location) {

	l = Location{1, 0}

	var newLine bool
	var positionAccurate bool
	for i, d := range s.data {
		switch d {
		case '\n', '\r':
			newLine = true
		default:
			if newLine {
				newLine = false
				l.Column = 1
				l.Line++
			} else {
				l.Column++
			}
		}
		if i == pos {
			positionAccurate = true
			break
		}
	}
	if !positionAccurate {
		l = Location{0, 0}
	}
	return
}

func (s *source) Peek(size int) (d Data) {
	d, _ = s.peek(size)
	return
}

func (s *source) peek(size int) (d Data, newPosition int) {
	d = Data{}
	if size < 1 {
		d.Range.Start = Location{0, 0}
		d.Range.End = Location{0, 0}
		return
	}

	peekFrom := s.data[s.position:]

	if size > len(peekFrom) {
		size = len(peekFrom)
	}
	d.Str = peekFrom[0:size]
	d.Range.Start = s.locationAtPosition(s.position)
	d.Range.End = s.locationAtPosition(s.position + (size - 1))
	newPosition = s.position + size

	return
}

func (s *source) Read(size int) (d Data) {
	d, s.position = s.peek(size)
	return
}

func (s *source) Name() string {
	return s.name
}

func FromReadCloser(r io.ReadCloser, options ...Option) (s Source, err error) {

	if r == nil {
		return nil, ErrNilReadCloser
	}

	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return
	}

	if err = r.Close(); err != nil {
		return
	}

	sp := &source{data: strings.Join(lines, "\n")}
	for _, option := range options {
		err = option(sp)
		if err != nil {
			return
		}
	}
	s = sp
	return
}

func FromString(data string, options ...Option) (s Source, err error) {
	return FromReadCloser(ioutil.NopCloser(bytes.NewBufferString(data)), options...)
}

type Option func(s *source) error

func WithName(name string) Option {
	return func(s *source) error {
		if name == "" {
			return ErrEmptySourceName
		}
		s.name = name
		return nil
	}
}

// Data is a fragment of the source with the range
type Data struct {
	Str   string
	Range Range
}

func (d Data) String() string {
	return fmt.Sprintf("Data:%q %v", d.Str, d.Range)
}

// Range designates a range within the source
type Range struct {
	Start Location
	End   Location
}

// Location designates a location within the source
type Location struct {
	Line   int
	Column int
}

var (
	ErrNilReadCloser       = fmt.Errorf("nil readCloser")
	ErrEmptySourceName     = fmt.Errorf("empty name")
	ErrDuplicateSourceName = fmt.Errorf("empty name")
)
