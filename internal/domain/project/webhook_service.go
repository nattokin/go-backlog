package project

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
)

// WebhookService handles webhook-related Backlog API calls for a project.
type WebhookService struct {
	method *client.Method
}

// List returns a list of webhooks in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-webhooks/
func (s *WebhookService) List(ctx context.Context, projectIDOrKey string) ([]*model.Webhook, error) {
	var ves validation.Errors
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
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Add adds a new webhook to a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-webhook/
func (s *WebhookService) Add(ctx context.Context, projectIDOrKey, name, hookURL string, opts ...*option.APIParamOption) (*model.Webhook, error) {
	optSvc := &option.OptionService{}
	form := url.Values{}
	validTypes := []option.APIParamOptionType{
		option.ParamName,
		option.ParamDescription,
		option.ParamHookURL,
		option.ParamAllEvent,
		option.ParamActivityTypeIDs,
	}
	options := append(
		[]*option.APIParamOption{
			optSvc.WithName(name),
			optSvc.WithHookURL(hookURL),
		},
		opts...,
	)

	var ves validation.Errors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, validTypes, options...)); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "webhooks")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := &model.Webhook{}
	if err := client.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
}

// One returns a single webhook.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-webhook/
func (s *WebhookService) One(ctx context.Context, projectIDOrKey string, webhookID int) (*model.Webhook, error) {
	var ves validation.Errors
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
	if err := client.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
}

// Update updates a webhook.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-webhook/
func (s *WebhookService) Update(ctx context.Context, projectIDOrKey string, webhookID int, opt *option.APIParamOption, opts ...*option.APIParamOption) (*model.Webhook, error) {
	form := url.Values{}
	options := append([]*option.APIParamOption{opt}, opts...)

	var ves validation.Errors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateWebhookID(webhookID); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, []option.APIParamOptionType{
		option.ParamName,
		option.ParamDescription,
		option.ParamHookURL,
		option.ParamAllEvent,
		option.ParamActivityTypeIDs,
	}, options...)); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "webhooks", strconv.Itoa(webhookID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := &model.Webhook{}
	if err := client.DecodeResponse(resp, v); err != nil {
		return nil, err
	}
	return v, nil
}

// Delete deletes a webhook from a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-webhook/
func (s *WebhookService) Delete(ctx context.Context, projectIDOrKey string, webhookID int) (*model.Webhook, error) {
	var ves validation.Errors
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
	if err := client.DecodeResponse(resp, v); err != nil {
		return nil, err
	}
	return v, nil
}

func NewWebhookService(method *client.Method) *WebhookService {
	return &WebhookService{method: method}
}
