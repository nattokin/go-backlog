package core

import (
	"io"
	"mime/multipart"
)

// ──────────────────────────────────────────────────────────────
//  Wrapper interface for I/O abstractions
// ──────────────────────────────────────────────────────────────

type Wrapper interface {
	Copy(dst io.Writer, src io.Reader) error
	NewMultipartWriter(w io.Writer) MultipartWriter
}

type MultipartWriter interface {
	CreateFormFile(fieldname, filename string) (io.Writer, error)
	FormDataContentType() string
	Close() error
}

// ──────────────────────────────────────────────────────────────
//  Default wrapper implementations
// ──────────────────────────────────────────────────────────────

type DefaultWrapper struct{}

func (*DefaultWrapper) Copy(dst io.Writer, src io.Reader) error {
	_, err := io.Copy(dst, src)
	return err
}

func (*DefaultWrapper) NewMultipartWriter(w io.Writer) MultipartWriter {
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
