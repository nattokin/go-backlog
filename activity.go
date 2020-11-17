package backlog

import (
	"encoding/json"
	"errors"
	"strconv"
)

func getActivityList(get clientGet, spath string, options ...*ActivityOption) ([]*Activity, error) {
	validOptions := []optionType{optionActivityTypeIDs, optionMinID, optionMaxID, optionCount, optionOrder}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	params := newRequestParams()
	for _, option := range options {
		if err := option.set(params); err != nil {
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

	Option *ActivityOptionService
}

// List returns a list of activities in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-recent-updates
func (s *ProjectActivityService) List(target ProjectIDOrKeyGetter, options ...*ActivityOption) ([]*Activity, error) {
	projectIDOrKey, err := target.getProjectIDOrKey()
	if err != nil {
		return nil, err
	}

	spath := "projects/" + projectIDOrKey + "/activities"
	return getActivityList(s.method.Get, spath, options...)
}

// SpaceActivityService has methods for activitys in your space.
type SpaceActivityService struct {
	method *method

	Option *ActivityOptionService
}

// List returns a list of activities in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-recent-updates
func (s *SpaceActivityService) List(options ...*ActivityOption) ([]*Activity, error) {
	spath := "space/activities"
	return getActivityList(s.method.Get, spath, options...)
}

// UserActivityService has methods for user activitys.
type UserActivityService struct {
	method *method

	Option *ActivityOptionService
}

// List returns a list of user activities.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-recent-updates
func (s *UserActivityService) List(userID int, options ...*ActivityOption) ([]*Activity, error) {
	if userID < 1 {
		return nil, errors.New("userID must be greater than 1")
	}

	spath := "users/" + strconv.Itoa(userID) + "/activities"
	return getActivityList(s.method.Get, spath, options...)
}
