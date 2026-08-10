package project

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// WebhookService handles webhook-related Backlog API calls for a project.
type WebhookService struct {
	method *core.Method
}

// List returns a list of webhooks in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-webhooks/
func (s *WebhookService) List(ctx context.Context, projectIDOrKey string) ([]*model.Webhook, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
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
func (s *WebhookService) Add(ctx context.Context, projectIDOrKey, name, hookURL string, opts ...*core.APIParamOption) (*model.Webhook, error) {
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
		[]*core.APIParamOption{
			option.WithName(name),
			option.WithHookURL(hookURL),
		},
		opts...,
	)

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if err := core.MergeValidationErrors(ves, core.ApplyOptions(form, validTypes, options...)); err != nil {
		return nil, err
	}

	allEvent := form.Get("allEvent")
	activityIDs := form["activityTypeId[]"]
	var postVes core.ValidationErrors
	if allEvent == "false" && len(activityIDs) == 0 {
		postVes = append(postVes, core.NewValidationError("activityTypeIds", "activityTypeIds is required when allEvent is false"))
	}
	if allEvent == "" && len(activityIDs) == 0 {
		postVes = append(postVes, core.NewValidationError("activityTypeIds", "requires WithAllEvent(true) or WithActivityTypeIDs"))
	}
	if len(postVes) > 0 {
		return nil, postVes
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

// One returns a single webhook.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-webhook/
func (s *WebhookService) One(ctx context.Context, projectIDOrKey string, webhookID int) (*model.Webhook, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateWebhookID(webhookID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
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
func (s *WebhookService) Update(ctx context.Context, projectIDOrKey string, webhookID int, option *core.APIParamOption, opts ...*core.APIParamOption) (*model.Webhook, error) {
	form := url.Values{}
	options := append([]*core.APIParamOption{option}, opts...)

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateWebhookID(webhookID); ve != nil {
		ves = append(ves, ve)
	}
	if err := core.MergeValidationErrors(ves, core.ApplyOptions(form, []core.APIParamOptionType{
		core.ParamName,
		core.ParamDescription,
		core.ParamHookURL,
		core.ParamAllEvent,
		core.ParamActivityTypeIDs,
	}, options...)); err != nil {
		return nil, err
	}

	allEvent := form.Get("allEvent")
	activityIDs := form["activityTypeId[]"]
	var postVes core.ValidationErrors
	if allEvent == "false" && len(activityIDs) == 0 {
		postVes = append(postVes, core.NewValidationError("activityTypeIds", "activityTypeIds is required when allEvent is false"))
	}
	if allEvent == "true" && len(activityIDs) > 0 {
		postVes = append(postVes, core.NewValidationError("activityTypeIds", "activityTypeIds cannot be specified when allEvent is true"))
	}
	if len(postVes) > 0 {
		return nil, postVes
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
func (s *WebhookService) Delete(ctx context.Context, projectIDOrKey string, webhookID int) (*model.Webhook, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateWebhookID(webhookID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
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

func NewWebhookService(method *core.Method) *WebhookService {
	return &WebhookService{method: method}
}
