package context

import (
	"fmt"
	"strings"
)

type Context struct {
	values   map[string]interface{}
	parent   *Context
	resolved map[string]bool
}

func New(values map[string]interface{}, parent ...*Context) (ctx *Context) {

	ctx = &Context{values: values}
	if ctx.values == nil {
		ctx.values = map[string]interface{}{}
	}
	for _, p := range parent {
		ctx.parent = p
		break
	}
	//fmt.Println(ctx)
	return
}

func (c *Context) String() string {
	if c == nil {
		return "Context: <nil>"
	}
	return fmt.Sprintf("Context: %s\n %s", c.values, c.parent)
}

func (c *Context) markResolved(key string) {
	if c.resolved == nil {
		c.clearResolved()
	}
	c.resolved[key] = true
}

func (c *Context) wasResolved(key string) bool {
	if c.resolved == nil {
		return false
	}

	first, _, ok := maybeSplitParts(key)
	if !ok {
		return false
	}

	return c.resolved[first]
}

func (c *Context) clearResolved() {
	c.resolved = map[string]bool{}
}

func (c *Context) Lookup(key string) (i interface{}, ok bool) {

	if i, ok = c.lookup(key, c.values); ok {
		return
	}

	if c.wasResolved(key) || c.parent == nil {
		return nil, false
	}

	return c.parent.Lookup(key)
}

func (c *Context) lookup(keyOuter string, vars map[string]interface{}) (i interface{}, ok bool) {

	key := strings.TrimSpace(keyOuter)

	first, remainder, ok := maybeSplitParts(key)

	if ok {
		if nestedMap, nestedOk := vars[first].(map[string]interface{}); nestedOk {
			c.markResolved(first)
			return c.lookup(remainder, nestedMap)
		}
	}

	i, ok = vars[key]
	return
}

func maybeSplitParts(key string) (first string, remainder string, ok bool) {
	if strings.Index(key, ".") > 0 {
		splitParts := strings.SplitN(key, ".", 2)
		first = splitParts[0]
		remainder = splitParts[1]
		ok = true
	}
	return
}
