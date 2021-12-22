package backlog

import (
	"encoding/json"
	"errors"
	"path"
	"strconv"
)

// WikiID is ID of Wiki.
type WikiID int

func (i WikiID) validate() error {
	if i < 1 {
		return newValidationError("wikiID must not be less than 1")
	}
	return nil
}

func (i WikiID) String() string {
	return strconv.Itoa(int(i))
}

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
func (s *WikiService) All(project ProjectIDOrKeyGetter, options ...*QueryOption) ([]*Wiki, error) {
	projectIDOrKey, err := project.getProjectIDOrKey()
	if err != nil {
		return nil, err
	}

	validOptions := []queryType{queryKeyword}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	query := NewQueryParams()
	for _, option := range options {
		if err := option.set(query); err != nil {
			return nil, err
		}
	}
	query.Set("projectIdOrKey", projectIDOrKey)

	resp, err := s.method.Get("wikis", query)
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
func (s *WikiService) Count(project ProjectIDOrKeyGetter) (int, error) {
	projectIDOrKey, err := project.getProjectIDOrKey()
	if err != nil {
		return 0, err
	}

	query := NewQueryParams()
	query.Set("projectIdOrKey", projectIDOrKey)

	resp, err := s.method.Get("wikis/count", query)
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
	wID := WikiID(wikiID)
	if err := wID.validate(); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", wID.String())
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
	pID := ProjectID(projectID)
	if err := pID.validate(); err != nil {
		return nil, err
	}

	form := NewFormParams()
	if err := withFormName(name).set(form); err != nil {
		return nil, err
	}
	if err := withFormContent(content).set(form); err != nil {
		return nil, err
	}

	validOptions := []formType{formMailNotify}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	for _, option := range options {
		if err := option.set(form); err != nil {
			return nil, err
		}
	}
	form.Set("projectId", pID.String())

	resp, err := s.method.Post("wikis", form)
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
	wID := WikiID(wikiID)
	if err := wID.validate(); err != nil {
		return nil, err
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

	form := NewFormParams()
	for _, option := range options {
		if err := option.set(form); err != nil {
			return nil, err
		}
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID))
	resp, err := s.method.Patch(spath, form)
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
	wID := WikiID(wikiID)
	if err := wID.validate(); err != nil {
		return nil, err
	}

	validOptions := []formType{formMailNotify}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	form := NewFormParams()
	for _, option := range options {
		if err := option.set(form); err != nil {
			return nil, err
		}
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID))
	resp, err := s.method.Delete(spath, form)
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
