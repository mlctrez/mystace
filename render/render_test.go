package render

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/mlctrez/mystace/context"
	"github.com/mlctrez/mystace/internal/mocks"
	"github.com/mlctrez/mystace/internal/spec"
	"github.com/mlctrez/mystace/internal/testify"
	"github.com/mlctrez/mystace/lexer"
	"github.com/mlctrez/mystace/source"
)

func TestNew(t *testing.T) {

	_, require := testify.New(t)

	r := New()
	require.NotNil(r)
	require.NotNil(r.(*render).sources)

}

func TestRender_AddSource(t *testing.T) {

	_, require := testify.New(t)

	src, err := source.FromString("one")
	require.Nil(err)
	require.NotNil(src)

	r := New()
	err = r.AddSource(src)
	require.ErrorIs(source.ErrEmptySourceName, err)

	src, err = source.FromString("one", source.WithName("one"))
	require.Nil(err)
	require.NotNil(src)

	require.Nil(r.AddSource(src))
	require.ErrorIs(r.AddSource(src), source.ErrDuplicateSourceName)

}

func TestRender_Writer(t *testing.T) {
	_, require := testify.New(t)
	r := New()
	writer := &bytes.Buffer{}
	r.Writer(writer)

	require.Equal(writer, r.(*render).writer)

}

func TestRender_Render(t *testing.T) {
	_, require := testify.New(t)
	r := New()
	err := r.Render("foo", context.New(nil))
	require.ErrorIs(err, ErrNoWriter)

	buf := &bytes.Buffer{}
	r.Writer(buf)

	err = r.Render("foo", context.New(nil))
	require.ErrorIs(err, ErrSourceNameNotFound)

	var src source.Source
	src, err = source.FromString("content", source.WithName("foo"))
	require.Nil(err)

	err = r.AddSource(src)
	require.Nil(err)

	err = r.Render("foo", context.New(nil))
	require.Nil(err)

	// testing bad token source
	src, err = source.FromString("{{", source.WithName("bad"))
	require.Nil(err)

	err = r.AddSource(src)
	require.Nil(err)

	err = r.Render("bad", context.New(nil))
	require.NotNil(err)

}

func Test_render(t *testing.T) {
	_, require := testify.New(t)

	buf := &bytes.Buffer{}
	r := render{writer: buf}
	tokens := []lexer.Token{{Data: source.Data{Str: "simple"}}}

	err := r.render(tokens, context.New(nil))
	require.Nil(err)
	require.Equal("simple", buf.String())

}

func Test_render_writerErr(t *testing.T) {
	_, require := testify.New(t)

	buf := &mocks.BadWriter{WriteErr: mocks.ErrBadWriterMockError}
	r := render{writer: buf}
	tokens := []lexer.Token{{Data: source.Data{Str: "simple"}}}

	err := r.render(tokens, context.New(nil))
	require.ErrorIs(err, mocks.ErrBadWriterMockError)

}

func TestRender_MustacheSpecs(t *testing.T) {
	_, require := testify.New(t)

	specFiles := []string{
		"../mustache/specs/comments.json",
		//"../mustache/specs/delimiters.json",
		"../mustache/specs/interpolation.json",
		//"../mustache/specs/inverted.json",
		//"../mustache/specs/partials.json",
		"../mustache/specs/sections.json",
	}

	for _, file := range specFiles {

		if !strings.HasSuffix(file, ".json") {
			continue
		}

		spec, err := spec.Read(file)
		require.Nil(err)

		for i, test := range spec.Tests {

			if test.Name != "Deeply Nested Contexts" {
				//continue
			}

			// since we normalize all line endings to \n, fix up spec tests \r\n -> \n
			test.Expected = strings.ReplaceAll(test.Expected, "\r", "")

			var src source.Source
			name := fmt.Sprintf("interpolation%d", i)
			src, err = source.FromString(test.Template, source.WithName(name))
			require.Nil(err)

			r := New()
			buf := &bytes.Buffer{}
			r.Writer(buf)
			err = r.AddSource(src)
			require.Nil(err)

			vars := make(map[string]interface{})
			if tryVars, ok := test.Data.(map[string]interface{}); ok {
				vars = tryVars
			} else {
				vars["."] = test.Data
			}

			err = r.Render(name, context.New(vars))
			message := fmt.Sprintf("specFile=%s testName=%s template=%q data=%v", file, test.Name, test.Template, test.Data)

			require.Nil(err, message)
			require.Equal(test.Expected, buf.String(), message)

		}
	}
}

func TestRender_modsAt(t *testing.T) {
	require := testify.Require(t)
	require.True(true)

	src, err := source.FromString(`some data{{!comment}}other data
more data
{{!comment}}
more more data
`)
	require.Nil(err)
	var tokens []lexer.Token
	tokens, err = lexer.New(src).Parse()
	require.Nil(err)

	// out of range should not blow up
	require.False(canRemoveWhitespace(tokens, 1, -1))
	require.False(canRemoveWhitespace(tokens, 1, 10))

	// comments on same line should not be treated as whitespace removal
	require.False(canRemoveWhitespace(tokens, 0, 1))
	require.False(canRemoveWhitespace(tokens, 2, 1))

	// comment on different line
	require.True(canRemoveWhitespace(tokens, 2, 3))
	require.True(canRemoveWhitespace(tokens, 4, 3))

	// comparison with non comments
	require.False(canRemoveWhitespace(tokens, 0, 0))

}

/*

{{#a}}
	{{one}}
	{{#b}}
		{{one}}{{two}}{{one}}
		{{#c}}
			{{one}}{{two}}{{three}}{{two}}{{one}}
			{{#d}}
				{{one}}{{two}}{{three}}{{four}}{{three}}{{two}}{{one}}
				{{#five}}
					{{one}}{{two}}{{three}}{{four}}{{five}}{{four}}{{three}}{{two}}{{one}}
					{{one}}{{two}}{{three}}{{four}}{{.}}6{{.}}{{four}}{{three}}{{two}}{{one}}
					{{one}}{{two}}{{three}}{{four}}{{five}}{{four}}{{three}}{{two}}{{one}}
				{{/five}}
				{{one}}{{two}}{{three}}{{four}}{{three}}{{two}}{{one}}
			{{/d}}
			{{one}}{{two}}{{three}}{{two}}{{one}}
		{{/c}}
		{{one}}{{two}}{{one}}
	{{/b}}
	{{one}}
{{/a}}


*/
