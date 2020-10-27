package backlog_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	backlog "github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestNewClientError(t *testing.T) {
	msg := "test error massage"
	e := backlog.ExportNewClientError(msg)

	assert.Equal(t, msg, e.Error())
}

func TestNewClient(t *testing.T) {
	cases := map[string]struct {
		url       string
		token     string
		wantError bool
	}{
		"no-error": {
			url:       "https://test.backlog.com",
			token:     "test",
			wantError: false,
		},
		"url-token-empty": {
			url:       "",
			token:     "",
			wantError: true,
		},
		"url-empty": {
			url:       "",
			token:     "test",
			wantError: true,
		},
		"token-empty": {
			url:       "https://test.backlog.com",
			token:     "",
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			c, err := backlog.NewClient(tc.url, tc.token)

			switch {
			case tc.wantError:
				assert.Error(t, err)
				assert.Nil(t, c)
			case !tc.wantError:
				assert.Nil(t, err)
				assert.NotNil(t, c)
			}

			if c == nil {
				return
			}

			assert.Equal(t, tc.url, c.ExportURL().String())
			assert.Equal(t, tc.token, c.ExportToken())
		})
	}
}

func TestClient_NewReqest(t *testing.T) {
	reader := bytes.NewReader([]byte("test"))
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	cases := map[string]struct {
		method    string
		spath     string
		params    *backlog.ExportRequestParams
		body      io.Reader
		wantError bool
	}{
		"method-get": {
			method:    http.MethodGet,
			spath:     "get",
			params:    backlog.ExportNewRequestParams(),
			body:      reader,
			wantError: false,
		},
		"method-post": {
			method:    http.MethodPost,
			spath:     "post",
			params:    backlog.ExportNewRequestParams(),
			body:      reader,
			wantError: false,
		},
		"method-patch": {
			method:    http.MethodPatch,
			spath:     "patch",
			params:    backlog.ExportNewRequestParams(),
			body:      reader,
			wantError: false,
		},
		"method-delete": {
			method:    http.MethodDelete,
			spath:     "delete",
			params:    backlog.ExportNewRequestParams(),
			body:      reader,
			wantError: false,
		},
		"method-empty": {
			method:    "",
			spath:     "nothing",
			params:    backlog.ExportNewRequestParams(),
			body:      reader,
			wantError: false,
		},
		"method-eroor": {
			method:    "@error",
			spath:     "nothing",
			params:    backlog.ExportNewRequestParams(),
			body:      reader,
			wantError: true,
		},
		"spath-empty": {
			method:    http.MethodGet,
			spath:     "",
			params:    backlog.ExportNewRequestParams(),
			body:      reader,
			wantError: true,
		},
		"params-empty": {
			method:    http.MethodGet,
			spath:     "test",
			params:    nil,
			body:      reader,
			wantError: false,
		},
		"body-empty": {
			method:    http.MethodGet,
			spath:     "test",
			params:    backlog.ExportNewRequestParams(),
			body:      nil,
			wantError: false,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			request, err := backlog.ExportClientNewReqest(c, tc.method, tc.spath, tc.params, tc.body)

			switch {
			case tc.wantError:
				assert.Error(t, err)
				assert.Nil(t, request)
			case !tc.wantError:
				assert.Nil(t, err)
				assert.NotNil(t, request)
			}
		})
	}

}

func TestClient_Do(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")

	user := &backlog.User{
		ID:          1,
		UserID:      "admin",
		Name:        "admin",
		RoleType:    1,
		Lang:        "ja",
		MailAddress: "test@example",
	}
	now := time.Now()

	want := backlog.Wiki{
		ID:          1,
		ProjectID:   1,
		Name:        "Home",
		Content:     "test",
		Tags:        []*backlog.Tag{},
		Attachments: []*backlog.Attachment{},
		SharedFiles: []*backlog.SharedFile{},
		Stars:       []*backlog.Star{},
		CreatedUser: user,
		Created:     now,
		UpdatedUser: user,
		Updated:     now,
	}

	bs, _ := json.Marshal(want)
	body := ioutil.NopCloser(bytes.NewReader(bs))

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       body,
		}

		return resp, nil
	})
	c.ExportSetHTTPClient(httpClient)

	req, _ := backlog.ExportClientNewReqest(
		c, http.MethodGet, "test",
		backlog.ExportNewRequestParams(),
		bytes.NewReader([]byte("test")),
	)

	res, err := backlog.ExportClientDo(c, req)
	assert.Nil(t, err)

	defer res.Body.Close()

	wiki := backlog.Wiki{}
	json.NewDecoder(res.Body).Decode(&wiki)

	assert.Equal(t, want.ID, wiki.ID)
	assert.Equal(t, want.Name, wiki.Name)
	assert.Equal(t, want.CreatedUser.Name, wiki.CreatedUser.Name)
}

