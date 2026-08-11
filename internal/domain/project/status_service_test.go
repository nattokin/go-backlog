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

		wantLen                int
		wantErrType            error
		wantValidationErrCount int
	}{
		"success-key": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/statuses", spath)
				return mock.NewResponse(fixture.Status.ListJSON), nil
			},
			wantLen: 2,
		},
		"success-id": {
			projectIDOrKey: "6",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/6/statuses", spath)
				return mock.NewResponse(fixture.Status.ListJSON), nil
			},
			wantLen: 2,
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			wantValidationErrCount: 1,
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
				return mock.NewResponse(fixture.InvalidJSON), nil
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

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, statuses)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
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

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			name:           "Open",
			color:          "#ed8077",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/statuses", spath)
				assert.Equal(t, "Open", form.Get("name"))
				assert.Equal(t, "#ed8077", form.Get("color"))
				return mock.NewResponse(fixture.Status.SingleJSON), nil
			},
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			name:                   "Open",
			color:                  "#ed8077",
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			name:                   "Open",
			color:                  "#ed8077",
			wantValidationErrCount: 1,
		},

		"error-validation-name-empty": {
			projectIDOrKey:         "TEST",
			name:                   "",
			color:                  "#ed8077",
			wantValidationErrCount: 1,
		},
		"error-validation-color-empty": {
			projectIDOrKey:         "TEST",
			name:                   "Open",
			color:                  "",
			wantValidationErrCount: 1,
		},
		"error-validation-name-and-color-empty": {
			projectIDOrKey:         "TEST",
			name:                   "",
			color:                  "",
			wantValidationErrCount: 2,
		},

		"error-validation-all": {
			projectIDOrKey:         "",
			name:                   "",
			color:                  "",
			wantValidationErrCount: 3,
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
				return mock.NewResponse(fixture.InvalidJSON), nil
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

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, status)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
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
		option         *core.APIParamOption
		opts           []*core.APIParamOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success": {
			projectIDOrKey: "TEST",
			statusID:       1,
			option:         o.WithName("Open Updated"),
			opts:           []*core.APIParamOption{o.WithColor("#f5ab35")},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/statuses/1", spath)
				assert.Equal(t, "Open Updated", form.Get("name"))
				assert.Equal(t, "#f5ab35", form.Get("color"))
				return mock.NewResponse(fixture.Status.SingleJSON), nil
			},
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			statusID:               1,
			option:                 o.WithName("Open"),
			wantValidationErrCount: 1,
		},
		"error-validation-statusID-zero": {
			projectIDOrKey:         "TEST",
			statusID:               0,
			option:                 o.WithName("Open"),
			wantValidationErrCount: 1,
		},

		"error-validation-fixed-option": {
			projectIDOrKey:         "TEST",
			statusID:               1,
			option:                 o.WithName(""),
			wantValidationErrCount: 1,
		},

		"error-validation-opt-single": {
			projectIDOrKey:         "TEST",
			statusID:               1,
			option:                 o.WithName("Open"),
			opts:                   []*core.APIParamOption{o.WithColor("")},
			wantValidationErrCount: 1,
		},

		"error-validation-all": {
			projectIDOrKey:         "",
			statusID:               0,
			option:                 o.WithName(""),
			opts:                   []*core.APIParamOption{o.WithColor("")},
			wantValidationErrCount: 4,
		},

		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "TEST",
			statusID:               1,
			option:                 o.WithName("Open"),
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			statusID:               0,
			option:                 o.WithName(""),
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		"error-option-invalid-type-with-valid-values": {
			projectIDOrKey: "TEST",
			statusID:       1,
			option:         mock.NewInvalidTypeOption(),
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			projectIDOrKey: "",
			statusID:       0,
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
				return mock.NewResponse(fixture.InvalidJSON), nil
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

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, status)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, status)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
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

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey:     "TEST",
			statusID:           1,
			substituteStatusID: 2,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/statuses/1", spath)
				assert.Equal(t, "2", form.Get("substituteStatusId"))
				return mock.NewResponse(fixture.Status.SingleJSON), nil
			},
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			statusID:               1,
			substituteStatusID:     2,
			wantValidationErrCount: 1,
		},
		"error-validation-statusID-zero": {
			projectIDOrKey:         "TEST",
			statusID:               0,
			substituteStatusID:     2,
			wantValidationErrCount: 1,
		},
		"error-validation-substituteStatusID-zero": {
			projectIDOrKey:         "TEST",
			statusID:               1,
			substituteStatusID:     0,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			statusID:               0,
			substituteStatusID:     0,
			wantValidationErrCount: 3,
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
				return mock.NewResponse(fixture.InvalidJSON), nil
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

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, status)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
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

		wantLen                int
		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			statusIDs:      []int{2, 1},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/statuses/updateDisplayOrder", spath)
				assert.Equal(t, []string{"2", "1"}, form["statusId[]"])
				return mock.NewResponse(fixture.Status.ListJSON), nil
			},
			wantLen: 2,
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			statusIDs:              []int{1, 2},
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			statusIDs:              []int{1, 2},
			wantValidationErrCount: 1,
		},
		"error-validation-statusIDs-empty": {
			projectIDOrKey:         "TEST",
			statusIDs:              []int{},
			wantValidationErrCount: 1,
		},
		"error-validation-statusID-zero": {
			projectIDOrKey:         "TEST",
			statusIDs:              []int{1, 0},
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			statusIDs:              []int{0},
			wantValidationErrCount: 2,
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
				return mock.NewResponse(fixture.InvalidJSON), nil
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

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, statuses)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, statuses)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, statuses)
			assert.Len(t, statuses, tc.wantLen)
		})
	}
}
