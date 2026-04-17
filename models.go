package backlog

import (
	"time"

	"github.com/nattokin/go-backlog/internal/model"
)

// Attachment represents an attached file.
type Attachment struct {
	ID          int
	Name        string
	Size        int
	CreatedUser *User
	Created     time.Time
}

// Activity represents a recent update or change in the project or space.
type Activity struct {
	ID            int
	Project       *Project
	Type          int
	Content       *ActivityContent
	Notifications []*Notification
	CreatedUser   *User
}

// ActivityContent represents the detailed content of an activity.
type ActivityContent struct {
	ID          int
	KeyID       int
	Summary     string
	Description string
	Comment     *Comment
}

// ChangeLog represents a history of changes made to an issue.
type ChangeLog struct {
	Field         string
	NewValue      string
	OriginalValue string
}

// Comment represents any one comment.
type Comment struct {
	ID            int
	Content       string
	ChangeLogs    []*ChangeLog
	CreatedUser   *User
	Created       time.Time
	Updated       time.Time
	Stars         []*Star
	Notifications []*Notification
}

// CustomField represents a custom field defined in the project.
type CustomField struct {
	ID                     int
	TypeID                 int
	Name                   string
	Description            string
	Required               bool
	ApplicableIssueTypeIDs []int
	AllowAddItem           bool
	Items                  []*CustomFieldItem
}

// CustomFieldItem represents one of Items in CustomField.
type CustomFieldItem struct {
	ID           int
	Name         string
	DisplayOrder int
}

// DiskUsageBase represents base of disk usage.
type DiskUsageBase struct {
	Issue      int
	Wiki       int
	File       int
	Subversion int
	Git        int
	GitLFS     int
}

// Licence represents licence.
type Licence struct {
	Active                            bool
	AttachmentLimit                   int
	AttachmentLimitPerFile            int
	AttachmentNumLimit                int
	Attribute                         bool
	AttributeLimit                    int
	Burndown                          bool
	CommentLimit                      int
	ComponentLimit                    int
	FileSharing                       bool
	Gantt                             bool
	Git                               bool
	IssueLimit                        int
	LicenceTypeID                     int
	LimitDate                         time.Time
	NulabAccount                      bool
	ParentChildIssue                  bool
	PostIssueByMail                   bool
	ProjectGroup                      bool
	ProjectLimit                      int
	PullRequestAttachmentLimitPerFile int
	PullRequestAttachmentNumLimit     int
	RemoteAddress                     bool
	RemoteAddressLimit                int
	StartedOn                         time.Time
	StorageLimit                      int64
	Subversion                        bool
	SubversionExternal                bool
	UserLimit                         int
	VersionLimit                      int
	WikiAttachment                    bool
	WikiAttachmentLimitPerFile        int
	WikiAttachmentNumLimit            int
}

// Notification represents some notification.
type Notification struct {
	ID                  int
	AlreadyRead         bool
	Reason              int
	ResourceAlreadyRead bool
	Project             *Project
	Issue               *Issue
	Comment             *Comment
	PullRequest         *PullRequest
	PullRequestComment  *Comment
	Sender              *User
	Created             time.Time
}

// SharedFile represents a file shared within the project or space.
type SharedFile struct {
	ID          int
	Type        string
	Dir         string
	Name        string
	Size        int
	CreatedUser *User
	Created     time.Time
	UpdatedUser *User
	Updated     time.Time
}

// Star represents any Star.
type Star struct {
	ID        int
	Comment   string
	URL       string
	Title     string
	Presenter *User
	Created   time.Time
}

// Status represents any status.
type Status struct {
	ID   int
	Name string
}

// Tag represents one of tags in Wiki.
type Tag struct {
	ID   int
	Name string
}

// Team represents team.
type Team struct {
	ID           int
	Name         string
	Members      []*User
	DisplayOrder int
	CreatedUser  *User
	Created      time.Time
	UpdatedUser  *User
	Updated      time.Time
}

// Version represents any version.
type Version struct {
	ID             int
	ProjectID      int
	Name           string
	Description    string
	StartDate      time.Time
	ReleaseDueDate time.Time
	Archived       bool
	DisplayOrder   int
}

// WatchingItem represents an item of watching list.
type WatchingItem struct {
	ID                  int
	ResourceAlreadyRead bool
	Note                string
	Type                string
	Issue               *Issue
	LastContentUpdated  time.Time
	Created             time.Time
	Updated             time.Time
}

