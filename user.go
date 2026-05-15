package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/user"
	"github.com/nattokin/go-backlog/internal/model"
)

// User represents user.
type User struct {
	ID          int
	UserID      string
	Name        string
	RoleType    Role
	Lang        string
	MailAddress string
}

// ──────────────────────────────────────────────────────────────
//  UserService
// ──────────────────────────────────────────────────────────────

// UserService has methods for user.
type UserService struct {
	base *user.Service

	Activity *UserActivityService
	Star     *UserStarService

	Option *UserOptionService
}

// List returns all users in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-list
func (s *UserService) List(ctx context.Context) ([]*User, error) {
	v, err := s.base.List(ctx)
	return usersFromModel(v), convertError(err)
}

// One returns a user in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user
func (s *UserService) One(ctx context.Context, id int) (*User, error) {
	v, err := s.base.One(ctx, id)
	return userFromModel(v), convertError(err)
}

// Own returns your own user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-own-user
func (s *UserService) Own(ctx context.Context) (*User, error) {
	v, err := s.base.Own(ctx)
	return userFromModel(v), convertError(err)
}

// Add adds a user to your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-user
func (s *UserService) Add(ctx context.Context, userID, password, name, mailAddress string, roleType Role) (*User, error) {
	v, err := s.base.Add(ctx, userID, password, name, mailAddress, int(roleType))
	return userFromModel(v), convertError(err)
}

// Update updates a user in your space.
//
// At least one option is required. This method supports options returned by
// methods in "*Client.User.Option", such as:
//   - WithMailAddress
//   - WithName
//   - WithPassword
//   - WithRoleType
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-user
func (s *UserService) Update(ctx context.Context, id int, option RequestOption, opts ...RequestOption) (*User, error) {
	v, err := s.base.Update(ctx, id, option, toCoreOptions(opts)...)
	return userFromModel(v), convertError(err)
}

// Delete deletes a user from your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-user
func (s *UserService) Delete(ctx context.Context, id int) (*User, error) {
	v, err := s.base.Delete(ctx, id)
	return userFromModel(v), convertError(err)
}

// Icon returns the icon image of a user.
// The caller is responsible for closing FileData.Body after use.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-icon
func (s *UserService) Icon(ctx context.Context, id int) (*FileData, error) {
	v, err := s.base.Icon(ctx, id)
	return fileDataFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  UserActivityService
// ──────────────────────────────────────────────────────────────

// UserActivityService handles communication with the user activities-related methods of the Backlog API.
type UserActivityService struct {
	base *user.ActivityService

	Option *ActivityOptionService
}

// List returns a list of user activities.
//
// This method supports options returned by methods in "*Client.User.Activity.Option",
// such as:
//   - WithActivityTypeIDs
//   - WithCount
//   - WithMaxID
//   - WithMinID
//   - WithOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-recent-updates
func (s *UserActivityService) List(ctx context.Context, userID int, opts ...RequestOption) ([]*Activity, error) {
	v, err := s.base.List(ctx, userID, toCoreOptions(opts)...)
	return activitiesFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  UserOptionService
// ──────────────────────────────────────────────────────────────

// UserOptionService provides a domain-specific set of option builders
// for operations within the UserService.
type UserOptionService struct {
	base *core.OptionService
}

// WithMailAddress sets the mail address of a user.
func (s *UserOptionService) WithMailAddress(mail string) RequestOption {
	return s.base.WithMailAddress(mail)
}

// WithName sets the name of a user.
func (s *UserOptionService) WithName(name string) RequestOption {
	return s.base.WithName(name)
}

// WithPassword sets the password of a user.
func (s *UserOptionService) WithPassword(password string) RequestOption {
	return s.base.WithPassword(password)
}

// WithRoleType sets the role type of a user.
func (s *UserOptionService) WithRoleType(role Role) RequestOption {
	return s.base.WithRoleType(int(role))
}

// WithSendMail sets whether to send a mail notification.
func (s *UserOptionService) WithSendMail(enabled bool) RequestOption {
	return s.base.WithSendMail(enabled)
}

// WithUserID sets the user ID.
func (s *UserOptionService) WithUserID(id int) RequestOption {
	return s.base.WithUserID(id)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newUserService(method *core.Method, option *core.OptionService) *UserService {
	return &UserService{
		base: user.NewService(method),

		Activity: newUserActivityService(method, option),
		Star:     newUserStarService(method, option),

		Option: newUserOptionService(option),
	}
}

func newUserActivityService(method *core.Method, option *core.OptionService) *UserActivityService {
	return &UserActivityService{
		base:   user.NewActivityService(method),
		Option: newActivityOptionService(option),
	}
}

func newUserOptionService(option *core.OptionService) *UserOptionService {
	return &UserOptionService{
		base: option,
	}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func userFromModel(m *model.User) *User {
	if m == nil {
		return nil
	}
	return &User{
		ID:          m.ID,
		UserID:      m.UserID,
		Name:        m.Name,
		RoleType:    Role(m.RoleType),
		Lang:        m.Lang,
		MailAddress: m.MailAddress,
	}
}

func usersFromModel(ms []*model.User) []*User {
	result := make([]*User, len(ms))
	for i, v := range ms {
		result[i] = userFromModel(v)
	}
	return result
}
