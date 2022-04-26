package testify

import (
	"testing"

	assertN "github.com/stretchr/testify/assert"
	requireN "github.com/stretchr/testify/require"
)

func New(t *testing.T) (assert *assertN.Assertions, require *requireN.Assertions) {
	return Assert(t), Require(t)
}

func Assert(t *testing.T) (assert *assertN.Assertions) {
	return assertN.New(t)
}

func Require(t *testing.T) (require *requireN.Assertions) {
	return requireN.New(t)
}
