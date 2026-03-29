package backlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWikiService_All(t *testing.T) {
	const testWiki1ID = 112
	const testWiki2ID = 115
	const testWiki1Name = "test1"
	const testWiki2Name = "test2"

	o := newWikiOptionService()

	cases := map[string]struct {
		// Input arguments
		projectIDOrKey string
		options        []*QueryOption // Variable arguments

		// HTTP mock settings
		httpStatus int
		httpBody   string
		httpError  error

		// API request assertion
		wantSpath               string
		wantQueryProjectIDOrKey string
		expectAPICall           bool

		// Expected results (Error handling)
		wantError   bool
		wantErrType error

		// Expected results (Success case only)
		wantIDs   []int
		wantNames []string
	}{
		"success-project-id": {
			projectIDOrKey:          "103",
			httpStatus:              http.StatusOK,
			httpBody:                testdataWikiListJSON,
			wantSpath:               "wikis",
			wantQueryProjectIDOrKey: "103",
			expectAPICall:           true,
			wantIDs:                 []int{testWiki1ID, testWiki2ID},
			wantNames:               []string{testWiki1Name, testWiki2Name},
		},
		"success-with-options": {
			projectIDOrKey: "PRJ_KEY",
			options: []*QueryOption{
				o.WithQueryKeyword("test"),
			},
			httpStatus:              http.StatusOK,
			httpBody:                testdataWikiListJSON,
			wantSpath:               "wikis",
			wantQueryProjectIDOrKey: "PRJ_KEY",
			expectAPICall:           true,
			wantIDs:                 []int{testWiki1ID, testWiki2ID},
			wantNames:               []string{testWiki1Name, testWiki2Name},
		},
		"validation-error-key-empty": {
			projectIDOrKey: "",
			expectAPICall:  false,
			wantError:      true,
			wantErrType:    &ValidationError{},
		},
		"validation-error-invalid-option-type": {
			projectIDOrKey: "PRJ",
			options: []*QueryOption{{
				t:         queryCount,
				checkFunc: nil,
				setFunc: func(p *QueryParams) error {
					return nil
				},
			}},
			expectAPICall: false,
			wantError:     true,
			wantErrType:   &InvalidQueryOptionError{},
		},
		"validation-error-option-set-fail": {
			projectIDOrKey: "PRJ",
			options:        []*QueryOption{newQueryOptionWithSetError(queryKeyword)},
			expectAPICall:  false,
			wantError:      true,
		},
		"client-error-network-failure": {
			projectIDOrKey:          "1",
			options:                 []*QueryOption{},
			httpError:               errors.New("network error"),
			wantSpath:               "wikis",
			wantQueryProjectIDOrKey: "1",
			expectAPICall:           true,
			wantError:               true,
		},
		"api-error-invalid-json": {
			projectIDOrKey:          "1",
			options:                 []*QueryOption{},
			httpStatus:              http.StatusOK,
			httpBody:                testdataInvalidJSON,
			wantSpath:               "wikis",
			wantQueryProjectIDOrKey: "1",
			expectAPICall:           true,
			wantError:               true,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			calledAPICall := false
			s := newWikiService()
			s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
				calledAPICall = true

				if tc.expectAPICall {
					assert.Equal(t, tc.wantSpath, spath)
					assert.Equal(t, tc.wantQueryProjectIDOrKey, query.Get("projectIdOrKey"))
				}

				resp := &http.Response{
					StatusCode: tc.httpStatus,
					Body:       io.NopCloser(bytes.NewReader([]byte(tc.httpBody))),
				}
				return resp, tc.httpError
			}

			wikis, err := s.All(tc.projectIDOrKey, tc.options...)

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
				assert.Nil(t, wikis)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, wikis)

				assert.Len(t, wikis, len(tc.wantIDs))
				for i := range wikis {
					assert.Equal(t, tc.wantIDs[i], wikis[i].ID)
					assert.Equal(t, tc.wantNames[i], wikis[i].Name)
				}
			}
		})
	}
}

