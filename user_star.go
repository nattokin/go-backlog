package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/domain/user"
	"github.com/nattokin/go-backlog/internal/option"
)

// UserStarService handles communication with the user star-related methods of the Backlog API.
type UserStarService struct {
	base *user.StarService

	// Option provides the options accepted by [UserStarService] methods.
	Option *UserStarOptionService
}

// List returns a list of stars received by the user with the given ID.
//
// This method supports options returned by methods in [UserStarOptionService],
// such as:
//   - [UserStarOptionService.WithCount]
//   - [UserStarOptionService.WithMaxID]
//   - [UserStarOptionService.WithMinID]
//   - [UserStarOptionService.WithOrder]
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-received-star-list
func (s *UserStarService) List(ctx context.Context, userID int, opts ...RequestOption) ([]*Star, error) {
	v, err := s.base.List(ctx, userID, toInnerOptions(opts)...)
	return starsFromModel(v), convertError(err)
}

// Count returns the number of stars received by the user with the given ID.
//
// This method supports options returned by methods in [UserStarOptionService],
// such as:
//   - [UserStarOptionService.WithSince]
//   - [UserStarOptionService.WithUntil]
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-user-received-stars
func (s *UserStarService) Count(ctx context.Context, userID int, opts ...RequestOption) (int, error) {
	v, err := s.base.Count(ctx, userID, toInnerOptions(opts)...)
	return v, convertError(err)
}

// UserStarOptionService provides a domain-specific set of option builders
// for operations within the UserStarService.
type UserStarOptionService struct {
	base *option.OptionService
}

// WithCount sets the number of stars to retrieve (1-100).
func (s *UserStarOptionService) WithCount(count int) RequestOption {
	return &requestOption{opt: s.base.WithCount(count)}
}

// WithMaxID sets the maximum ID to filter results.
func (s *UserStarOptionService) WithMaxID(id int) RequestOption {
	return &requestOption{opt: s.base.WithMaxID(id)}
}

// WithMinID sets the minimum ID to filter results.
func (s *UserStarOptionService) WithMinID(id int) RequestOption {
	return &requestOption{opt: s.base.WithMinID(id)}
}

// WithOrder sets the sort order of results.
func (s *UserStarOptionService) WithOrder(order Order) RequestOption {
	return &requestOption{opt: s.base.WithOrder(string(order))}
}

// WithSince counts only stars received on or after the given date.
// The date must be in "yyyy-MM-dd" format.
func (s *UserStarOptionService) WithSince(date string) RequestOption {
	return &requestOption{opt: s.base.WithSince(date)}
}

// WithUntil counts only stars received on or before the given date.
// The date must be in "yyyy-MM-dd" format.
func (s *UserStarOptionService) WithUntil(date string) RequestOption {
	return &requestOption{opt: s.base.WithUntil(date)}
}

func newUserStarService(method *client.Method, option *option.OptionService) *UserStarService {
	return &UserStarService{
		base:   user.NewStarService(method),
		Option: newUserStarOptionService(option),
	}
}

func newUserStarOptionService(option *option.OptionService) *UserStarOptionService {
	return &UserStarOptionService{
		base: option,
	}
}
