package render

import (
	"bytes"
	"testing"

	"github.com/mlctrez/mystace/internal/mocks"
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
	err := r.Render("foo", map[string]interface{}{})
	require.ErrorIs(err, ErrNoWriter)

	buf := &bytes.Buffer{}
	r.Writer(buf)

	err = r.Render("foo", map[string]interface{}{})
	require.ErrorIs(err, ErrSourceNameNotFound)

	var src source.Source
	src, err = source.FromString("content", source.WithName("foo"))
	require.Nil(err)

	err = r.AddSource(src)
	require.Nil(err)

	err = r.Render("foo", map[string]interface{}{})
	require.Nil(err)

	// testing bad token source
	src, err = source.FromString("{{", source.WithName("bad"))
	require.Nil(err)

	err = r.AddSource(src)
	require.Nil(err)

	err = r.Render("bad", map[string]interface{}{})
	require.NotNil(err)

}

func Test_render(t *testing.T) {
	_, require := testify.New(t)

	buf := &bytes.Buffer{}
	r := render{writer: buf}
	tokens := []lexer.Token{{Data: source.Data{Str: "simple"}}}

	err := r.render(tokens, map[string]interface{}{})
	require.Nil(err)
	require.Equal("simple", buf.String())

}

func Test_render_writerErr(t *testing.T) {
	_, require := testify.New(t)

	buf := &mocks.BadWriter{WriteErr: mocks.ErrBadWriterMockError}
	r := render{writer: buf}
	tokens := []lexer.Token{{Data: source.Data{Str: "simple"}}}

	err := r.render(tokens, map[string]interface{}{})
	require.ErrorIs(err, mocks.ErrBadWriterMockError)

}