func TestClient_Do_httpClientError(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	emsg := "http client error"
	// If http.Client.Do return error, Shuld return error.
	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New(emsg)
	})
	c.ExportSetHTTPClient(httpClient)

	req, _ := backlog.ExportClientNewReqest(
		c, http.MethodGet, "test",
		backlog.ExportNewRequestParams(),
		bytes.NewReader([]byte("test")),
	)
	_, err := backlog.ExportClientDo(c, req)
	assert.Error(t, err, emsg)

}

func TestClient_Do_errorResponse(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")

	apiErrors := backlog.APIResponseError{
		Errors: []*backlog.Error{
			{
				Message:  "No project.",
				Code:     6,
				MoreInfo: "",
			},
		},
	}

	bs, _ := json.Marshal(apiErrors)
	body := ioutil.NopCloser(bytes.NewReader(bs))

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: 404,
			Header:     header,
			Body:       body,
		}

		return resp, nil
	})
	c.ExportSetHTTPClient(httpClient)

	req, _ := backlog.ExportClientNewReqest(
		c, http.MethodGet, "test",
		backlog.ExportNewRequestParams(),
		bytes.NewReader([]byte("test")),
	)

	_, err := backlog.ExportClientDo(c, req)
	assert.Error(t, err, apiErrors.Error())
}

func TestClient_Get(t *testing.T) {
	baseURL := "https://test.backlog.com"
	apiKey := "apikey"
	spath := "spath"
	want := struct {
		method      string
		url         string
		contentType string
	}{
		method:      http.MethodGet,
		url:         baseURL + "/api/v2/" + spath + "?apiKey=" + apiKey,
		contentType: "",
	}

	c, _ := backlog.NewClient(baseURL, apiKey)

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())
		assert.Equal(t, want.contentType, req.Header.Get("Content-Type"))
		return &http.Response{StatusCode: http.StatusOK}, nil
	})
	c.ExportSetHTTPClient(httpClient)

	res, _ := backlog.ExportClientGet(c, spath, nil)
	statusCode := res.ExportGetHTTPResponse().StatusCode
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestClient_Get_newRequestError(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	_, err := backlog.ExportClientGet(c, "", backlog.ExportNewRequestParams())
	assert.Error(t, err)
}

func TestClient_Post(t *testing.T) {
	baseURL := "https://test.backlog.com"
	apiKey := "apikey"
	spath := "spath"
	want := struct {
		method      string
		url         string
		contentType string
	}{
		method:      http.MethodPost,
		url:         baseURL + "/api/v2/" + spath + "?apiKey=" + apiKey,
		contentType: "application/x-www-form-urlencoded",
	}

	c, _ := backlog.NewClient(baseURL, apiKey)

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())
		assert.Equal(t, want.contentType, req.Header.Get("Content-Type"))
		return &http.Response{StatusCode: http.StatusOK}, nil
	})
	c.ExportSetHTTPClient(httpClient)

	res, _ := backlog.ExportClientPost(c, spath, nil)
	assert.Equal(t, http.StatusOK, res.ExportGetHTTPResponse().StatusCode)
}

func TestClient_Post_newRequestError(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	_, err := backlog.ExportClientPost(c, "", backlog.ExportNewRequestParams())
	assert.Error(t, err)
}

func TestClient_Patch(t *testing.T) {
	baseURL := "https://test.backlog.com"
	apiKey := "apikey"
	spath := "spath"
	want := struct {
		method      string
		url         string
		contentType string
	}{
		method:      http.MethodPatch,
		url:         baseURL + "/api/v2/" + spath + "?apiKey=" + apiKey,
		contentType: "application/x-www-form-urlencoded",
	}

	c, _ := backlog.NewClient(baseURL, apiKey)

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())
		assert.Equal(t, want.contentType, req.Header.Get("Content-Type"))
		content := []byte{}
		req.Body.Read(content)
		assert.Empty(t, content)
		return &http.Response{StatusCode: http.StatusOK}, nil
	})
	c.ExportSetHTTPClient(httpClient)

	params := backlog.ExportNewRequestParams()
	params.Set("key", "value")

	res, _ := backlog.ExportClientPatch(c, spath, params)
	assert.Equal(t, http.StatusOK, res.ExportGetHTTPResponse().StatusCode)
}

func TestClient_Patch_emptyParams(t *testing.T) {
	baseURL := "https://test.backlog.com"
	apiKey := "apikey"
	spath := "spath"

	c, _ := backlog.NewClient(baseURL, apiKey)

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		return &http.Response{StatusCode: http.StatusOK}, nil
	})
	c.ExportSetHTTPClient(httpClient)

	res, _ := backlog.ExportClientPatch(c, spath, nil)
	assert.Equal(t, http.StatusOK, res.ExportGetHTTPResponse().StatusCode)
}

func TestClient_Patch_newRequestError(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	_, err := backlog.ExportClientPatch(c, "", backlog.ExportNewRequestParams())
	assert.Error(t, err)
}

