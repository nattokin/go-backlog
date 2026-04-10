package backlog

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
)

func validateWikiID(wikiID int) error {
	if wikiID < 1 {
		return newValidationError("wikiID must not be less than 1")
	}
	return nil
}

// WikiService handles communication with the wiki-related methods of the Backlog API.
type WikiService struct {
	method *core.Method

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
func (s *WikiService) All(ctx context.Context, projectIDOrKey string, opts ...RequestOption) ([]*Wiki, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	query := url.Values{}
	validTypes := []apiParamOptionType{paramKeyword}
	if err := applyOptions(query, validTypes, opts...); err != nil {
		return nil, err
	}

	query.Set("projectIdOrKey", projectIDOrKey)

	resp, err := s.method.Get(ctx, "wikis", query)
	if err != nil {
		return nil, err
	}

	v := []*Wiki{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Count returns the number of wikis in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-wiki-page
func (s *WikiService) Count(ctx context.Context, projectIDOrKey string) (int, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return 0, err
	}

	query := url.Values{}
	query.Set("projectIdOrKey", projectIDOrKey)

	resp, err := s.method.Get(ctx, "wikis/count", query)
	if err != nil {
		return 0, err
	}

	v := map[string]int{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return 0, err
	}

	return v["count"], nil
}

// One returns a specific wiki by ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page
func (s *WikiService) One(ctx context.Context, wikiID int) (*Wiki, error) {
	if err := validateWikiID(wikiID); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID))
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := Wiki{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
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
func (s *WikiService) Create(ctx context.Context, projectID int, name, content string, opts ...RequestOption) (*Wiki, error) {
	if err := validateProjectID(projectID); err != nil {
		return nil, err
	}

	form := url.Values{}
	validTypes := []apiParamOptionType{paramName, paramContent, paramMailNotify}
	options := append([]RequestOption{s.Option.base.WithName(name), s.Option.base.WithContent(content)}, opts...)
	if err := applyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	form.Set("projectId", strconv.Itoa(projectID))

	resp, err := s.method.Post(ctx, "wikis", form)
	if err != nil {
		return nil, err
	}

	v := Wiki{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
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
func (s *WikiService) Update(ctx context.Context, wikiID int, option RequestOption, opts ...RequestOption) (*Wiki, error) {
	if err := validateWikiID(wikiID); err != nil {
		return nil, err
	}

	form := url.Values{}
	validTypes := []apiParamOptionType{paramName, paramContent, paramMailNotify}
	options := append([]RequestOption{option}, opts...)

	if !hasRequiredOption(options, []apiParamOptionType{paramName, paramContent}) {
		return nil, newValidationError("requires an option to modify wiki content or name (WithName or WithContent)")
	}

	if err := applyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := Wiki{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete removes a wiki by ID.
//
// This method supports options returned by methods in "*Client.Wiki.Option",
// such as:
//   - WithMailNotify
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-wiki-page
func (s *WikiService) Delete(ctx context.Context, wikiID int, opts ...RequestOption) (*Wiki, error) {
	if err := validateWikiID(wikiID); err != nil {
		return nil, err
	}

	form := url.Values{}
	validTypes := []apiParamOptionType{paramMailNotify}
	if err := applyOptions(form, validTypes, opts...); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID))
	resp, err := s.method.Delete(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := Wiki{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}
