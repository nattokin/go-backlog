package backlog_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWikiService_All(t *testing.T) {
	t.Parallel()

	projectIDOrKey := "103"

	want := struct {
		spath          string
		projectIDOrKey string
		idList         []int
		nameList       []string
	}{
		spath:          "wikis",
		projectIDOrKey: projectIDOrKey,
		idList:         []int{112, 115},
		nameList:       []string{"test1", "test2"},
	}

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.projectIDOrKey, query.Get("projectIdOrKey"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiListJSON))),
			}
			return resp, nil
		},
	})

	wikis, err := s.All(projectIDOrKey)
	assert.NoError(t, err)
	count := len(wikis)
	assert.Equal(t, len(want.idList), count)
	for i := 0; i < count; i++ {
		assert.Equal(t, want.idList[i], wikis[i].ID)
		assert.Equal(t, want.nameList[i], wikis[i].Name)
	}
}

func TestWikiService_All_param(t *testing.T) {
	option := &backlog.WikiOptionService{}
	cases := map[string]struct {
		projectIdOrKey string
		keywordOption  *backlog.QueryOption
		content        string
		mailNotify     bool
		wantError      bool
	}{
		"valid-1": {
			projectIdOrKey: "1",
			keywordOption:  option.WithQueryKeyword("test"),
			wantError:      false,
		},
		"valid-2": {
			projectIdOrKey: "TEST",
			keywordOption:  option.WithQueryKeyword(""),
			wantError:      false,
		},
		"invalid-ProjectID": {
			projectIdOrKey: "0",
			keywordOption:  option.WithQueryKeyword("test"),
			wantError:      true,
		},
		"invalid-ProjectKey": {
			projectIdOrKey: "",
			keywordOption:  option.WithQueryKeyword("test"),
			wantError:      true,
		},
		"invalid-option": {
			projectIdOrKey: "TEST",
			keywordOption: backlog.ExportNewQueryOption(0, func(p *backlog.QueryParams) error {
				return nil
			}),
			wantError: true,
		},
		"option-error": {
			projectIdOrKey: "TEST",
			keywordOption: backlog.ExportNewQueryOption(backlog.ExportQueryKeyword, func(p *backlog.QueryParams) error {
				return errors.New("error")
			}),
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			s := &backlog.WikiService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiListJSON))),
					}
					return resp, nil
				},
			})

			if wikis, err := s.All(tc.projectIdOrKey, tc.keywordOption); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, wikis, 2)
			}
		})

	}
}

func TestWikiService_All_param_error(t *testing.T) {
	t.Parallel()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiListJSON))),
			}
			return resp, nil
		},
	})
	{
		wikis, err := s.All("0")
		assert.Error(t, err)
		assert.Nil(t, wikis)
	}
	{
		wikis, err := s.All("")
		assert.Error(t, err)
		assert.Nil(t, wikis)
	}
}

func TestWikiService_All_clientError(t *testing.T) {
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})

	wikis, err := s.All("TEST")
	assert.Error(t, err)
	assert.Nil(t, wikis)
}

func TestWikiService_All_invalidJson(t *testing.T) {
	t.Parallel()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

	wikis, err := s.All("TEST")
	assert.Error(t, err)
	assert.Nil(t, wikis)
}

