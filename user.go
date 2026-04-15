package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/user"
)

// ──────────────────────────────────────────────────────────────
//  UserService
// ──────────────────────────────────────────────────────────────

// UserService has methods for user.
type UserService struct {
	base *user.Service

	Activity *UserActivityService
	Option   *UserOptionService
}

// All returns all users in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-list
func (s *UserService) All(ctx context.Context) ([]*model.User, error) {
	return s.base.All(ctx)
}

// One returns a user in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user
func (s *UserService) One(ctx context.Context, id int) (*model.User, error) {
	return s.base.One(ctx, id)
}

// Own returns your own user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-own-user
func (s *UserService) Own(ctx context.Context) (*model.User, error) {
	return s.base.Own(ctx)
}

// Add adds a user to your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-user
func (s *UserService) Add(ctx context.Context, userID, password, name, mailAddress string, roleType model.Role) (*model.User, error) {
	return s.base.Add(ctx, userID, password, name, mailAddress, roleType)
}

// Update updates a user in your space.
//
// This method supports options returned by methods in "*Client.User.Option",
// such as:
//   - WithMailAddress
//   - WithName
//   - WithPassword
//   - WithRoleType
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-user
func (s *UserService) Update(ctx context.Context, id int, opts ...RequestOption) (*model.User, error) {
	return s.base.Update(ctx, id, toCoreOptions(opts)...)
}

// Delete deletes a user from your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-user
func (s *UserService) Delete(ctx context.Context, id int) (*model.User, error) {
	return s.base.Delete(ctx, id)
}

// ──────────────────────────────────────────────────────────────
//  UserActivityService
// ──────────────────────────────────────────────────────────────

// UserActivityService handles communication with the user activities-related methods of the Backlog API.
type UserActivityService struct {
	base *activity.UserService

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
func (s *UserActivityService) List(ctx context.Context, userID int, opts ...RequestOption) ([]*model.Activity, error) {
	return s.base.List(ctx, userID, toCoreOptions(opts)...)
}

// ProjectUserService has methods for user of project.
type ProjectUserService struct {
	base *user.ProjectService
}

// All returns all users in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-user-list
func (s *ProjectUserService) All(ctx context.Context, projectIDOrKey string, excludeGroupMembers bool) ([]*model.User, error) {
	return s.base.All(ctx, projectIDOrKey, excludeGroupMembers)
}

// Add adds a user to the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project-user
func (s *ProjectUserService) Add(ctx context.Context, projectIDOrKey string, userID int) (*model.User, error) {
	return s.base.Add(ctx, projectIDOrKey, userID)
}

// Delete deletes a user from the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project-user
func (s *ProjectUserService) Delete(ctx context.Context, projectIDOrKey string, userID int) (*model.User, error) {
	return s.base.Delete(ctx, projectIDOrKey, userID)
}

// AddAdmin adds a admin user to the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project-administrator
func (s *ProjectUserService) AddAdmin(ctx context.Context, projectIDOrKey string, userID int) (*model.User, error) {
	return s.base.AddAdmin(ctx, projectIDOrKey, userID)
}

// AdminAll returns a list of all admin users in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-project-administrators
func (s *ProjectUserService) AdminAll(ctx context.Context, projectIDOrKey string) ([]*model.User, error) {
	return s.base.AdminAll(ctx, projectIDOrKey)
}

// DeleteAdmin removes an admin user from the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project-administrator
func (s *ProjectUserService) DeleteAdmin(ctx context.Context, projectIDOrKey string, userID int) (*model.User, error) {
	return s.base.DeleteAdmin(ctx, projectIDOrKey, userID)
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
func (s *UserOptionService) WithRoleType(role model.Role) RequestOption {
	return s.base.WithRoleType(role)
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
		base:     user.NewService(method),
		Activity: newUserActivityService(method, option),
		Option:   newUserOptionService(option),
	}
}

func newUserActivityService(method *core.Method, option *core.OptionService) *UserActivityService {
	return &UserActivityService{
		base:   activity.NewUserService(method),
		Option: newActivityOptionService(option),
	}
}

func newProjectUserService(method *core.Method, option *core.OptionService) *ProjectUserService {
	return &ProjectUserService{
		base: user.NewProjectService(method),
	}
}

func newUserOptionService(option *core.OptionService) *UserOptionService {
	return &UserOptionService{
		base: option,
	}
}
