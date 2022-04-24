package mocks

import (
	"fmt"
	"io"
)

var _ io.ReadCloser = (*BadReader)(nil)

var (
	ErrBadReaderMockError = fmt.Errorf("a test error for mocking")
	ErrBadWriterMockError = fmt.Errorf("a mock write error")
)

type BadReader struct {
	ReadErr  error
	CloseErr error
}

func (b *BadReader) Read(p []byte) (n int, err error) {
	if b.ReadErr != nil {
		return 0, b.ReadErr
	}

	b.ReadErr = io.EOF

	if len(p) > 10 {
		for i := 0; i < 10; i++ {
			p[i] = 1
		}
		n = 10
	} else {
		for i := 0; i < len(p); i++ {
			p[i] = 2
		}
		n = len(p)
	}

	return
}

func (b *BadReader) Close() error {
	return b.CloseErr
}

var _ io.Writer = (*BadWriter)(nil)

type BadWriter struct {
	WriteErr error
}

func (b *BadWriter) Write(p []byte) (n int, err error) {
	if b.WriteErr != nil {
		err = b.WriteErr
	}
	return
}
