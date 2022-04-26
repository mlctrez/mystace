package lexer

import (
	"fmt"
	"strings"

	"github.com/mlctrez/mystace/source"
)

type Token struct {
	Data source.Data
}

func (t Token) Line() int {
	return t.Data.Range.Start.Line
}

func (t Token) IsChar() bool {
	return !strings.HasPrefix(t.Data.Str, "{{")
}

func (t Token) IsTwoBracket() bool {
	return strings.HasPrefix(t.Data.Str, "{{") && !t.IsThreeBracket()
}

func (t Token) IsThreeBracket() bool {
	return strings.HasPrefix(t.Data.Str, "{{{")
}

func (t Token) String() string {
	return fmt.Sprintf("Token: %s", t.Data)
}

// Value returns the character data or the data inside brackets
func (t Token) Value() (mods Modifiers, value string) {
	if t.IsChar() {
		value = t.Data.Str
		return
	}
	prefix := "{{"
	suffix := "}}"

	if t.IsThreeBracket() {
		prefix = "{{{"
		suffix = "}}}"
	}
	value = strings.TrimSuffix(strings.TrimPrefix(t.Data.Str, prefix), suffix)
	for _, modifier := range AllModifiers {
		modStr := string(modifier)
		if strings.HasPrefix(value, modStr) {
			mods = append(mods, modifier)
			value = strings.TrimPrefix(value, modStr)
		}
	}
	return
}

type Modifier string

const (
	HashModifier     Modifier = "#"
	AmpModifier      Modifier = "&"
	ImportModifier   Modifier = ">"
	TildeModifier    Modifier = "~"
	CloseModifier    Modifier = "/"
	CommentModifier  Modifier = "!"
	InvertedModifier Modifier = "^"
)

type Modifiers []Modifier

func (mods Modifiers) HasModifier(list ...Modifier) bool {
	for _, m := range mods {
		for _, mod := range list {
			if m == mod {
				return true
			}
		}
	}
	return false
}

var (
	AllModifiers = []Modifier{HashModifier, AmpModifier, ImportModifier, TildeModifier, CloseModifier, CommentModifier, InvertedModifier}
)
