package lexer

import (
	"testing"

	"github.com/mlctrez/mystace/internal/testify"
	"github.com/mlctrez/mystace/source"
)

func makeToken(s string) Token {
	return Token{Data: source.Data{Str: s}}
}

func TestToken_Line(t *testing.T) {
	require := testify.Require(t)
	token := makeToken("a")
	token.Data.Range.Start.Line = 10
	require.Equal(10, token.Line())
}

func TestToken_IsChar(t *testing.T) {
	require := testify.Require(t)

	require.True(makeToken("a").IsChar())
	require.True(makeToken("{a}").IsChar())
	require.False(makeToken("{{a}}").IsChar())

}

func TestToken_IsTwoBracket(t *testing.T) {
	require := testify.Require(t)

	require.True(makeToken("{{a}}").IsTwoBracket())
	require.False(makeToken("{{{a}}}").IsTwoBracket())
}

func TestToken_IsThreeBracket(t *testing.T) {
	require := testify.Require(t)

	require.False(makeToken("{{a}}").IsThreeBracket())
	require.True(makeToken("{{{a}}}").IsThreeBracket())

}

func TestToken_Value(t *testing.T) {
	require := testify.Require(t)

	var expectedMods Modifiers

	mods, v := makeToken("a").Value()
	require.Equal("a", v)
	require.Equal(expectedMods, mods)

	mods, v = makeToken("{a}").Value()
	require.Equal("{a}", v)
	require.Equal(expectedMods, mods)

	mods, v = makeToken("{{a}}").Value()
	require.Equal("a", v)
	require.Equal(expectedMods, mods)

	mods, v = makeToken("{{{a}}}").Value()
	require.Equal("a", v)
	require.Equal(expectedMods, mods)

	mods, v = makeToken("{{&a}}").Value()
	require.Equal("a", v)
	expectedMods = Modifiers{AmpModifier}
	require.Equal(expectedMods, mods)

}

func TestModifiers_HasModifier(t *testing.T) {
	require := testify.Require(t)

	require.True(Modifiers{HashModifier}.HasModifier(HashModifier))

	var empty Modifiers
	require.False(empty.HasModifier(HashModifier))

	require.False(Modifiers{HashModifier}.HasModifier(AmpModifier))

}
