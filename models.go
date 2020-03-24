package backlog

import (
	"time"
)

// Activity represents
type Activity struct {
	ID            int              `json:"id,omitempty"`
	Project       *Project         `json:"project,omitempty"`
	Type          int              `json:"type,omitempty"`
	Content       *ActivityContent `json:"content,omitempty"`
	Notifications []*Notification  `json:"notifications,omitempty"`
	CreatedUser   *User            `json:"createdUser,omitempty"`
}

// ActivityContent represents content of Backlog activity.
type ActivityContent struct {
	ID          int      `json:"id,omitempty"`
	KeyID       int      `json:"key_id,omitempty"`
	Summary     string   `json:"summary,omitempty"`
	Description string   `json:"description,omitempty"`
	Comment     *Comment `json:"comment,omitempty"`
}

// Attachment represents one of attachments.
type Attachment struct {
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Size        int       `json:"size,omitempty"`
	CreatedUser *User     `json:"createdUser,omitempty"`
	Created     time.Time `json:"created,omitempty"`
}

// Category represents category of Backlog.
type Category []struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
}

// ChangeLog reprements one of ChangeLogs.
type ChangeLog struct {
	Field         string `json:"field,omitempty"`
	NewValue      string `json:"newValue,omitempty"`
	OriginalValue string `json:"originalValue,omitempty"`
}

// Comment reprements comment of Backlog.
type Comment struct {
	ID            int             `json:"id,omitempty"`
	Content       string          `json:"content,omitempty"`
	ChangeLogs    []*ChangeLog    `json:"changeLog,omitempty"`
	CreatedUser   *User           `json:"createdUser,omitempty"`
	Created       time.Time       `json:"created,omitempty"`
	Updated       time.Time       `json:"updated,omitempty"`
	Stars         *Star           `json:"stars,omitempty"`
	Notifications []*Notification `json:"notifications,omitempty"`
}

// CustomField represents custom field of Backlog.
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

// DiskUsageSpace represents space's disk usage.
type DiskUsageSpace struct {
	DiskUsageBase
	Capacity int                 `json:"capacity,omitempty"`
	Details  []*DiskUsageProject `json:"details,omitempty"`
}

// DiskUsageProject represents project's disk usage.
type DiskUsageProject struct {
	DiskUsageBase
	ProjectID int `json:"projectId,omitempty"`
}

// Error represents one of Backlog API response errors.
type Error struct {
	Message  string `json:"message,omitempty"`
	Code     int    `json:"code,omitempty"`
	MoreInfo string `json:"moreInfo,omitempty"`
}

// Licence represents licence of space.
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

// Notification represents notification of Backlog.
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

// Issue represents Backlog Issue.
type Issue struct {
	ID             int            `json:"id,omitempty"`
	ProjectID      int            `json:"projectId,omitempty"`
	IssueKey       string         `json:"issueKey,omitempty"`
	KeyID          int            `json:"keyId,omitempty"`
	IssueType      *IssueType     `json:"issueType,omitempty"`
	Summary        string         `json:"summary,omitempty"`
	Description    string         `json:"description,omitempty"`
	Resolutions    []*Resolution  `json:"resolutions,omitempty"`
	Priority       *Priority      `json:"priority,omitempty"`
	Status         *Status        `json:"status,omitempty"`
	Assignee       *User          `json:"assignee,omitempty"`
	Category       *Category      `json:"category,omitempty"`
	Versions       *Version       `json:"versions,omitempty"`
	Milestone      *Version       `json:"milestone,omitempty"`
	StartDate      time.Time      `json:"startDate,omitempty"`
	DueDate        time.Time      `json:"dueDate,omitempty"`
	EstimatedHours int            `json:"estimatedHours,omitempty"`
	ActualHours    int            `json:"actualHours,omitempty"`
	ParentIssueID  int            `json:"parentIssueId,omitempty"`
	CreatedUser    *User          `json:"createdUser,omitempty"`
	Created        time.Time      `json:"created,omitempty"`
	UpdatedUser    *User          `json:"updatedUser,omitempty"`
	Updated        time.Time      `json:"updated,omitempty"`
	CustomFields   []*CustomField `json:"customFields,omitempty"`
	Attachments    []*Attachment  `json:"attachments,omitempty"`
	SharedFiles    []*SharedFile  `json:"sharedFiles,omitempty"`
	Stars          []*Star        `json:"stars,omitempty"`
}

