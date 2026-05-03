package backlog_test

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
	"github.com/nattokin/go-backlog/internal/project"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestIssueTypeService_All(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantLen     int
		wantErrType error
	}{
		"success-projectIDOrKey-key": {
			projectIDOrKey: "TEST",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/issueTypes", spath)
				assert.Nil(t, query)
				return mock.NewJSONResponse(fixture.IssueType.ListJSON), nil
			},

			wantLen:     2,
			wantErrType: nil,
		},
		"success-projectIDOrKey-id": {
			projectIDOrKey: "6",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/6/issueTypes", spath)
				return mock.NewJSONResponse(fixture.IssueType.ListJSON), nil
			},

			wantLen:     2,
			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}
			s := project.NewIssueTypeService(method)

			issueTypes, err := s.All(context.Background(), tc.projectIDOrKey)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, issueTypes)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, issueTypes)
			assert.Len(t, issueTypes, tc.wantLen)
		})
	}
}

func TestIssueTypeService_Create(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		name           string
		color          string

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			name:           "Bug",
			color:          "#e30000",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/issueTypes", spath)
				assert.Equal(t, "Bug", form.Get("name"))
				assert.Equal(t, "#e30000", form.Get("color"))
				return mock.NewJSONResponse(fixture.IssueType.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			name:           "Bug",
			color:          "#e30000",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-name-empty": {
			projectIDOrKey: "TEST",
			name:           "",
			color:          "#e30000",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-color-empty": {
			projectIDOrKey: "TEST",
			name:           "Bug",
			color:          "",
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			name:           "Bug",
			color:          "#e30000",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			name:           "Bug",
			color:          "#e30000",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockPostFn != nil {
				method.Post = tc.mockPostFn
			}
			s := project.NewIssueTypeService(method)

			issueType, err := s.Create(context.Background(), tc.projectIDOrKey, tc.name, tc.color)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, issueType)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, issueType)
			assert.Equal(t, 1, issueType.ID)
			assert.Equal(t, "Bug", issueType.Name)
			assert.Equal(t, "#e30000", issueType.Color)
		})
	}
}

func TestIssueTypeService_Update(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		issueTypeID    int
		option         core.RequestOption
		opts           []core.RequestOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			issueTypeID:    1,
			option:         o.WithName("Bug Updated"),
			opts:           []core.RequestOption{o.WithColor("#990000")},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/issueTypes/1", spath)
				assert.Equal(t, "Bug Updated", form.Get("name"))
				assert.Equal(t, "#990000", form.Get("color"))
				return mock.NewJSONResponse(fixture.IssueType.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"success-option-only": {
			projectIDOrKey: "TEST",
			issueTypeID:    1,
			option:         o.WithName("Bug Updated"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/issueTypes/1", spath)
				assert.Equal(t, "Bug Updated", form.Get("name"))
				return mock.NewJSONResponse(fixture.IssueType.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			issueTypeID:    1,
			option:         o.WithName("Bug Updated"),
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-issueTypeID-zero": {
			projectIDOrKey: "TEST",
			issueTypeID:    0,
			option:         o.WithName("Bug Updated"),
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: "TEST",
			issueTypeID:    1,
			option:         mock.NewInvalidTypeOption(),
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			issueTypeID:    1,
			option:         o.WithName("Bug Updated"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			issueTypeID:    1,
			option:         o.WithName("Bug Updated"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockPatchFn != nil {
				method.Patch = tc.mockPatchFn
			}
			s := project.NewIssueTypeService(method)

			issueType, err := s.Update(context.Background(), tc.projectIDOrKey, tc.issueTypeID, tc.option, tc.opts...)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, issueType)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, issueType)
			assert.Equal(t, 1, issueType.ID)
		})
	}
}

func TestIssueTypeService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey        string
		issueTypeID           int
		substituteIssueTypeID int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey:        "TEST",
			issueTypeID:           1,
			substituteIssueTypeID: 2,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/issueTypes/1", spath)
				assert.Equal(t, "2", form.Get("substituteIssueTypeId"))
				return mock.NewJSONResponse(fixture.IssueType.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:        "",
			issueTypeID:           1,
			substituteIssueTypeID: 2,
			wantErrType:           &core.ValidationError{},
		},
		"error-validation-issueTypeID-zero": {
			projectIDOrKey:        "TEST",
			issueTypeID:           0,
			substituteIssueTypeID: 2,
			wantErrType:           &core.ValidationError{},
		},
		"error-validation-substituteIssueTypeID-zero": {
			projectIDOrKey:        "TEST",
			issueTypeID:           1,
			substituteIssueTypeID: 0,
			wantErrType:           &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey:        "TEST",
			issueTypeID:           1,
			substituteIssueTypeID: 2,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey:        "TEST",
			issueTypeID:           1,
			substituteIssueTypeID: 2,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockDeleteFn != nil {
				method.Delete = tc.mockDeleteFn
			}
			s := project.NewIssueTypeService(method)

			issueType, err := s.Delete(context.Background(), tc.projectIDOrKey, tc.issueTypeID, tc.substituteIssueTypeID)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, issueType)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, issueType)
			assert.Equal(t, 1, issueType.ID)
		})
	}
}
