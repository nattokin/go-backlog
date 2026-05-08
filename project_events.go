package backlog

import (
	"context"
	"time"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/project"
)

// Webhook represents webhook of Backlog.
type Webhook struct {
	ID              int
	Name            string
	Description     string
	HookURL         string
	AllEvent        bool
	ActivityTypeIDs []int
	CreatedUser     *User
	Created         time.Time
	UpdatedUser     *User
	Updated         time.Time
}

// ──────────────────────────────────────────────────────────────
//  ProjectActivityService
// ──────────────────────────────────────────────────────────────

// ProjectActivityService handles communication with the project activities-related methods of the Backlog API.
type ProjectActivityService struct {
	base *project.ActivityService

	Option *ActivityOptionService
}

// List returns a list of activities in the project.
//
// This method supports options returned by methods in "*Client.Project.Activity.Option",
// such as:
//   - WithActivityTypeIDs
//   - WithCount
//   - WithMaxID
//   - WithMinID
//   - WithOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-recent-updates
func (s *ProjectActivityService) List(ctx context.Context, projectIDOrKey string, opts ...RequestOption) ([]*Activity, error) {
	v, err := s.base.List(ctx, projectIDOrKey, toCoreOptions(opts)...)
	return activitiesFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  ProjectWebhookService
// ──────────────────────────────────────────────────────────────

// ProjectWebhookService handles communication with the project webhook-related methods of the Backlog API.
type ProjectWebhookService struct {
	base   *project.WebhookService
	Option *ProjectWebhookOptionService
}

// All returns a list of webhooks in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-webhooks
func (s *ProjectWebhookService) All(ctx context.Context, projectIDOrKey string) ([]*Webhook, error) {
	v, err := s.base.List(ctx, projectIDOrKey)
	return webhooksFromModel(v), convertError(err)
}

// Create adds a webhook to the project.
//
// This method supports options returned by methods in "*Client.Project.Webhook.Option",
// such as:
//   - WithActivityTypeIDs
//   - WithAllEvent
//   - WithDescription
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-webhook
func (s *ProjectWebhookService) Create(ctx context.Context, projectIDOrKey, name, hookURL string, opts ...RequestOption) (*Webhook, error) {
	v, err := s.base.Add(ctx, projectIDOrKey, name, hookURL, toCoreOptions(opts)...)
	return webhookFromModel(v), convertError(err)
}

// One returns a webhook in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-webhook
func (s *ProjectWebhookService) One(ctx context.Context, projectIDOrKey string, webhookID int) (*Webhook, error) {
	v, err := s.base.Get(ctx, projectIDOrKey, webhookID)
	return webhookFromModel(v), convertError(err)
}

// Update updates a webhook.
//
// At least one option is required. This method supports options returned by
// methods in "*Client.Project.Webhook.Option", such as:
//   - WithActivityTypeIDs
//   - WithAllEvent
//   - WithDescription
//   - WithHookURL
//   - WithName
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-webhook
func (s *ProjectWebhookService) Update(ctx context.Context, projectIDOrKey string, webhookID int, option RequestOption, opts ...RequestOption) (*Webhook, error) {
	v, err := s.base.Update(ctx, projectIDOrKey, webhookID, option, toCoreOptions(opts)...)
	return webhookFromModel(v), convertError(err)
}

// Delete deletes a webhook.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-webhook
func (s *ProjectWebhookService) Delete(ctx context.Context, projectIDOrKey string, webhookID int) (*Webhook, error) {
	v, err := s.base.Delete(ctx, projectIDOrKey, webhookID)
	return webhookFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  ProjectWebhookOptionService
// ──────────────────────────────────────────────────────────────

// ProjectWebhookOptionService provides a domain-specific set of option builders
// for operations within the ProjectWebhookService.
type ProjectWebhookOptionService struct {
	base *core.OptionService
}

// WithActivityTypeIDs sets activity type IDs for webhook events.
func (s *ProjectWebhookOptionService) WithActivityTypeIDs(typeIDs []int) RequestOption {
	return s.base.WithActivityTypeIDs(typeIDs)
}

// WithAllEvent sets whether the webhook receives all events.
func (s *ProjectWebhookOptionService) WithAllEvent(enabled bool) RequestOption {
	return s.base.WithAllEvent(enabled)
}

// WithDescription sets the webhook description.
func (s *ProjectWebhookOptionService) WithDescription(description string) RequestOption {
	return s.base.WithDescription(description)
}

// WithHookURL sets the webhook URL.
func (s *ProjectWebhookOptionService) WithHookURL(hookURL string) RequestOption {
	return s.base.WithHookURL(hookURL)
}

// WithName sets the webhook name.
func (s *ProjectWebhookOptionService) WithName(name string) RequestOption {
	return s.base.WithName(name)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newProjectActivityService(method *core.Method, option *core.OptionService) *ProjectActivityService {
	return &ProjectActivityService{
		base:   project.NewActivityService(method),
		Option: newActivityOptionService(option),
	}
}

func newProjectWebhookService(method *core.Method, option *core.OptionService) *ProjectWebhookService {
	return &ProjectWebhookService{
		base:   project.NewWebhookService(method),
		Option: newWebhookOptionService(option),
	}
}

func newWebhookOptionService(option *core.OptionService) *ProjectWebhookOptionService {
	return &ProjectWebhookOptionService{base: option}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func webhookFromModel(m *model.Webhook) *Webhook {
	if m == nil {
		return nil
	}

	return &Webhook{
		ID:              m.ID,
		Name:            m.Name,
		Description:     m.Description,
		HookURL:         m.HookURL,
		AllEvent:        m.AllEvent,
		ActivityTypeIDs: m.ActivityTypeIDs,
		CreatedUser:     userFromModel(m.CreatedUser),
		Created:         m.Created,
		UpdatedUser:     userFromModel(m.UpdatedUser),
		Updated:         m.Updated,
	}
}

func webhooksFromModel(ms []*model.Webhook) []*Webhook {
	if ms == nil {
		return nil
	}

	result := make([]*Webhook, len(ms))
	for i, v := range ms {
		result[i] = webhookFromModel(v)
	}
	return result
}