// IssueType represents type of Issue.
type IssueType struct {
	ID           int    `json:"id,omitempty"`
	ProjectID    int    `json:"projectId,omitempty"`
	Name         string `json:"name,omitempty"`
	Color        string `json:"color,omitempty"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
}

// Priority represents priority of Backlog.
type Priority struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Project represents Backlog Project.
type Project struct {
	ID                                int    `json:"id,omitempty"`
	ProjectKey                        string `json:"projectKey,omitempty"`
	Name                              string `json:"name,omitempty"`
	ChartEnabled                      bool   `json:"chartEnabled,omitempty"`
	SubtaskingEnabled                 bool   `json:"subtaskingEnabled,omitempty"`
	ProjectLeaderCanEditProjectLeader bool   `json:"projectLeaderCanEditProjectLeader,omitempty"`
	TextFormattingRule                string `json:"textFormattingRule,omitempty"`
	Archived                          bool   `json:"archived,omitempty"`
}

// PullRequest represents pull request of Backlog git.
type PullRequest struct {
	ID           int           `json:"id,omitempty"`
	ProjectID    int           `json:"projectId,omitempty"`
	RepositoryID int           `json:"repositoryId,omitempty"`
	Number       int           `json:"number,omitempty"`
	Summary      string        `json:"summary,omitempty"`
	Description  string        `json:"description,omitempty"`
	Base         string        `json:"base,omitempty"`
	Branch       string        `json:"branch,omitempty"`
	Status       *Status       `json:"status,omitempty"`
	Assignee     *User         `json:"assignee,omitempty"`
	Issue        *Issue        `json:"issue,omitempty"`
	BaseCommit   interface{}   `json:"baseCommit,omitempty"`
	BranchCommit interface{}   `json:"branchCommit,omitempty"`
	CloseAt      time.Time     `json:"closeAt,omitempty"`
	MergeAt      time.Time     `json:"mergeAt,omitempty"`
	CreatedUser  *User         `json:"createdUser,omitempty"`
	Created      time.Time     `json:"created,omitempty"`
	UpdatedUser  *User         `json:"updatedUser,omitempty"`
	Updated      time.Time     `json:"updated,omitempty"`
	Attachments  []*Attachment `json:"attachments,omitempty"`
	Stars        []*Star       `json:"stars,omitempty"`
}

// Repository represents repository of Backlog git.
type Repository struct {
	ID           int       `json:"id,omitempty"`
	ProjectID    int       `json:"projectId,omitempty"`
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	HookURL      string    `json:"hookUrl,omitempty"`
	HTTPURL      string    `json:"httpUrl,omitempty"`
	SSHURL       string    `json:"sshUrl,omitempty"`
	DisplayOrder int       `json:"displayOrder,omitempty"`
	PushedAt     time.Time `json:"pushedAt,omitempty"`
	CreatedUser  *User     `json:"createdUser,omitempty"`
	Created      time.Time `json:"created,omitempty"`
	UpdatedUser  *User     `json:"updatedUser,omitempty"`
	Updated      time.Time `json:"updated,omitempty"`
}

// Resolution represents
type Resolution struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// SharedFile represents one of SharedFiles.
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

// Space represents Space of Backlog.
type Space struct {
	SpaceKey           string    `json:"spaceKey,omitempty"`
	Name               string    `json:"name,omitempty"`
	OwnerID            int       `json:"ownerId,omitempty"`
	Lang               string    `json:"lang,omitempty"`
	Timezone           string    `json:"timezone,omitempty"`
	ReportSendTime     string    `json:"reportSendTime,omitempty"`
	TextFormattingRule string    `json:"textFormattingRule,omitempty"`
	Created            time.Time `json:"created,omitempty"`
	Updated            time.Time `json:"updated,omitempty"`
}

// SpaceNotification represents a notification of Space.
type SpaceNotification struct {
	Content string    `json:"content,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
}

// Star represents any Star.
type Star struct {
	ID        int       `json:"id,omitempty"`
	Comment   *Comment  `json:"comment,omitempty"`
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

// User represents user.
type User struct {
	ID          int    `json:"id,omitempty"`
	UserID      string `json:"userId,omitempty"`
	Name        string `json:"name,omitempty"`
	RoleType    int    `json:"roleType,omitempty"`
	Lang        string `json:"lang,omitempty"`
	MailAddress string `json:"mailAddress,omitempty"`
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

// Wiki represents Backlog Wiki.
type Wiki struct {
	ID          int           `json:"id,omitempty"`
	ProjectID   int           `json:"projectId,omitempty"`
	Name        string        `json:"name,omitempty"`
	Content     string        `json:"content,omitempty"`
	Tags        []*Tag        `json:"tags,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
	SharedFiles []*SharedFile `json:"sharedFiles,omitempty"`
	Stars       []*Star       `json:"stars,omitempty"`
	CreatedUser *User         `json:"createdUser,omitempty"`
	Created     time.Time     `json:"created,omitempty"`
	UpdatedUser *User         `json:"updatedUser,omitempty"`
	Updated     time.Time     `json:"updated,omitempty"`
}

// WikiHistory reprements history of Wiki.
type WikiHistory struct {
	PageID      int       `json:"pageId,omitempty"`
	Version     int       `json:"version,omitempty"`
	Name        string    `json:"name,omitempty"`
	Content     string    `json:"content,omitempty"`
	CreatedUser *User     `json:"createdUser,omitempty"`
	Created     time.Time `json:"created,omitempty"`
}
