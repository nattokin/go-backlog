package backlog

import (
	"io"
	"net/http"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newHTTPClientMock(f roundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(f),
	}
}

func newClientMock(baseURL, token string, f roundTripFunc) (*Client, error) {
	c, err := NewClient(baseURL, token)
	if err != nil {
		return nil, err
	}

	httpClient := newHTTPClientMock(f)
	c.httpClient = httpClient

	return c, nil
}

type mockWrapper struct {
	createErr error
	copyErr   error
	closeErr  error
}

func (w mockWrapper) NewMultipartWriter(_ io.Writer) multipartWriter {
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
