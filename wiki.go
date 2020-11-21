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
// This method can specify the options returned by methods in "*Client.Wiki.Option".
//
// Use the following methods:
//   WithQueryKeyword
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-list
func (s *WikiService) All(target ProjectIDOrKeyGetter, options ...*QueryOption) ([]*Wiki, error) {
	projectIDOrKey, err := target.getProjectIDOrKey()
	if err != nil {
		return nil, err
	}

	validOptions := []queryType{queryKeyword}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	params := NewQueryParams()
	for _, option := range options {
		if err := option.set(params); err != nil {
			return nil, err
		}
	}
	params.Set("projectIdOrKey", projectIDOrKey)

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
	projectIDOrKey, err := target.getProjectIDOrKey()
	if err != nil {
		return 0, err
	}

	params := NewQueryParams()
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
// This method can specify the options returned by methods in "*Client.Wiki.Option".
//
// Use the following methods:
//   WithFormMailNotify
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-wiki-page
func (s *WikiService) Create(projectID int, name, content string, options ...*FormOption) (*Wiki, error) {
	if projectID < 1 {
		return nil, fmt.Errorf("projectID must not be less than 1")
	}
	if name == "" {
		return nil, errors.New("name must not be empty")
	}
	if content == "" {
		return nil, errors.New("content must not be empty")
	}

	validOptions := []formType{formMailNotify}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	params := NewFormParams()
	for _, option := range options {
		if err := option.set(params); err != nil {
			return nil, err
		}
	}
	params.Set("projectId", strconv.Itoa(projectID))
	params.Set("name", name)
	params.Set("content", content)

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
// This method can specify the options returned by methods in "*Client.Wiki.Option".
//
// Use the following methods:
//   WithFormName
//   WithFormContent
//   WithFormMailNotify
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-wiki-page
func (s *WikiService) Update(wikiID int, options ...*FormOption) (*Wiki, error) {
	if wikiID < 1 {
		return nil, fmt.Errorf("wikiID must not be less than 1")
	}

	if options == nil {
		return nil, errors.New("requires one or more options")
	}

	validOptions := []formType{formName, formContent, formMailNotify}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	params := NewFormParams()
	for _, option := range options {
		if err := option.set(params); err != nil {
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
// This method can specify the options returned by methods in "*Client.Wiki.Option".
//
// Use the following methods:
//   WithFormMailNotify
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-wiki-page
func (s *WikiService) Delete(wikiID int, options ...*FormOption) (*Wiki, error) {
	if wikiID < 1 {
		return nil, fmt.Errorf("wikiID must not be less than 1")
	}

	validOptions := []formType{formMailNotify}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	params := NewFormParams()
	for _, option := range options {
		if err := option.set(params); err != nil {
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
