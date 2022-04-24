package testify

import "testing"

func TestNew(t *testing.T) {
	assert, require := New(t)
	require.NotNil(assert)
	require.NotNil(require)
}
