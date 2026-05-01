package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/project"
	"github.com/nattokin/go-backlog/internal/sharedfile"
	"github.com/nattokin/go-backlog/internal/user"
)

// ──────────────────────────────────────────────────────────────
//  ProjectCategoryService
// ──────────────────────────────────────────────────────────────

// ProjectCategoryService handles communication with the project category-related methods of the Backlog API.
type ProjectCategoryService struct {
	base *project.CategoryService
}

// All returns a list of categories in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-category-list
func (s *ProjectCategoryService) All(ctx context.Context, projectIDOrKey string) ([]*Category, error) {
	v, err := s.base.All(ctx, projectIDOrKey)
	return categoriesFromModel(v), convertError(err)
}

// Create adds a new category to a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-category
func (s *ProjectCategoryService) Create(ctx context.Context, projectIDOrKey string, name string) (*Category, error) {
	v, err := s.base.Create(ctx, projectIDOrKey, name)
	return categoryFromModel(v), convertError(err)
}

// Update updates a category in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-category
func (s *ProjectCategoryService) Update(ctx context.Context, projectIDOrKey string, categoryID int, name string) (*Category, error) {
	v, err := s.base.Update(ctx, projectIDOrKey, categoryID, name)
	return categoryFromModel(v), convertError(err)
}

// Delete deletes a category from a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-category
func (s *ProjectCategoryService) Delete(ctx context.Context, projectIDOrKey string, categoryID int) (*Category, error) {
	v, err := s.base.Delete(ctx, projectIDOrKey, categoryID)
	return categoryFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  ProjectSharedFileService
// ──────────────────────────────────────────────────────────────

// ProjectSharedFileService handles communication with the project shared-file-related methods of the Backlog API.
type ProjectSharedFileService struct {
	base *sharedfile.ProjectService
}

// List returns a list of shared files in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-shared-files
func (s *ProjectSharedFileService) List(ctx context.Context, projectIDOrKey string) ([]*SharedFile, error) {
	v, err := s.base.List(ctx, projectIDOrKey)
	return sharedFilesFromModel(v), convertError(err)
}

// ProjectUserService has methods for user of project.
type ProjectUserService struct {
	base *user.ProjectService
}

// All returns all users in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-user-list
func (s *ProjectUserService) All(ctx context.Context, projectIDOrKey string, excludeGroupMembers bool) ([]*User, error) {
	v, err := s.base.All(ctx, projectIDOrKey, excludeGroupMembers)
	return usersFromModel(v), convertError(err)
}

// Add adds a user to the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project-user
func (s *ProjectUserService) Add(ctx context.Context, projectIDOrKey string, userID int) (*User, error) {
	v, err := s.base.Add(ctx, projectIDOrKey, userID)
	return userFromModel(v), convertError(err)
}

// Delete deletes a user from the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project-user
func (s *ProjectUserService) Delete(ctx context.Context, projectIDOrKey string, userID int) (*User, error) {
	v, err := s.base.Delete(ctx, projectIDOrKey, userID)
	return userFromModel(v), convertError(err)
}

// AddAdmin adds a admin user to the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project-administrator
func (s *ProjectUserService) AddAdmin(ctx context.Context, projectIDOrKey string, userID int) (*User, error) {
	v, err := s.base.AddAdmin(ctx, projectIDOrKey, userID)
	return userFromModel(v), convertError(err)
}

// AdminAll returns a list of all admin users in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-project-administrators
func (s *ProjectUserService) AdminAll(ctx context.Context, projectIDOrKey string) ([]*User, error) {
	v, err := s.base.AdminAll(ctx, projectIDOrKey)
	return usersFromModel(v), convertError(err)
}

// DeleteAdmin removes an admin user from the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project-administrator
func (s *ProjectUserService) DeleteAdmin(ctx context.Context, projectIDOrKey string, userID int) (*User, error) {
	v, err := s.base.DeleteAdmin(ctx, projectIDOrKey, userID)
	return userFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newProjectCategoryService(method *core.Method) *ProjectCategoryService {
	return &ProjectCategoryService{
		base: project.NewCategoryService(method),
	}
}

func newProjectSharedFileService(method *core.Method) *ProjectSharedFileService {
	return &ProjectSharedFileService{
		base: sharedfile.NewProjectService(method),
	}
}
