package project_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/project"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestVersionService_All(t *testing.T) {
	option := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		opts           []core.RequestOption
		wantErrType    error
		wantLen        int
		mockGetFn      func(context.Context, string, url.Values) (*http.Response, error)
	}{
		"success": {
			projectIDOrKey: "TEST",
			wantLen:        2,
			opts: []core.RequestOption{
				option.WithArchived(true),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/versions", spath)
				assert.Equal(t, "true", query.Get("archived"))
				return mock.NewJSONResponse(fixture.Version.ListJSON), nil
			},
		},
		"error-project-empty": {
			projectIDOrKey: "",
			wantErrType:    &core.ValidationError{},
			mockGetFn:      mock.NewUnexpectedGetFn(t),
		},
		"error-option-invalid-type": {
			projectIDOrKey: "TEST",
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		"error-client": {
			projectIDOrKey: "TEST",
			wantErrType:    errors.New(""),
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
			projectIDOrKey: "TEST",
			wantErrType:    &json.SyntaxError{},
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := mock.NewMethod(t)
			m.Get = tc.mockGetFn
			s := project.NewVersionService(m)
			got, err := s.List(context.Background(), tc.projectIDOrKey, tc.opts...)
			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, got, tc.wantLen)
		})
	}
}

func TestVersionService_Add(t *testing.T) {
	o := &core.OptionService{}
	date := "2025-01-01"

	cases := map[string]struct {
		projectIDOrKey string
		name           string
		opts           []core.RequestOption
		wantErrType    error
		wantID         int
		mockPostFn     func(context.Context, string, url.Values) (*http.Response, error)
	}{
		"success": {
			projectIDOrKey: "TEST",
			name:           "v1",
			opts:           []core.RequestOption{o.WithDescription("desc"), o.WithStartDate(date), o.WithReleaseDueDate(date)},
			wantID:         fixture.Version.Single.ID,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/versions", spath)
				assert.Equal(t, "v1", form.Get("name"))
				return mock.NewJSONResponse(fixture.Version.SingleJSON), nil
			},
		},
		"error-name-empty": {
			projectIDOrKey: "TEST",
			name:           "",
			opts:           []core.RequestOption{o.WithDescription("desc")},
			wantErrType:    &core.ValidationError{},
		},
		"error-project-empty": {
			projectIDOrKey: "",
			name:           "v1",
			opts:           []core.RequestOption{o.WithDescription("desc")},
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: "TEST",
			name:           "v1",
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-set-failed": {
			projectIDOrKey: "TEST",
			name:           "v1",
			opts:           []core.RequestOption{mock.NewFailingSetOption(core.ParamDescription)},
			wantErrType:    errors.New(""),
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			name:           "v1",
			wantErrType:    errors.New(""),
			mockPostFn:     func(context.Context, string, url.Values) (*http.Response, error) { return nil, errors.New("network") },
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			name:           "v1",
			wantErrType:    &json.SyntaxError{},
			mockPostFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := mock.NewMethod(t)
			if tc.mockPostFn != nil {
				m.Post = tc.mockPostFn
			}
			s := project.NewVersionService(m)
			got, err := s.Add(context.Background(), tc.projectIDOrKey, tc.name, tc.opts...)
			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestVersionService_Update(t *testing.T) {
	o := &core.OptionService{}
	date := "2025-01-01"

	cases := map[string]struct {
		projectIDOrKey string
		versionID      int
		option         core.RequestOption
		opts           []core.RequestOption
		wantErrType    error
		wantID         int
		mockPatchFn    func(context.Context, string, url.Values) (*http.Response, error)
	}{
		"success": {
			projectIDOrKey: "TEST",
			versionID:      1,
			option:         o.WithName("name"),
			opts:           []core.RequestOption{o.WithReleaseDueDate(date)},
			wantID:         fixture.Version.Single.ID,
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/versions/1", spath)
				return mock.NewJSONResponse(fixture.Version.SingleJSON), nil
			},
		},
		"error-versionID-negative": {
			projectIDOrKey: "TEST",
			versionID:      -1,
			option:         o.WithName("name"),
			wantErrType:    &core.ValidationError{},
			mockPatchFn:    mock.NewUnexpectedPatchFn(t),
		},
		"error-option-invalid-type": {
			projectIDOrKey: "TEST",
			versionID:      1,
			option:         mock.NewInvalidTypeOption(),
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-set-failed": {
			projectIDOrKey: "TEST",
			versionID:      1,
			option:         mock.NewFailingSetOption(core.ParamArchived),
			wantErrType:    errors.New(""),
		},
		"error-project-empty": {
			projectIDOrKey: "",
			versionID:      1,
			option:         o.WithName("name"),
			wantErrType:    &core.ValidationError{},
			mockPatchFn:    mock.NewUnexpectedPatchFn(t),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			versionID:      1,
			option:         o.WithName("name"),
			wantErrType:    &json.SyntaxError{},
			mockPatchFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
		"error-versionID-zero": {
			projectIDOrKey: "TEST",
			versionID:      0,
			option:         o.WithName("name"),
			wantErrType:    &core.ValidationError{},
			mockPatchFn:    mock.NewUnexpectedPatchFn(t),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := mock.NewMethod(t)
			if tc.mockPatchFn != nil {
				m.Patch = tc.mockPatchFn
			}
			s := project.NewVersionService(m)
			got, err := s.Update(context.Background(), tc.projectIDOrKey, tc.versionID, tc.option, tc.opts...)
			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestVersionService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		versionID      int
		wantErrType    error
		wantID         int
		mockDeleteFn   func(context.Context, string, url.Values) (*http.Response, error)
	}{
		"success": {
			projectIDOrKey: "TEST", versionID: 1, wantID: fixture.Version.Single.ID,
			mockDeleteFn: func(ctx context.Context, spath string, _ url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/versions/1", spath)
				return mock.NewJSONResponse(fixture.Version.SingleJSON), nil
			},
		},
		"error-versionID-zero": {
			projectIDOrKey: "TEST", versionID: 0, wantErrType: &core.ValidationError{}, mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},
		"error-versionID-negative": {
			projectIDOrKey: "TEST",
			versionID:      -1,
			wantErrType:    &core.ValidationError{},
			mockDeleteFn:   mock.NewUnexpectedDeleteFn(t),
		},
		"error-project-empty": {
			projectIDOrKey: "", versionID: 1, wantErrType: &core.ValidationError{}, mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},
		"error-client-network": {
			projectIDOrKey: "TEST", versionID: 1, wantErrType: errors.New(""),
			mockDeleteFn: func(context.Context, string, url.Values) (*http.Response, error) { return nil, errors.New("network") },
		},
		"error-invalid-json": {
			projectIDOrKey: "TEST", versionID: 1, wantErrType: &json.SyntaxError{},
			mockDeleteFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := mock.NewMethod(t)
			m.Delete = tc.mockDeleteFn
			s := project.NewVersionService(m)
			got, err := s.Delete(context.Background(), tc.projectIDOrKey, tc.versionID)
			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}
