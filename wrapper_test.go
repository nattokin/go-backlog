package backlog

import (
	"io"

	"github.com/nattokin/go-backlog/internal/core"
)

type mockWrapper struct {
	createErr error
	copyErr   error
	closeErr  error
}

func (w mockWrapper) NewMultipartWriter(_ io.Writer) core.MultipartWriter {
	return &mockMultipartWriter{wrapper: w}
}

func (w mockWrapper) Copy(_ io.Writer, _ io.Reader) error {
	return w.copyErr
}

type mockMultipartWriter struct {
	wrapper mockWrapper
}

func (mw *mockMultipartWriter) CreateFormFile(fieldname, filename string) (io.Writer, error) {
	if mw.wrapper.createErr != nil {
		return nil, mw.wrapper.createErr
	}
	return io.Discard, nil
}
func (mw *mockMultipartWriter) FormDataContentType() string { return "mock/type" }
func (mw *mockMultipartWriter) Close() error                { return mw.wrapper.closeErr }
