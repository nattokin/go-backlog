package backlog

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// WikiService has methods for Wiki.
type WikiService struct {
	method *method

	Attachment *WikiAttachmentService
	Option     *WikiOptionService
}

// All Wiki in project is gotten.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-list
func (s *WikiService) All(target ProjectIDOrKeyGetter) ([]*Wiki, error) {
	return s.Search(target, "")
}

// Search returns wikis by keyword from within the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-list
func (s *WikiService) Search(target ProjectIDOrKeyGetter, keyword string) ([]*Wiki, error) {
	params := newRequestParams()
	projectIDOrKey, err := target.getProjectIDOrKey()
	if err != nil {
		return nil, err
	}
	params.Set("projectIdOrKey", projectIDOrKey)
	if keyword != "" {
		params.Set("keyword", keyword)
	}
	resp, err := s.method.Get("wikis", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := []*Wiki{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}

// Count returns the number of wikis in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-wiki-page
func (s *WikiService) Count(target ProjectIDOrKeyGetter) (int, error) {
	params := newRequestParams()
	projectIDOrKey, err := target.getProjectIDOrKey()
	if err != nil {
		return 0, err
	}
	params.Set("projectIdOrKey", projectIDOrKey)
	resp, err := s.method.Get("wikis/count", params)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	v := map[string]int{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}

	return v["count"], nil
}

// One returns one of the wiki by ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page
func (s *WikiService) One(wikiID int) (*Wiki, error) {
	if wikiID < 1 {
		return nil, fmt.Errorf("wikiID must not be less than 1")
	}

	spath := "wikis/" + strconv.Itoa(wikiID)
	resp, err := s.method.Get(spath, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := Wiki{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Create creates a new Wiki for the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-wiki-page
func (s *WikiService) Create(projectID int, name, content string, options ...WikiOption) (*Wiki, error) {
	if projectID < 1 {
		return nil, fmt.Errorf("projectID must not be less than 1")
	}
	if name == "" {
		return nil, errors.New("name must not be empty")
	}
	if content == "" {
		return nil, errors.New("content must not be empty")
	}
	params := newRequestParams()
	params.Set("projectId", strconv.Itoa(projectID))
	params.Set("name", name)
	params.Set("content", content)

	for _, option := range options {
		if err := option(params); err != nil {
			return nil, err
		}
	}

	resp, err := s.method.Post("wikis", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := Wiki{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update a wiki.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-wiki-page
func (s *WikiService) Update(wikiID int, options ...WikiOption) (*Wiki, error) {
	if wikiID < 1 {
		return nil, fmt.Errorf("wikiID must not be less than 1")
	}

	if options == nil {
		return nil, errors.New("requires one or more options")
	}

	params := newRequestParams()
	for _, option := range options {
		if err := option(params); err != nil {
			return nil, err
		}
	}

	spath := "wikis/" + strconv.Itoa(wikiID)
	resp, err := s.method.Patch(spath, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := Wiki{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete a wiki by ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-wiki-page
func (s *WikiService) Delete(wikiID int, options ...WikiOption) (*Wiki, error) {
	if wikiID < 1 {
		return nil, fmt.Errorf("wikiID must not be less than 1")
	}
	params := newRequestParams()
	for _, option := range options {
		if err := option(params); err != nil {
			return nil, err
		}
	}

	spath := "wikis/" + strconv.Itoa(wikiID)
	resp, err := s.method.Delete(spath, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := Wiki{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}
