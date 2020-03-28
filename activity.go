package backlog

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// ActivityOptionService has methods to make functional option for ActivityService.
type ActivityOptionService struct {
}

// ActivityOption is type of functional option for ActivityService.
type ActivityOption func(p *requestParams) error

// WithActivityTypeIDs returns option. the option sets `activityTypeId` for user.
func (*ActivityOptionService) WithActivityTypeIDs(typeIDs []int) ActivityOption {
	return func(p *requestParams) error {
		for _, id := range typeIDs {
			if id < 1 || 26 < id {
				return errors.New("activityTypeId must be between 1 and 26")
			}
			p.Add("activityTypeId[]", strconv.Itoa(id))
		}
		return nil
	}
}

// WithMinID returns option. the option sets `minId` for user.
func (*ActivityOptionService) WithMinID(minID int) ActivityOption {
	return func(p *requestParams) error {
		if minID < 1 {
			return errors.New("minId must be greater than 1")
		}
		p.Set("minId", strconv.Itoa(minID))
		return nil
	}
}

// WithMaxID returns option. the option sets `maxId` for user.
func (*ActivityOptionService) WithMaxID(maxID int) ActivityOption {
	return func(p *requestParams) error {
		if maxID < 1 {
			return errors.New("maxId must be greater than 1")
		}
		p.Set("maxId", strconv.Itoa(maxID))
		return nil
	}
}

// WithCount returns option. the option sets `count` for user.
func (*ActivityOptionService) WithCount(count int) ActivityOption {
	return func(p *requestParams) error {
		if count < 1 || 100 < count {
			return errors.New("count must be between 1 and 100")
		}
		p.Set("count", strconv.Itoa(count))
		return nil
	}
}

// WithOrder returns option. the option sets `order` for user.
func (*ActivityOptionService) WithOrder(order string) ActivityOption {
	return func(p *requestParams) error {
		if order != OrderAsc && order != OrderDesc {
			return fmt.Errorf("order must be only '%s' or '%s'", OrderAsc, OrderDesc)
		}
		p.Set("order", order)
		return nil
	}
}

type baseActivityService struct {
	clientMethod *clientMethod

	Option *ActivityOptionService
}

func newBaseActivityService(cm *clientMethod) *baseActivityService {
	return &baseActivityService{
		clientMethod: cm,
		Option:       &ActivityOptionService{},
	}
}

func (s *baseActivityService) getList(spath string, options ...ActivityOption) ([]*Activity, error) {
	params := newRequestParams()
	for _, option := range options {
		if err := option(params); err != nil {
			return nil, err
		}
	}

	resp, err := s.clientMethod.Get(spath, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := []*Activity{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}

// ActivityService has methods for Activitys.
type ActivityService struct {
	*baseActivityService
}

func newActivityService(cm *clientMethod) *ActivityService {
	return &ActivityService{
		baseActivityService: newBaseActivityService(cm),
	}
}

// ProjectActivityService has methods for activitys of the project.
type ProjectActivityService struct {
	*baseActivityService
}

func newProjectActivityService(cm *clientMethod) *ProjectActivityService {
	return &ProjectActivityService{
		baseActivityService: newBaseActivityService(cm),
	}
}

// List returns a list of activities in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-recent-updates
func (s *ProjectActivityService) List(projectIDOrKey string, options ...ActivityOption) ([]*Activity, error) {
	if projectIDOrKey == "" {
		return nil, errors.New("projectIDOrKey must not be empty")
	}

	spath := "projects/" + projectIDOrKey + "/activities"
	return s.getList(spath, options...)
}

// SpaceActivityService has methods for activitys in your space.
type SpaceActivityService struct {
	*baseActivityService
}

func newSpaceActivityService(cm *clientMethod) *SpaceActivityService {
	return &SpaceActivityService{
		baseActivityService: newBaseActivityService(cm),
	}
}

// List returns a list of activities in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-recent-updates
func (s *SpaceActivityService) List(options ...ActivityOption) ([]*Activity, error) {
	spath := "space/activities"
	return s.getList(spath, options...)
}

// UserActivityService has methods for user activitys.
type UserActivityService struct {
	*baseActivityService
}

func newUserActivityService(cm *clientMethod) *UserActivityService {
	return &UserActivityService{
		baseActivityService: newBaseActivityService(cm),
	}
}

// List returns a list of user activities.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-recent-updates
func (s *UserActivityService) List(id int, options ...ActivityOption) ([]*Activity, error) {
	if id < 1 {
		return nil, errors.New("id must be greater than 1")
	}

	spath := "users/" + strconv.Itoa(id) + "/activities"
	return s.getList(spath, options...)
}
