package webhook

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// Service handles communication with the webhook related methods of the Backlog API.
type Service struct {
	method *core.Method
}

// List returns a list of webhooks in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-webhooks/
func (s *Service) List(ctx context.Context, projectIDOrKey string) ([]*model.Webhook, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "webhooks")
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := []*model.Webhook{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Add adds a new webhook to a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-webhook/
func (s *Service) Add(ctx context.Context, projectIDOrKey, name, hookURL string, opts ...core.RequestOption) (*model.Webhook, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamName,
		core.ParamDescription,
		core.ParamHookURL,
		core.ParamAllEvent,
		core.ParamActivityTypeIDs,
	}
	options := append(
		[]core.RequestOption{
			option.WithName(name),
			option.WithHookURL(hookURL),
		},
		opts...,
	)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	allEvent := form.Get("allEvent")
	activityIDs := form["activityTypeId[]"]
	if allEvent == "false" && len(activityIDs) == 0 {
		return nil, core.NewValidationError("activityTypeIds is required when allEvent is false")
	}
	if allEvent == "" && len(activityIDs) == 0 {
		return nil, core.NewValidationError("requires WithAllEvent(true) or WithActivityTypeIDs")
	}

	spath := path.Join("projects", projectIDOrKey, "webhooks")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := &model.Webhook{}
	if err := core.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
}

// Get returns information about a specific webhook.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-webhook/
func (s *Service) Get(ctx context.Context, projectIDOrKey string, webhookID int) (*model.Webhook, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateWebhookID(webhookID); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "webhooks", strconv.Itoa(webhookID))
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := &model.Webhook{}
	if err := core.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
}

// Update updates a webhook.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-webhook/
func (s *Service) Update(ctx context.Context, projectIDOrKey string, webhookID int, opts ...core.RequestOption) (*model.Webhook, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateWebhookID(webhookID); err != nil {
		return nil, err
	}

	if !core.HasRequiredOption(opts, []core.APIParamOptionType{
		core.ParamName,
		core.ParamDescription,
		core.ParamHookURL,
		core.ParamAllEvent,
		core.ParamActivityTypeIDs,
	}) {
		return nil, core.NewValidationError("requires at least one webhook update option")
	}

	form := url.Values{}
	if err := core.ApplyOptions(form, []core.APIParamOptionType{
		core.ParamName,
		core.ParamDescription,
		core.ParamHookURL,
		core.ParamAllEvent,
		core.ParamActivityTypeIDs,
	}, opts...); err != nil {
		return nil, err
	}

	allEvent := form.Get("allEvent")
	activityIDs := form["activityTypeId[]"]
	if allEvent == "false" && len(activityIDs) == 0 {
		return nil, core.NewValidationError("activityTypeIds is required when allEvent is false")
	}
	if allEvent == "true" && len(activityIDs) > 0 {
		return nil, core.NewValidationError("activityTypeIds cannot be specified when allEvent is true")
	}

	spath := path.Join("projects", projectIDOrKey, "webhooks", strconv.Itoa(webhookID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := &model.Webhook{}
	if err := core.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
}

// Delete deletes a webhook from a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-webhook/
func (s *Service) Delete(ctx context.Context, projectIDOrKey string, webhookID int) (*model.Webhook, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateWebhookID(webhookID); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "webhooks", strconv.Itoa(webhookID))
	resp, err := s.method.Delete(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := &model.Webhook{}
	if err := core.DecodeResponse(resp, v); err != nil {
		return nil, err
	}
	return v, nil
}

// NewService creates and returns a new webhook Service.
func NewService(method *core.Method) *Service {
	return &Service{method: method}
}
