package backlog

import (
	"encoding/json"
	"path"
	"strconv"
)

func validateWikiID(wikiID int) error {
	if wikiID < 1 {
		return newValidationError("wikiID must not be less than 1")
	}
	return nil
}

// WikiService handles communication with the wiki-related methods of the Backlog API.
type WikiService struct {
	method *method

	Attachment *WikiAttachmentService
	Option     *WikiOptionService
}

// All returns a list of all wikis in the specified project.
//
// This method supports options returned by methods in "*Client.Wiki.Option".
//
// Use the following methods:
//
//	WithQueryKeyword
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-list
func (s *WikiService) All(projectIDOrKey string, opts ...*QueryOption) ([]*Wiki, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	validOptions := []queryType{queryKeyword}
	for _, option := range opts {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	o := s.Option.support.query
	query := NewQueryParams()
	err := o.applyOptions(query, opts...)
	if err != nil {
		return nil, err
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
func (s *WikiService) Count(projectIDOrKey string) (int, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
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

// One returns a specific wiki by ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page
func (s *WikiService) One(wikiID int) (*Wiki, error) {
	if err := validateWikiID(wikiID); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID))
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
// This method supports options returned by methods in "*Client.Wiki.Option".
//
// Use the following methods:
//
//	WithFormMailNotify
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-wiki-page
func (s *WikiService) Create(projectID int, name, content string, opts ...*FormOption) (*Wiki, error) {
	if err := validateProjectID(projectID); err != nil {
		return nil, err
	}

	validOptions := []formType{formMailNotify}
	for _, option := range opts {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	o := s.Option.support.form
	form := NewFormParams()
	err := o.applyOptions(form, append(opts, o.WithName(name), o.WithContent(content))...)
	if err != nil {
		return nil, err
	}

	form.Set("projectId", strconv.Itoa(projectID))

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

// Update modifies an existing wiki page.
//
// This method requires at least one option to modify the page's name or content.
// The initial option is passed as a mandatory argument (`option`), and any
// additional options are passed via the variadic argument (`opts`).
//
// Internally, the method validates that at least one of WithFormName or WithFormContent is provided.
//
// This method supports options returned by methods in "*Client.Wiki.Option".
//
// Use the following methods:
//
//	WithFormName
//	WithFormContent
//	WithFormMailNotify
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-wiki-page
func (s *WikiService) Update(wikiID int, option *FormOption, opts ...*FormOption) (*Wiki, error) {
	if err := validateWikiID(wikiID); err != nil {
		return nil, err
	}

	options := append([]*FormOption{option}, opts...)

	for _, option := range options {
		if err := option.validate([]formType{formName, formContent, formMailNotify}); err != nil {
			return nil, err
		}
	}

	if !hasRequiredFormOption(options, []formType{formName, formContent}) {
		return nil, newValidationError("requires an option to modify wiki content or name (WithFormName or WithFormContent)")
	}

	o := s.Option.support.form
	form := NewFormParams()
	err := o.applyOptions(form, options...)
	if err != nil {
		return nil, err
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

// Delete removes a wiki by ID.
//
// This method supports options returned by methods in "*Client.Wiki.Option".
//
// Use the following methods:
//
//	WithFormMailNotify
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-wiki-page
func (s *WikiService) Delete(wikiID int, opts ...*FormOption) (*Wiki, error) {
	if err := validateWikiID(wikiID); err != nil {
		return nil, err
	}

	validOptions := []formType{formMailNotify}
	for _, option := range opts {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	o := s.Option.support.form
	form := NewFormParams()
	err := o.applyOptions(form, opts...)
	if err != nil {
		return nil, err
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
