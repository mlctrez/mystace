package testify

import (
	"testing"

	assertN "github.com/stretchr/testify/assert"
	requireN "github.com/stretchr/testify/require"
)

func New(t *testing.T) (assert *assertN.Assertions, require *requireN.Assertions) {
	return assertN.New(t), requireN.New(t)
}
