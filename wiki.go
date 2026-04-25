package backlog

import (
	"context"
	"time"

	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/sharedfile"
	"github.com/nattokin/go-backlog/internal/star"
	"github.com/nattokin/go-backlog/internal/wiki"
)

// ──────────────────────────────────────────────────────────────
//  Wiki models
// ──────────────────────────────────────────────────────────────

// Wiki represents Backlog Wiki.
type Wiki struct {
	ID          int
	ProjectID   int
	Name        string
	Content     string
	Tags        []*Tag
	Attachments []*Attachment
	SharedFiles []*SharedFile
	Stars       []*Star
	CreatedUser *User
	Created     time.Time
	UpdatedUser *User
	Updated     time.Time
}

// WikiHistory represents a version history entry for a wiki page.
type WikiHistory struct {
	PageID      int
	Version     int
	Name        string
	Content     string
	CreatedUser *User
	Created     time.Time
}

// ──────────────────────────────────────────────────────────────
//  WikiService
// ──────────────────────────────────────────────────────────────

// WikiService handles communication with the wiki-related methods of the Backlog API.
type WikiService struct {
	base *wiki.Service

	Attachment *WikiAttachmentService
	SharedFile *WikiSharedFileService
	Star       *WikiStarService

	Option *WikiOptionService
}

// All returns a list of all wiki pages in the project.
//
// This method supports options returned by methods in "*Client.Wiki.Option",
// such as:
//   - WithKeyword
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-list
func (s *WikiService) All(ctx context.Context, projectIDOrKey string, opts ...RequestOption) ([]*Wiki, error) {
	v, err := s.base.All(ctx, projectIDOrKey, toCoreOptions(opts)...)
	return wikisFromModel(v), convertError(err)
}

// Count returns the number of wiki pages in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-wiki-page
func (s *WikiService) Count(ctx context.Context, projectIDOrKey string) (int, error) {
	v, err := s.base.Count(ctx, projectIDOrKey)
	return v, convertError(err)
}

// One returns a wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page
func (s *WikiService) One(ctx context.Context, wikiID int) (*Wiki, error) {
	v, err := s.base.One(ctx, wikiID)
	return wikiFromModel(v), convertError(err)
}

// Create creates a new wiki page.
//
// This method supports options returned by methods in "*Client.Wiki.Option",
// such as:
//   - WithMailNotify
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-wiki-page
func (s *WikiService) Create(ctx context.Context, projectID int, name, content string, opts ...RequestOption) (*Wiki, error) {
	v, err := s.base.Create(ctx, projectID, name, content, toCoreOptions(opts)...)
	return wikiFromModel(v), convertError(err)
}

// Update updates a wiki page.
//
// This method supports options returned by methods in "*Client.Wiki.Option",
// such as:
//   - WithContent
//   - WithMailNotify
//   - WithName
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-wiki-page
func (s *WikiService) Update(ctx context.Context, wikiID int, option RequestOption, opts ...RequestOption) (*Wiki, error) {
	v, err := s.base.Update(ctx, wikiID, option, toCoreOptions(opts)...)
	return wikiFromModel(v), convertError(err)
}

