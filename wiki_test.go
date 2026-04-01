package backlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
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
		projectIDOrKey string
		options        []*QueryOption

		mockGetFn func(spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantIDs     []int
		wantNames   []string
	}{
		"success-project-id": {
			projectIDOrKey: "103",

			mockGetFn: func(spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "103", query.Get("projectIdOrKey"))

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiListJSON))),
				}, nil
			},

			wantIDs:   []int{testWiki1ID, testWiki2ID},
			wantNames: []string{testWiki1Name, testWiki2Name},
		},
		"success-with-options": {
			projectIDOrKey: "PRJ_KEY",
			options: []*QueryOption{
				o.WithQueryKeyword("test"),
			},

			mockGetFn: func(spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "PRJ_KEY", query.Get("projectIdOrKey"))
				assert.Equal(t, "test", query.Get("keyword"))

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiListJSON))),
				}, nil
			},

			wantIDs:   []int{testWiki1ID, testWiki2ID},
			wantNames: []string{testWiki1Name, testWiki2Name},
		},
		"validation-error-key-empty": {
			projectIDOrKey: "",
			wantErrType:    &ValidationError{},
		},
		"validation-error-invalid-option-type": {
			projectIDOrKey: "PRJ",
			options: []*QueryOption{{
				t:         queryCount,
				checkFunc: nil,
				setFunc: func(p url.Values) error {
					return nil
				},
			}},
			wantErrType: &InvalidOptionError[queryType]{},
		},
		"validation-error-option-set-fail": {
			projectIDOrKey: "PRJ",
			options:        []*QueryOption{newQueryOptionWithSetError(queryKeyword)},
			wantErrType:    errors.New(""),
		},
		"client-error-network-failure": {
			projectIDOrKey: "1",

			mockGetFn: func(spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "1", query.Get("projectIdOrKey"))
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},
		"api-error-invalid-json": {
			projectIDOrKey: "1",

			mockGetFn: func(spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "1", query.Get("projectIdOrKey"))

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newWikiService()

			// default: unexpected API call
			s.method.Get = newUnexpectedGetFn(t)

			if tc.mockGetFn != nil {
				s.method.Get = tc.mockGetFn
			}

			wikis, err := s.All(tc.projectIDOrKey, tc.options...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, wikis)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, wikis)

			assert.Len(t, wikis, len(tc.wantIDs))
			for i := range wikis {
				assert.Equal(t, tc.wantIDs[i], wikis[i].ID)
				assert.Equal(t, tc.wantNames[i], wikis[i].Name)
			}
		})
	}
}

func TestWikiService_Count(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockGetFn func(spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantCount   int
	}{
		"success-project-id": {
			projectIDOrKey: "103",

			mockGetFn: func(spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/count", spath)
				assert.Equal(t, "103", query.Get("projectIdOrKey"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"count": 34}`))),
				}, nil
			},

			wantCount: 34,
		},
		"success-project-key": {
			projectIDOrKey: "PRJ_KEY",

			mockGetFn: func(spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/count", spath)
				assert.Equal(t, "PRJ_KEY", query.Get("projectIdOrKey"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"count": 10}`))),
				}, nil
			},

			wantCount: 10,
		},
		"validation-error-key-empty": {
			projectIDOrKey: "",
			wantErrType:    &ValidationError{},
		},
		"client-error-network-failure": {
			projectIDOrKey: "1",

			mockGetFn: func(spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/count", spath)
				assert.Equal(t, "1", query.Get("projectIdOrKey"))
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},
		"api-error-invalid-json": {
			projectIDOrKey: "1",

			mockGetFn: func(spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/count", spath)
				assert.Equal(t, "1", query.Get("projectIdOrKey"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newWikiService()
			s.method.Get = newUnexpectedGetFn(t)
			if tc.mockGetFn != nil {
				s.method.Get = tc.mockGetFn
			}

			count, err := s.Count(tc.projectIDOrKey)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Equal(t, 0, count)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantCount, count)
		})
	}
}

