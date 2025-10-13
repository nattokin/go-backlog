package backlog_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWikiService_All(t *testing.T) {
	type testCase struct {
		// Input arguments
		projectIDOrKey string
		options        []*backlog.QueryOption // Variable arguments

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
	}

	// Constants extracted from testdataWikiListJSON
	const testWiki1ID = 112
	const testWiki2ID = 115
	const testWiki1Name = "test1"
	const testWiki2Name = "test2"

	o := backlog.ExportNewWikiOptionService()
	cases := map[string]testCase{
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
			options: []*backlog.QueryOption{
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
		// 1. Validation Error: projectIDOrKey is empty (validateProjectIDOrKey cover)
		"validation-error-key-empty": {
			projectIDOrKey: "",
			expectAPICall:  false,
			wantError:      true,
			wantErrType:    &backlog.ValidationError{},
		},
		// 2. Validation Error: Invalid Option Type (option.validate cover)
		"validation-error-invalid-option-type": {
			projectIDOrKey: "PRJ",
			options: []*backlog.QueryOption{backlog.ExportNewQueryOption(
				backlog.ExportQueryCount,
				nil,
				func(p *backlog.QueryParams) error {
					return nil
				},
			)}, // Invalid option for WikiService.All
			expectAPICall: false,
			wantError:     true,
			wantErrType:   &backlog.InvalidQueryOptionError{},
		},
		// 3. Option Set Error (option.set cover)
		"validation-error-option-set-fail": {
			projectIDOrKey: "PRJ",
			options: []*backlog.QueryOption{
				backlog.ExportNewQueryOption(
					backlog.ExportQueryKeyword,
					nil,
					func(p *backlog.QueryParams) error {
						return errors.New("error during option set")
					},
				),
			},
			expectAPICall: false,
			wantError:     true,
		},
		// --- Existing Failure Cases ---
		"client-error-network-failure": {
			projectIDOrKey:          "1",
			options:                 []*backlog.QueryOption{},
			httpError:               errors.New("network error"),
			wantSpath:               "wikis",
			wantQueryProjectIDOrKey: "1",
			expectAPICall:           true,
			wantError:               true,
		},
		"api-error-invalid-json": {
			projectIDOrKey:          "1",
			options:                 []*backlog.QueryOption{},
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
			s := backlog.ExportNewWikiService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
					calledAPICall = true

					// Assert the API request when expected
					if tc.expectAPICall {
						assert.Equal(t, tc.wantSpath, spath)
						assert.Equal(t, tc.wantQueryProjectIDOrKey, query.Get("projectIdOrKey"))
					}

					resp := &http.Response{
						StatusCode: tc.httpStatus,
						Body:       io.NopCloser(bytes.NewReader([]byte(tc.httpBody))),
					}
					return resp, tc.httpError
				},
			})

			wikis, err := s.All(tc.projectIDOrKey, tc.options...)

			if tc.expectAPICall {
				assert.True(t, calledAPICall)
			} else {
				assert.False(t, calledAPICall)
			}

			// Assert result
			if tc.wantError {
				assert.Error(t, err)
				if tc.wantErrType != nil {
					assert.IsType(t, tc.wantErrType, err)
				}
				assert.Nil(t, wikis)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, wikis)

				// Assert the list content
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
	t.Parallel()

	type testCase struct {
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
	}

	cases := map[string]testCase{
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
		// --- Failure Cases ---
		"validation-error-key-empty": {
			projectIDOrKey: "",
			expectAPICall:  false,
			wantError:      true,
			wantErrType:    &backlog.ValidationError{},
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
			httpBody:       testdataInvalidJSON, // JSON structure does not match `{"count": N}`
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
			s := backlog.ExportNewWikiService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
					calledAPICall = true

					// Assert the API request when expected
					if tc.expectAPICall {
						assert.Equal(t, tc.wantSpath, spath)
						assert.Equal(t, tc.wantQueryParam, query.Get("projectIdOrKey"))
					}

					resp := &http.Response{
						StatusCode: tc.httpStatus,
						Body:       io.NopCloser(bytes.NewReader([]byte(tc.httpBody))),
					}
					return resp, tc.httpError
				},
			})

			count, err := s.Count(tc.projectIDOrKey)

			if tc.expectAPICall {
				assert.True(t, calledAPICall)
			} else {
				assert.False(t, calledAPICall)
			}

			// Assert result
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
	t.Parallel()

	type testCase struct {
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
		wantID   int
		wantName string
	}

	cases := map[string]testCase{
		"success-normal": {
			wikiID:        34,
			httpStatus:    http.StatusOK,
			httpBody:      testdataWikiMaximumJSON,
			wantSpath:     "wikis/34",
			expectAPICall: true,
			wantID:        34,
			wantName:      "Maximum Wiki Page",
		},
		// --- Failure Cases ---
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
			s := backlog.ExportNewWikiService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
					calledAPICall = true

					// Assert the API request when expected
					if tc.expectAPICall {
						assert.Equal(t, tc.wantSpath, spath)
						assert.Nil(t, query) // One() does not use query params
					}

					resp := &http.Response{
						StatusCode: tc.httpStatus,
						Body:       io.NopCloser(bytes.NewReader([]byte(tc.httpBody))),
					}
					return resp, tc.httpError
				},
			})

			wiki, err := s.One(tc.wikiID)

			if tc.expectAPICall {
				assert.True(t, calledAPICall)
			} else {
				assert.False(t, calledAPICall)
			}

			// Assert result
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
				assert.Equal(t, tc.wantName, wiki.Name)
			}
		})
	}
}

