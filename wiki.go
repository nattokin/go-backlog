package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/wiki"
)

// ──────────────────────────────────────────────────────────────
//  WikiService
// ──────────────────────────────────────────────────────────────

// WikiService handles communication with the wiki-related methods of the Backlog API.
type WikiService struct {
	base *wiki.WikiService

	Attachment *WikiAttachmentService
	Option     *WikiOptionService
}

// All returns a list of all wikis in the specified project.
//
// This method supports options returned by methods in "*Client.Wiki.Option",
// such as:
//   - WithKeyword
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-list
func (s *WikiService) All(ctx context.Context, projectIDOrKey string, opts ...core.RequestOption) ([]*model.Wiki, error) {
	return s.base.All(ctx, projectIDOrKey, opts...)
}

// Count returns the number of wikis in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-wiki-page
func (s *WikiService) Count(ctx context.Context, projectIDOrKey string) (int, error) {
	return s.base.Count(ctx, projectIDOrKey)
}

// One returns a specific wiki by ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page
func (s *WikiService) One(ctx context.Context, wikiID int) (*model.Wiki, error) {
	return s.base.One(ctx, wikiID)
}

// Create creates a new Wiki for the project.
//
// This method supports options returned by methods in "*Client.Wiki.Option",
// such as:
//   - WithContent
//   - WithMailNotify
//   - WithName
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-wiki-page
func (s *WikiService) Create(ctx context.Context, projectID int, name, content string, opts ...core.RequestOption) (*model.Wiki, error) {
	return s.base.Create(ctx, projectID, name, content, opts...)
}

// Update modifies an existing wiki page.
//
// This method supports options returned by methods in "*Client.Wiki.Option",
// such as:
//   - WithContent
//   - WithMailNotify
//   - WithName
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-wiki-page
func (s *WikiService) Update(ctx context.Context, wikiID int, option core.RequestOption, opts ...core.RequestOption) (*model.Wiki, error) {
	return s.base.Update(ctx, wikiID, option, opts...)
}

// Delete removes a wiki by ID.
//
// This method supports options returned by methods in "*Client.Wiki.Option",
// such as:
//   - WithMailNotify
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-wiki-page
func (s *WikiService) Delete(ctx context.Context, wikiID int, opts ...core.RequestOption) (*model.Wiki, error) {
	return s.base.Delete(ctx, wikiID, opts...)
}

// ──────────────────────────────────────────────────────────────
//  WikiAttachmentService
// ──────────────────────────────────────────────────────────────

// WikiAttachmentService handles communication with the wiki attachment-related methods of the Backlog API.
type WikiAttachmentService struct {
	base *attachment.WikiAttachmentService
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newWikiService(method *core.Method, option *core.OptionService) *WikiService {
	return &WikiService{
		base:       wiki.NewWikiService(method, option),
		Attachment: newWikiAttachmentService(method),
		Option:     newWikiOptionService(option),
	}
}

func newWikiAttachmentService(method *core.Method) *WikiAttachmentService {
	return &WikiAttachmentService{
		base: attachment.NewWikiAttachmentService(method),
	}
}
