package lexer

import (
	"fmt"
	"strings"

	"github.com/mlctrez/mystace/source"
)

var (
	ErrMissingEndToken  = fmt.Errorf("missing end token }}")
	ErrMaxLoopsExceeded = fmt.Errorf("max loops exceeded")
)

const (
	DefaultMaxLoops = 500
)

type Lexer interface {
	Parse() (tokens []Token, err error)
}

type lexer struct {
	source   source.Source
	maxLoops int
}

func New(source source.Source, options ...Option) Lexer {
	l := &lexer{source: source}
	for _, option := range options {
		option(l)
	}
	return l
}

func (l *lexer) Parse() (tokens []Token, err error) {

	if l.maxLoops == 0 {
		l.maxLoops = DefaultMaxLoops
	}

	loops := 0
	for {
		loops++
		if loops > l.maxLoops {
			err = ErrMaxLoopsExceeded
			return
		}
		var peek source.Data
		if peek = l.source.Peek(100); peek.Str == "" {
			break
		}

		if start := strings.Index(peek.Str, "{{"); start > 0 {
			tokens = append(tokens, Token{Data: l.source.Read(start)})
			continue
		}
		if start := strings.Index(peek.Str, "{{"); start == 0 {
			end := strings.Index(peek.Str, "}}")
			if end == -1 {
				err = ErrMissingEndToken
				return
			}
			if len(peek.Str) > end+2 && peek.Str[end:end+3] == "}}}" {
				end++
			}
			tokens = append(tokens, Token{Data: l.source.Read(end + 2)})
			continue
		}
		if start := strings.Index(peek.Str, "{{"); start < 0 {
			tokens = append(tokens, Token{Data: l.source.Read(len(peek.Str))})
		}

	}

	return
}

type Option func(s *lexer) error

func WithMaxLoops(maxLoops int) Option {
	return func(s *lexer) error {
		s.maxLoops = maxLoops
		return nil
	}
}