func TestWikiService_Create(t *testing.T) {
	t.Parallel()

	o := backlog.ExportNewWikiOptionService()

	type testCase struct {
		// Input arguments
		projectID int
		name      string
		content   string
		options   []*backlog.FormOption

		// HTTP mock settings
		httpStatus int
		httpBody   string
		httpError  error

		// API request assertion
		wantSpath      string
		wantProjectID  string
		wantMailNotify string
		expectAPICall  bool

		// Expected results (Error handling)
		wantError   bool
		wantErrType error

		// Expected results (Success case only)
		wantID      int
		wantName    string
		wantContent string
	}

	cases := map[string]testCase{
		"success-minimum": {
			projectID:     56,
			name:          "Minimum Wiki Page",
			content:       "This is a minimal wiki page.",
			options:       []*backlog.FormOption{},
			httpStatus:    http.StatusOK,
			httpBody:      testdataWikiMinimumJSON,
			wantSpath:     "wikis",
			wantProjectID: "56",
			expectAPICall: true,
			wantID:        34,
			wantName:      "Minimum Wiki Page",
			wantContent:   "This is a minimal wiki page.",
		},
		"success-with-mailNotify": {
			projectID:      56,
			name:           "Minimum Wiki Page",
			content:        "This is a minimal wiki page.",
			options:        []*backlog.FormOption{o.WithFormMailNotify(true)},
			httpStatus:     http.StatusOK,
			httpBody:       testdataWikiMinimumJSON,
			wantSpath:      "wikis",
			wantProjectID:  "56",
			wantMailNotify: "true",
			expectAPICall:  true,
			wantID:         34,
			wantName:       "Minimum Wiki Page",
			wantContent:    "This is a minimal wiki page.",
		},
		// --- Failure Cases (omitting success-only fields) ---
		"validation-error-projectID-zero": {
			projectID:     0,
			name:          "Test",
			content:       "test",
			options:       []*backlog.FormOption{},
			expectAPICall: false,
			wantError:     true,
			wantErrType:   &backlog.ValidationError{},
		},
		"validation-error-name-empty": {
			projectID:     1,
			name:          "",
			content:       "test",
			options:       []*backlog.FormOption{},
			expectAPICall: false,
			wantError:     true,
			wantErrType:   &backlog.ValidationError{},
		},
		"validation-error-content-empty": {
			projectID:     1,
			name:          "Test",
			content:       "",
			options:       []*backlog.FormOption{},
			expectAPICall: false,
			wantError:     true,
			wantErrType:   &backlog.ValidationError{},
		},
		"validation-error-option-set-fail": {
			projectID: 1,
			name:      "Test",
			content:   "content",
			options: []*backlog.FormOption{
				backlog.ExportNewFormOption(
					backlog.ExportFormMailNotify,
					nil,
					func(p *backlog.ExportRequestParams) error {
						return errors.New("error during option set")
					},
				),
			},
			expectAPICall: false,
			wantError:     true,
		},
		"validation-error-invalid-option-type": {
			projectID: 1,
			name:      "Test",
			content:   "content",
			options: []*backlog.FormOption{backlog.ExportNewFormOption(
				backlog.ExportFormMailAddress,
				nil,
				func(p *backlog.ExportRequestParams) error {
					return nil
				},
			)},
			expectAPICall: false,
			wantError:     true,
			wantErrType:   &backlog.InvalidFormOptionError{},
		},
		"client-error-network-failure": {
			projectID:     1,
			name:          "Test",
			content:       "content",
			httpError:     errors.New("network error"),
			wantSpath:     "wikis",
			wantProjectID: "1",
			expectAPICall: true,
			wantError:     true,
		},
		"api-error-invalid-json": {
			projectID:     1,
			name:          "Test",
			content:       "content",
			httpStatus:    http.StatusOK,
			httpBody:      testdataInvalidJSON,
			wantSpath:     "wikis",
			wantProjectID: "1",
			expectAPICall: true,
			wantError:     true,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			calledAPICall := false
			s := backlog.ExportNewWikiService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
					calledAPICall = true

					if tc.expectAPICall {
						assert.Equal(t, tc.wantSpath, spath)
						assert.Equal(t, tc.wantProjectID, form.Get("projectId"))

						// 💡 修正箇所: tc.wantName/wantContent -> tc.name/content に変更
						assert.Equal(t, tc.name, form.Get("name"))
						assert.Equal(t, tc.content, form.Get("content"))

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

			wiki, err := s.Create(tc.projectID, tc.name, tc.content, tc.options...)

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
				assert.Equal(t, tc.wantName, wiki.Name)
				assert.Equal(t, tc.wantContent, wiki.Content)
			}
		})
	}
}

