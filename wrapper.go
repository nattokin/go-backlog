package backlog

import (
	"io"
	"mime/multipart"
)

// ──────────────────────────────────────────────────────────────
//  Wrapper interface for I/O abstractions
// ──────────────────────────────────────────────────────────────

type wrapper interface {
	Copy(dst io.Writer, src io.Reader) error
	NewMultipartWriter(w io.Writer) multipartWriter
}

type multipartWriter interface {
	CreateFormFile(fieldname, filename string) (io.Writer, error)
	FormDataContentType() string
	Close() error
}

// ──────────────────────────────────────────────────────────────
//  Default wrapper implementations
// ──────────────────────────────────────────────────────────────

type defaultWrapper struct{}

func (*defaultWrapper) Copy(dst io.Writer, src io.Reader) error {
	_, err := io.Copy(dst, src)
	return err
}

func (*defaultWrapper) NewMultipartWriter(w io.Writer) multipartWriter {
	return &defaultMultipartWriter{multipart.NewWriter(w)}
}

type defaultMultipartWriter struct {
	*multipart.Writer
}

func (mw *defaultMultipartWriter) CreateFormFile(fieldname, filename string) (io.Writer, error) {
	return mw.Writer.CreateFormFile(fieldname, filename)
}

func (mw *defaultMultipartWriter) FormDataContentType() string {
	return mw.Writer.FormDataContentType()
}

func (mw *defaultMultipartWriter) Close() error {
	return mw.Writer.Close()
}
