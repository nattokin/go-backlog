package backlog

import (
	"io"
	"net/http"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newHTTPClientMock(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(fn),
	}
}

func newClientMock(baseURL, token string, fn roundTripFunc) (*Client, error) {
	c, err := NewClient(baseURL, token)
	if err != nil {
		return nil, err
	}

	httpClient := newHTTPClientMock(fn)
	c.httpClient = httpClient

	return c, nil
}

type mockWrapper struct {
	createErr error
	copyErr   error
	closeErr  error
}

func (m mockWrapper) NewMultipartWriter(_ io.Writer) multipartWriter {
	return &mockMultipartWriter{m: m}
}

func (m mockWrapper) Copy(_ io.Writer, _ io.Reader) error {
	return m.copyErr
}

type mockMultipartWriter struct {
	m mockWrapper
}

func (mw *mockMultipartWriter) CreateFormFile(fieldname, filename string) (io.Writer, error) {
	if mw.m.createErr != nil {
		return nil, mw.m.createErr
	}
	return io.Discard, nil
}
func (mw *mockMultipartWriter) FormDataContentType() string { return "mock/type" }
func (mw *mockMultipartWriter) Close() error                { return mw.m.closeErr }
