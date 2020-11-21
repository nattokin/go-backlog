package backlog

import (
	"encoding/json"
	"errors"
	"io"
	"path"
	"strconv"
)

// AttachmentID implements IssueIDOrKeyGetter interface.
type AttachmentID int

func (i AttachmentID) validate() error {
	if i < 1 {
		return errors.New("attachmentID must not be less than 1")
	}
	return nil
}

func (i AttachmentID) String() string {
	return strconv.Itoa(int(i))
}

// SpaceAttachmentService hs methods for attachment.
type SpaceAttachmentService struct {
	method *method
}

// Upload uploads a any file to the space.
//
// File's path and name are must not empty.
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

func listAttachments(get clientGet, spath string) ([]*Attachment, error) {
	resp, err := get(spath, nil)
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

func removeAttachment(delete clientDelete, spath string) (*Attachment, error) {
	resp, err := delete(spath, nil)
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

// WikiAttachmentService hs methods for attachment file of wiki.
type WikiAttachmentService struct {
	method *method
}

// Attach attachs files uploaded to space to the wiki.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/attach-file-to-wiki
func (s *WikiAttachmentService) Attach(wikiID int, attachmentIDs []int) ([]*Attachment, error) {
	wID := WikiID(wikiID)
	if err := wID.validate(); err != nil {
		return nil, err
	}
	if len(attachmentIDs) == 0 {
		return nil, errors.New("attachmentIDs is must not empty")
	}

	form := NewFormParams()
	for _, id := range attachmentIDs {
		aID := AttachmentID(id)
		if err := aID.validate(); err != nil {
			return nil, err
		}
		form.Add("attachmentId[]", aID.String())
	}

	spath := path.Join("wikis/", wID.String(), "/attachments")
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
	wID := WikiID(wikiID)
	if err := wID.validate(); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", wID.String(), "attachments")
	return listAttachments(s.method.Get, spath)
}

// Remove removes a file attached to the wiki.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-wiki-attachment
func (s *WikiAttachmentService) Remove(wikiID, attachmentID int) (*Attachment, error) {
	wID := WikiID(wikiID)
	if err := wID.validate(); err != nil {
		return nil, err
	}
	aID := AttachmentID(attachmentID)
	if err := aID.validate(); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", wID.String(), "attachments", aID.String())
	return removeAttachment(s.method.Delete, spath)
}

// IssueAttachmentService hs methods for attachment file of issue.
type IssueAttachmentService struct {
	method *method
}

// List returns a list of all attachments in the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-issue-attachments
func (s *IssueAttachmentService) List(issue IssueIDOrKeyGetter) ([]*Attachment, error) {
	issueIDOrKey, err := issue.getIssueIDOrKey()
	if err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "attachments")
	return listAttachments(s.method.Get, spath)
}

// Remove removes a file attached to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue-attachment
func (s *IssueAttachmentService) Remove(issue IssueIDOrKeyGetter, attachmentID int) (*Attachment, error) {
	issueIDOrKey, err := issue.getIssueIDOrKey()
	if err != nil {
		return nil, err
	}
	aID := AttachmentID(attachmentID)
	if err := aID.validate(); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "attachments", aID.String())
	return removeAttachment(s.method.Delete, spath)
}

// PullRequestAttachmentService hs methods for attachment file of pull request.
type PullRequestAttachmentService struct {
	method *method
}

// List returns a list of all attachments in the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-pull-request-attachment
func (s *PullRequestAttachmentService) List(project ProjectIDOrKeyGetter, repository RepositoryIDOrKeyGetter, prNumber int) ([]*Attachment, error) {
	projectIDOrKey, err := project.getProjectIDOrKey()
	if err != nil {
		return nil, err
	}
	repoIDOrName, err := repository.getRepositoryIDOrKey()
	if err != nil {
		return nil, err
	}
	prNum := PRNumber(prNumber)
	if err := prNum.validate(); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", prNum.String(), "attachments")
	return listAttachments(s.method.Get, spath)
}

// Remove removes a file attached to the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-pull-request-attachments
func (s *PullRequestAttachmentService) Remove(project ProjectIDOrKeyGetter, repository RepositoryIDOrKeyGetter, prNumber int, attachmentID int) (*Attachment, error) {
	projectIDOrKey, err := project.getProjectIDOrKey()
	if err != nil {
		return nil, err
	}
	repoIDOrName, err := repository.getRepositoryIDOrKey()
	if err != nil {
		return nil, err
	}
	prNum := PRNumber(prNumber)
	if err := prNum.validate(); err != nil {
		return nil, err
	}
	aID := AttachmentID(attachmentID)
	if err := aID.validate(); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", prNum.String(), "attachments", aID.String())
	return removeAttachment(s.method.Delete, spath)
}