func TestWikiService_Count(t *testing.T) {
	cases := map[string]struct {
		// Input arguments
		projectIDOrKey string

		// HTTP mock settings
		httpStatus int
		httpBody   string
		httpError  error

		// API request assertion
		wantSpath      string
		wantQueryParam string
		expectAPICall  bool

		// Expected results (Error handling)
		wantError   bool
		wantErrType error

		// Expected results (Success case only)
		wantCount int
	}{
		"success-project-id": {
			projectIDOrKey: "103",
			httpStatus:     http.StatusOK,
			httpBody:       `{"count": 34}`,
			wantSpath:      "wikis/count",
			wantQueryParam: "103",
			expectAPICall:  true,
			wantCount:      34,
		},
		"success-project-key": {
			projectIDOrKey: "PRJ_KEY",
			httpStatus:     http.StatusOK,
			httpBody:       `{"count": 10}`,
			wantSpath:      "wikis/count",
			wantQueryParam: "PRJ_KEY",
			expectAPICall:  true,
			wantCount:      10,
		},
		"validation-error-key-empty": {
			projectIDOrKey: "",
			expectAPICall:  false,
			wantError:      true,
			wantErrType:    &ValidationError{},
		},
		"client-error-network-failure": {
			projectIDOrKey: "1",
			httpError:      errors.New("network error"),
			wantSpath:      "wikis/count",
			wantQueryParam: "1",
			expectAPICall:  true,
			wantError:      true,
		},
		"api-error-invalid-json": {
			projectIDOrKey: "1",
			httpStatus:     http.StatusOK,
			httpBody:       testdataInvalidJSON,
			wantSpath:      "wikis/count",
			wantQueryParam: "1",
			expectAPICall:  true,
			wantError:      true,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			calledAPICall := false
			s := newWikiService()
			s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
				calledAPICall = true

				if tc.expectAPICall {
					assert.Equal(t, tc.wantSpath, spath)
					assert.Equal(t, tc.wantQueryParam, query.Get("projectIdOrKey"))
				}

				resp := &http.Response{
					StatusCode: tc.httpStatus,
					Body:       io.NopCloser(bytes.NewReader([]byte(tc.httpBody))),
				}
				return resp, tc.httpError
			}

			count, err := s.Count(tc.projectIDOrKey)

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
				assert.Equal(t, 0, count)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantCount, count)
			}
		})
	}
}

func TestWikiService_One(t *testing.T) {
	cases := map[string]struct {
		// Input arguments
		wikiID int

		// HTTP mock settings
		httpStatus int
		httpBody   string
		httpError  error

		// API request assertion
		wantSpath     string
		expectAPICall bool

		// Expected results (Error handling)
		wantError   bool
		wantErrType error

		// Expected results (Success case only)
		wantWikiID   int
		wantWikiName string
	}{
		"success-normal": {
			wikiID:        34,
			httpStatus:    http.StatusOK,
			httpBody:      testdataWikiMaximumJSON,
			wantSpath:     "wikis/34",
			expectAPICall: true,
			wantWikiID:    34,
			wantWikiName:  "Maximum Wiki Page",
		},

		"validation-error-id-zero": {
			wikiID:        0,
			expectAPICall: false,
			wantError:     true,
			wantErrType:   &ValidationError{},
		},
		"validation-error-id-negative": {
			wikiID:        -1,
			expectAPICall: false,
			wantError:     true,
			wantErrType:   &ValidationError{},
		},
		"client-error-network-failure": {
			wikiID:        1,
			httpError:     errors.New("network error"),
			wantSpath:     "wikis/1",
			expectAPICall: true,
			wantError:     true,
		},
		"api-error-invalid-json": {
			wikiID:        1,
			httpStatus:    http.StatusOK,
			httpBody:      testdataInvalidJSON,
			wantSpath:     "wikis/1",
			expectAPICall: true,
			wantError:     true,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			calledAPICall := false
			s := newWikiService()

			s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
				calledAPICall = true

				if tc.expectAPICall {
					assert.Equal(t, tc.wantSpath, spath)
					assert.Nil(t, query)
				}

				resp := &http.Response{
					StatusCode: tc.httpStatus,
					Body:       io.NopCloser(bytes.NewReader([]byte(tc.httpBody))),
				}
				return resp, tc.httpError
			}

			wiki, err := s.One(tc.wikiID)

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
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, wiki)

			assert.Equal(t, tc.wantWikiID, wiki.ID)
			assert.Equal(t, tc.wantWikiName, wiki.Name)
		})
	}
}

