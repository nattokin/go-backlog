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
		opts           []RequestOption

		mockGetFn func(spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantIDs     []int
		wantNames   []string
	}{
		"success-projectIDOrKey-id": {
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

		"success-projectIDOrKey-key-with-options": {
			projectIDOrKey: "PRJ_KEY",
			opts: []RequestOption{
				o.WithKeyword("test"),
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

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &ValidationError{},
		},

		"error-option-invalid-type": {
			projectIDOrKey: "invalid",
			opts:           []RequestOption{newInvalidTypeOption()},
			wantErrType:    &InvalidOptionKeyError{},
		},

		"error-option-set-failed": {
			projectIDOrKey: "PRJ",
			opts:           []RequestOption{newFailingSetOption(paramKeyword)},
			wantErrType:    errors.New(""),
		},

		"error-client-network": {
			projectIDOrKey: "1",

			mockGetFn: func(spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "1", query.Get("projectIdOrKey"))
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},

		"error-response-invalid-json": {
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
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newWikiService()

			// default: unexpected API call
			s.method.Get = newUnexpectedGetFn(t)

			if tc.mockGetFn != nil {
				s.method.Get = tc.mockGetFn
			}

			wikis, err := s.All(tc.projectIDOrKey, tc.opts...)

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
		"success-projectIDOrKey-id": {
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
		"success-projectIDOrKey-key": {
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
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "1",

			mockGetFn: func(spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/count", spath)
				assert.Equal(t, "1", query.Get("projectIdOrKey"))
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
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
		"success-wikiID-normal": {
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
		"error-validation-wikiID-zero": {
			wikiID:      0,
			wantErrType: &ValidationError{},
		},
		"error-validation-wikiID-negative": {
			wikiID:      -1,
			wantErrType: &ValidationError{},
		},
		"error-client-network": {
			wikiID: 1,

			mockGetFn: func(spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1", spath)
				assert.Nil(t, query)
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
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
		opts      []RequestOption

		mockPostFn func(spath string, form url.Values) (*http.Response, error)

		wantWiki    *Wiki
		wantErrType error
	}{
		"success-projectID-name-content-minimum": {
			projectID: 56,
			name:      "Minimum Wiki Page",
			content:   "This is a minimal wiki page.",

			mockPostFn: func(spath string, form url.Values) (*http.Response, error) {
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
		"success-projectID-name-content-withMailNotify": {
			projectID: 56,
			name:      "Minimum Wiki Page",
			content:   "This is a minimal wiki page.",
			opts:      []RequestOption{o.WithMailNotify(true)},

			mockPostFn: func(spath string, form url.Values) (*http.Response, error) {
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
		"error-validation-projectID-zero": {
			projectID:   0,
			name:        "Test",
			content:     "test",
			wantErrType: &ValidationError{},
		},
		"error-validation-name-empty": {
			projectID:   1,
			name:        "",
			content:     "test",
			wantErrType: &ValidationError{},
		},
		"error-validation-content-empty": {
			projectID:   1,
			name:        "Test",
			content:     "",
			wantErrType: &ValidationError{},
		},
		"error-option-invalid-type": {
			projectID:   1,
			name:        "Test",
			content:     "content",
			opts:        []RequestOption{newInvalidTypeOption()},
			wantErrType: &InvalidOptionKeyError{},
		},
		"error-option-set-faild": {
			projectID:   1,
			name:        "Test",
			content:     "content",
			opts:        []RequestOption{newFailingSetOption(paramMailNotify)},
			wantErrType: errors.New(""),
		},
		"error-client-network": {
			projectID: 1,
			name:      "Test",
			content:   "content",

			mockPostFn: func(spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "1", form.Get("projectId"))
				assert.Equal(t, "Test", form.Get("name"))
				assert.Equal(t, "content", form.Get("content"))
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectID: 1,
			name:      "Test",
			content:   "content",

			mockPostFn: func(spath string, form url.Values) (*http.Response, error) {
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
		option RequestOption
		opts   []RequestOption

		mockPatchFn func(spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantWiki    *Wiki
	}{
		"success-wikiID-name-only": {
			wikiID: 34,
			option: o.WithName("New Page Name"),

			mockPatchFn: func(spath string, form url.Values) (*http.Response, error) {
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
		"success-wikiID-content-only": {
			wikiID: 34,
			option: o.WithContent("Full Options Content"),

			mockPatchFn: func(spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Equal(t, "Full Options Content", form.Get("content"))
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
		"success-wikiID-mailNotify-name": {
			wikiID: 34,
			option: o.WithMailNotify(true),
			opts: []RequestOption{
				o.WithName("Full Options Name"),
			},

			mockPatchFn: func(spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Equal(t, "Full Options Name", form.Get("name"))
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
		"success-wikiID-full-options": {
			wikiID: 34,
			option: o.WithName("Full Options Name"),
			opts: []RequestOption{
				o.WithContent("Full Options Content"),
				o.WithMailNotify(true),
			},

			mockPatchFn: func(spath string, form url.Values) (*http.Response, error) {
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
		"error-validation-required-option": {
			wikiID:      12,
			option:      o.WithMailNotify(true),
			wantErrType: &ValidationError{},
		},
		"error-validation-wikiID-zero": {
			wikiID:      0,
			option:      o.WithName("New Name"),
			wantErrType: &ValidationError{},
		},
		"error-option-invalid-type": {
			wikiID: 12,
			option: newInvalidTypeOption(),
			opts: []RequestOption{
				o.WithName("New Name"),
			},
			wantErrType: &InvalidOptionKeyError{},
		},
		"error-option-set-faild": {
			wikiID:      12,
			option:      newFailingSetOption(paramName),
			wantErrType: errors.New(""),
		},
		"error-client-network": {
			wikiID: 13,
			option: o.WithName("New Name"),

			mockPatchFn: func(spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/13", spath)
				assert.Equal(t, "New Name", form.Get("name"))
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			wikiID: 14,
			option: o.WithName("New Name"),

			mockPatchFn: func(spath string, form url.Values) (*http.Response, error) {
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

	cases := map[string]struct {
		wikiID int
		opts   []RequestOption

		mockDeleteFn func(spath string, form url.Values) (*http.Response, error)

		wantWikiID  int
		wantErrType error
	}{
		"success-wikiID-withMailNotify": {
			wikiID: 34,
			opts:   []RequestOption{o.WithMailNotify(true)},

			mockDeleteFn: func(spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Equal(t, "true", form.Get("mailNotify"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiMaximumJSON))),
				}, nil
			},

			wantWikiID: 34,
		},
		"success-wikiID-no-option": {
			wikiID: 1,

			mockDeleteFn: func(spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiMaximumJSON))),
				}, nil
			},

			wantWikiID: 34,
		},
		"error-validation-wikiID-zero": {
			wikiID:      0,
			wantErrType: &ValidationError{},
		},
		"error-validation-wikiID-negative": {
			wikiID:      -1,
			wantErrType: &ValidationError{},
		},
		"error-option-set-faild": {
			wikiID:      1,
			opts:        []RequestOption{newFailingSetOption(paramMailNotify)},
			wantErrType: errors.New(""),
		},
		"error-option-invalid-type": {
			wikiID:      1,
			opts:        []RequestOption{newInvalidTypeOption()},
			wantErrType: &InvalidOptionKeyError{},
		},
		"error-client-network": {
			wikiID: 34,

			mockDeleteFn: func(spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			wikiID: 34,

			mockDeleteFn: func(spath string, form url.Values) (*http.Response, error) {
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