func TestClient_Delete(t *testing.T) {
	baseURL := "https://test.backlog.com"
	apiKey := "apikey"
	spath := "spath"
	want := struct {
		method      string
		url         string
		contentType string
	}{
		method:      http.MethodDelete,
		url:         baseURL + "/api/v2/" + spath + "?apiKey=" + apiKey,
		contentType: "application/x-www-form-urlencoded",
	}

	c, _ := backlog.NewClient(baseURL, apiKey)

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()

		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())
		assert.Equal(t, want.contentType, req.Header.Get("Content-Type"))
		content := []byte{}
		req.Body.Read(content)
		assert.Empty(t, content)
		return &http.Response{StatusCode: http.StatusOK}, nil
	})
	c.ExportSetHTTPClient(httpClient)

	params := backlog.ExportNewRequestParams()
	params.Set("key", "value")

	res, _ := backlog.ExportClientDelete(c, spath, params)
	assert.Equal(t, http.StatusOK, res.ExportGetHTTPResponse().StatusCode)
}

func TestClient_Delete_emptyParams(t *testing.T) {
	baseURL := "https://test.backlog.com"
	apiKey := "apikey"
	spath := "spath"

	c, _ := backlog.NewClient(baseURL, apiKey)

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		return &http.Response{StatusCode: http.StatusOK}, nil
	})
	c.ExportSetHTTPClient(httpClient)

	res, _ := backlog.ExportClientDelete(c, spath, nil)
	assert.Equal(t, http.StatusOK, res.ExportGetHTTPResponse().StatusCode)
}

func TestClient_Delete_newRequestError(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	_, err := backlog.ExportClientDelete(c, "", backlog.ExportNewRequestParams())
	assert.Error(t, err)
}

func TestClient_Uploade(t *testing.T) {
	baseURL := "https://test.backlog.com"
	apiKey := "apikey"
	spath := "spath"
	fpath := "testdata/testfile"
	want := struct {
		method      string
		url         string
		contentType string
	}{
		method:      http.MethodPost,
		url:         baseURL + "/api/v2/" + spath + "?apiKey=" + apiKey,
		contentType: "multipart/form-data",
	}

	c, _ := backlog.NewClient(baseURL, apiKey)

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())

		contentType := req.Header.Get("Content-Type")
		if strings.Index(contentType, want.contentType) != 0 {
			t.Errorf("want: %s, got: %s", want.contentType, contentType)
		}
		return &http.Response{StatusCode: http.StatusOK}, nil
	})
	c.ExportSetHTTPClient(httpClient)

	res, err := backlog.ExportClientUploade(c, spath, fpath, "fname")
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.ExportGetHTTPResponse().StatusCode)
}

func TestClient_Uploade_newRequestError(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	_, err := backlog.ExportClientUploade(c, "", "testdata/testfile", "fname")
	assert.NotNil(t, err)
}

func TestClient_Uploade_emptyFilePath(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	_, err := backlog.ExportClientUploade(c, "spath", "", "fname")
	assert.Error(t, err, "file's path and name is required")
}

func TestClient_Uploade_emptyFileName(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	_, err := backlog.ExportClientUploade(c, "spath", "fpath", "")
	assert.Error(t, err, "file's path and name is required")
}

func TestCeckResponseError(t *testing.T) {
	cases := map[string]struct {
		statusCode int
		wantError  bool
	}{
		"199": {
			statusCode: 199,
			wantError:  true,
		},
		"200": {
			statusCode: 200,
			wantError:  false,
		},
		"299": {
			statusCode: 299,
			wantError:  false,
		},
		"300": {
			statusCode: 300,
			wantError:  true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			body := ioutil.NopCloser(bytes.NewReader([]byte(`{"errors":[{"message": "No project.","code": 6,"moreInfo": ""}]}`)))

			resp := &backlog.ExportResponse{
				Response: &http.Response{
					StatusCode: tc.statusCode,
					Body:       body,
				},
				Error: &backlog.APIResponseError{},
			}

			if r, err := backlog.ExportCeckResponseError(resp); tc.wantError {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, resp, r)
			}
		})
	}
}

func TestCeckResponseError_emptyBody(t *testing.T) {
	resp := &backlog.ExportResponse{
		Response: &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       nil,
		},
		Error: &backlog.APIResponseError{},
	}
	_, err := backlog.ExportCeckResponseError(resp)
	assert.Error(t, err, "response body is empty")
}

func TestCeckResponseError_invalidJSON(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(`{{"errors":[{"message": "No project.","code": 6,"moreInfo": ""}]}`)))

	resp := &backlog.ExportResponse{
		Response: &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       body,
		},
		Error: &backlog.APIResponseError{},
	}
	want := &json.SyntaxError{}
	if _, err := backlog.ExportCeckResponseError(resp); err == nil {
		assert.NotNil(t, err)
	} else {
		assert.Equal(t, reflect.TypeOf(want), reflect.TypeOf(err))
	}
}
