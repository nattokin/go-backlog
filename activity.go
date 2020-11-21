package backlog

import (
	"encoding/json"
	"errors"
	"path"
	"strconv"
)

func getActivityList(get clientGet, spath string, options ...*QueryOption) ([]*Activity, error) {
	validOptions := []queryType{queryActivityTypeIDs, queryMinID, queryMaxID, queryCount, queryOrder}
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

	resp, err := get(spath, query)
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
func (s *ProjectActivityService) List(project ProjectIDOrKeyGetter, options ...*QueryOption) ([]*Activity, error) {
	projectIDOrKey, err := project.getProjectIDOrKey()
	if err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "activities")
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
	return getActivityList(s.method.Get, "space/activities", options...)
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

	spath := path.Join("users", strconv.Itoa(userID), "activities")
	return getActivityList(s.method.Get, spath, options...)
}
