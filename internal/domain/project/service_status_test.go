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

func TestStatusService_List(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantLen     int
		wantErrType error
	}{
		"success-projectIDOrKey-key": {
			projectIDOrKey: "TEST",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/statuses", spath)
				assert.Nil(t, query)
				return mock.NewJSONResponse(fixture.Status.ListJSON), nil
			},

			wantLen:     2,
			wantErrType: nil,
		},
		"success-projectIDOrKey-id": {
			projectIDOrKey: "6",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/6/statuses", spath)
				return mock.NewJSONResponse(fixture.Status.ListJSON), nil
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
			s := project.NewStatusService(method)

			statuses, err := s.List(context.Background(), tc.projectIDOrKey)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, statuses)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, statuses)
			assert.Len(t, statuses, tc.wantLen)
		})
	}
}

func TestStatusService_Create(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		name           string
		color          string

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			name:           "Open",
			color:          "#ed8077",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/statuses", spath)
				assert.Equal(t, "Open", form.Get("name"))
				assert.Equal(t, "#ed8077", form.Get("color"))
				return mock.NewJSONResponse(fixture.Status.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			name:           "Open",
			color:          "#ed8077",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-name-empty": {
			projectIDOrKey: "TEST",
			name:           "",
			color:          "#ed8077",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-color-empty": {
			projectIDOrKey: "TEST",
			name:           "Open",
			color:          "",
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			name:           "Open",
			color:          "#ed8077",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			name:           "Open",
			color:          "#ed8077",

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
			s := project.NewStatusService(method)

			status, err := s.Create(context.Background(), tc.projectIDOrKey, tc.name, tc.color)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, status)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, status)
			assert.Equal(t, 1, status.ID)
			assert.Equal(t, "Open", status.Name)
			assert.Equal(t, "#ed8077", status.Color)
		})
	}
}

func TestStatusService_Update(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		statusID       int
		option         core.RequestOption
		opts           []core.RequestOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			statusID:       1,
			option:         o.WithName("Open Updated"),
			opts: []core.RequestOption{
				o.WithColor("#f5ab35"),
			},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/statuses/1", spath)
				assert.Equal(t, "Open Updated", form.Get("name"))
				assert.Equal(t, "#f5ab35", form.Get("color"))
				return mock.NewJSONResponse(fixture.Status.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			statusID:       1,
			option:         o.WithName("Open"),
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-statusID-zero": {
			projectIDOrKey: "TEST",
			statusID:       0,
			option:         o.WithName("Open"),
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: "TEST",
			statusID:       1,
			option:         mock.NewInvalidTypeOption(),
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			statusID:       1,
			option:         o.WithName("Open"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			statusID:       1,
			option:         o.WithName("Open"),

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
			s := project.NewStatusService(method)

			status, err := s.Update(context.Background(), tc.projectIDOrKey, tc.statusID, tc.option, tc.opts...)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, status)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, status)
			assert.Equal(t, 1, status.ID)
		})
	}
}

func TestStatusService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey     string
		statusID           int
		substituteStatusID int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey:     "TEST",
			statusID:           1,
			substituteStatusID: 2,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/statuses/1", spath)
				assert.Equal(t, "2", form.Get("substituteStatusId"))
				return mock.NewJSONResponse(fixture.Status.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:     "",
			statusID:           1,
			substituteStatusID: 2,
			wantErrType:        &core.ValidationError{},
		},
		"error-validation-statusID-zero": {
			projectIDOrKey:     "TEST",
			statusID:           0,
			substituteStatusID: 2,
			wantErrType:        &core.ValidationError{},
		},
		"error-validation-substituteStatusID-zero": {
			projectIDOrKey:     "TEST",
			statusID:           1,
			substituteStatusID: 0,
			wantErrType:        &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey:     "TEST",
			statusID:           1,
			substituteStatusID: 2,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey:     "TEST",
			statusID:           1,
			substituteStatusID: 2,

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
			s := project.NewStatusService(method)

			status, err := s.Delete(context.Background(), tc.projectIDOrKey, tc.statusID, tc.substituteStatusID)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, status)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, status)
			assert.Equal(t, 1, status.ID)
		})
	}
}

func TestStatusService_UpdateOrder(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		statusIDs      []int

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantLen     int
		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			statusIDs:      []int{2, 1},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/statuses/updateDisplayOrder", spath)
				assert.Equal(t, []string{"2", "1"}, form["statusId[]"])
				return mock.NewJSONResponse(fixture.Status.ListJSON), nil
			},

			wantLen:     2,
			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			statusIDs:      []int{1, 2},
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-statusIDs-empty": {
			projectIDOrKey: "TEST",
			statusIDs:      []int{},
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-statusID-zero": {
			projectIDOrKey: "TEST",
			statusIDs:      []int{1, 0},
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			statusIDs:      []int{1, 2},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			statusIDs:      []int{1, 2},

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
			s := project.NewStatusService(method)

			statuses, err := s.UpdateOrder(context.Background(), tc.projectIDOrKey, tc.statusIDs)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, statuses)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, statuses)
			assert.Len(t, statuses, tc.wantLen)
		})
	}
}
