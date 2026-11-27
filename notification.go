package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/domain/notification"
	"github.com/nattokin/go-backlog/internal/option"
)

// ──────────────────────────────────────────────────────────────
//  NotificationService
// ──────────────────────────────────────────────────────────────

// NotificationService handles communication with the notification-related methods of the Backlog API.
// All endpoints are scoped to the authenticated user, so no userID argument is needed.
type NotificationService struct {
	base *notification.Service

	// Option provides the options accepted by [NotificationService] methods.
	Option *NotificationOptionService
}

// List returns a list of notifications for the authenticated user.
//
// This method supports options returned by methods in [NotificationOptionService],
// such as:
//   - [NotificationOptionService.WithCount]
//   - [NotificationOptionService.WithMaxID]
//   - [NotificationOptionService.WithMinID]
//   - [NotificationOptionService.WithOrder]
//   - [NotificationOptionService.WithSenderID]
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-notification
func (s *NotificationService) List(ctx context.Context, opts ...RequestOption) ([]*Notification, error) {
	v, err := s.base.List(ctx, toInnerOptions(opts)...)
	return notificationsFromModel(v), convertError(err)
}

// Count returns the number of notifications for the authenticated user.
//
// This method supports options returned by methods in [NotificationOptionService],
// such as:
//   - [NotificationOptionService.WithAlreadyRead]
//   - [NotificationOptionService.WithResourceAlreadyRead]
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-notification
func (s *NotificationService) Count(ctx context.Context, opts ...RequestOption) (int, error) {
	v, err := s.base.Count(ctx, toInnerOptions(opts)...)
	return v, convertError(err)
}

// ResetUnreadCount resets the unread notification count of the authenticated user.
// It returns the count in the response.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/reset-unread-notification-count
func (s *NotificationService) ResetUnreadCount(ctx context.Context) (int, error) {
	v, err := s.base.ResetUnreadCount(ctx)
	return v, convertError(err)
}

// MarkAsRead marks the notification with the given ID as read.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/read-notification
func (s *NotificationService) MarkAsRead(ctx context.Context, id int) error {
	return convertError(s.base.MarkAsRead(ctx, id))
}

// ──────────────────────────────────────────────────────────────
//  NotificationOptionService
// ──────────────────────────────────────────────────────────────

// NotificationOptionService provides a domain-specific set of option builders
// for operations within the NotificationService.
type NotificationOptionService struct {
	base *option.OptionService
}

// WithAlreadyRead counts only notifications that are read (true) or unread (false)
// since the user last checked.
func (s *NotificationOptionService) WithAlreadyRead(enabled bool) RequestOption {
	return &requestOption{opt: s.base.WithAlreadyRead(enabled)}
}

// WithCount sets the number of notifications to retrieve (1-100).
func (s *NotificationOptionService) WithCount(count int) RequestOption {
	return &requestOption{opt: s.base.WithCount(count)}
}

// WithMaxID sets the maximum notification ID to filter results.
func (s *NotificationOptionService) WithMaxID(id int) RequestOption {
	return &requestOption{opt: s.base.WithMaxID(id)}
}

// WithMinID sets the minimum notification ID to filter results.
func (s *NotificationOptionService) WithMinID(id int) RequestOption {
	return &requestOption{opt: s.base.WithMinID(id)}
}

// WithOrder sets the sort order of results.
func (s *NotificationOptionService) WithOrder(order Order) RequestOption {
	return &requestOption{opt: s.base.WithOrder(string(order))}
}

// WithResourceAlreadyRead counts only notifications that are read (true) or unread (false).
func (s *NotificationOptionService) WithResourceAlreadyRead(enabled bool) RequestOption {
	return &requestOption{opt: s.base.WithResourceAlreadyRead(enabled)}
}

// WithSenderID filters notifications by the ID of the user who sent them.
func (s *NotificationOptionService) WithSenderID(id int) RequestOption {
	return &requestOption{opt: s.base.WithSenderID(id)}
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newNotificationService(method *client.Method, option *option.OptionService) *NotificationService {
	return &NotificationService{
		base:   notification.NewService(method),
		Option: newNotificationOptionService(option),
	}
}

func newNotificationOptionService(option *option.OptionService) *NotificationOptionService {
	return &NotificationOptionService{
		base: option,
	}
}
