package mocks

import (
	"testing"

	"github.com/mlctrez/mystace/internal/testify"
)

func TestBadReader_Read(t *testing.T) {

	_, require := testify.New(t)

	makeByteArray := func(len int, v byte) (d []byte) {
		d = make([]byte, len)
		for i := 0; i < len; i++ {
			d[i] = v
		}
		return d
	}

	p := make([]byte, 100)
	read, err := (&BadReader{}).Read(p)
	require.Nil(err)
	require.Equal(10, read)
	require.Equal(makeByteArray(10, 1), p[0:10])

	p = make([]byte, 5)
	read, err = (&BadReader{}).Read(p)
	require.Nil(err)
	require.Equal(5, read)
	require.Equal(makeByteArray(5, 2), p)

	p = make([]byte, 5)
	read, err = (&BadReader{ReadErr: ErrBadReaderMockError}).Read(p)
	require.Equal(0, read)
	require.ErrorIs(err, ErrBadReaderMockError)

}

func TestBadReader_Close(t *testing.T) {
	_, require := testify.New(t)

	err := (&BadReader{CloseErr: ErrBadReaderMockError}).Close()
	require.ErrorIs(err, ErrBadReaderMockError)

}
