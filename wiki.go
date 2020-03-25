package backlog

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// WikiService has methods for Wiki.
type WikiService struct {
	clientMethod *clientMethod

	Attachment *WikiAttachmentService
	option     *wikiOptionService
}

func newWikiService(cm *clientMethod) *WikiService {
	as := &WikiAttachmentService{
		baseAttachmentService: &baseAttachmentService{
			clientMethod: cm,
		},
	}
	return &WikiService{
		clientMethod: cm,
		Attachment:   as,
		option:       &wikiOptionService{},
	}
}

// wikiOptionService has methods to make functional option for WikiService.
type wikiOptionService struct {
}

// wikiOption is type of functional option for WikiService.
type wikiOption func(p *requestParams) error

// WithName returns option. the option sets `name` for wiki.
func (*wikiOptionService) WithName(name string) wikiOption {
	return func(p *requestParams) error {
		if name == "" {
			return errors.New("[*wikiOptionService.WithName] name must not be empty")
		}
		p.Set("name", name)
		return nil
	}
}

// WithContent returns option. the option sets `content` for wiki.
func (*wikiOptionService) WithContent(content string) wikiOption {
	return func(p *requestParams) error {
		if content == "" {
			return errors.New("[*wikiOptionService.WithContent] content must not be empty")
		}
		p.Set("content", content)
		return nil
	}
}

// WithMailNotify returns option. the option sets `mailNotify` true for wiki.
func (*wikiOptionService) WithMailNotify() wikiOption {
	return func(p *requestParams) error {
		p.Set("mailNotify", "true")
		return nil
	}
}

// All Wiki in project is gotten.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-list
func (s *WikiService) All(projectIDOrKey string) ([]*Wiki, error) {
	return s.Search(projectIDOrKey, "")
}

// Search returns wikis by keyword from within the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-list
func (s *WikiService) Search(projectIDOrKey, keyword string) ([]*Wiki, error) {
	params := newRequestParams()
	params.Set("projectIDOrKey", projectIDOrKey)
	if keyword != "" {
		params.Set("keyword", keyword)
	}
	resp, err := s.clientMethod.Get("wikis", params)
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
	params := newRequestParams()
	params.Set("projectIDOrKey", projectIDOrKey)
	resp, err := s.clientMethod.Get("wikis/count", params)
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
	spath := "wikis/" + strconv.Itoa(wikiID)
	resp, err := s.clientMethod.Get(spath, nil)
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
func (s *WikiService) Create(projectID int, name, content string, mailNotify bool) (*Wiki, error) {
	if mailNotify {
		return s.create(projectID, name, content, s.option.WithMailNotify())
	}
	return s.create(projectID, name, content)
}

func (s *WikiService) create(projectID int, name, content string, options ...wikiOption) (*Wiki, error) {
	if projectID == 0 {
		return nil, errors.New("projectID must not be zero")
	}
	if name == "" {
		return nil, errors.New("name is requierd")
	}
	if content == "" {
		return nil, errors.New("content is requierd")
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

	resp, err := s.clientMethod.Post("wikis", params)
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
func (s *WikiService) Update(wikiID int, name, content string, mailNotify bool) (*Wiki, error) {
	options := []wikiOption{}
	if name != "" {
		options = append(options, s.option.WithName(name))
	}
	if content != "" {
		options = append(options, s.option.WithContent(content))
	}
	if mailNotify {
		options = append(options, s.option.WithMailNotify())
	}
	return s.update(wikiID, options...)
}

func (s *WikiService) update(wikiID int, options ...wikiOption) (*Wiki, error) {
	if wikiID <= 0 {
		return nil, fmt.Errorf("wikiID must be 1 or more: %d", wikiID)
	}

	params := newRequestParams()
	for _, option := range options {
		if err := option(params); err != nil {
			return nil, err
		}
	}

	spath := "wikis/" + strconv.Itoa(wikiID)
	resp, err := s.clientMethod.Patch(spath, params)
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
func (s *WikiService) Delete(wikiID int, mailNotify bool) (*Wiki, error) {
	if mailNotify {
		return s.delete(wikiID, s.option.WithMailNotify())
	}
	return s.delete(wikiID)
}

func (s *WikiService) delete(wikiID int, options ...wikiOption) (*Wiki, error) {
	if wikiID <= 0 {
		return nil, fmt.Errorf("wikiID must be 1 or more: %d", wikiID)
	}
	params := newRequestParams()
	for _, option := range options {
		if err := option(params); err != nil {
			return nil, err
		}
	}

	spath := "wikis/" + strconv.Itoa(wikiID)
	resp, err := s.clientMethod.Delete(spath, params)
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
