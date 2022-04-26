package context

import (
	"testing"

	"github.com/mlctrez/mystace/internal/testify"
)

func TestNew(t *testing.T) {
	require := testify.Require(t)
	require.True(true)

	ctx := New(nil)
	require.NotNil(ctx)
	require.Nil(ctx.parent)

	values := map[string]interface{}{"foo": "bar"}
	ctx = New(values)
	require.NotNil(ctx)
	require.Nil(ctx.parent)
	require.Equal(values, ctx.values)

	parentOne := New(nil)
	parentTwo := New(nil)

	ctx = New(nil, parentOne, parentTwo)
	require.NotNil(ctx)
	require.Equal(parentOne, ctx.parent)

}

func TestContext_Lookup(t *testing.T) {
	require := testify.Require(t)
	require.True(true)

	c := New(nil)
	lookup, ok := c.Lookup("foo")
	require.False(ok)
	require.Nil(lookup)

	expected := "bar"
	c = New(map[string]interface{}{"foo": expected})
	lookup, ok = c.Lookup("foo")
	require.True(ok)
	require.Equal(expected, lookup)

	expected = "nested value"
	c = New(map[string]interface{}{"nested": map[string]interface{}{"value": expected}})
	lookup, ok = c.Lookup("nested.value")
	require.True(ok)
	require.Equal(expected, lookup)

	expected = "from child"
	c = New(
		map[string]interface{}{"foo": expected},
		New(map[string]interface{}{"foo": "from parent"}),
	)

	lookup, ok = c.Lookup("foo")
	require.True(ok)
	require.Equal(expected, lookup)

	expected = "from parent"
	c = New(
		map[string]interface{}{"foo_does_not_exist": "from child"},
		New(map[string]interface{}{"foo": expected}),
	)
	lookup, ok = c.Lookup("foo")
	require.True(ok)
	require.Equal(expected, lookup)

}

func TestContext_Lookup_NoShadow(t *testing.T) {
	require := testify.Require(t)
	require.True(true)

	// data=map[a:map[b:map[]] b:map[c:ERROR]]

	values := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{},
		},
		"b": map[string]interface{}{
			"c": "ERROR",
		},
	}
	c := New(values)

	var oi interface{}
	var ook bool

	if i, ok := c.Lookup("a"); ok {
		switch it := i.(type) {
		case map[string]interface{}:
			nc := &Context{values: it, parent: c}
			oi, ook = nc.Lookup("b.c")
		}
	}
	require.False(ook)
	require.Equal(nil, oi)

}

func Test_lookup(t *testing.T) {
	require := testify.Require(t)
	require.True(true)

	c := New(nil)

	expectString := "string at root"
	lookup, ok := c.lookup(".", map[string]interface{}{".": expectString})
	require.True(ok)
	require.Equal(expectString, lookup)

	expectString = "value for key"
	lookup, ok = c.lookup("key", map[string]interface{}{"key": expectString})
	require.True(ok)
	require.Equal(expectString, lookup)

	expectString = "value for levelTwo"
	lookup, ok = c.lookup("levelOne.levelTwo", map[string]interface{}{"levelOne": map[string]interface{}{
		"levelTwo": expectString,
	}})
	require.True(ok)
	require.Equal(expectString, lookup)

	lookup, ok = c.lookup("levelOne.doesnotexist", map[string]interface{}{"levelOne": map[string]interface{}{
		"levelTwo": expectString,
	}})
	require.False(ok)
	require.Nil(lookup)

}

func TestContext_String(t *testing.T) {
	require := testify.Require(t)
	require.True(true)

	require.Equal("Context: <nil>", (*Context)(nil).String())
	require.Equal("Context: map[]\n Context: <nil>", New(nil).String())
}

func TestContext_wasResolved(t *testing.T) {
	require := testify.Require(t)
	require.True(true)

	ctx := New(map[string]interface{}{
		"a": map[string]interface{}{"b": "c"},
		"d": "e",
	},
		New(map[string]interface{}{
			"a": map[string]interface{}{"b": "ERRROR"},
			"d": "ERROR",
		}),
	)

	lookup, ok := ctx.Lookup("d")
	require.True(ok)
	require.Equal("e", lookup)
	require.False(ctx.wasResolved("d"))

	lookup, ok = ctx.Lookup("a.b")
	require.True(ok)
	require.Equal("c", lookup)
	require.True(ctx.wasResolved("a.b"))
	require.False(ctx.wasResolved("a"))

}
