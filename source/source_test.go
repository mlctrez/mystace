package source

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/mlctrez/mystace/internal/mocks"
	"github.com/mlctrez/mystace/internal/testify"
)

func TestFromString(t *testing.T) {
	_, require := testify.New(t)

	src, err := FromString("")
	require.Nil(err)
	require.Equal("", src.(*source).data)

	src, err = FromString("a")
	require.Nil(err)
	require.Equal("a", src.(*source).data)

	src, err = FromString("😄")
	require.Nil(err)
	require.Equal("😄", src.(*source).data)

	src, err = FromString("has newline at end\n")
	require.Nil(err)
	require.Equal("has newline at end\n", src.(*source).data)

}

func TestFromReadCloser(t *testing.T) {
	_, require := testify.New(t)

	src, err := FromReadCloser(nil)
	require.NotNil(err)
	require.ErrorIs(err, ErrNilReadCloser)

	src, err = FromReadCloser(&mocks.BadReader{ReadErr: mocks.ErrBadReaderMockError})
	require.ErrorIs(err, mocks.ErrBadReaderMockError)

	src, err = FromReadCloser(&mocks.BadReader{CloseErr: mocks.ErrBadReaderMockError})
	require.ErrorIs(err, mocks.ErrBadReaderMockError)

	src, err = FromReadCloser(ioutil.NopCloser(&bytes.Buffer{}))
	require.Nil(err)
	require.Equal("", src.(*source).data)

}

func TestWithName(t *testing.T) {
	_, require := testify.New(t)
	src, err := FromString("", WithName("templateName"))
	require.Nil(err)
	require.NotNil(src)
	require.Equal("templateName", src.Name())

	src, err = FromString("", WithName(""))
	require.Nil(src)
	require.ErrorIs(err, ErrEmptySourceName)
}

func TestSource_Peek(t *testing.T) {
	_, require := testify.New(t)
	src, err := FromString("0123456789")
	require.Nil(err)
	require.NotNil(src)

	negPos := src.Peek(-1)
	require.Equal("", negPos.Str)
	require.Equal(Location{0, 0}, negPos.Range.Start)
	require.Equal(Location{0, 0}, negPos.Range.End)

	zeroPos := src.Peek(0)
	require.Equal("", zeroPos.Str)
	require.Equal(Location{0, 0}, zeroPos.Range.Start)
	require.Equal(Location{0, 0}, zeroPos.Range.End)

	require.Equal("0", src.Peek(1).Str)
	require.Equal("01", src.Peek(2).Str)
	require.Equal("0123456789", src.Peek(200).Str)

	src.(*source).position = 1

	require.Equal("1", src.Peek(1).Str)
	require.Equal("12", src.Peek(2).Str)
	require.Equal("123456789", src.Peek(200).Str)

	src.(*source).position = 9
	require.Equal("9", src.Peek(1).Str)

	src.(*source).position = 10
	require.Equal("", src.Peek(1).Str)

}

func TestSource_PeekRange(t *testing.T) {
	assert, require := testify.New(t)
	src, err := FromString("0123456789\nABCDEFGHI")
	require.Nil(err)
	require.NotNil(src)

	p := src.Peek(1)
	assert.Equal(Location{1, 1}, p.Range.Start)
	assert.Equal(Location{1, 1}, p.Range.End)

	p = src.Peek(14)
	assert.Equal(Location{1, 1}, p.Range.Start)
	assert.Equal(Location{2, 3}, p.Range.End)

	src.(*source).position = 5

	p = src.Peek(100)
	assert.Equal(Location{1, 6}, p.Range.Start)
	assert.Equal(Location{2, 9}, p.Range.End)
}

func TestSource_locationAtPosition(t *testing.T) {
	assert, _ := testify.New(t)

	s := &source{data: "012345\n890"}

	assert.Equal(Location{1, 1}, s.locationAtPosition(0))
	assert.Equal(Location{1, 2}, s.locationAtPosition(1))
	assert.Equal(Location{2, 2}, s.locationAtPosition(8))
	assert.Equal(Location{0, 0}, s.locationAtPosition(100))

	s = &source{data: ""}
	assert.Equal(Location{0, 0}, s.locationAtPosition(1))
	s = &source{data: "abcdefg"}
	assert.Equal(Location{0, 0}, s.locationAtPosition(-1))

}

func TestSource_Read(t *testing.T) {
	_, require := testify.New(t)

	s := &source{data: "012345\n890"}
	require.NotNil(s)

	require.Equal("01", s.Read(2).Str)
	require.Equal(2, s.position)
	require.Equal("23", s.Read(2).Str)
	require.Equal(4, s.position)
	require.Equal("45\n890", s.Read(20).Str)
	require.Equal(10, s.position)

	eofRead := s.Read(1)
	require.Equal(10, s.position)
	require.Equal("", eofRead.Str)
	require.Equal(Location{0, 0}, eofRead.Range.Start)

}

func TestData_String(t *testing.T) {
	_, require := testify.New(t)

	s := &source{data: "some data \n with newline"}
	d := s.Read(100)

	require.Equal("Data:\"some data \\n with newline\" {{1 1} {2 13}}", d.String())

}
