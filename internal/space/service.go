package space

import (
	"context"
	"net/url"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

type Service struct {
	method *core.Method
}

// One returns information about your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-space
func (s *Service) One(ctx context.Context) (*model.Space, error) {
	resp, err := s.method.Get(ctx, "space", nil)
	if err != nil {
		return nil, err
	}

	v := model.Space{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// DiskUsage returns information about the disk usage of your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-space-disk-usage
func (s *Service) DiskUsage(ctx context.Context) (*model.DiskUsageSpace, error) {
	resp, err := s.method.Get(ctx, "space/diskUsage", nil)
	if err != nil {
		return nil, err
	}

	v := model.DiskUsageSpace{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Notification returns the space notification.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-space-notification
func (s *Service) Notification(ctx context.Context) (*model.SpaceNotification, error) {
	resp, err := s.method.Get(ctx, "space/notification", nil)
	if err != nil {
		return nil, err
	}

	v := model.SpaceNotification{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// UpdateNotification updates the space notification.
//
// content must not be empty.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-space-notification
func (s *Service) UpdateNotification(ctx context.Context, content string) (*model.SpaceNotification, error) {
	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamContent}
	if err := core.ApplyOptions(form, validTypes, option.WithContent(content)); err != nil {
		return nil, err
	}

	resp, err := s.method.Put(ctx, "space/notification", form)
	if err != nil {
		return nil, err
	}

	v := model.SpaceNotification{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewService(method *core.Method) *Service {
	return &Service{
		method: method,
	}
}