func TestWikiService_Update(t *testing.T) {
	t.Parallel()

	o := backlog.ExportNewWikiOptionService()

	type testCase struct {
		// Input arguments
		wikiID int
		option *backlog.FormOption   // Required option argument
		opts   []*backlog.FormOption // Variable arguments (additional options)

		// HTTP mock settings
		httpStatus int
		httpBody   string
		httpError  error

		// API request assertion
		wantSpath          string
		expectAPICall      bool
		wantFormName       string
		wantFormContent    string
		wantFormMailNotify string

		// Expected results (Error handling)
		wantError   bool
		wantErrType error

		// Expected results (Success case only)
		wantID      int
		wantName    string
		wantContent string
	}

	cases := map[string]testCase{
		"success-name-only": {
			wikiID:        34,
			option:        o.WithFormName("New Page Name"),
			opts:          []*backlog.FormOption{},
			httpStatus:    http.StatusOK,
			httpBody:      testdataWikiMaximumJSON,
			wantSpath:     "wikis/34",
			expectAPICall: true,
			wantFormName:  "New Page Name",
			wantID:        34,
			wantName:      "Maximum Wiki Page",
			wantContent:   "This is a muximal wiki page.",
		},
		// Test Viewpoint 1: Specify all available options
		"success-full-options": {
			wikiID: 34,
			option: o.WithFormName("Full Options Name"),
			opts: []*backlog.FormOption{
				o.WithFormContent("Full Options Content"),
				o.WithFormMailNotify(true),
			},
			httpStatus:         http.StatusOK,
			httpBody:           testdataWikiMaximumJSON,
			wantSpath:          "wikis/34",
			expectAPICall:      true,
			wantFormName:       "Full Options Name",
			wantFormContent:    "Full Options Content",
			wantFormMailNotify: "true",
			wantID:             34,
			wantName:           "Maximum Wiki Page",
			wantContent:        "This is a muximal wiki page.",
		},
		// Test Viewpoint 2: Verify correct handling when mandatory option is in opts...
		"success-option-opts-split": {
			wikiID: 35,
			option: o.WithFormMailNotify(true), // Non-mandatory option in the required argument slot
			opts: []*backlog.FormOption{
				o.WithFormName("Split Option Name"), // Mandatory option in the variadic argument slot
			},
			httpStatus:         http.StatusOK,
			httpBody:           testdataWikiMaximumJSON,
			wantSpath:          "wikis/35",
			expectAPICall:      true,
			wantFormName:       "Split Option Name",
			wantFormMailNotify: "true",
			wantID:             34,
			wantName:           "Maximum Wiki Page",
			wantContent:        "This is a muximal wiki page.",
		},
		// --- Failure Cases (omitting success-only fields) ---
		"validation-error-required-option": {
			wikiID: 12,
			// All provided options (option and opts...) do not set mandatory fields (name/content)
			option:        o.WithFormMailNotify(true),
			opts:          []*backlog.FormOption{},
			expectAPICall: false,
			wantError:     true,
			wantErrType:   &backlog.ValidationError{},
		},
		"validation-error-invalid-wikiID": {
			wikiID:        0,
			option:        o.WithFormName("New Name"),
			opts:          []*backlog.FormOption{},
			expectAPICall: false,
			wantError:     true,
			wantErrType:   &backlog.ValidationError{},
		},
		"validation-error-invalid-option-type": {
			wikiID: 12,
			option: o.WithFormName("New Name"),
			opts: []*backlog.FormOption{backlog.ExportNewFormOption(
				backlog.ExportFormMailAddress,
				nil,
				func(p *backlog.ExportRequestParams) error {
					return nil
				},
			)},
			expectAPICall: false,
			wantError:     true,
			wantErrType:   &backlog.InvalidFormOptionError{},
		},
		"validation-error-option-set-fail": {
			wikiID: 12,
			option: o.WithFormName("New Name"),
			opts: []*backlog.FormOption{
				backlog.ExportNewFormOption(
					backlog.ExportFormMailNotify,
					nil,
					func(p *backlog.ExportRequestParams) error {
						return &backlog.ValidationError{Message: "error during option set"}
					},
				),
			},
			expectAPICall: false,
			wantError:     true,
			wantErrType:   &backlog.ValidationError{},
		},
		"client-error-network-failure": {
			wikiID:        13,
			option:        o.WithFormName("New Name"),
			opts:          []*backlog.FormOption{},
			httpError:     errors.New("network error"),
			wantSpath:     "wikis/13",
			expectAPICall: true,
			wantFormName:  "New Name",
			wantError:     true,
		},
		"api-error-invalid-json": {
			wikiID:        14,
			option:        o.WithFormName("New Name"),
			opts:          []*backlog.FormOption{},
			httpStatus:    http.StatusOK,
			httpBody:      testdataInvalidJSON,
			wantSpath:     "wikis/14",
			expectAPICall: true,
			wantFormName:  "New Name",
			wantError:     true,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			calledAPICall := false
			s := backlog.ExportNewWikiService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Patch: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
					calledAPICall = true

					// Assert the API request when expected
					if tc.expectAPICall {
						assert.Equal(t, tc.wantSpath, spath)

						// Assert form payload
						assert.Equal(t, tc.wantFormName, form.Get("name"))
						assert.Equal(t, tc.wantFormContent, form.Get("content"))
						if tc.wantFormMailNotify != "" {
							assert.Equal(t, tc.wantFormMailNotify, form.Get("mailNotify"))
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

			wiki, err := s.Update(tc.wikiID, tc.option, tc.opts...)

			if tc.expectAPICall {
				assert.True(t, calledAPICall)
			} else {
				assert.False(t, calledAPICall)
			}

			// Assert result
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
				assert.Equal(t, tc.wantName, wiki.Name)
				assert.Equal(t, tc.wantContent, wiki.Content)
			}
		})
	}
}