// Delete deletes a wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-wiki-page
func (s *WikiService) Delete(ctx context.Context, wikiID int, opts ...RequestOption) (*Wiki, error) {
	v, err := s.base.Delete(ctx, wikiID, toCoreOptions(opts)...)
	return wikiFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  WikiAttachmentService
// ──────────────────────────────────────────────────────────────

// WikiAttachmentService handles communication with the wiki attachment-related methods of the Backlog API.
type WikiAttachmentService struct {
	base *attachment.WikiService
}

// Attach attaches files uploaded to the space to the specified wiki.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/attach-file-to-wiki
func (s *WikiAttachmentService) Attach(ctx context.Context, wikiID int, attachmentIDs []int) ([]*Attachment, error) {
	v, err := s.base.Attach(ctx, wikiID, attachmentIDs)
	return attachmentsFromModel(v), convertError(err)
}

// List returns a list of files attached to the wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-wiki-attachments
func (s *WikiAttachmentService) List(ctx context.Context, wikiID int) ([]*Attachment, error) {
	v, err := s.base.List(ctx, wikiID)
	return attachmentsFromModel(v), convertError(err)
}

// Remove removes an attachment from the wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-wiki-attachment
func (s *WikiAttachmentService) Remove(ctx context.Context, wikiID, attachmentID int) (*Attachment, error) {
	v, err := s.base.Remove(ctx, wikiID, attachmentID)
	return attachmentFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  WikiSharedFileService
// ──────────────────────────────────────────────────────────────

// WikiSharedFileService handles communication with the wiki shared-file-related methods of the Backlog API.
type WikiSharedFileService struct {
	base *sharedfile.WikiService
}

// List returns a list of shared files linked to the wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-shared-files-on-wiki
func (s *WikiSharedFileService) List(ctx context.Context, wikiID int) ([]*SharedFile, error) {
	v, err := s.base.List(ctx, wikiID)
	return sharedFilesFromModel(v), convertError(err)
}

// Link links shared files to the wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/link-shared-files-to-wiki
func (s *WikiSharedFileService) Link(ctx context.Context, wikiID int, fileIDs []int) ([]*SharedFile, error) {
	v, err := s.base.Link(ctx, wikiID, fileIDs)
	return sharedFilesFromModel(v), convertError(err)
}

// Unlink removes a shared file link from the wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-link-to-shared-file-from-wiki
func (s *WikiSharedFileService) Unlink(ctx context.Context, wikiID, fileID int) (*SharedFile, error) {
	v, err := s.base.Unlink(ctx, wikiID, fileID)
	return sharedFileFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  WikiStarService
// ──────────────────────────────────────────────────────────────

// WikiStarService handles communication with the wiki star-related methods of the Backlog API.
type WikiStarService struct {
	base *star.WikiService
	star *StarService
}

// List returns a list of stars on the wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-star
func (s *WikiStarService) List(ctx context.Context, wikiID int) ([]*Star, error) {
	v, err := s.base.List(ctx, wikiID)
	return starsFromModel(v), convertError(err)
}

// Add adds a star to the wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-star
func (s *WikiStarService) Add(ctx context.Context, wikiID int) error {
	return s.star.Add(ctx, s.star.Option.WithWikiID(wikiID))
}

// Remove removes a star by its ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-star
func (s *WikiStarService) Remove(ctx context.Context, starID int) error {
	return s.star.Remove(ctx, starID)
}

// ──────────────────────────────────────────────────────────────
//  WikiOptionService
// ──────────────────────────────────────────────────────────────

// WikiOptionService provides a domain-specific set of option builders
// for operations within the WikiService.
type WikiOptionService struct {
	base *core.OptionService
}

// WithKeyword filters wiki pages by keyword.
func (s *WikiOptionService) WithKeyword(keyword string) RequestOption {
	return s.base.WithKeyword(keyword)
}

// WithContent sets the content of a wiki page.
func (s *WikiOptionService) WithContent(content string) RequestOption {
	return s.base.WithContent(content)
}

// WithMailNotify sets whether to send a mail notification.
func (s *WikiOptionService) WithMailNotify(enabled bool) RequestOption {
	return s.base.WithMailNotify(enabled)
}

// WithName sets the name of a wiki page.
func (s *WikiOptionService) WithName(name string) RequestOption {
	return s.base.WithName(name)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newWikiService(method *core.Method, option *core.OptionService) *WikiService {
	return &WikiService{
		base:       wiki.NewService(method),
		Attachment: newWikiAttachmentService(method),
		Option:     newWikiOptionService(option),
		SharedFile: newWikiSharedFileService(method),
		Star:       newWikiStarService(method, option),
	}
}

func newWikiAttachmentService(method *core.Method) *WikiAttachmentService {
	return &WikiAttachmentService{
		base: attachment.NewWikiService(method),
	}
}

func newWikiSharedFileService(method *core.Method) *WikiSharedFileService {
	return &WikiSharedFileService{
		base: sharedfile.NewWikiService(method),
	}
}

func newWikiStarService(method *core.Method, option *core.OptionService) *WikiStarService {
	return &WikiStarService{
		base: star.NewWikiService(method),
		star: newStarService(method, option),
	}
}

func newWikiOptionService(option *core.OptionService) *WikiOptionService {
	return &WikiOptionService{
		base: option,
	}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func wikiFromModel(m *model.Wiki) *Wiki {
	if m == nil {
		return nil
	}
	tags := make([]*Tag, len(m.Tags))
	for i, v := range m.Tags {
		tags[i] = tagFromModel(v)
	}
	attachments := make([]*Attachment, len(m.Attachments))
	for i, v := range m.Attachments {
		attachments[i] = attachmentFromModel(v)
	}
	sharedFiles := make([]*SharedFile, len(m.SharedFiles))
	for i, v := range m.SharedFiles {
		sharedFiles[i] = sharedFileFromModel(v)
	}
	stars := make([]*Star, len(m.Stars))
	for i, v := range m.Stars {
		stars[i] = starFromModel(v)
	}
	return &Wiki{
		ID:          m.ID,
		ProjectID:   m.ProjectID,
		Name:        m.Name,
		Content:     m.Content,
		Tags:        tags,
		Attachments: attachments,
		SharedFiles: sharedFiles,
		Stars:       stars,
		CreatedUser: userFromModel(m.CreatedUser),
		Created:     m.Created,
		UpdatedUser: userFromModel(m.UpdatedUser),
		Updated:     m.Updated,
	}
}

func wikisFromModel(ms []*model.Wiki) []*Wiki {
	result := make([]*Wiki, len(ms))
	for i, v := range ms {
		result[i] = wikiFromModel(v)
	}
	return result
}
