package wiki_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
	"github.com/nattokin/go-backlog/internal/wiki"
)

func TestWikiService_All(t *testing.T) {
	const testWiki1ID = 112
	const testWiki2ID = 115
	const testWiki1Name = "test1"
	const testWiki2Name = "test2"

	o := wiki.NewWikiOptionService(&core.OptionService{})

	cases := map[string]struct {
		projectIDOrKey string
		opts           []core.RequestOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantIDs     []int
		wantNames   []string
	}{
		"success-projectIDOrKey-id": {
			projectIDOrKey: "103",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "103", query.Get("projectIdOrKey"))

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.ListJSON))),
				}, nil
			},

			wantIDs:   []int{testWiki1ID, testWiki2ID},
			wantNames: []string{testWiki1Name, testWiki2Name},
		},

		"success-projectIDOrKey-key-with-options": {
			projectIDOrKey: "PRJ_KEY",
			opts: []core.RequestOption{
				o.WithKeyword("test"),
			},

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "PRJ_KEY", query.Get("projectIdOrKey"))
				assert.Equal(t, "test", query.Get("keyword"))

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.ListJSON))),
				}, nil
			},

			wantIDs:   []int{testWiki1ID, testWiki2ID},
			wantNames: []string{testWiki1Name, testWiki2Name},
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &core.ValidationError{},
		},

		"error-option-invalid-type": {
			projectIDOrKey: "invalid",
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		"error-option-set-failed": {
			projectIDOrKey: "PRJ",
			opts:           []core.RequestOption{mock.NewFailingSetOption(core.ParamKeyword)},
			wantErrType:    errors.New(""),
		},

		"error-client-network": {
			projectIDOrKey: "1",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "1", query.Get("projectIdOrKey"))
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},

		"error-response-invalid-json": {
			projectIDOrKey: "1",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "1", query.Get("projectIdOrKey"))

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{Get: mock.NewUnexpectedGetFn(t)}
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}

			s := wiki.NewWikiService(method, nil)

			wikis, err := s.All(context.Background(), tc.projectIDOrKey, tc.opts...)

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

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantCount   int
	}{
		"success-projectIDOrKey-id": {
			projectIDOrKey: "103",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
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

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
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
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "1",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/count", spath)
				assert.Equal(t, "1", query.Get("projectIdOrKey"))
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "1",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/count", spath)
				assert.Equal(t, "1", query.Get("projectIdOrKey"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{Get: mock.NewUnexpectedGetFn(t)}
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}

			s := wiki.NewWikiService(method, nil)

			count, err := s.Count(context.Background(), tc.projectIDOrKey)

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

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType  error
		wantWikiID   int
		wantWikiName string
	}{
		"success-wikiID-normal": {
			wikiID: 34,

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Nil(t, query)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.MaximumJSON))),
				}, nil
			},

			wantWikiID:   34,
			wantWikiName: "Maximum Wiki Page",
		},
		"error-validation-wikiID-zero": {
			wikiID:      0,
			wantErrType: &core.ValidationError{},
		},
		"error-validation-wikiID-negative": {
			wikiID:      -1,
			wantErrType: &core.ValidationError{},
		},
		"error-client-network": {
			wikiID: 1,

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1", spath)
				assert.Nil(t, query)
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			wikiID: 1,

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1", spath)
				assert.Nil(t, query)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{Get: mock.NewUnexpectedGetFn(t)}
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}

			s := wiki.NewWikiService(method, nil)

			wiki, err := s.One(context.Background(), tc.wikiID)

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
	o := wiki.NewWikiOptionService(&core.OptionService{})

	cases := map[string]struct {
		projectID int
		name      string
		content   string
		opts      []core.RequestOption

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantWiki    *model.Wiki
		wantErrType error
	}{
		"success-projectID-name-content-minimum": {
			projectID: 56,
			name:      "Minimum Wiki Page",
			content:   "This is a minimal wiki page.",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "56", form.Get("projectId"))
				assert.Equal(t, "Minimum Wiki Page", form.Get("name"))
				assert.Equal(t, "This is a minimal wiki page.", form.Get("content"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewReader(
						[]byte(fixture.Wiki.MinimumJSON),
					)),
				}, nil
			},

			wantWiki: &model.Wiki{
				ID:      34,
				Name:    "Minimum Wiki Page",
				Content: "This is a minimal wiki page.",
			},
		},
		"success-projectID-name-content-withMailNotify": {
			projectID: 56,
			name:      "Minimum Wiki Page",
			content:   "This is a minimal wiki page.",
			opts:      []core.RequestOption{o.WithMailNotify(true)},

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "56", form.Get("projectId"))
				assert.Equal(t, "Minimum Wiki Page", form.Get("name"))
				assert.Equal(t, "This is a minimal wiki page.", form.Get("content"))
				assert.Equal(t, "true", form.Get("mailNotify"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewReader(
						[]byte(fixture.Wiki.MinimumJSON),
					)),
				}, nil
			},

			wantWiki: &model.Wiki{
				ID:      34,
				Name:    "Minimum Wiki Page",
				Content: "This is a minimal wiki page.",
			},
		},
		"error-validation-projectID-zero": {
			projectID:   0,
			name:        "Test",
			content:     "test",
			wantErrType: &core.ValidationError{},
		},
		"error-validation-name-empty": {
			projectID:   1,
			name:        "",
			content:     "test",
			wantErrType: &core.ValidationError{},
		},
		"error-validation-content-empty": {
			projectID:   1,
			name:        "Test",
			content:     "",
			wantErrType: &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectID:   1,
			name:        "Test",
			content:     "content",
			opts:        []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"error-option-set-faild": {
			projectID:   1,
			name:        "Test",
			content:     "content",
			opts:        []core.RequestOption{mock.NewFailingSetOption(core.ParamMailNotify)},
			wantErrType: errors.New(""),
		},
		"error-client-network": {
			projectID: 1,
			name:      "Test",
			content:   "content",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
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

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "1", form.Get("projectId"))
				assert.Equal(t, "Test", form.Get("name"))
				assert.Equal(t, "content", form.Get("content"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewReader(
						[]byte(fixture.InvalidJSON),
					)),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{Post: mock.NewUnexpectedPostFn(t)}
			if tc.mockPostFn != nil {
				method.Post = tc.mockPostFn
			}

			s := wiki.NewWikiService(method, nil)

			wiki, err := s.Create(context.Background(), tc.projectID, tc.name, tc.content, tc.opts...)

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
	o := wiki.NewWikiOptionService(&core.OptionService{})

	cases := map[string]struct {
		wikiID int
		option core.RequestOption
		opts   []core.RequestOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantWiki    *model.Wiki
	}{
		"success-wikiID-name-only": {
			wikiID: 34,
			option: o.WithName("New Page Name"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Equal(t, "New Page Name", form.Get("name"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.MaximumJSON))),
				}, nil
			},

			wantWiki: &model.Wiki{
				ID:      34,
				Name:    "Maximum Wiki Page",
				Content: "This is a muximal wiki page.",
			},
		},
		"success-wikiID-content-only": {
			wikiID: 34,
			option: o.WithContent("Full Options Content"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Equal(t, "Full Options Content", form.Get("content"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.MaximumJSON))),
				}, nil
			},

			wantWiki: &model.Wiki{
				ID:      34,
				Name:    "Maximum Wiki Page",
				Content: "This is a muximal wiki page.",
			},
		},
		"success-wikiID-mailNotify-name": {
			wikiID: 34,
			option: o.WithMailNotify(true),
			opts: []core.RequestOption{
				o.WithName("Full Options Name"),
			},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Equal(t, "Full Options Name", form.Get("name"))
				assert.Equal(t, "true", form.Get("mailNotify"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.MaximumJSON))),
				}, nil
			},

			wantWiki: &model.Wiki{
				ID:      34,
				Name:    "Maximum Wiki Page",
				Content: "This is a muximal wiki page.",
			},
		},
		"success-wikiID-full-options": {
			wikiID: 34,
			option: o.WithName("Full Options Name"),
			opts: []core.RequestOption{
				o.WithContent("Full Options Content"),
				o.WithMailNotify(true),
			},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Equal(t, "Full Options Name", form.Get("name"))
				assert.Equal(t, "Full Options Content", form.Get("content"))
				assert.Equal(t, "true", form.Get("mailNotify"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.MaximumJSON))),
				}, nil
			},

			wantWiki: &model.Wiki{
				ID:      34,
				Name:    "Maximum Wiki Page",
				Content: "This is a muximal wiki page.",
			},
		},
		"error-validation-required-option": {
			wikiID:      12,
			option:      o.WithMailNotify(true),
			wantErrType: &core.ValidationError{},
		},
		"error-validation-wikiID-zero": {
			wikiID:      0,
			option:      o.WithName("New Name"),
			wantErrType: &core.ValidationError{},
		},
		"error-option-invalid-type": {
			wikiID: 12,
			option: mock.NewInvalidTypeOption(),
			opts: []core.RequestOption{
				o.WithName("New Name"),
			},
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"error-option-set-faild": {
			wikiID:      12,
			option:      mock.NewFailingSetOption(core.ParamName),
			wantErrType: errors.New(""),
		},
		"error-client-network": {
			wikiID: 13,
			option: o.WithName("New Name"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/13", spath)
				assert.Equal(t, "New Name", form.Get("name"))
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			wikiID: 14,
			option: o.WithName("New Name"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/14", spath)
				assert.Equal(t, "New Name", form.Get("name"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{Patch: mock.NewUnexpectedPatchFn(t)}
			if tc.mockPatchFn != nil {
				method.Patch = tc.mockPatchFn
			}

			s := wiki.NewWikiService(method, nil)

			wiki, err := s.Update(context.Background(), tc.wikiID, tc.option, tc.opts...)

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
	o := wiki.NewWikiOptionService(&core.OptionService{})

	cases := map[string]struct {
		wikiID int
		opts   []core.RequestOption

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantWikiID  int
		wantErrType error
	}{
		"success-wikiID-withMailNotify": {
			wikiID: 34,
			opts:   []core.RequestOption{o.WithMailNotify(true)},

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				assert.Equal(t, "true", form.Get("mailNotify"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.MaximumJSON))),
				}, nil
			},

			wantWikiID: 34,
		},
		"success-wikiID-no-option": {
			wikiID: 1,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Wiki.MaximumJSON))),
				}, nil
			},

			wantWikiID: 34,
		},
		"error-validation-wikiID-zero": {
			wikiID:      0,
			wantErrType: &core.ValidationError{},
		},
		"error-validation-wikiID-negative": {
			wikiID:      -1,
			wantErrType: &core.ValidationError{},
		},
		"error-option-set-faild": {
			wikiID:      1,
			opts:        []core.RequestOption{mock.NewFailingSetOption(core.ParamMailNotify)},
			wantErrType: errors.New(""),
		},
		"error-option-invalid-type": {
			wikiID:      1,
			opts:        []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"error-client-network": {
			wikiID: 34,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			wikiID: 34,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{Delete: mock.NewUnexpectedDeleteFn(t)}
			if tc.mockDeleteFn != nil {
				method.Delete = tc.mockDeleteFn
			}

			s := wiki.NewWikiService(method, nil)

			wiki, err := s.Delete(context.Background(), tc.wikiID, tc.opts...)

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

// TestWikiService_contextPropagation verifies that the context passed to each
// WikiService method is correctly relayed to the underlying method call.
// A sentinel value is embedded in the context and its pointer identity is
// asserted inside the mock to catch any ctx substitution (e.g. context.Background()).
func TestWikiService_contextPropagation(t *testing.T) {
	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	o := wiki.NewWikiOptionService(&core.OptionService{})

	cases := []struct {
		name string
		call func(t *testing.T, m *core.Method)
	}{
		{"All", func(t *testing.T, m *core.Method) {
			m.Get = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s := wiki.NewWikiService(m, nil)
			s.All(ctx, "TEST") //nolint:errcheck
		}},
		{"Count", func(t *testing.T, m *core.Method) {
			m.Get = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s := wiki.NewWikiService(m, nil)
			s.Count(ctx, "TEST") //nolint:errcheck
		}},
		{"One", func(t *testing.T, m *core.Method) {
			m.Get = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s := wiki.NewWikiService(m, nil)
			s.One(ctx, 1) //nolint:errcheck
		}},
		{"Create", func(t *testing.T, m *core.Method) {
			m.Post = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s := wiki.NewWikiService(m, nil)
			s.Create(ctx, 1, "name", "content") //nolint:errcheck
		}},
		{"Update", func(t *testing.T, m *core.Method) {
			m.Patch = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s := wiki.NewWikiService(m, nil)
			s.Update(ctx, 1, o.WithName("n")) //nolint:errcheck
		}},
		{"Delete", func(t *testing.T, m *core.Method) {
			m.Delete = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s := wiki.NewWikiService(m, nil)
			s.Delete(ctx, 1) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
