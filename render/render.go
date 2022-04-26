package render

import (
	"fmt"
	"html"
	"io"
	"reflect"
	"strings"

	"github.com/mlctrez/mystace/context"
	"github.com/mlctrez/mystace/lexer"
	"github.com/mlctrez/mystace/source"
)

type Render interface {
	AddSource(src source.Source) (err error)
	Writer(writer io.Writer)
	Render(name string, context *context.Context) (err error)
}

type render struct {
	writer  io.Writer
	sources map[string]source.Source
}

func New() Render {
	return &render{
		sources: make(map[string]source.Source),
	}
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

func (r *render) Render(name string, context *context.Context) (err error) {

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
		err = r.render(tokens, context)
	}
	return
}

// canRemoveWhitespace determins if a comment exists at position <at> is on a different line
func canRemoveWhitespace(tokens []lexer.Token, current int, at int) bool {
	currentToken := tokens[current]
	if at > -1 && at < len(tokens) {
		maybe := tokens[at]
		mods, _ := maybe.Value()
		if mods.HasModifier(lexer.CommentModifier, lexer.HashModifier) {
			if maybe.Line() != currentToken.Line() {
				return true
			}
			_, value := currentToken.Value()
			if current > at {
				return len(value) > 1 && strings.HasPrefix(value, "\n")
			}
			return strings.TrimSpace(value) == ""
		}
	}
	return false
}

func (r *render) render(tokens []lexer.Token, ctx *context.Context) (err error) {

	totalTokens := len(tokens)
	for i := 0; i < totalTokens; i++ {

		token := tokens[i]
		mods, value := token.Value()

		if mods.HasModifier(lexer.CommentModifier) {
			continue
		}

		if token.IsChar() {

			if canRemoveWhitespace(tokens, i, i-1) && strings.HasPrefix(value, "\n") {
				value = strings.TrimPrefix(value, "\n")
			}
			if canRemoveWhitespace(tokens, i, i+1) && strings.HasSuffix(value, " ") {
				value = strings.TrimRight(value, " ")
			}

			if _, err = r.writer.Write([]byte(value)); err != nil {
				return
			}
			continue
		}

		if token.IsThreeBracket() {
			if v, ok := ctx.Lookup(value); ok {
				if err = r.writeValue(v, false); err != nil {
					return
				}
			}
			continue
		}

		if token.IsTwoBracket() {
			if mods.HasModifier(lexer.ImportModifier) {
				err = fmt.Errorf("implement me")
				return
			}

			if mods.HasModifier(lexer.HashModifier, lexer.InvertedModifier) {
				// find matching close token for *this* token
				openModifiers := 1
				nextToken := -1
				for j := i + 1; j < totalTokens; j++ {
					aheadMods, aheadValue := tokens[j].Value()
					if aheadMods.HasModifier(lexer.HashModifier, lexer.InvertedModifier) {
						openModifiers++
					}
					if aheadMods.HasModifier(lexer.CloseModifier) {
						openModifiers--
						if openModifiers == 0 && aheadValue == value {
							nextToken = j
							break
						}
					}
				}
				if nextToken == -1 {
					return fmt.Errorf("unable to find close for %s", token)
				}
				if value == "if" {
					return fmt.Errorf("if not implemented yet %s", token)
				}

				nestedTokens := tokens[i+1 : nextToken]
				i = nextToken

				if len(nestedTokens) > 0 {
					if nestedTokens[0].Data.Str == "\n" {
						nestedTokens = nestedTokens[1:]
					}
				}

				if len(nestedTokens) > 2 {
					if nestedTokens[len(nestedTokens)-1].Data.Str == "\n" {
						nestedTokens = nestedTokens[:len(nestedTokens)-1]
					}
				}

				if v, ok := ctx.Lookup(value); ok {
					if mods.HasModifier(lexer.HashModifier) {
						switch vv := v.(type) {
						case nil:
						case bool:
							if vv {
								err = r.render(nestedTokens, ctx)
							}
						case map[string]interface{}:
							nc := context.New(vv, ctx)
							err = r.render(nestedTokens, nc)
						case string, float64:
							newValues := map[string]interface{}{".": vv}
							nc := context.New(newValues, ctx)
							err = r.render(nestedTokens, nc)
						case []interface{}:
							for _, nm := range vv {
								if inm, oknm := nm.(map[string]interface{}); oknm {
									nc := context.New(inm, ctx)
									err = r.render(nestedTokens, nc)
									if err != nil {
										break
									}
								}
							}
						default:
							err = fmt.Errorf("hash missing type %s at value %s", reflect.TypeOf(v), value)
						}
					}
					if mods.HasModifier(lexer.InvertedModifier) {
						switch vv := v.(type) {
						case bool:
							if !vv {
								err = r.render(nestedTokens, ctx)
							}
						case nil:
							err = r.render(nestedTokens, ctx)
						case map[string]interface{}:
						case []interface{}:
							if len(vv) == 0 {
								err = r.render(nestedTokens, ctx)
							}
						default:
							err = fmt.Errorf("inverted add check for type %s", reflect.TypeOf(v))
						}
					}
				} else {
					err = fmt.Errorf("missing var for %s", token)
				}

				if err != nil {
					return err
				}
				continue
			}

			escaping := true
			if mods.HasModifier(lexer.AmpModifier) {
				escaping = false
			}

			if v, ok := ctx.Lookup(value); ok {
				if err = r.writeValue(v, escaping); err != nil {
					return
				}
			}
			continue
		}
		return fmt.Errorf("unhandled case")
	}
	return nil
}

func (r *render) writeValue(v interface{}, escape bool) (err error) {
	switch vt := v.(type) {
	case string:
		if escape {
			_, err = r.writer.Write([]byte(htmlEscape(vt)))
		} else {
			_, err = r.writer.Write([]byte(vt))
		}
	case float64:
		_, err = r.writer.Write([]byte(formatFloat(vt)))
	case nil:
	default:
		err = fmt.Errorf("unhandled type %s", reflect.TypeOf(vt))
	}
	return
}

func formatFloat(f float64) string {
	fv := fmt.Sprintf("%1.2f", f)
	return strings.TrimSuffix(fv, ".00")
}

func htmlEscape(s string) (result string) {
	return strings.ReplaceAll(html.EscapeString(s), "&#34;", "&quot;")
}
