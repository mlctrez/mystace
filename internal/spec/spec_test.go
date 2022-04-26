package spec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mlctrez/mystace/internal/testify"
)

func TestRead(t *testing.T) {
	_, require := testify.New(t)

	read, err := Read("../../mustache/specs/interpolation.json")
	require.Nil(err)
	require.NotNil(read)

	require.Equal("No Interpolation", read.Tests[0].Name)

	read, err = Read("../../mustache/specs/not_found.json")
	require.ErrorIs(err, os.ErrNotExist)

}

func TestRead_All(t *testing.T) {
	_, require := testify.New(t)

	glob, err := filepath.Glob("../../mustache/specs/*.json")
	require.Nil(err)
	for _, s := range glob {
		var read File
		read, err = Read(s)
		require.Nil(err)
		require.NotNil(read)
	}

}
