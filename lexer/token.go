package lexer

import (
	"fmt"
	"strings"

	"github.com/mlctrez/mystace/source"
)

type Token struct {
	Data source.Data
}

func (t Token) IsChar() bool {
	return !strings.HasPrefix(t.Data.Str, "{{")
}

func (t Token) String() string {
	return fmt.Sprintf("Token: %s", t.Data)
}