func TestWikiService_Create(t *testing.T) {
	o := newWikiOptionService()

	type wantAPI struct {
		called bool
		spath  string
		form   map[string]string
	}

	cases := map[string]struct {
		projectID int
		name      string
		content   string
		opts      []*FormOption

		httpStatus int
		httpBody   string
		httpError  error

		wantAPI     wantAPI
		wantWiki    *Wiki
		wantErrType error
	}{
		"success-minimum": {
			projectID: 56,
			name:      "Minimum Wiki Page",
			content:   "This is a minimal wiki page.",

			httpStatus: http.StatusOK,
			httpBody:   testdataWikiMinimumJSON,

			wantAPI: wantAPI{
				called: true,
				spath:  "wikis",
				form: map[string]string{
					"projectId": "56",
					"name":      "Minimum Wiki Page",
					"content":   "This is a minimal wiki page.",
				},
			},

			wantWiki: &Wiki{
				ID:      34,
				Name:    "Minimum Wiki Page",
				Content: "This is a minimal wiki page.",
			},
		},

		"success-with-mailNotify": {
			projectID: 56,
			name:      "Minimum Wiki Page",
			content:   "This is a minimal wiki page.",
			opts:      []*FormOption{o.WithFormMailNotify(true)},

			httpStatus: http.StatusOK,
			httpBody:   testdataWikiMinimumJSON,

			wantAPI: wantAPI{
				called: true,
				spath:  "wikis",
				form: map[string]string{
					"projectId":  "56",
					"name":       "Minimum Wiki Page",
					"content":    "This is a minimal wiki page.",
					"mailNotify": "true",
				},
			},

			wantWiki: &Wiki{
				ID:      34,
				Name:    "Minimum Wiki Page",
				Content: "This is a minimal wiki page.",
			},
		},

		"validation-error-projectID-zero": {
			projectID: 0,
			name:      "Test",
			content:   "test",

			wantAPI: wantAPI{
				called: false,
			},

			wantErrType: &ValidationError{},
		},

		"validation-error-name-empty": {
			projectID: 1,
			name:      "",
			content:   "test",

			wantAPI: wantAPI{
				called: false,
			},

			wantErrType: &ValidationError{},
		},

		"validation-error-content-empty": {
			projectID: 1,
			name:      "Test",
			content:   "",

			wantAPI: wantAPI{
				called: false,
			},

			wantErrType: &ValidationError{},
		},

		"validation-error-option-set-fail": {
			projectID: 1,
			name:      "Test",
			content:   "content",
			opts:      []*FormOption{newFormOptionWithSetError(formMailNotify)},

			wantAPI: wantAPI{
				called: false,
			},

			wantErrType: errors.New(""),
		},

		"validation-error-invalid-option-type": {
			projectID: 1,
			name:      "Test",
			content:   "content",
			opts: []*FormOption{
				{
					formMailAddress,
					nil,
					func(p *FormParams) error {
						return nil
					},
				},
			},

			wantAPI: wantAPI{
				called: false,
			},

			wantErrType: &InvalidFormOptionError{},
		},

		"client-error-network-failure": {
			projectID: 1,
			name:      "Test",
			content:   "content",

			httpError: errors.New("network error"),

			wantAPI: wantAPI{
				called: true,
				spath:  "wikis",
				form: map[string]string{
					"projectId": "1",
					"name":      "Test",
					"content":   "content",
				},
			},

			wantErrType: errors.New(""),
		},

		"api-error-invalid-json": {
			projectID: 1,
			name:      "Test",
			content:   "content",

			httpStatus: http.StatusOK,
			httpBody:   testdataInvalidJSON,

			wantAPI: wantAPI{
				called: true,
				spath:  "wikis",
				form: map[string]string{
					"projectId": "1",
					"name":      "Test",
					"content":   "content",
				},
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			called := false
			s := newWikiService()

			s.method.Post = func(spath string, form *FormParams) (*http.Response, error) {
				called = true

				if tc.wantAPI.called {
					assert.Equal(t, tc.wantAPI.spath, spath)

					for k, v := range tc.wantAPI.form {
						assert.Equal(t, v, form.Get(k))
					}
				}

				return &http.Response{
					StatusCode: tc.httpStatus,
					Body:       io.NopCloser(bytes.NewReader([]byte(tc.httpBody))),
				}, tc.httpError
			}

			wiki, err := s.Create(tc.projectID, tc.name, tc.content, tc.opts...)

			assert.Equal(t, tc.wantAPI.called, called)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, wiki)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, wiki)

			assert.Equal(t, tc.wantWiki.ID, wiki.ID)
			assert.Equal(t, tc.wantWiki.Name, wiki.Name)
			assert.Equal(t, tc.wantWiki.Content, wiki.Content)
		})
	}
}

