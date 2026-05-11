package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/user"
)

// UserStarService handles communication with the user star-related methods of the Backlog API.
type UserStarService struct {
	base *user.StarService

	Option *UserStarOptionService
}

// List returns a list of stars received by the user with the given ID.
//
// This method supports options returned by methods in "*Client.User.Star.Option",
// such as:
//   - WithCount
//   - WithMaxID
//   - WithMinID
//   - WithOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-received-star-list
func (s *UserStarService) List(ctx context.Context, userID int, opts ...RequestOption) ([]*Star, error) {
	v, err := s.base.List(ctx, userID, toCoreOptions(opts)...)
	return starsFromModel(v), convertError(err)
}

// Count returns the number of stars received by the user with the given ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-user-received-stars
func (s *UserStarService) Count(ctx context.Context, userID int) (int, error) {
	v, err := s.base.Count(ctx, userID)
	return v, convertError(err)
}

// UserStarOptionService provides a domain-specific set of option builders
// for operations within the UserStarService.
type UserStarOptionService struct {
	base *core.OptionService
}

// WithCount sets the number of results to return.
func (s *UserStarOptionService) WithCount(count int) RequestOption {
	return s.base.WithCount(count)
}

// WithMaxID sets the maximum ID to filter results.
func (s *UserStarOptionService) WithMaxID(id int) RequestOption {
	return s.base.WithMaxID(id)
}

// WithMinID sets the minimum ID to filter results.
func (s *UserStarOptionService) WithMinID(id int) RequestOption {
	return s.base.WithMinID(id)
}

// WithOrder sets the sort order of results.
func (s *UserStarOptionService) WithOrder(order Order) RequestOption {
	return s.base.WithOrder(string(order))
}

func newUserStarService(method *core.Method, option *core.OptionService) *UserStarService {
	return &UserStarService{
		base:   user.NewStarService(method),
		Option: newUserStarOptionService(option),
	}
}

func newUserStarOptionService(option *core.OptionService) *UserStarOptionService {
	return &UserStarOptionService{
		base: option,
	}
}
