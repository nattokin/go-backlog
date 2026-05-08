package core

import (
	"io"
	"mime/multipart"
)

// Wrapper abstracts I/O operations used in multipart uploads to allow test injection.
type Wrapper interface {
	Copy(dst io.Writer, src io.Reader) error
	NewMultipartWriter(w io.Writer) MultipartWriter
}

type MultipartWriter interface {
	CreateFormFile(fieldname, filename string) (io.Writer, error)
	FormDataContentType() string
	Close() error
}

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
