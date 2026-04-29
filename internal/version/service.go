package version

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// Service provides Version/Milestone API operations.
type Service struct {
	method *core.Method
}

// All returns versions/milestones in a project.
func (s *Service) All(ctx context.Context, projectIDOrKey string, opts ...core.RequestOption) ([]*model.Version, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	query := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamArchived,
		core.ParamAll,
	}
	if err := core.ApplyOptions(query, validTypes, opts...); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "versions")
	resp, err := s.method.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.Version{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Add adds a version/milestone.
func (s *Service) Add(ctx context.Context, projectIDOrKey, name string, opts ...core.RequestOption) (*model.Version, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, core.NewValidationError("name is required")
	}

	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamName,
		core.ParamDescription,
		core.ParamStartDate,
		core.ParamReleaseDueDate,
	}
	options := append([]core.RequestOption{option.WithName(name)}, opts...)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "versions")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Version{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates a version/milestone.
func (s *Service) Update(ctx context.Context, projectIDOrKey string, versionID int, opts ...core.RequestOption) (*model.Version, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateVersionID(versionID); err != nil {
		return nil, err
	}

	if !core.HasRequiredOption(opts, []core.APIParamOptionType{
		core.ParamName,
		core.ParamDescription,
		core.ParamStartDate,
		core.ParamReleaseDueDate,
		core.ParamArchived,
	}) {
		return nil, core.NewValidationError(
			"requires an option to modify version fields",
		)
	}

	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamName,
		core.ParamDescription,
		core.ParamStartDate,
		core.ParamReleaseDueDate,
		core.ParamArchived,
	}
	if err := core.ApplyOptions(form, validTypes, opts...); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "versions", strconv.Itoa(versionID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Version{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete deletes a version/milestone.
func (s *Service) Delete(ctx context.Context, projectIDOrKey string, versionID int) (*model.Version, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateVersionID(versionID); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "versions", strconv.Itoa(versionID))
	resp, err := s.method.Delete(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Version{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

// NewService returns a Version service.
func NewService(method *core.Method) *Service {
	return &Service{method: method}
}