func TestWikiService_Delete(t *testing.T) {
	t.Parallel()

	o := backlog.ExportNewWikiOptionService()
	projectOption := backlog.ExportNewProjectOptionService() // For testing InvalidFormOptionError

	type testCase struct {
		// Input arguments
		wikiID int
		opts   []*backlog.FormOption // Variable arguments

		// HTTP mock settings
		httpStatus int
		httpBody   string
		httpError  error

		// API request assertion
		wantSpath      string
		wantMailNotify string
		expectAPICall  bool

		// Expected results (Error handling)
		wantError   bool
		wantErrType error

		// Expected results (Success case only)
		wantID int
	}

	cases := map[string]testCase{
		"success-with-option": {
			wikiID:         34,
			opts:           []*backlog.FormOption{o.WithFormMailNotify(true)},
			httpStatus:     http.StatusOK,
			httpBody:       testdataWikiMaximumJSON,
			wantSpath:      "wikis/34",
			wantMailNotify: "true",
			expectAPICall:  true,
			wantID:         34,
		},
		"success-no-option": {
			wikiID:        1,
			opts:          []*backlog.FormOption{},
			httpStatus:    http.StatusOK,
			httpBody:      testdataWikiMaximumJSON,
			wantSpath:     "wikis/1",
			expectAPICall: true,
			wantID:        34, // ID from JSON body
		},
		// --- Failure Cases (omitting success-only fields) ---
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
		"validation-error-option-set-fail": {
			wikiID: 1,
			opts: []*backlog.FormOption{
				backlog.ExportNewFormOption(
					backlog.ExportFormMailNotify,
					nil,
					func(p *backlog.ExportRequestParams) error {
						return errors.New("error during option set")
					},
				),
			},
			expectAPICall: false,
			wantError:     true,
		},
		"validation-error-invalid-option-type": {
			wikiID: 1,
			opts: []*backlog.FormOption{
				// Use ProjectOptionService to correctly trigger InvalidFormOptionError in WikiService
				projectOption.WithFormKey("Invalid Option"),
			},
			expectAPICall: false,
			wantError:     true,
			wantErrType:   &backlog.InvalidFormOptionError{},
		},
		"client-error-network-failure": {
			wikiID:        34,
			opts:          []*backlog.FormOption{},
			httpError:     errors.New("network error"),
			wantSpath:     "wikis/34",
			expectAPICall: true,
			wantError:     true,
		},
		"api-error-invalid-json": {
			wikiID:        34,
			opts:          []*backlog.FormOption{},
			httpStatus:    http.StatusOK,
			httpBody:      testdataInvalidJSON,
			wantSpath:     "wikis/34",
			expectAPICall: true,
			wantError:     true,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			calledAPICall := false
			s := backlog.ExportNewWikiService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
					calledAPICall = true

					// Assert the API request when expected
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

			wiki, err := s.Delete(tc.wikiID, tc.opts...)

			if tc.expectAPICall {
				assert.True(t, calledAPICall)
			} else {
				assert.False(t, calledAPICall)
			}

			// Assert result
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