func TestWikiService_Update(t *testing.T) {
	o := newWikiOptionService()

	type wantAPI struct {
		called bool
		spath  string
		form   map[string]string
	}

	cases := map[string]struct {
		wikiID int
		option *FormOption
		opts   []*FormOption

		httpStatus int
		httpBody   string
		httpError  error

		wantAPI     wantAPI
		wantErrType error
		wantWiki    *Wiki
	}{
		"success-name-only": {
			wikiID: 34,
			option: o.WithFormName("New Page Name"),

			httpStatus: http.StatusOK,
			httpBody:   testdataWikiMaximumJSON,

			wantAPI: wantAPI{
				called: true,
				spath:  "wikis/34",
				form: map[string]string{
					"name": "New Page Name",
				},
			},

			wantWiki: &Wiki{
				ID:      34,
				Name:    "Maximum Wiki Page",
				Content: "This is a muximal wiki page.",
			},
		},

		"success-full-options": {
			wikiID: 34,
			option: o.WithFormName("Full Options Name"),
			opts: []*FormOption{
				o.WithFormContent("Full Options Content"),
				o.WithFormMailNotify(true),
			},

			httpStatus: http.StatusOK,
			httpBody:   testdataWikiMaximumJSON,

			wantAPI: wantAPI{
				called: true,
				spath:  "wikis/34",
				form: map[string]string{
					"name":       "Full Options Name",
					"content":    "Full Options Content",
					"mailNotify": "true",
				},
			},

			wantWiki: &Wiki{
				ID:      34,
				Name:    "Maximum Wiki Page",
				Content: "This is a muximal wiki page.",
			},
		},

		"validation-error-required-option": {
			wikiID: 12,
			option: o.WithFormMailNotify(true),

			wantAPI: wantAPI{
				called: false,
			},

			wantErrType: &ValidationError{},
		},

		"validation-error-invalid-wikiID": {
			wikiID: 0,
			option: o.WithFormName("New Name"),

			wantAPI: wantAPI{
				called: false,
			},

			wantErrType: &ValidationError{},
		},

		"client-error-network-failure": {
			wikiID: 13,
			option: o.WithFormName("New Name"),

			httpError: errors.New("network error"),

			wantAPI: wantAPI{
				called: true,
				spath:  "wikis/13",
				form: map[string]string{
					"name": "New Name",
				},
			},

			wantErrType: errors.New(""),
		},

		"api-error-invalid-json": {
			wikiID: 14,
			option: o.WithFormName("New Name"),

			httpStatus: http.StatusOK,
			httpBody:   testdataInvalidJSON,

			wantAPI: wantAPI{
				called: true,
				spath:  "wikis/14",
				form: map[string]string{
					"name": "New Name",
				},
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			called := false
			s := newWikiService()

			s.method.Patch = func(spath string, form *FormParams) (*http.Response, error) {
				called = true

				if tc.wantAPI.called {
					assert.Equal(t, tc.wantAPI.spath, spath)

					for k, v := range tc.wantAPI.form {
						assert.Equal(t, v, form.Get(k))
					}
				}

				return &http.Response{
					StatusCode: tc.httpStatus,
					Body:       io.NopCloser(bytes.NewReader([]byte(tc.httpBody))),
				}, tc.httpError
			}

			wiki, err := s.Update(tc.wikiID, tc.option, tc.opts...)

			assert.Equal(t, tc.wantAPI.called, called)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, wiki)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, wiki)

			assert.Equal(t, tc.wantWiki.ID, wiki.ID)
			assert.Equal(t, tc.wantWiki.Name, wiki.Name)
			assert.Equal(t, tc.wantWiki.Content, wiki.Content)
		})
	}
}

