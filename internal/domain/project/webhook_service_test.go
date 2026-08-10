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

func TestWebhookService_List(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockGetFn func(context.Context, string, url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks", spath)
				return mock.NewResponse("[" + fixture.Webhook.AllEventJSON + "," + fixture.Webhook.ActivityTypesJSON + "]"), nil
			},
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			wantValidationErrCount: 1,
		},

		// --- other errors ---
		"error-client-network": {
			projectIDOrKey: "TEST",
			wantErrType:    errors.New(""),
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			wantErrType:    &json.SyntaxError{},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				m.Get = tc.mockGetFn
			}
			s := project.NewWebhookService(m)
			got, err := s.List(context.Background(), tc.projectIDOrKey)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.Len(t, got, 2)
			assert.Equal(t, fixture.Webhook.AllEvent.ID, got[0].ID)
		})
	}
}

func TestWebhookService_Add(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		name           string
		hookURL        string
		opts           []*core.APIParamOption

		mockPostFn func(context.Context, string, url.Values) (*http.Response, error)

		wantID                 int
		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success-all-event-false-with-activity-types": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts:           []*core.APIParamOption{o.WithAllEvent(false), o.WithActivityTypeIDs([]int{1, 2})},
			wantID:         fixture.Webhook.ActivityTypes.ID,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "false", form.Get("allEvent"))
				assert.Equal(t, []string{"1", "2"}, form["activityTypeId[]"])
				return mock.NewResponse(fixture.Webhook.ActivityTypesJSON), nil
			},
		},
		"success-all-event-true": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts:           []*core.APIParamOption{o.WithAllEvent(true)},
			wantID:         fixture.Webhook.AllEvent.ID,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "true", form.Get("allEvent"))
				return mock.NewResponse(fixture.Webhook.AllEventJSON), nil
			},
		},
		"success-activity-types-only": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts:           []*core.APIParamOption{o.WithActivityTypeIDs([]int{1, 2})},
			wantID:         fixture.Webhook.ActivityTypes.ID,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"1", "2"}, form["activityTypeId[]"])
				return mock.NewResponse(fixture.Webhook.ActivityTypesJSON), nil
			},
		},

		// --- validation errors: argument only ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			name:                   "webhook",
			hookURL:                "https://example.com/webhook",
			opts:                   []*core.APIParamOption{o.WithAllEvent(true)},
			wantValidationErrCount: 1,
		},

		// --- validation errors: fixed options only ---
		"error-validation-name-empty": {
			projectIDOrKey:         "TEST",
			name:                   "",
			hookURL:                "https://example.com/webhook",
			opts:                   []*core.APIParamOption{o.WithAllEvent(true)},
			wantValidationErrCount: 1,
		},
		"error-validation-hookURL-empty": {
			projectIDOrKey:         "TEST",
			name:                   "webhook",
			hookURL:                "",
			opts:                   []*core.APIParamOption{o.WithAllEvent(true)},
			wantValidationErrCount: 1,
		},
		"error-validation-name-and-hookURL-empty": {
			projectIDOrKey:         "TEST",
			name:                   "",
			hookURL:                "",
			opts:                   []*core.APIParamOption{o.WithAllEvent(true)},
			wantValidationErrCount: 2,
		},

		// --- validation errors: all ---
		"error-validation-all": {
			projectIDOrKey:         "",
			name:                   "",
			hookURL:                "",
			opts:                   []*core.APIParamOption{o.WithAllEvent(true)},
			wantValidationErrCount: 3,
		},

		// --- application-level validation errors ---
		"error-all-event-false-without-activity-types": {
			projectIDOrKey:         "TEST",
			name:                   "webhook",
			hookURL:                "https://example.com/webhook",
			opts:                   []*core.APIParamOption{o.WithAllEvent(false)},
			wantValidationErrCount: 1,
		},
		"error-no-options": {
			projectIDOrKey:         "TEST",
			name:                   "webhook",
			hookURL:                "https://example.com/webhook",
			opts:                   []*core.APIParamOption{},
			wantValidationErrCount: 1,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "TEST",
			name:                   "webhook",
			hookURL:                "https://example.com/webhook",
			opts:                   []*core.APIParamOption{o.WithAllEvent(true), nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			name:                   "",
			hookURL:                "",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type-with-valid-values": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			projectIDOrKey: "",
			name:           "",
			hookURL:        "",
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-option-set-failed": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts:           []*core.APIParamOption{mock.NewFailingSetOption(core.ParamAllEvent)},
			wantErrType:    errors.New(""),
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts:           []*core.APIParamOption{o.WithAllEvent(true)},
			wantErrType:    errors.New(""),
			mockPostFn:     func(context.Context, string, url.Values) (*http.Response, error) { return nil, errors.New("network") },
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts:           []*core.APIParamOption{o.WithAllEvent(true)},
			wantErrType:    &json.SyntaxError{},
			mockPostFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
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
			s := project.NewWebhookService(m)
			got, err := s.Add(context.Background(), tc.projectIDOrKey, tc.name, tc.hookURL, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, got)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestWebhookService_One(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		webhookID      int

		mockGetFn func(context.Context, string, url.Values) (*http.Response, error)

		wantID                 int
		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			wantID:         fixture.Webhook.AllEvent.ID,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks/1", spath)
				return mock.NewResponse(fixture.Webhook.AllEventJSON), nil
			},
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			webhookID:              1,
			wantValidationErrCount: 1,
		},
		"error-validation-webhookID-zero": {
			projectIDOrKey:         "TEST",
			webhookID:              0,
			wantValidationErrCount: 1,
		},
		"error-validation-webhookID-negative": {
			projectIDOrKey:         "TEST",
			webhookID:              -1,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			webhookID:              0,
			wantValidationErrCount: 2,
		},

		// --- other errors ---
		"error-client-network": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			wantErrType:    errors.New(""),
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			wantErrType:    &json.SyntaxError{},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				m.Get = tc.mockGetFn
			}
			s := project.NewWebhookService(m)
			got, err := s.One(context.Background(), tc.projectIDOrKey, tc.webhookID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestWebhookService_Update(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		webhookID      int
		opt            *core.APIParamOption
		opts           []*core.APIParamOption

		mockPatchFn func(context.Context, string, url.Values) (*http.Response, error)

		wantID                 int
		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success-activity-types-only": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            o.WithActivityTypeIDs([]int{1, 2}),
			wantID:         fixture.Webhook.ActivityTypes.ID,
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks/1", spath)
				return mock.NewResponse(fixture.Webhook.ActivityTypesJSON), nil
			},
		},
		"success-all-event-true": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            o.WithAllEvent(true),
			wantID:         fixture.Webhook.AllEvent.ID,
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "true", form.Get("allEvent"))
				return mock.NewResponse(fixture.Webhook.AllEventJSON), nil
			},
		},
		"success-all-event-true-with-activity-types": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            o.WithAllEvent(true),
			opts:           []*core.APIParamOption{o.WithActivityTypeIDs([]int{1, 2})},
			wantID:         fixture.Webhook.AllEvent.ID,
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "true", form.Get("allEvent"))
				assert.Equal(t, []string{"1", "2"}, form["activityTypeId[]"])
				return mock.NewResponse(fixture.Webhook.AllEventJSON), nil
			},
		},

		// --- validation errors: argument only ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			webhookID:              1,
			opt:                    o.WithAllEvent(true),
			wantValidationErrCount: 1,
		},
		"error-validation-webhookID-zero": {
			projectIDOrKey:         "TEST",
			webhookID:              0,
			opt:                    o.WithAllEvent(true),
			wantValidationErrCount: 1,
		},
		"error-validation-webhookID-negative": {
			projectIDOrKey:         "TEST",
			webhookID:              -1,
			opt:                    o.WithAllEvent(true),
			wantValidationErrCount: 1,
		},

		// --- validation errors: fixed option only ---
		"error-validation-fixed-option-hookURL-empty": {
			projectIDOrKey:         "TEST",
			webhookID:              1,
			opt:                    o.WithHookURL(""),
			opts:                   []*core.APIParamOption{o.WithAllEvent(true)},
			wantValidationErrCount: 1,
		},
		"error-validation-fixed-option-name-empty": {
			projectIDOrKey:         "TEST",
			webhookID:              1,
			opt:                    o.WithName(""),
			opts:                   []*core.APIParamOption{o.WithAllEvent(true)},
			wantValidationErrCount: 1,
		},

		// --- validation errors: all ---
		"error-validation-all": {
			projectIDOrKey:         "",
			webhookID:              0,
			opt:                    o.WithName(""),
			opts:                   []*core.APIParamOption{o.WithHookURL("")},
			wantValidationErrCount: 4,
		},

		// --- application-level validation errors ---
		"error-all-event-false-without-activity-types": {
			projectIDOrKey:         "TEST",
			webhookID:              1,
			opt:                    o.WithAllEvent(false),
			wantValidationErrCount: 1,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "TEST",
			webhookID:              1,
			opt:                    o.WithAllEvent(true),
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			webhookID:              0,
			opt:                    o.WithName(""),
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type-with-valid-values": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            o.WithAllEvent(true),
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			projectIDOrKey: "",
			webhookID:      0,
			opt:            mock.NewInvalidTypeOption(),
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-option-set-failed": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            mock.NewFailingSetOption(core.ParamAllEvent),
			wantErrType:    errors.New(""),
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            o.WithAllEvent(true),
			wantErrType:    errors.New(""),
			mockPatchFn:    func(context.Context, string, url.Values) (*http.Response, error) { return nil, errors.New("network") },
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            o.WithAllEvent(true),
			wantErrType:    &json.SyntaxError{},
			mockPatchFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := mock.NewMethod(t)
			if tc.mockPatchFn != nil {
				m.Patch = tc.mockPatchFn
			}
			s := project.NewWebhookService(m)
			got, err := s.Update(context.Background(), tc.projectIDOrKey, tc.webhookID, tc.opt, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, got)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestWebhookService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		webhookID      int

		mockDeleteFn func(context.Context, string, url.Values) (*http.Response, error)

		wantID                 int
		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			wantID:         fixture.Webhook.AllEvent.ID,
			mockDeleteFn: func(ctx context.Context, spath string, _ url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks/1", spath)
				return mock.NewResponse(fixture.Webhook.AllEventJSON), nil
			},
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			webhookID:              1,
			wantValidationErrCount: 1,
		},
		"error-validation-webhookID-zero": {
			projectIDOrKey:         "TEST",
			webhookID:              0,
			wantValidationErrCount: 1,
		},
		"error-validation-webhookID-negative": {
			projectIDOrKey:         "TEST",
			webhookID:              -1,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			webhookID:              0,
			wantValidationErrCount: 2,
		},

		// --- other errors ---
		"error-client-network": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			wantErrType:    errors.New(""),
			mockDeleteFn:   func(context.Context, string, url.Values) (*http.Response, error) { return nil, errors.New("network") },
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			wantErrType:    &json.SyntaxError{},
			mockDeleteFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := mock.NewMethod(t)
			if tc.mockDeleteFn != nil {
				m.Delete = tc.mockDeleteFn
			}
			s := project.NewWebhookService(m)
			got, err := s.Delete(context.Background(), tc.projectIDOrKey, tc.webhookID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}
