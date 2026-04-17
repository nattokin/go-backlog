package model

import "time"

// Attachment represents an attached file.
type Attachment struct {
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Size        int       `json:"size,omitempty"`
	CreatedUser *User     `json:"createdUser,omitempty"`
	Created     time.Time `json:"created,omitempty"`
}

// ChangeLog represents a history of changes made to an issue.
type ChangeLog struct {
	Field         string `json:"field,omitempty"`
	NewValue      string `json:"newValue,omitempty"`
	OriginalValue string `json:"originalValue,omitempty"`
}

// Comment represents any one comment.
type Comment struct {
	ID            int             `json:"id,omitempty"`
	Content       string          `json:"content,omitempty"`
	ChangeLogs    []*ChangeLog    `json:"changeLog,omitempty"`
	CreatedUser   *User           `json:"createdUser,omitempty"`
	Created       time.Time       `json:"created,omitempty"`
	Updated       time.Time       `json:"updated,omitempty"`
	Stars         []*Star         `json:"stars,omitempty"`
	Notifications []*Notification `json:"notifications,omitempty"`
}

// CustomField represents a custom field defined in the project.
type CustomField struct {
	ID                     int                `json:"id,omitempty"`
	TypeID                 int                `json:"typeId,omitempty"`
	Name                   string             `json:"name,omitempty"`
	Description            string             `json:"description,omitempty"`
	Required               bool               `json:"required,omitempty"`
	ApplicableIssueTypeIDs []int              `json:"applicableIssueTypes,omitempty"`
	AllowAddItem           bool               `json:"allowAddItem,omitempty"`
	Items                  []*CustomFieldItem `json:"items,omitempty"`
}

// CustomFieldItem represents one of Items in CustomField.
type CustomFieldItem struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
}

// DiskUsageBase represents base of disk usage.
type DiskUsageBase struct {
	Issue      int `json:"issue,omitempty"`
	Wiki       int `json:"wiki,omitempty"`
	File       int `json:"file,omitempty"`
	Subversion int `json:"subversion,omitempty"`
	Git        int `json:"git,omitempty"`
	GitLFS     int `json:"gitLFS,omitempty"`
}

// Licence represents licence.
type Licence struct {
	Active                            bool      `json:"active,omitempty"`
	AttachmentLimit                   int       `json:"attachmentLimit,omitempty"`
	AttachmentLimitPerFile            int       `json:"attachmentLimitPerFile,omitempty"`
	AttachmentNumLimit                int       `json:"attachmentNumLimit,omitempty"`
	Attribute                         bool      `json:"attribute,omitempty"`
	AttributeLimit                    int       `json:"attributeLimit,omitempty"`
	Burndown                          bool      `json:"burndown,omitempty"`
	CommentLimit                      int       `json:"commentLimit,omitempty"`
	ComponentLimit                    int       `json:"componentLimit,omitempty"`
	FileSharing                       bool      `json:"fileSharing,omitempty"`
	Gantt                             bool      `json:"gantt,omitempty"`
	Git                               bool      `json:"git,omitempty"`
	IssueLimit                        int       `json:"issueLimit,omitempty"`
	LicenceTypeID                     int       `json:"licenceTypeId,omitempty"`
	LimitDate                         time.Time `json:"limitDate,omitempty"`
	NulabAccount                      bool      `json:"nulabAccount,omitempty"`
	ParentChildIssue                  bool      `json:"parentChildIssue,omitempty"`
	PostIssueByMail                   bool      `json:"postIssueByMail,omitempty"`
	ProjectGroup                      bool      `json:"projectGroup,omitempty"`
	ProjectLimit                      int       `json:"projectLimit,omitempty"`
	PullRequestAttachmentLimitPerFile int       `json:"pullRequestAttachmentLimitPerFile,omitempty"`
	PullRequestAttachmentNumLimit     int       `json:"pullRequestAttachmentNumLimit,omitempty"`
	RemoteAddress                     bool      `json:"remoteAddress,omitempty"`
	RemoteAddressLimit                int       `json:"remoteAddressLimit,omitempty"`
	StartedOn                         time.Time `json:"startedOn,omitempty"`
	StorageLimit                      int64     `json:"storageLimit,omitempty"`
	Subversion                        bool      `json:"subversion,omitempty"`
	SubversionExternal                bool      `json:"subversionExternal,omitempty"`
	UserLimit                         int       `json:"userLimit,omitempty"`
	VersionLimit                      int       `json:"versionLimit,omitempty"`
	WikiAttachment                    bool      `json:"wikiAttachment,omitempty"`
	WikiAttachmentLimitPerFile        int       `json:"wikiAttachmentLimitPerFile,omitempty"`
	WikiAttachmentNumLimit            int       `json:"wikiAttachmentNumLimit,omitempty"`
}