func TestWikiService_Count(t *testing.T) {
	t.Parallel()

	projectKey := "TEST"

	want := struct {
		spath      string
		projectKey string
		count      int
	}{
		spath:      "wikis/count",
		projectKey: projectKey,
		count:      5,
	}

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.projectKey, query.Get("projectIdOrKey"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"count":5}`)),
			}
			return resp, nil
		},
	})

	count, err := s.Count(projectKey)
	require.NoError(t, err)
	assert.Equal(t, want.count, count)
}

func TestWikiService_Count_param_error(t *testing.T) {
	t.Parallel()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiListJSON))),
			}
			return resp, nil
		},
	})

	{
		count, err := s.Count("0")
		assert.Error(t, err)
		assert.Zero(t, count)

	}
	{
		count, err := s.Count("")
		assert.Error(t, err)
		assert.Zero(t, count)
	}
}

func TestWikiService_Count_clientError(t *testing.T) {
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})

	count, err := s.Count("TEST")
	assert.Error(t, err)
	assert.Zero(t, count)
}

func TestWikiService_Count_invalidJson(t *testing.T) {
	t.Parallel()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

	count, err := s.Count("TEST")
	assert.Error(t, err)
	assert.Zero(t, count)
}

func TestWikiService_One(t *testing.T) {
	t.Parallel()

	wikiID := 34

	want := struct {
		spath  string
		wikiID int
		name   string
	}{
		spath:  "wikis/" + strconv.Itoa(wikiID),
		wikiID: wikiID,
		name:   "Maximum Wiki Page",
	}

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Nil(t, query)

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiMaximumJSON))),
			}
			return resp, nil
		},
	})

	wiki, err := s.One(wikiID)
	require.NoError(t, err)
	assert.Equal(t, want.wikiID, wiki.ID)
	assert.Equal(t, want.name, wiki.Name)
}

func TestWikiService_One_param(t *testing.T) {
	cases := map[string]struct {
		wikiID    int
		wantError bool
	}{
		"valid": {
			wikiID:    1,
			wantError: false,
		},
		"invalid-1": {
			wikiID:    0,
			wantError: true,
		},
		"invalid-2": {
			wikiID:    -1,
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			s := &backlog.WikiService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiMaximumJSON))),
					}
					return resp, nil
				},
			})

			if wiki, err := s.One(tc.wikiID); tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, wiki)
			} else {
				require.NoError(t, err)
				assert.Equal(t, 34, wiki.ID)
			}
		})

	}
}

func TestWikiService_One_clientError(t *testing.T) {
	t.Parallel()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})

	wiki, err := s.One(1)
	assert.Error(t, err)
	assert.Nil(t, wiki)
}

func TestWikiService_One_invalidJson(t *testing.T) {
	t.Parallel()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

	wiki, err := s.One(1)
	assert.Error(t, err)
	assert.Nil(t, wiki)
}

func TestWikiService_Create(t *testing.T) {
	t.Parallel()

	wikiID := 34
	name := "Minimum Wiki Page"
	content := "This is a minimal wiki page."

	want := struct {
		wikiID     int
		spath      string
		name       string
		content    string
		mailNotify string
	}{
		wikiID:     wikiID,
		spath:      "wikis",
		name:       name,
		content:    content,
		mailNotify: "true",
	}

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.NotNil(t, form)
			assert.Equal(t, want.name, form.Get("name"))
			assert.Equal(t, want.content, form.Get("content"))
			assert.Equal(t, want.mailNotify, form.Get("mailNotify"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiMinimumJSON))),
			}
			return resp, nil
		},
	})

	wiki, err := s.Create(56, name, content, s.Option.WithFormMailNotify(true))
	require.NoError(t, err)
	assert.Equal(t, want.wikiID, wiki.ID)
	assert.Equal(t, want.name, wiki.Name)
	assert.Equal(t, want.content, wiki.Content)
}

func TestWikiService_Create_param(t *testing.T) {
	cases := map[string]struct {
		projectID  int
		name       string
		content    string
		mailNotify bool
		wantError  bool
	}{
		"no_error_1": {
			projectID: 1,
			name:      "Test",
			content:   "test",
			wantError: false,
		},
		"no_error_2": {
			projectID: 100,
			name:      "Test Name",
			content:   "test content",
			wantError: false,
		},
		"projectId_zero": {
			projectID: 0,
			name:      "Test",
			content:   "test",
			wantError: true,
		},
		"name_empty": {
			projectID: 1,
			name:      "",
			content:   "test",
			wantError: true,
		},
		"content_empty": {
			projectID: 1,
			name:      "Test",
			content:   "",
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			s := &backlog.WikiService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiMinimumJSON))),
					}
					return resp, nil
				},
			})

			if resp, err := s.Create(tc.projectID, tc.name, tc.content); tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})

	}
}

func TestWikiService_Create_clientError(t *testing.T) {
	t.Parallel()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})

	resp, err := s.Create(1, "name", "test")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestWikiService_Create_invalidJson(t *testing.T) {
	t.Parallel()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

	wiki, err := s.Create(1, "name", "test")
	assert.Error(t, err)
	assert.Nil(t, wiki)
}

func TestWikiService_Create_OptionError(t *testing.T) {
	t.Parallel()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiMinimumJSON))),
			}
			return resp, nil
		},
	})

	errorOption := backlog.ExportNewFormOption(backlog.ExportFormMailNotify, func(p *backlog.ExportRequestParams) error {
		return errors.New("error")
	})

	wiki, err := s.Create(1, "name", "content", errorOption)
	assert.Error(t, err)
	assert.Nil(t, wiki)
}

func TestWikiService_Create_invalidOption(t *testing.T) {
	t.Parallel()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiMinimumJSON))),
			}
			return resp, nil
		},
	})

	invalidOption := backlog.ExportNewFormOption(0, func(p *backlog.ExportRequestParams) error {
		return nil
	})

	wiki, err := s.Create(1, "name", "content", invalidOption)
	assert.IsType(t, &backlog.InvalidFormOptionError{}, err)
	assert.Nil(t, wiki)
}

func TestWikiService_Update(t *testing.T) {
	t.Parallel()

	option := &backlog.WikiOptionService{}

	type testCase struct {
		wikiID int
		option *backlog.FormOption
		opts   []*backlog.FormOption

		httpStatus int
		httpBody   string
		httpError  error

		wantError   bool
		wantErrType interface{}
		wantID      int
		wantName    string
	}

	cases := map[string]testCase{
		"success-name-only": {
			wikiID:     34,
			option:     option.WithFormName("New Page Name"),
			opts:       []*backlog.FormOption{},
			httpStatus: http.StatusOK,
			httpBody:   testdataWikiMaximumJSON,
			wantError:  false,
			wantID:     34,
			wantName:   "Maximum Wiki Page",
		},
		"success-content-and-notify": {
			wikiID:     34,
			option:     option.WithFormContent("Updated content."),
			opts:       []*backlog.FormOption{option.WithFormMailNotify(true)},
			httpStatus: http.StatusOK,
			httpBody:   testdataWikiMaximumJSON,
			wantError:  false,
			wantID:     34,
			wantName:   "Maximum Wiki Page",
		},
		"validation-error-required-option": {
			wikiID:      12,
			option:      option.WithFormMailNotify(true),
			opts:        []*backlog.FormOption{},
			httpStatus:  http.StatusBadRequest,
			httpBody:    "",
			wantError:   true,
			wantErrType: &backlog.ValidationError{},
		},
		"validation-error-required-content": {
			wikiID:      12,
			option:      option.WithFormMailNotify(true),
			opts:        []*backlog.FormOption{},
			httpStatus:  http.StatusBadRequest,
			httpBody:    "",
			wantError:   true,
			wantErrType: &backlog.ValidationError{},
		},
		"validation-error-invalid-wikiID": {
			wikiID:      0,
			option:      option.WithFormName("New Name"),
			httpStatus:  http.StatusBadRequest,
			httpBody:    "",
			wantError:   true,
			wantErrType: &backlog.ValidationError{},
		},
		"validation-error-invalid-option-type": {
			wikiID: 12,
			option: option.WithFormMailNotify(false),
			opts: []*backlog.FormOption{
				(&backlog.ProjectOptionService{}).WithFormKey("Invalid Option"),
			},
			wantError:   true,
			wantErrType: &backlog.InvalidFormOptionError{},
			wantID:      12,
		},
		"client-error-network-failure": {
			wikiID:      13,
			option:      option.WithFormName("New Name"),
			httpStatus:  http.StatusOK,
			httpBody:    "",
			httpError:   errors.New("network error"),
			wantError:   true,
			wantErrType: nil,
		},
		"api-error-invalid-json": {
			wikiID:      14,
			option:      option.WithFormName("New Name"),
			httpStatus:  http.StatusOK,
			httpBody:    testdataInvalidJSON,
			wantError:   true,
			wantErrType: nil,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := &backlog.WikiService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Patch: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: tc.httpStatus,
						Body:       io.NopCloser(bytes.NewReader([]byte(tc.httpBody))),
					}
					return resp, tc.httpError
				},
			})

			wiki, err := s.Update(tc.wikiID, tc.option, tc.opts...)
			if tc.wantError {
				assert.Error(t, err)
				if tc.wantErrType != nil {
					assert.IsType(t, tc.wantErrType, err)
				}
			} else {
				assert.NoError(t, err)
				require.NotNil(t, wiki)
				assert.Equal(t, tc.wantID, wiki.ID)
				assert.Equal(t, tc.wantName, wiki.Name)
			}
		})
	}
}

func TestWikiService_Delete(t *testing.T) {
	t.Parallel()

	option := &backlog.WikiOptionService{}

	type testCase struct {
		wikiID         int
		options        []*backlog.FormOption
		httpStatus     int
		httpBody       string
		httpError      error
		wantSpath      string
		wantMailNotify string
		expectAPICall  bool
		wantError      bool
		wantErrType    interface{}
		wantID         int
	}

	cases := map[string]testCase{
		"success-with-option": {
			wikiID:         34,
			options:        []*backlog.FormOption{option.WithFormMailNotify(true)},
			httpStatus:     http.StatusOK,
			httpBody:       testdataWikiMaximumJSON,
			wantSpath:      "wikis/34",
			wantMailNotify: "true",
			expectAPICall:  true,
			wantError:      false,
			wantID:         34,
		},
		"success-no-option": {
			wikiID:         1,
			options:        []*backlog.FormOption{},
			httpStatus:     http.StatusOK,
			httpBody:       testdataWikiMaximumJSON,
			wantSpath:      "wikis/1",
			wantMailNotify: "",
			expectAPICall:  true,
			wantError:      false,
			wantID:         34,
		},
		"validation-error-id-zero": {
			wikiID:        0,
			expectAPICall: false,
			wantError:     true,
			wantErrType:   &backlog.ValidationError{},
		},
		"validation-error-id-negative": {
			wikiID:        -1,
			expectAPICall: false,
			wantError:     true,
			wantErrType:   &backlog.ValidationError{},
		},
		"client-error-network-failure": {
			wikiID:        34,
			options:       []*backlog.FormOption{},
			httpStatus:    http.StatusOK,
			httpError:     errors.New("network error"),
			wantSpath:     "wikis/34",
			expectAPICall: true,
			wantError:     true,
			wantID:        0,
		},
		"api-error-invalid-json": {
			wikiID:        34,
			options:       []*backlog.FormOption{},
			httpStatus:    http.StatusOK,
			httpBody:      testdataInvalidJSON,
			wantSpath:     "wikis/34",
			expectAPICall: true,
			wantError:     true,
			wantID:        0,
		},
		"validation-error-option-set-fail": {
			wikiID: 1,
			options: []*backlog.FormOption{
				backlog.ExportNewFormOption(
					backlog.ExportFormMailNotify,
					func(p *backlog.ExportRequestParams) error {
						return errors.New("error during option set")
					},
				),
			},
			expectAPICall: false,
			wantError:     true,
			wantID:        0,
		},
		"validation-error-invalid-option-type": {
			wikiID: 1,
			options: []*backlog.FormOption{
				option.WithFormName("Invalid Option"),
			},
			expectAPICall: false,
			wantError:     true,
			wantErrType:   &backlog.InvalidFormOptionError{},
			wantID:        0,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			calledAPICall := false
			s := &backlog.WikiService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
					calledAPICall = true

					if tc.expectAPICall {
						assert.Equal(t, tc.wantSpath, spath)
						if tc.wantMailNotify != "" {
							assert.Equal(t, tc.wantMailNotify, form.Get("mailNotify"))
						} else {
							assert.Empty(t, form.Get("mailNotify"))
						}
					}

					resp := &http.Response{
						StatusCode: tc.httpStatus,
						Body:       io.NopCloser(bytes.NewReader([]byte(tc.httpBody))),
					}
					return resp, tc.httpError
				},
			})

			wiki, err := s.Delete(tc.wikiID, tc.options...)

			if tc.expectAPICall {
				assert.True(t, calledAPICall)
			} else {
				assert.False(t, calledAPICall)
			}

			if tc.wantError {
				assert.Error(t, err)
				if tc.wantErrType != nil {
					assert.IsType(t, tc.wantErrType, err)
				}
				assert.Nil(t, wiki)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, wiki)
				assert.Equal(t, tc.wantID, wiki.ID)
			}
		})
	}
}
