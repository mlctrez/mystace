package lexer

import (
	"testing"

	"github.com/mlctrez/mystace/internal/testify"
	"github.com/mlctrez/mystace/source"
	log "github.com/sirupsen/logrus"
)

func TestLexer_Parse(t *testing.T) {
	_, require := testify.New(t)

	log.SetLevel(log.DebugLevel)
	s, err := source.FromString("{{aaa}}bbb{{{ccc}}}\n\nddd{{eee}}fff")
	require.Nil(err)

	var parse []Token
	parse, err = New(s).Parse()
	require.Nil(err)
	require.NotNil(parse)
	require.Equal("{{aaa}}", parse[0].data.Str)
	require.Equal("bbb", parse[1].data.Str)
	require.Equal("{{{ccc}}}", parse[2].data.Str)
	require.Equal("fff", parse[5].data.Str)

	s, err = source.FromString("{{")
	require.Nil(err)

	parse, err = New(s).Parse()
	require.ErrorIs(ErrMissingEndToken, err)

	s, err = source.FromString("")
	require.Nil(err)

	parse, err = New(s).Parse()
	require.Nil(err)
	var expected []Token
	require.Equal(expected, parse)

}

func TestToken_String(t *testing.T) {
	_, require := testify.New(t)

	src, err := source.FromString("some data \n with newline")
	require.Nil(err)

	token := Token{data: src.Read(200)}

	require.Equal("Token: Data:\"some data \\n with newline\" {{1 1} {2 13}}", token.String())

}

func TestWithMaxLoops(t *testing.T) {

	_, require := testify.New(t)

	src, err := source.FromString("some data \n with newline")
	require.Nil(err)

	tokens, err := New(src, WithMaxLoops(-1)).Parse()
	require.Nil(tokens)
	require.ErrorIs(ErrMaxLoopsExceeded, err)

}