// Notification represents some notification.
type Notification struct {
	ID                  int          `json:"id,omitempty"`
	AlreadyRead         bool         `json:"alreadyRead,omitempty"`
	Reason              int          `json:"reason,omitempty"`
	ResourceAlreadyRead bool         `json:"resourceAlreadyRead,omitempty"`
	Project             *Project     `json:"project,omitempty"`
	Issue               *Issue       `json:"issue,omitempty"`
	Comment             *Comment     `json:"comment,omitempty"`
	PullRequest         *PullRequest `json:"pullRequest,omitempty"`
	PullRequestComment  *Comment     `json:"pullRequestComment,omitempty"`
	Sender              *User        `json:"sender,omitempty"`
	Created             time.Time    `json:"created,omitempty"`
}

// SharedFile represents a file shared within the project or space.
type SharedFile struct {
	ID          int       `json:"id,omitempty"`
	Type        string    `json:"type,omitempty"`
	Dir         string    `json:"dir,omitempty"`
	Name        string    `json:"name,omitempty"`
	Size        int       `json:"size,omitempty"`
	CreatedUser *User     `json:"createdUser,omitempty"`
	Created     time.Time `json:"created,omitempty"`
	UpdatedUser *User     `json:"updatedUser,omitempty"`
	Updated     time.Time `json:"updated,omitempty"`
}

// Star represents any Star.
type Star struct {
	ID        int       `json:"id,omitempty"`
	Comment   string    `json:"comment,omitempty"`
	URL       string    `json:"url,omitempty"`
	Title     string    `json:"title,omitempty"`
	Presenter *User     `json:"presenter,omitempty"`
	Created   time.Time `json:"created,omitempty"`
}

// Status represents any status.
type Status struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Tag represents one of tags in Wiki.
type Tag struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Team represents team.
type Team struct {
	ID           int       `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Members      []*User   `json:"members,omitempty"`
	DisplayOrder int       `json:"displayOrder,omitempty"`
	CreatedUser  *User     `json:"createdUser,omitempty"`
	Created      time.Time `json:"created,omitempty"`
	UpdatedUser  *User     `json:"updatedUser,omitempty"`
	Updated      time.Time `json:"updated,omitempty"`
}

// Version represents any version.
type Version struct {
	ID             int       `json:"id,omitempty"`
	ProjectID      int       `json:"projectId,omitempty"`
	Name           string    `json:"name,omitempty"`
	Description    string    `json:"description,omitempty"`
	StartDate      time.Time `json:"startDate,omitempty"`
	ReleaseDueDate time.Time `json:"releaseDueDate,omitempty"`
	Archived       bool      `json:"archived,omitempty"`
	DisplayOrder   int       `json:"displayOrder,omitempty"`
}

// WatchingItem represents an item of watching list.
type WatchingItem struct {
	ID                  int       `json:"id,omitempty"`
	ResourceAlreadyRead bool      `json:"resourceAlreadyRead,omitempty"`
	Note                string    `json:"note,omitempty"`
	Type                string    `json:"type,omitempty"`
	Issue               *Issue    `json:"issue,omitempty"`
	LastContentUpdated  time.Time `json:"lastContentUpdated,omitempty"`
	Created             time.Time `json:"created,omitempty"`
	Updated             time.Time `json:"updated,omitempty"`
}

// Webhook represents webhook of Backlog.
type Webhook struct {
	ID              int       `json:"id,omitempty"`
	Name            string    `json:"name,omitempty"`
	Description     string    `json:"description,omitempty"`
	HookURL         string    `json:"hookUrl,omitempty"`
	AllEvent        bool      `json:"allEvent,omitempty"`
	ActivityTypeIds []int     `json:"activityTypeIds,omitempty"`
	CreatedUser     *User     `json:"createdUser,omitempty"`
	Created         time.Time `json:"created,omitempty"`
	UpdatedUser     *User     `json:"updatedUser,omitempty"`
	Updated         time.Time `json:"updated,omitempty"`
}