func TestWikiService_Delete(t *testing.T) {
	o := newWikiOptionService()
	projectOption := newProjectOptionService()

	type wantAPI struct {
		called bool
		spath  string
		form   map[string]string
	}

	cases := map[string]struct {
		wikiID int
		opts   []*FormOption

		httpStatus int
		httpBody   string
		httpError  error

		wantAPI     wantAPI
		wantWikiID  int
		wantErrType error
	}{
		"success-with-option": {
			wikiID: 34,
			opts:   []*FormOption{o.WithFormMailNotify(true)},

			httpStatus: http.StatusOK,
			httpBody:   testdataWikiMaximumJSON,

			wantAPI: wantAPI{
				called: true,
				spath:  "wikis/34",
				form: map[string]string{
					"mailNotify": "true",
				},
			},

			wantWikiID: 34,
		},

		"success-no-option": {
			wikiID: 1,

			httpStatus: http.StatusOK,
			httpBody:   testdataWikiMaximumJSON,

			wantAPI: wantAPI{
				called: true,
				spath:  "wikis/1",
			},

			wantWikiID: 34,
		},

		"validation-error-id-zero": {
			wikiID: 0,

			wantAPI: wantAPI{
				called: false,
			},

			wantErrType: &ValidationError{},
		},

		"validation-error-id-negative": {
			wikiID: -1,

			wantAPI: wantAPI{
				called: false,
			},

			wantErrType: &ValidationError{},
		},

		"validation-error-option-set-fail": {
			wikiID: 1,
			opts:   []*FormOption{newFormOptionWithSetError(formMailNotify)},

			wantAPI: wantAPI{
				called: false,
			},

			wantErrType: errors.New(""),
		},

		"validation-error-invalid-option-type": {
			wikiID: 1,
			opts: []*FormOption{
				projectOption.WithFormKey("Invalid Option"),
			},

			wantAPI: wantAPI{
				called: false,
			},

			wantErrType: &InvalidFormOptionError{},
		},

		"client-error-network-failure": {
			wikiID:    34,
			httpError: errors.New("network error"),

			wantAPI: wantAPI{
				called: true,
				spath:  "wikis/34",
			},

			wantErrType: errors.New(""),
		},

		"api-error-invalid-json": {
			wikiID:     34,
			httpStatus: http.StatusOK,
			httpBody:   testdataInvalidJSON,

			wantAPI: wantAPI{
				called: true,
				spath:  "wikis/34",
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			called := false
			s := newWikiService()

			s.method.Delete = func(spath string, form *FormParams) (*http.Response, error) {
				called = true

				if tc.wantAPI.called {
					assert.Equal(t, tc.wantAPI.spath, spath)

					for k, v := range tc.wantAPI.form {
						assert.Equal(t, v, form.Get(k))
					}
				}

				return &http.Response{
					StatusCode: tc.httpStatus,
					Body:       io.NopCloser(bytes.NewReader([]byte(tc.httpBody))),
				}, tc.httpError
			}

			wiki, err := s.Delete(tc.wikiID, tc.opts...)

			assert.Equal(t, tc.wantAPI.called, called)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, wiki)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, wiki)

			assert.Equal(t, tc.wantWikiID, wiki.ID)
		})
	}
}