func TestWikiService_One(t *testing.T) {
	cases := map[string]struct {
		wikiID int

		mockGetFn func(spath string, query url.Values) (*http.Response, error)

		wantErrType  error
		wantWikiID   int
		wantWikiName string
	}{
		"success-normal": {
			wikiID: 34,

			mockGetFn: func(spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Nil(t, query)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiMaximumJSON))),
				}, nil
			},

			wantWikiID:   34,
			wantWikiName: "Maximum Wiki Page",
		},
		"validation-error-id-zero": {
			wikiID:      0,
			wantErrType: &ValidationError{},
		},
		"validation-error-id-negative": {
			wikiID:      -1,
			wantErrType: &ValidationError{},
		},
		"client-error-network-failure": {
			wikiID: 1,

			mockGetFn: func(spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1", spath)
				assert.Nil(t, query)
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},
		"api-error-invalid-json": {
			wikiID: 1,

			mockGetFn: func(spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1", spath)
				assert.Nil(t, query)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newWikiService()
			s.method.Get = newUnexpectedGetFn(t)
			if tc.mockGetFn != nil {
				s.method.Get = tc.mockGetFn
			}

			wiki, err := s.One(tc.wikiID)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, wiki)
				assert.IsType(t, tc.wantErrType, err)
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

	cases := map[string]struct {
		projectID int
		name      string
		content   string
		opts      []*FormOption

		mockPostFn func(spath string, form *FormParams) (*http.Response, error)

		wantWiki    *Wiki
		wantErrType error
	}{
		"success-minimum": {
			projectID: 56,
			name:      "Minimum Wiki Page",
			content:   "This is a minimal wiki page.",

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "56", form.Get("projectId"))
				assert.Equal(t, "Minimum Wiki Page", form.Get("name"))
				assert.Equal(t, "This is a minimal wiki page.", form.Get("content"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewReader(
						[]byte(testdataWikiMinimumJSON),
					)),
				}, nil
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

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "56", form.Get("projectId"))
				assert.Equal(t, "Minimum Wiki Page", form.Get("name"))
				assert.Equal(t, "This is a minimal wiki page.", form.Get("content"))
				assert.Equal(t, "true", form.Get("mailNotify"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewReader(
						[]byte(testdataWikiMinimumJSON),
					)),
				}, nil
			},

			wantWiki: &Wiki{
				ID:      34,
				Name:    "Minimum Wiki Page",
				Content: "This is a minimal wiki page.",
			},
		},
		"validation-error-projectID-zero": {
			projectID:   0,
			name:        "Test",
			content:     "test",
			wantErrType: &ValidationError{},
		},
		"validation-error-name-empty": {
			projectID:   1,
			name:        "",
			content:     "test",
			wantErrType: &ValidationError{},
		},
		"validation-error-content-empty": {
			projectID:   1,
			name:        "Test",
			content:     "",
			wantErrType: &ValidationError{},
		},
		"validation-error-invalid-option-type": {
			projectID: 1,
			name:      "Test",
			content:   "content",
			opts: []*FormOption{
				{
					formMailAddress,
					nil,
					func(p *FormParams) error { return nil },
				},
			},
			wantErrType: &InvalidOptionError[formType]{},
		},
		"validation-error-option-set-fail": {
			projectID:   1,
			name:        "Test",
			content:     "content",
			opts:        []*FormOption{newFormOptionWithSetError(formMailNotify)},
			wantErrType: errors.New(""),
		},
		"client-error-network-failure": {
			projectID: 1,
			name:      "Test",
			content:   "content",

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "1", form.Get("projectId"))
				assert.Equal(t, "Test", form.Get("name"))
				assert.Equal(t, "content", form.Get("content"))
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},
		"api-error-invalid-json": {
			projectID: 1,
			name:      "Test",
			content:   "content",

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "1", form.Get("projectId"))
				assert.Equal(t, "Test", form.Get("name"))
				assert.Equal(t, "content", form.Get("content"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewReader(
						[]byte(testdataInvalidJSON),
					)),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newWikiService()
			s.method.Post = newUnexpectedPostFn(t)

			if tc.mockPostFn != nil {
				s.method.Post = tc.mockPostFn
			}

			wiki, err := s.Create(tc.projectID, tc.name, tc.content, tc.opts...)

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

	cases := map[string]struct {
		wikiID int
		option *FormOption
		opts   []*FormOption

		mockPatchFn func(spath string, form *FormParams) (*http.Response, error)

		wantErrType error
		wantWiki    *Wiki
	}{
		"success-name-only": {
			wikiID: 34,
			option: o.WithFormName("New Page Name"),

			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Equal(t, "New Page Name", form.Get("name"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiMaximumJSON))),
				}, nil
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

			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Equal(t, "Full Options Name", form.Get("name"))
				assert.Equal(t, "Full Options Content", form.Get("content"))
				assert.Equal(t, "true", form.Get("mailNotify"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiMaximumJSON))),
				}, nil
			},

			wantWiki: &Wiki{
				ID:      34,
				Name:    "Maximum Wiki Page",
				Content: "This is a muximal wiki page.",
			},
		},
		"validation-error-required-option": {
			wikiID:      12,
			option:      o.WithFormMailNotify(true),
			wantErrType: &ValidationError{},
		},
		"validation-error-invalid-wikiID": {
			wikiID:      0,
			option:      o.WithFormName("New Name"),
			wantErrType: &ValidationError{},
		},
		"validation-error-invalid-option-type": {
			wikiID: 12,
			option: &FormOption{
				t:         formRoleType,
				checkFunc: nil,
				setFunc: func(p *FormParams) error {
					return nil
				},
			},
			wantErrType: &InvalidOptionError[formType]{},
		},
		"validation-error-option-set-fail": {
			wikiID:      12,
			option:      newFormOptionWithSetError(formName),
			wantErrType: errors.New(""),
		},
		"client-error-network-failure": {
			wikiID: 13,
			option: o.WithFormName("New Name"),

			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "wikis/13", spath)
				assert.Equal(t, "New Name", form.Get("name"))
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},
		"api-error-invalid-json": {
			wikiID: 14,
			option: o.WithFormName("New Name"),

			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "wikis/14", spath)
				assert.Equal(t, "New Name", form.Get("name"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newWikiService()
			s.method.Patch = newUnexpectedPatchFn(t)

			if tc.mockPatchFn != nil {
				s.method.Patch = tc.mockPatchFn
			}

			wiki, err := s.Update(tc.wikiID, tc.option, tc.opts...)

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

	cases := map[string]struct {
		wikiID int
		opts   []*FormOption

		mockDeleteFn func(spath string, form *FormParams) (*http.Response, error)

		wantWikiID  int
		wantErrType error
	}{
		"success-with-option": {
			wikiID: 34,
			opts:   []*FormOption{o.WithFormMailNotify(true)},

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Equal(t, "true", form.Get("mailNotify"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiMaximumJSON))),
				}, nil
			},

			wantWikiID: 34,
		},
		"success-no-option": {
			wikiID: 1,

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "wikis/1", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiMaximumJSON))),
				}, nil
			},

			wantWikiID: 34,
		},
		"validation-error-id-zero": {
			wikiID:      0,
			wantErrType: &ValidationError{},
		},
		"validation-error-id-negative": {
			wikiID:      -1,
			wantErrType: &ValidationError{},
		},
		"validation-error-option-set-fail": {
			wikiID:      1,
			opts:        []*FormOption{newFormOptionWithSetError(formMailNotify)},
			wantErrType: errors.New(""),
		},
		"validation-error-invalid-option-type": {
			wikiID: 1,
			opts: []*FormOption{
				projectOption.WithFormKey("Invalid Option"),
			},
			wantErrType: &InvalidOptionError[formType]{},
		},
		"client-error-network-failure": {
			wikiID: 34,

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},
		"api-error-invalid-json": {
			wikiID: 34,

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newWikiService()
			s.method.Delete = newUnexpectedDeleteFn(t)

			if tc.mockDeleteFn != nil {
				s.method.Delete = tc.mockDeleteFn
			}

			wiki, err := s.Delete(tc.wikiID, tc.opts...)

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
