package backlog

import (
	"encoding/json"
	"errors"
	"strconv"
)

// ActivityService has methods for Activitys.
type ActivityService struct {
	method *method
	Option *ActivityOptionService
}

func getActivityList(get clientGet, spath string, options ...ActivityOption) ([]*Activity, error) {
	params := newRequestParams()
	for _, option := range options {
		if err := option(params); err != nil {
			return nil, err
		}
	}

	resp, err := get(spath, params)
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

// ProjectActivityService has methods for activitys of the project.
type ProjectActivityService struct {
	method *method
}

// List returns a list of activities in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-recent-updates
func (s *ProjectActivityService) List(projectIDOrKey string, options ...ActivityOption) ([]*Activity, error) {
	if projectIDOrKey == "" {
		return nil, errors.New("projectIDOrKey must not be empty")
	}

	spath := "projects/" + projectIDOrKey + "/activities"
	return getActivityList(s.method.Get, spath, options...)
}

// List returns a list of activities in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-recent-updates
func (s *SpaceActivityService) List(options ...ActivityOption) ([]*Activity, error) {
	spath := "space/activities"
	return getActivityList(s.method.Get, spath, options...)
}

// UserActivityService has methods for user activitys.
type UserActivityService struct {
	method *method
}

// List returns a list of user activities.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-recent-updates
func (s *UserActivityService) List(id int, options ...ActivityOption) ([]*Activity, error) {
	if id < 1 {
		return nil, errors.New("id must be greater than 1")
	}

	spath := "users/" + strconv.Itoa(id) + "/activities"
	return getActivityList(s.method.Get, spath, options...)
}
