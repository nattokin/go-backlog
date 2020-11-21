package backlog

import (
	"encoding/json"
	"errors"
	"strconv"
)

func getActivityList(get clientGet, spath string, options ...*QueryOption) ([]*Activity, error) {
	validOptions := []queryType{queryActivityTypeIDs, queryMinID, queryMaxID, queryCount, queryOrder}
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
// This method can specify the options returned by methods in "*Client.Project.Activity.Option".
//
// Use the following methods:
//   WithQueryActivityTypeIDs
//   WithQueryMinID
//   WithQueryMaxID
//   WithQueryCount
//   WithQueryOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-recent-updates
func (s *ProjectActivityService) List(target ProjectIDOrKeyGetter, options ...*QueryOption) ([]*Activity, error) {
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
// This method can specify the options returned by methods in "*Client.Space.Activity.Option".
//
// Use the following methods:
//   WithQueryActivityTypeIDs
//   WithQueryMinID
//   WithQueryMaxID
//   WithQueryCount
//   WithQueryOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-recent-updates
func (s *SpaceActivityService) List(options ...*QueryOption) ([]*Activity, error) {
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
// This method can specify the options returned by methods in "*Client.User.Activity.Option".
//
// Use the following methods:
//   WithQueryActivityTypeIDs
//   WithQueryMinID
//   WithQueryMaxID
//   WithQueryCount
//   WithQueryOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-recent-updates
func (s *UserActivityService) List(userID int, options ...*QueryOption) ([]*Activity, error) {
	if userID < 1 {
		return nil, errors.New("userID must be greater than 1")
	}

	spath := "users/" + strconv.Itoa(userID) + "/activities"
	return getActivityList(s.method.Get, spath, options...)
}
