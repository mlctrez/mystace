package render

import (
	"fmt"
	"io"

	"github.com/mlctrez/mystace/lexer"
	"github.com/mlctrez/mystace/source"
)

type Render interface {
	AddSource(src source.Source) (err error)
	Writer(writer io.Writer)
	Render(name string, vars map[string]interface{}) (err error)
}

type render struct {
	writer  io.Writer
	sources map[string]source.Source
}

func (r *render) AddSource(src source.Source) (err error) {
	name := src.Name()
	if name == "" {
		err = source.ErrEmptySourceName
		return
	}
	if _, ok := r.sources[name]; ok {
		err = fmt.Errorf("name %q : %w", name, source.ErrDuplicateSourceName)
		return
	}
	r.sources[name] = src
	return
}

var (
	ErrSourceNameNotFound = fmt.Errorf("source name not found")
	ErrNoWriter           = fmt.Errorf("no writer")
)

func (r *render) Writer(writer io.Writer) {
	r.writer = writer
}

func (r *render) Render(name string, vars map[string]interface{}) (err error) {

	if r.writer == nil {
		err = ErrNoWriter
		return
	}

	if s, ok := r.sources[name]; !ok {
		err = ErrSourceNameNotFound
	} else {
		var tokens []lexer.Token
		if tokens, err = lexer.New(s).Parse(); err != nil {
			return
		}
		err = r.render(tokens, vars)
	}
	return
}

func (r *render) render(tokens []lexer.Token, vars map[string]interface{}) (err error) {
	for _, token := range tokens {
		if token.IsChar() {
			_, err = r.writer.Write([]byte(token.Data.Str))
			if err != nil {
				return
			}
		}
	}
	return nil
}

func New() Render {
	return &render{
		sources: make(map[string]source.Source),
	}
}
