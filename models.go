package backlog

import (
	"io"

	"github.com/nattokin/go-backlog/internal/model"
)

// Attachment represents a file attached to an issue, wiki page, or pull request.
type Attachment struct {
	ID          int
	Name        string
	Size        int
	CreatedUser *User
	Created     Timestamp
}

// Activity represents a recent update or change in the project or space.
type Activity struct {
	ID            int
	Project       *Project
	Type          int
	Content       *ActivityContent
	Notifications []*Notification
	CreatedUser   *User
	Created       Timestamp
}

// ActivityContent represents the detailed content of an activity.
type ActivityContent struct {
	ID          int
	KeyID       int
	Summary     string
	Description string
	Comment     *Comment
}

// Category represents a project category.
type Category struct {
	ID           int
	Name         string
	DisplayOrder int
}

// ChangeLog represents a history of changes made to an issue.
type ChangeLog struct {
	Field          string
	NewValue       string
	OriginalValue  string
	AttachmentInfo *struct {
		ID   int
		Name string
	}
	AttributeInfo *struct {
		ID     int
		TypeID string
	}
	NotificationInfo *struct {
		Type string
	}
}

// Comment represents a comment on an issue or pull request.
type Comment struct {
	ID            int
	Content       string
	ChangeLogs    []*ChangeLog
	CreatedUser   *User
	Created       Timestamp
	Updated       Timestamp
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

// CustomFieldItem represents one item in a list type [CustomField].
type CustomFieldItem struct {
	ID           int
	Name         string
	DisplayOrder int
}

// FileData represents a downloaded binary file with its metadata.
// Body must be closed by the caller after use.
type FileData struct {
	Body        io.ReadCloser
	Filename    string
	ContentType string
}

// Licence represents the licence information for a Backlog space.
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
	LimitDate                         Timestamp
	NulabAccount                      bool
	ParentChildIssue                  bool
	PostIssueByMail                   bool
	ProjectGroup                      bool
	ProjectLimit                      int
	PullRequestAttachmentLimitPerFile int
	PullRequestAttachmentNumLimit     int
	RemoteAddress                     bool
	RemoteAddressLimit                int
	StartedOn                         Timestamp
	StorageLimit                      int64
	Subversion                        bool
	SubversionExternal                bool
	UserLimit                         int
	VersionLimit                      int
	WikiAttachment                    bool
	WikiAttachmentLimitPerFile        int
	WikiAttachmentNumLimit            int
}

// Notification represents a notification delivered to a user.
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
	User                *User
	Created             Timestamp
}

// SharedFile represents a file shared within the project or space.
type SharedFile struct {
	ID          int
	Type        string
	Dir         string
	Name        string
	Size        int
	CreatedUser *User
	Created     Timestamp
	UpdatedUser *User
	Updated     Timestamp
}

// Star represents a star added to an issue, wiki page, or pull request.
type Star struct {
	ID        int
	Comment   string
	URL       string
	Title     string
	Presenter *User
	Created   Timestamp
}

// Status represents a project status that can be assigned to issues.
type Status struct {
	ID           int
	ProjectID    int
	Name         string
	Color        string
	DisplayOrder int
}

// Tag represents a tag attached to a wiki page.
type Tag struct {
	ID   int
	Name string
}

// Team represents a team within a space.
type Team struct {
	ID           int
	Name         string
	Members      []*User
	DisplayOrder int
	CreatedUser  *User
	Created      Timestamp
	UpdatedUser  *User
	Updated      Timestamp
}

// Version represents a version or milestone in a project.
type Version struct {
	ID             int
	ProjectID      int
	Name           string
	Description    string
	StartDate      Date
	ReleaseDueDate Date
	Archived       bool
	DisplayOrder   int
}

// Watching represents an item in a user's watching list.
type Watching struct {
	ID                 int
	AlreadyRead        bool
	Note               string
	Type               string
	Issue              *Issue
	LastContentUpdated Timestamp
	Created            Timestamp
	Updated            Timestamp
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
		Created:       Timestamp{m.Created},
	}
}

func activitiesFromModel(ms []*model.Activity) []*Activity {
	result := make([]*Activity, len(ms))
	for i, v := range ms {
		result[i] = activityFromModel(v)
	}
	return result
}

func categoryFromModel(m *model.Category) *Category {
	if m == nil {
		return nil
	}
	return &Category{
		ID:           m.ID,
		Name:         m.Name,
		DisplayOrder: m.DisplayOrder,
	}
}

func categoriesFromModel(ms []*model.Category) []*Category {
	if ms == nil {
		return nil
	}
	result := make([]*Category, len(ms))
	for i, v := range ms {
		result[i] = categoryFromModel(v)
	}
	return result
}

func changeLogFromModel(m *model.ChangeLog) *ChangeLog {
	if m == nil {
		return nil
	}
	out := &ChangeLog{
		Field:         m.Field,
		NewValue:      m.NewValue,
		OriginalValue: m.OriginalValue,
	}
	if m.AttachmentInfo != nil {
		out.AttachmentInfo = &struct {
			ID   int
			Name string
		}{
			ID:   m.AttachmentInfo.ID,
			Name: m.AttachmentInfo.Name,
		}
	}
	if m.AttributeInfo != nil {
		out.AttributeInfo = &struct {
			ID     int
			TypeID string
		}{
			ID:     m.AttributeInfo.ID,
			TypeID: m.AttributeInfo.TypeID,
		}
	}
	if m.NotificationInfo != nil {
		out.NotificationInfo = &struct {
			Type string
		}{
			Type: m.NotificationInfo.Type,
		}
	}
	return out
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
		Created:   Timestamp{m.Created},
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
		Created:     Timestamp{m.Created},
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

func customFieldsFromModel(m []*model.CustomField) []*CustomField {
	if m == nil {
		return nil
	}
	result := make([]*CustomField, len(m))
	for i, v := range m {
		result[i] = customFieldFromModel(v)
	}
	return result
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
		Created:       Timestamp{m.Created},
		Updated:       Timestamp{m.Updated},
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
		User:                userFromModel(m.User),
		Created:             Timestamp{m.Created},
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
		Created:     Timestamp{m.Created},
		UpdatedUser: userFromModel(m.UpdatedUser),
		Updated:     Timestamp{m.Updated},
	}
}

func sharedFilesFromModel(m []*model.SharedFile) []*SharedFile {
	if m == nil {
		return nil
	}
	result := make([]*SharedFile, len(m))
	for i, v := range m {
		result[i] = sharedFileFromModel(v)
	}
	return result
}

func statusFromModel(m *model.Status) *Status {
	if m == nil {
		return nil
	}
	return &Status{
		ID:           m.ID,
		ProjectID:    m.ProjectID,
		Name:         m.Name,
		Color:        m.Color,
		DisplayOrder: m.DisplayOrder,
	}
}

func statusesFromModel(ms []*model.Status) []*Status {
	if ms == nil {
		return nil
	}
	result := make([]*Status, len(ms))
	for i, v := range ms {
		result[i] = statusFromModel(v)
	}
	return result
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

func fileDataFromModel(m *model.FileData) *FileData {
	if m == nil {
		return nil
	}
	return &FileData{
		Body:        m.Body,
		Filename:    m.Filename,
		ContentType: m.ContentType,
	}
}