// Webhook represents webhook of Backlog.
type Webhook struct {
	ID              int
	Name            string
	Description     string
	HookURL         string
	AllEvent        bool
	ActivityTypeIds []int
	CreatedUser     *User
	Created         time.Time
	UpdatedUser     *User
	Updated         time.Time
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func activityContentFromModel(m *model.ActivityContent) *ActivityContent {
	if m == nil {
		return nil
	}
	return &ActivityContent{
		ID:          m.ID,
		KeyID:       m.KeyID,
		Summary:     m.Summary,
		Description: m.Description,
		Comment:     commentFromModel(m.Comment),
	}
}

func activityFromModel(m *model.Activity) *Activity {
	if m == nil {
		return nil
	}
	notifications := make([]*Notification, len(m.Notifications))
	for i, v := range m.Notifications {
		notifications[i] = notificationFromModel(v)
	}
	return &Activity{
		ID:            m.ID,
		Project:       projectFromModel(m.Project),
		Type:          m.Type,
		Content:       activityContentFromModel(m.Content),
		Notifications: notifications,
		CreatedUser:   userFromModel(m.CreatedUser),
	}
}

func activitiesFromModel(ms []*model.Activity) []*Activity {
	result := make([]*Activity, len(ms))
	for i, v := range ms {
		result[i] = activityFromModel(v)
	}
	return result
}

func changeLogFromModel(m *model.ChangeLog) *ChangeLog {
	if m == nil {
		return nil
	}
	return &ChangeLog{
		Field:         m.Field,
		NewValue:      m.NewValue,
		OriginalValue: m.OriginalValue,
	}
}

func starFromModel(m *model.Star) *Star {
	if m == nil {
		return nil
	}
	return &Star{
		ID:        m.ID,
		Comment:   m.Comment,
		URL:       m.URL,
		Title:     m.Title,
		Presenter: userFromModel(m.Presenter),
		Created:   m.Created,
	}
}

func starsFromModel(ms []*model.Star) []*Star {
	if ms == nil {
		return nil
	}
	result := make([]*Star, len(ms))
	for i, v := range ms {
		result[i] = starFromModel(v)
	}
	return result
}

func attachmentFromModel(m *model.Attachment) *Attachment {
	if m == nil {
		return nil
	}
	return &Attachment{
		ID:          m.ID,
		Name:        m.Name,
		Size:        m.Size,
		CreatedUser: userFromModel(m.CreatedUser),
		Created:     m.Created,
	}
}

func attachmentsFromModel(m []*model.Attachment) []*Attachment {
	if m == nil {
		return nil
	}
	result := make([]*Attachment, len(m))
	for i, v := range m {
		result[i] = attachmentFromModel(v)
	}
	return result
}

func customFieldItemFromModel(m *model.CustomFieldItem) *CustomFieldItem {
	if m == nil {
		return nil
	}
	return &CustomFieldItem{
		ID:           m.ID,
		Name:         m.Name,
		DisplayOrder: m.DisplayOrder,
	}
}

func customFieldFromModel(m *model.CustomField) *CustomField {
	if m == nil {
		return nil
	}
	items := make([]*CustomFieldItem, len(m.Items))
	for i, v := range m.Items {
		items[i] = customFieldItemFromModel(v)
	}
	return &CustomField{
		ID:                     m.ID,
		TypeID:                 m.TypeID,
		Name:                   m.Name,
		Description:            m.Description,
		Required:               m.Required,
		ApplicableIssueTypeIDs: m.ApplicableIssueTypeIDs,
		AllowAddItem:           m.AllowAddItem,
		Items:                  items,
	}
}

func commentFromModel(m *model.Comment) *Comment {
	if m == nil {
		return nil
	}
	changeLogs := make([]*ChangeLog, len(m.ChangeLogs))
	for i, v := range m.ChangeLogs {
		changeLogs[i] = changeLogFromModel(v)
	}
	notifications := make([]*Notification, len(m.Notifications))
	for i, v := range m.Notifications {
		notifications[i] = notificationFromModel(v)
	}
	return &Comment{
		ID:            m.ID,
		Content:       m.Content,
		ChangeLogs:    changeLogs,
		CreatedUser:   userFromModel(m.CreatedUser),
		Created:       m.Created,
		Updated:       m.Updated,
		Stars:         starsFromModel(m.Stars),
		Notifications: notifications,
	}
}

func notificationFromModel(m *model.Notification) *Notification {
	if m == nil {
		return nil
	}
	return &Notification{
		ID:                  m.ID,
		AlreadyRead:         m.AlreadyRead,
		Reason:              m.Reason,
		ResourceAlreadyRead: m.ResourceAlreadyRead,
		Project:             projectFromModel(m.Project),
		Issue:               issueFromModel(m.Issue),
		Comment:             commentFromModel(m.Comment),
		PullRequest:         pullRequestFromModel(m.PullRequest),
		PullRequestComment:  commentFromModel(m.PullRequestComment),
		Sender:              userFromModel(m.Sender),
		Created:             m.Created,
	}
}

func sharedFileFromModel(m *model.SharedFile) *SharedFile {
	if m == nil {
		return nil
	}
	return &SharedFile{
		ID:          m.ID,
		Type:        m.Type,
		Dir:         m.Dir,
		Name:        m.Name,
		Size:        m.Size,
		CreatedUser: userFromModel(m.CreatedUser),
		Created:     m.Created,
		UpdatedUser: userFromModel(m.UpdatedUser),
		Updated:     m.Updated,
	}
}

func statusFromModel(m *model.Status) *Status {
	if m == nil {
		return nil
	}
	return &Status{
		ID:   m.ID,
		Name: m.Name,
	}
}

func tagFromModel(m *model.Tag) *Tag {
	if m == nil {
		return nil
	}
	return &Tag{
		ID:   m.ID,
		Name: m.Name,
	}
}
