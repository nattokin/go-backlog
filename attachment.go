package backlog

import (
	"encoding/json"
	"errors"
	"strconv"
)

// SpaceAttachmentService hs methods for attachment.
type SpaceAttachmentService struct {
	method *method
}

// Upload uploads a any file to the space.
//
// File's path and name are must not empty.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/post-attachment-file
func (s *SpaceAttachmentService) Upload(fpath, fname string) (*Attachment, error) {
	spath := "space/attachment"
	resp, err := s.method.Upload(spath, fpath, fname)
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
	params := NewFormParams()
	for _, id := range attachmentIDs {
		params.Add("attachmentId[]", strconv.Itoa(id))
	}

	spath := "wikis/" + strconv.Itoa(wikiID) + "/attachments"
	resp, err := s.method.Post(spath, params)
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
	if wikiID < 1 {
		return nil, errors.New("wikiID must not be less than 1")
	}

	spath := "wikis/" + strconv.Itoa(wikiID) + "/attachments"
	return listAttachments(s.method.Get, spath)
}

// Remove removes a file attached to the wiki.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-wiki-attachment
func (s *WikiAttachmentService) Remove(wikiID, attachmentID int) (*Attachment, error) {
	if wikiID < 1 {
		return nil, errors.New("wikiID must not be less than 1")
	}
	if attachmentID < 1 {
		return nil, errors.New("attachmentID must not be less than 1")
	}

	spath := "wikis/" + strconv.Itoa(wikiID) + "/attachments/" + strconv.Itoa(attachmentID)
	return removeAttachment(s.method.Delete, spath)
}

// IssueAttachmentService hs methods for attachment file of issue.
type IssueAttachmentService struct {
	method *method
}

// List returns a list of all attachments in the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-issue-attachments
func (s *IssueAttachmentService) List(target IssueIDOrKeyGetter) ([]*Attachment, error) {
	issueIDOrKey, err := target.getIssueIDOrKey()
	if err != nil {
		return nil, err
	}

	spath := "issues/" + issueIDOrKey + "/attachments"
	return listAttachments(s.method.Get, spath)
}

// Remove removes a file attached to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue-attachment
func (s *IssueAttachmentService) Remove(target IssueIDOrKeyGetter, attachmentID int) (*Attachment, error) {
	issueIDOrKey, err := target.getIssueIDOrKey()
	if err != nil {
		return nil, err
	}
	if attachmentID < 1 {
		return nil, errors.New("attachmentID must not be less than 1")
	}

	spath := "issues/" + issueIDOrKey + "/attachments/" + strconv.Itoa(attachmentID)
	return removeAttachment(s.method.Delete, spath)
}

// PullRequestAttachmentService hs methods for attachment file of pull request.
type PullRequestAttachmentService struct {
	method *method
}

// List returns a list of all attachments in the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-pull-request-attachment
func (s *PullRequestAttachmentService) List(targetProject ProjectIDOrKeyGetter, targetRepository RepositoryIDOrKeyGetter, prNumber int) ([]*Attachment, error) {
	projectIDOrKey, err := targetProject.getProjectIDOrKey()
	if err != nil {
		return nil, err
	}
	repoIDOrName, err := targetRepository.getRepositoryIDOrKey()
	if err != nil {
		return nil, err
	}
	if prNumber < 1 {
		return nil, errors.New("prNumber must not be less than 1")
	}

	spath := "projects/" + projectIDOrKey + "/git/repositories/" + repoIDOrName + "/pullRequests/" + strconv.Itoa(prNumber) + "/attachments"
	return listAttachments(s.method.Get, spath)
}

// Remove removes a file attached to the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-pull-request-attachments
func (s *PullRequestAttachmentService) Remove(targetProject ProjectIDOrKeyGetter, targetRepository RepositoryIDOrKeyGetter, prNumber int, attachmentID int) (*Attachment, error) {
	projectIDOrKey, err := targetProject.getProjectIDOrKey()
	if err != nil {
		return nil, err
	}
	repoIDOrName, err := targetRepository.getRepositoryIDOrKey()
	if err != nil {
		return nil, err
	}
	if prNumber < 1 {
		return nil, errors.New("prNumber must not be less than 1")
	}
	if attachmentID < 1 {
		return nil, errors.New("attachmentID must not be less than 1")
	}

	spath := "projects/" + projectIDOrKey + "/git/repositories/" + repoIDOrName + "/pullRequests/" + strconv.Itoa(prNumber) + "/attachments" + strconv.Itoa(attachmentID)
	return removeAttachment(s.method.Delete, spath)
}
