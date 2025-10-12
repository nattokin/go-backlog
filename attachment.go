package backlog

import (
	"encoding/json"
	"errors"
	"io"
	"path"
	"strconv"
)

func validateAttachmentID(attachmentID int) error {
	if attachmentID < 1 {
		return newValidationError("attachmentID must not be less than 1")
	}
	return nil
}

// SpaceAttachmentService handles communication with the space attachment-related methods of the Backlog API.
type SpaceAttachmentService struct {
	method *method
}

// Upload uploads any file to the space.
//
// The file name must not be empty.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/post-attachment-file
func (s *SpaceAttachmentService) Upload(fileName string, r io.Reader) (*Attachment, error) {
	resp, err := s.method.Upload("space/attachment", fileName, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := Attachment{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

func listAttachments(m *method, spath string) ([]*Attachment, error) {
	resp, err := m.Get(spath, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := []*Attachment{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}

func removeAttachment(m *method, spath string) (*Attachment, error) {
	resp, err := m.Delete(spath, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := Attachment{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

// WikiAttachmentService handles communication with the wiki attachment-related methods of the Backlog API.
type WikiAttachmentService struct {
	method *method
}

// Attach attaches files uploaded to the space to the specified wiki.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/attach-file-to-wiki
func (s *WikiAttachmentService) Attach(wikiID int, attachmentIDs []int) ([]*Attachment, error) {
	if err := validateWikiID(wikiID); err != nil {
		return nil, err
	}
	if len(attachmentIDs) == 0 {
		return nil, errors.New("attachmentIDs must not be empty")
	}

	form := NewFormParams()
	for _, id := range attachmentIDs {
		if err := validateAttachmentID(id); err != nil {
			return nil, err
		}
		form.Add("attachmentId[]", strconv.Itoa(id))
	}

	spath := path.Join("wikis/", strconv.Itoa(wikiID), "/attachments")
	resp, err := s.method.Post(spath, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := []*Attachment{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}

// List returns a list of all attachments in the wiki.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-wiki-attachments
func (s *WikiAttachmentService) List(wikiID int) ([]*Attachment, error) {
	if err := validateWikiID(wikiID); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "attachments")
	return listAttachments(s.method, spath)
}

// Remove removes a file attached to the wiki.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-wiki-attachment
func (s *WikiAttachmentService) Remove(wikiID, attachmentID int) (*Attachment, error) {
	if err := validateWikiID(wikiID); err != nil {
		return nil, err
	}
	if err := validateAttachmentID(attachmentID); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "attachments", strconv.Itoa(attachmentID))
	return removeAttachment(s.method, spath)
}

// IssueAttachmentService handles communication with the issue attachment-related methods of the Backlog API.
type IssueAttachmentService struct {
	method *method
}

// List returns a list of all attachments in the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-issue-attachments
func (s *IssueAttachmentService) List(issueIDOrKey string) ([]*Attachment, error) {
	if err := validateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "attachments")
	return listAttachments(s.method, spath)
}

// Remove removes a file attached to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue-attachment
func (s *IssueAttachmentService) Remove(issueIDOrKey string, attachmentID int) (*Attachment, error) {
	if err := validateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}
	if err := validateAttachmentID(attachmentID); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "attachments", strconv.Itoa(attachmentID))
	return removeAttachment(s.method, spath)
}

// PullRequestAttachmentService handles communication with the pull request attachment-related methods of the Backlog API.
type PullRequestAttachmentService struct {
	method *method
}

// List returns a list of all attachments in the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-pull-request-attachment
func (s *PullRequestAttachmentService) List(projectIDOrKey string, repositoryIDOrName string, prNumber int) ([]*Attachment, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validateRepositoryIDOrName(repositoryIDOrName); err != nil {
		return nil, err
	}
	if err := validatePRNumber(prNumber); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repositoryIDOrName, "pullRequests", strconv.Itoa(prNumber), "attachments")
	return listAttachments(s.method, spath)
}

// Remove removes a file attached to the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-pull-request-attachments
func (s *PullRequestAttachmentService) Remove(projectIDOrKey string, repositoryIDOrName string, prNumber int, attachmentID int) (*Attachment, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validateRepositoryIDOrName(repositoryIDOrName); err != nil {
		return nil, err
	}
	if err := validatePRNumber(prNumber); err != nil {
		return nil, err
	}
	if err := validateAttachmentID(attachmentID); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repositoryIDOrName, "pullRequests", strconv.Itoa(prNumber), "attachments", strconv.Itoa(attachmentID))
	return removeAttachment(s.method, spath)
}
