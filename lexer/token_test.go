package lexer

import (
	"testing"

	"github.com/mlctrez/mystace/internal/testify"
	"github.com/mlctrez/mystace/source"
)

func TestToken_IsChar(t *testing.T) {
	_, require := testify.New(t)

	require.True(Token{Data: source.Data{Str: "a"}}.IsChar())
	require.True(Token{Data: source.Data{Str: "{a}"}}.IsChar())
	require.False(Token{Data: source.Data{Str: "{{a}}"}}.IsChar())

}
