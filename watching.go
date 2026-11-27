package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/domain/watching"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
)

// ──────────────────────────────────────────────────────────────
//  WatchingService
// ──────────────────────────────────────────────────────────────

// WatchingService handles communication with the watching-related methods of the Backlog API.
type WatchingService struct {
	base *watching.Service

	// Option provides the options accepted by [WatchingService] methods.
	Option *WatchingOptionService
}

// List returns a list of watchings of the user with the given ID.
//
// This method supports options returned by methods in [WatchingOptionService],
// such as:
//   - [WatchingOptionService.WithCount]
//   - [WatchingOptionService.WithIssueIDs]
//   - [WatchingOptionService.WithOffset]
//   - [WatchingOptionService.WithOrder]
//   - [WatchingOptionService.WithResourceAlreadyRead]
//   - [WatchingOptionService.WithSort]
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-watching-list
func (s *WatchingService) List(ctx context.Context, userID int, opts ...RequestOption) ([]*Watching, error) {
	v, err := s.base.List(ctx, userID, toInnerOptions(opts)...)
	return watchingsFromModel(v), convertError(err)
}

// Count returns the number of watchings of the user with the given ID.
//
// This method supports options returned by methods in [WatchingOptionService],
// such as:
//   - [WatchingOptionService.WithAlreadyRead]
//   - [WatchingOptionService.WithResourceAlreadyRead]
//
// If both options are given, the API ignores
// [WatchingOptionService.WithResourceAlreadyRead].
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-watching
func (s *WatchingService) Count(ctx context.Context, userID int, opts ...RequestOption) (int, error) {
	v, err := s.base.Count(ctx, userID, toInnerOptions(opts)...)
	return v, convertError(err)
}

// One returns the watching with the given ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-watching
func (s *WatchingService) One(ctx context.Context, watchingID int) (*Watching, error) {
	v, err := s.base.One(ctx, watchingID)
	return watchingFromModel(v), convertError(err)
}

// Create adds a watching for an issue.
// issueIDOrKey can be an issue ID or an issue key.
//
// This method supports options returned by methods in [WatchingOptionService],
// such as:
//   - [WatchingOptionService.WithNote]
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-watching
func (s *WatchingService) Create(ctx context.Context, issueIDOrKey string, opts ...RequestOption) (*Watching, error) {
	v, err := s.base.Create(ctx, issueIDOrKey, toInnerOptions(opts)...)
	return watchingFromModel(v), convertError(err)
}

// Update updates the note of the watching with the given ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-watching
func (s *WatchingService) Update(ctx context.Context, watchingID int, note string) (*Watching, error) {
	v, err := s.base.Update(ctx, watchingID, note)
	return watchingFromModel(v), convertError(err)
}

// Delete deletes the watching with the given ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-watching
func (s *WatchingService) Delete(ctx context.Context, watchingID int) (*Watching, error) {
	v, err := s.base.Delete(ctx, watchingID)
	return watchingFromModel(v), convertError(err)
}

// MarkAsRead marks the watching with the given ID as read.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/mark-watching-as-read
func (s *WatchingService) MarkAsRead(ctx context.Context, watchingID int) error {
	return convertError(s.base.MarkAsRead(ctx, watchingID))
}

// ──────────────────────────────────────────────────────────────
//  WatchingOptionService
// ──────────────────────────────────────────────────────────────

// WatchingOptionService provides a domain-specific set of option builders
// for operations within the WatchingService.
type WatchingOptionService struct {
	base *option.OptionService
}

// WithAlreadyRead counts only watchings that are read (true) or unread (false)
// since the user last checked them.
func (s *WatchingOptionService) WithAlreadyRead(enabled bool) RequestOption {
	return &requestOption{opt: s.base.WithAlreadyRead(enabled)}
}

// WithCount sets the number of watchings to retrieve (1-100).
func (s *WatchingOptionService) WithCount(count int) RequestOption {
	return &requestOption{opt: s.base.WithCount(count)}
}

// WithIssueIDs filters results by issue IDs.
func (s *WatchingOptionService) WithIssueIDs(ids []int) RequestOption {
	return &requestOption{opt: s.base.WithIssueIDs(ids)}
}

// WithNote sets the note of the watching.
func (s *WatchingOptionService) WithNote(note string) RequestOption {
	return &requestOption{opt: s.base.WithNote(note)}
}

// WithOffset sets the number of watchings to skip.
func (s *WatchingOptionService) WithOffset(offset int) RequestOption {
	return &requestOption{opt: s.base.WithOffset(offset)}
}

// WithOrder sets the sort order of results.
func (s *WatchingOptionService) WithOrder(order Order) RequestOption {
	return &requestOption{opt: s.base.WithOrder(string(order))}
}

// WithResourceAlreadyRead filters results to watchings whose issue is read
// (true) or unread (false).
func (s *WatchingOptionService) WithResourceAlreadyRead(enabled bool) RequestOption {
	return &requestOption{opt: s.base.WithResourceAlreadyRead(enabled)}
}

// WithSort sets the field to sort watching list results by.
func (s *WatchingOptionService) WithSort(sort WatchingSort) RequestOption {
	return &requestOption{opt: s.base.WithWatchingSort(string(sort))}
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newWatchingService(method *client.Method, option *option.OptionService) *WatchingService {
	return &WatchingService{
		base:   watching.NewService(method),
		Option: newWatchingOptionService(option),
	}
}

func newWatchingOptionService(option *option.OptionService) *WatchingOptionService {
	return &WatchingOptionService{
		base: option,
	}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func watchingFromModel(m *model.Watching) *Watching {
	if m == nil {
		return nil
	}
	return &Watching{
		ID:                 m.ID,
		AlreadyRead:        m.AlreadyRead,
		Note:               m.Note,
		Type:               m.Type,
		Issue:              issueFromModel(m.Issue),
		LastContentUpdated: Timestamp{m.LastContentUpdated},
		Created:            Timestamp{m.Created},
		Updated:            Timestamp{m.Updated},
	}
}

func watchingsFromModel(ms []*model.Watching) []*Watching {
	if ms == nil {
		return nil
	}
	result := make([]*Watching, len(ms))
	for i, v := range ms {
		result[i] = watchingFromModel(v)
	}
	return result
}
