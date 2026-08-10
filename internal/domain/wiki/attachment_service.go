package wiki

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/shared/attachment"
	"github.com/nattokin/go-backlog/internal/validate"
)

// AttachmentService handles attachment-related Backlog API calls for wiki pages.
// It delegates all HTTP operations to the shared attachment.Service and is
// responsible only for validation and spath construction.
type AttachmentService struct {
	base   *attachment.Service
	method *core.Method
}

// Attach attaches files uploaded to the space to the specified wiki.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/attach-file-to-wiki
func (s *AttachmentService) Attach(ctx context.Context, wikiID int, attachmentIDs []int) ([]*model.Attachment, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateWikiID(wikiID); ve != nil {
		ves = append(ves, ve)
	}
	if len(attachmentIDs) == 0 {
		ves = append(ves, core.NewValidationError("attachmentIDs", "attachmentIDs must not be empty"))
	} else {
		for _, id := range attachmentIDs {
			if id <= 0 {
				ves = append(ves, core.NewValidationError("attachmentIDs", "attachmentID must be greater than 0"))
				break
			}
		}
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "attachments")
	form := url.Values{}
	for _, id := range attachmentIDs {
		form.Add("attachmentId[]", strconv.Itoa(id))
	}

	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := []*model.Attachment{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// List returns a list of files attached to the wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-wiki-attachments
func (s *AttachmentService) List(ctx context.Context, wikiID int) ([]*model.Attachment, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateWikiID(wikiID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "attachments")
	return s.base.List(ctx, spath)
}

// Remove removes an attachment from the wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-wiki-attachment
func (s *AttachmentService) Remove(ctx context.Context, wikiID, attachmentID int) (*model.Attachment, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateWikiID(wikiID); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateAttachmentID(attachmentID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "attachments", strconv.Itoa(attachmentID))
	return s.base.Remove(ctx, spath)
}

// Download downloads a file attached to the wiki page.
// The caller is responsible for closing FileData.Body after use.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-attachment
func (s *AttachmentService) Download(ctx context.Context, wikiID, attachmentID int) (*model.FileData, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateWikiID(wikiID); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateAttachmentID(attachmentID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "attachments", strconv.Itoa(attachmentID))
	return s.base.Download(ctx, spath)
}

func NewAttachmentService(method *core.Method) *AttachmentService {
	return &AttachmentService{
		base:   attachment.NewService(method),
		method: method,
	}
}
