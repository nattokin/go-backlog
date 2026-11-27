// Package notification implements the Backlog Notification API service.
package notification

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
)

// Service handles notification-related Backlog API calls.
// All endpoints are scoped to the authenticated user.
type Service struct {
	method *client.Method
}

// List returns a list of notifications for the authenticated user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-notification
func (s *Service) List(ctx context.Context, opts ...*option.APIParamOption) ([]*model.Notification, error) {
	query := url.Values{}
	validOptionKeys := []option.APIParamOptionType{
		option.ParamMinID,
		option.ParamMaxID,
		option.ParamCount,
		option.ParamOrder,
		option.ParamSenderID,
	}
	if err := option.ApplyOptions(query, validOptionKeys, opts...); err != nil {
		return nil, err
	}

	resp, err := s.method.Get(ctx, "notifications", query)
	if err != nil {
		return nil, err
	}

	v := []*model.Notification{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Count returns the number of notifications for the authenticated user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-notification
func (s *Service) Count(ctx context.Context, opts ...*option.APIParamOption) (int, error) {
	query := url.Values{}
	validOptionKeys := []option.APIParamOptionType{
		option.ParamAlreadyRead,
		option.ParamResourceAlreadyRead,
	}
	if err := option.ApplyOptions(query, validOptionKeys, opts...); err != nil {
		return 0, err
	}

	resp, err := s.method.Get(ctx, "notifications/count", query)
	if err != nil {
		return 0, err
	}

	return decodeCount(resp)
}

// ResetUnreadCount resets the unread notification count of the authenticated user.
// It returns the count in the response.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/reset-unread-notification-count
func (s *Service) ResetUnreadCount(ctx context.Context) (int, error) {
	resp, err := s.method.Post(ctx, "notifications/markAsRead", nil)
	if err != nil {
		return 0, err
	}

	return decodeCount(resp)
}

// MarkAsRead marks the notification with the given ID as read.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/read-notification
func (s *Service) MarkAsRead(ctx context.Context, id int) error {
	var ves validation.Errors
	if ve := validate.ValidateNotificationID(id); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return ves
	}

	spath := client.PathJoin("notifications", strconv.Itoa(id), "markAsRead")
	if _, err := s.method.Post(ctx, spath, nil); err != nil {
		return err
	}

	return nil
}

// NewService returns a *Service that sends requests through method.
func NewService(method *client.Method) *Service {
	return &Service{method: method}
}

func decodeCount(resp *http.Response) (int, error) {
	var v struct {
		Count int `json:"count"`
	}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return 0, err
	}

	return v.Count, nil
}
