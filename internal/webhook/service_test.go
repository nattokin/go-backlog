package webhook_test

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
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
	"github.com/nattokin/go-backlog/internal/webhook"
)

func TestService_List(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		wantErrType    error
		mockGetFn      func(context.Context, string, url.Values) (*http.Response, error)
	}{
		"success": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks", spath)
				return mock.NewJSONResponse("[" + fixture.Webhook.AllEventJSON + "," + fixture.Webhook.ActivityTypesJSON + "]"), nil
			},
		},
		"error-client": {
			projectIDOrKey: "TEST",
			wantErrType:    errors.New(""),
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
			projectIDOrKey: "TEST",
			wantErrType:    &json.SyntaxError{},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
		"error-project-empty": {
			projectIDOrKey: "",
			wantErrType:    &core.ValidationError{},
			mockGetFn:      mock.NewUnexpectedGetFn(t),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := mock.NewMethod(t)
			m.Get = tc.mockGetFn
			s := webhook.NewService(m)

			got, err := s.List(context.Background(), tc.projectIDOrKey)
			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, got, 2)
			assert.Equal(t, fixture.Webhook.AllEvent.ID, got[0].ID)
		})
	}
}

func TestService_Add(t *testing.T) {
	option := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		name           string
		hookURL        string
		opts           []core.RequestOption
		wantErrType    error
		wantID         int
		mockPostFn     func(context.Context, string, url.Values) (*http.Response, error)
	}{
		"success-all-event-false-with-activity-types": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts: []core.RequestOption{
				option.WithAllEvent(false),
				option.WithActivityTypeIDs([]int{1, 2}),
			},
			wantID: fixture.Webhook.ActivityTypes.ID,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks", spath)
				assert.Equal(t, "false", form.Get("allEvent"))
				assert.Equal(t, []string{"1", "2"}, form["activityTypeId[]"])
				return mock.NewJSONResponse(fixture.Webhook.ActivityTypesJSON), nil
			},
		},
		"success-all-event-true": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts: []core.RequestOption{
				option.WithAllEvent(true),
			},
			wantID: fixture.Webhook.AllEvent.ID,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks", spath)
				assert.Equal(t, "webhook", form.Get("name"))
				assert.Equal(t, "https://example.com/webhook", form.Get("hookUrl"))
				assert.Equal(t, "true", form.Get("allEvent"))
				return mock.NewJSONResponse(fixture.Webhook.AllEventJSON), nil
			},
		},
		"success-activity-types-only": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts: []core.RequestOption{
				option.WithActivityTypeIDs([]int{1, 2}),
			},
			wantID: fixture.Webhook.ActivityTypes.ID,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks", spath)
				assert.Equal(t, []string{"1", "2"}, form["activityTypeId[]"])
				return mock.NewJSONResponse(fixture.Webhook.ActivityTypesJSON), nil
			},
		},
		"error-all-event-false-without-activity-types": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts: []core.RequestOption{
				option.WithAllEvent(false),
			},
			wantErrType: &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts: []core.RequestOption{
				option.WithAllEvent(true),
			},
			wantErrType: errors.New(""),
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		},
		"error-hookURL-empty": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "",
			opts: []core.RequestOption{
				option.WithAllEvent(true),
			},
			wantErrType: &core.ValidationError{},
		},
		"error-invalid-option-type": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-name-empty": {
			projectIDOrKey: "TEST",
			name:           "",
			hookURL:        "https://example.com/webhook",
			opts: []core.RequestOption{
				option.WithAllEvent(true),
			},
			wantErrType: &core.ValidationError{},
		},
		"error-no-options": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts:           []core.RequestOption{},
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-set-failed": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts:           []core.RequestOption{mock.NewFailingSetOption(core.ParamAllEvent)},
			wantErrType:    errors.New(""),
		},
		"error-project-empty": {
			projectIDOrKey: "",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			wantErrType:    &core.ValidationError{},
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			name:           "webhook",
			hookURL:        "https://example.com/webhook",
			opts: []core.RequestOption{
				option.WithAllEvent(true),
			},
			wantErrType: &json.SyntaxError{},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
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
			s := webhook.NewService(m)

			got, err := s.Add(context.Background(), tc.projectIDOrKey, tc.name, tc.hookURL, tc.opts...)
			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestService_Get(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		webhookID      int
		wantErrType    error
		wantID         int
		mockGetFn      func(context.Context, string, url.Values) (*http.Response, error)
	}{
		"success": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			wantID:         fixture.Webhook.AllEvent.ID,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks/1", spath)
				assert.Nil(t, query)
				return mock.NewJSONResponse(fixture.Webhook.AllEventJSON), nil
			},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			wantErrType:    errors.New(""),
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks/1", spath)
				return nil, errors.New("network error")
			},
		},
		"error-project-empty": {
			projectIDOrKey: "",
			webhookID:      1,
			wantErrType:    &core.ValidationError{},
			mockGetFn:      mock.NewUnexpectedGetFn(t),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			wantErrType:    &json.SyntaxError{},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks/1", spath)
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
		"error-webhookID-negative": {
			projectIDOrKey: "TEST",
			webhookID:      -1,
			wantErrType:    &core.ValidationError{},
			mockGetFn:      mock.NewUnexpectedGetFn(t),
		},
		"error-webhookID-zero": {
			projectIDOrKey: "TEST",
			webhookID:      0,
			wantErrType:    &core.ValidationError{},
			mockGetFn:      mock.NewUnexpectedGetFn(t),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Get = tc.mockGetFn
			s := webhook.NewService(method)

			got, err := s.Get(context.Background(), tc.projectIDOrKey, tc.webhookID)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestService_Update(t *testing.T) {
	option := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		webhookID      int
		opt            core.RequestOption
		opts           []core.RequestOption
		wantErrType    error
		wantID         int
		mockPatchFn    func(context.Context, string, url.Values) (*http.Response, error)
	}{
		"success-activity-types-only": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            option.WithActivityTypeIDs([]int{1, 2}),
			wantID:         fixture.Webhook.ActivityTypes.ID,
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks/1", spath)
				assert.Equal(t, []string{"1", "2"}, form["activityTypeId[]"])
				return mock.NewJSONResponse(fixture.Webhook.ActivityTypesJSON), nil
			},
		},
		"success-all-event-false-with-activity-types": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            option.WithAllEvent(false),
			opts:           []core.RequestOption{option.WithActivityTypeIDs([]int{1, 2})},
			wantID:         fixture.Webhook.ActivityTypes.ID,
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks/1", spath)
				assert.Equal(t, "false", form.Get("allEvent"))
				assert.Equal(t, []string{"1", "2"}, form["activityTypeId[]"])
				return mock.NewJSONResponse(fixture.Webhook.ActivityTypesJSON), nil
			},
		},
		"success-all-event-true": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            option.WithAllEvent(true),
			wantID:         fixture.Webhook.AllEvent.ID,
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks/1", spath)
				assert.Equal(t, "true", form.Get("allEvent"))
				return mock.NewJSONResponse(fixture.Webhook.AllEventJSON), nil
			},
		},
		"error-all-event-false-without-activity-types": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            option.WithAllEvent(false),
			wantErrType:    &core.ValidationError{},
		},
		"error-all-event-true-with-activity-types": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            option.WithAllEvent(true),
			opts:           []core.RequestOption{option.WithActivityTypeIDs([]int{1, 2})},
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            option.WithAllEvent(true),
			wantErrType:    errors.New(""),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		},
		"error-hookURL-empty": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            option.WithHookURL(""),
			opts:           []core.RequestOption{option.WithAllEvent(true)},
			wantErrType:    &core.ValidationError{},
		},
		"error-invalid-option-type": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            option.WithAllEvent(true),
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-name-empty": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            option.WithName(""),
			opts:           []core.RequestOption{option.WithAllEvent(true)},
			wantErrType:    &core.ValidationError{},
		},
		"error-option-set-failed": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            mock.NewFailingSetOption(core.ParamAllEvent),
			wantErrType:    errors.New(""),
		},
		"error-project-empty": {
			projectIDOrKey: "",
			webhookID:      1,
			opt:            option.WithAllEvent(true),
			wantErrType:    &core.ValidationError{},
			mockPatchFn:    mock.NewUnexpectedPatchFn(t),
		},
		"error-webhookID-zero": {
			projectIDOrKey: "TEST",
			webhookID:      0,
			opt:            option.WithAllEvent(true),
			wantErrType:    &core.ValidationError{},
			mockPatchFn:    mock.NewUnexpectedPatchFn(t),
		},
		"error-webhookID-negative": {
			projectIDOrKey: "TEST",
			webhookID:      -1,
			opt:            option.WithAllEvent(true),
			wantErrType:    &core.ValidationError{},
			mockPatchFn:    mock.NewUnexpectedPatchFn(t),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			opt:            option.WithAllEvent(true),
			wantErrType:    &json.SyntaxError{},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
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
			s := webhook.NewService(m)
			got, err := s.Update(context.Background(), tc.projectIDOrKey, tc.webhookID, tc.opt, tc.opts...)
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

func TestService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		webhookID      int
		wantID         int
		wantErrType    error
		mockDeleteFn   func(context.Context, string, url.Values) (*http.Response, error)
	}{
		"success": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			wantID:         fixture.Webhook.AllEvent.ID,
			mockDeleteFn: func(ctx context.Context, spath string, _ url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks/1", spath)
				return mock.NewJSONResponse(fixture.Webhook.AllEventJSON), nil
			},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			wantErrType:    errors.New(""),
			mockDeleteFn: func(ctx context.Context, spath string, _ url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks/1", spath)
				return nil, errors.New("network error")
			},
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			webhookID:      1,
			wantErrType:    &json.SyntaxError{},
			mockDeleteFn: func(ctx context.Context, spath string, _ url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/webhooks/1", spath)
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
		"error-project-empty": {
			projectIDOrKey: "",
			webhookID:      1,
			wantErrType:    &core.ValidationError{},
			mockDeleteFn:   mock.NewUnexpectedDeleteFn(t),
		},
		"error-webhookID-negative": {
			projectIDOrKey: "TEST",
			webhookID:      -1,
			wantErrType:    &core.ValidationError{},
			mockDeleteFn:   mock.NewUnexpectedDeleteFn(t),
		},
		"error-webhookID-zero": {
			projectIDOrKey: "TEST",
			webhookID:      0,
			wantErrType:    &core.ValidationError{},
			mockDeleteFn:   mock.NewUnexpectedDeleteFn(t),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := mock.NewMethod(t)
			m.Delete = tc.mockDeleteFn
			s := webhook.NewService(m)

			got, err := s.Delete(context.Background(), tc.projectIDOrKey, tc.webhookID)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}
