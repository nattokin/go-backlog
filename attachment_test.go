package backlog_test

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestSpaceAttachmentService_Upload(t *testing.T) {
	bj, err := os.Open("testdata/json/attachment_upload.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	fpath := "fpath"
	fname := "test.txt"

	want := struct {
		spath string
		fpath string
		fname string
		id    int
		name  string
		size  int
	}{
		spath: "space/attachment",
		fpath: fpath,
		fname: fname,
		id:    1,
		name:  fname,
		size:  8857,
	}
	s := &backlog.SpaceAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Upload: func(spath, fpath, fname string) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.fpath, fpath)
			assert.Equal(t, want.fname, fname)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	attachment, err := s.Upload(fpath, fname)
	assert.NoError(t, err)
	if assert.NotNil(t, attachment) {
		assert.Equal(t, want.id, attachment.ID)
		assert.Equal(t, want.name, attachment.Name)
		assert.Equal(t, want.size, attachment.Size)
	}
}

func TestSpaceAttachmentService_Upload_clientError(t *testing.T) {
	s := &backlog.SpaceAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Upload: func(spath, fpath, fname string) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	attachement, err := s.Upload("fpath", "fname")
	assert.Error(t, err)
	assert.Nil(t, attachement)
}

func TestSpaceAttachmentService_Upload_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.SpaceAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Upload: func(spath, fpath, fname string) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	attachement, err := s.Upload("fpath", "fname")
	assert.Error(t, err)
	assert.Nil(t, attachement)
}

func TestWikiAttachmentService_Attach(t *testing.T) {
	bj, err := os.Open("testdata/json/attachment_list.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	wikiID := 1234

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   "wikis/" + strconv.Itoa(wikiID) + "/attachments",
		id:      2,
		name:    "A.png",
		size:    196186,
		created: time.Date(2014, time.September, 11, 6, 26, 5, 0, time.UTC),
	}
	s := &backlog.WikiAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			v := *params.ExportURLValues()
			assert.Equal(t, []string{"2"}, v["attachmentId[]"])
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	attachments, err := s.Attach(wikiID, []int{2})
	assert.NoError(t, err)
	if assert.NotNil(t, attachments) {
		assert.Equal(t, want.id, attachments[0].ID)
		assert.Equal(t, want.name, attachments[0].Name)
		assert.Equal(t, want.size, attachments[0].Size)
		assert.Equal(t, want.size, attachments[0].Size)
		assert.ObjectsAreEqualValues(want.created, attachments[0].Created)
	}
}

func TestWikiAttachmentService_Attach_clientError(t *testing.T) {
	s := &backlog.WikiAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	attachments, err := s.Attach(1234, []int{2})
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestWikiAttachmentService_Attach_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	attachments, err := s.Attach(1234, []int{2})
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestWikiAttachmentService_List(t *testing.T) {
	bj, err := os.Open("testdata/json/attachment_list.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	wikiID := 1234

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   "wikis/" + strconv.Itoa(wikiID) + "/attachments",
		id:      2,
		name:    "A.png",
		size:    196186,
		created: time.Date(2014, time.September, 11, 6, 26, 5, 0, time.UTC),
	}
	s := &backlog.WikiAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	attachments, err := s.List(wikiID)
	assert.NoError(t, err)
	if assert.NotNil(t, attachments) {
		assert.Equal(t, want.id, attachments[0].ID)
		assert.Equal(t, want.name, attachments[0].Name)
		assert.Equal(t, want.size, attachments[0].Size)
		assert.Equal(t, want.size, attachments[0].Size)
		assert.ObjectsAreEqualValues(want.created, attachments[0].Created)
	}
}

func TestWikiAttachmentService_List_param(t *testing.T) {
	cases := map[string]struct {
		wikiID    int
		wantError bool
	}{
		"valid": {
			wikiID:    1,
			wantError: false,
		},
		"invalid-1": {
			wikiID:    0,
			wantError: true,
		},
		"invalid-2": {
			wikiID:    -1,
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			bj, err := os.Open("testdata/json/attachment_list.json")
			if err != nil {
				t.Fatal(err)
			}
			defer bj.Close()

			s := &backlog.WikiAttachmentService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if attachements, err := s.List(tc.wikiID); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, attachements, 2)
				assert.Equal(t, attachements[0].ID, 2)
				assert.Equal(t, attachements[1].ID, 5)
			}
		})
	}
}

func TestWikiAttachmentService_List_clientError(t *testing.T) {
	s := &backlog.WikiAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	attachments, err := s.List(1234)
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestWikiAttachmentService_List_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	attachments, err := s.List(1234)
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestWikiAttachmentService_Remove(t *testing.T) {
	bj, err := os.Open("testdata/json/attachment.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	wikiID := 1234
	attachmentID := 8

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   "wikis/" + strconv.Itoa(wikiID) + "/attachments/" + strconv.Itoa(attachmentID),
		id:      8,
		name:    "IMG0088.png",
		size:    5563,
		created: time.Date(2014, time.October, 28, 9, 24, 43, 0, time.UTC),
	}
	s := &backlog.WikiAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	attachments, err := s.Remove(wikiID, attachmentID)
	assert.NoError(t, err)
	if assert.NotNil(t, attachments) {
		assert.Equal(t, want.id, attachments.ID)
		assert.Equal(t, want.name, attachments.Name)
		assert.Equal(t, want.size, attachments.Size)
		assert.Equal(t, want.size, attachments.Size)
		assert.ObjectsAreEqualValues(want.created, attachments.Created)
	}
}

func TestWikiAttachmentService_Remove_param(t *testing.T) {
	bj, err := os.Open("testdata/json/attachment.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	cases := map[string]struct {
		wikiID       int
		attachmentID int
		wantError    bool
	}{
		"valid": {
			wikiID:       1,
			attachmentID: 8,
			wantError:    false,
		},
		"invalid-wikiID-1": {
			wikiID:       0,
			attachmentID: 8,
			wantError:    true,
		},
		"invalid-wikiID-2": {
			wikiID:       -1,
			attachmentID: 8,
			wantError:    true,
		},
		"invalid-attachmentID-1": {
			wikiID:       1,
			attachmentID: 0,
			wantError:    true,
		},
		"invalid-attachmentID-2": {
			wikiID:       1,
			attachmentID: -1,
			wantError:    true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := &backlog.WikiAttachmentService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Delete: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if attachements, err := s.Remove(tc.wikiID, tc.attachmentID); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, attachements)
			}
		})
	}
}

func TestWikiAttachmentService_Remove_clientError(t *testing.T) {
	s := &backlog.WikiAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	attachment, err := s.Remove(1234, 8)
	assert.Error(t, err)
	assert.Nil(t, attachment)
}

func TestWikiAttachmentService_Remove_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	attachment, err := s.Remove(1234, 8)
	assert.Error(t, err)
	assert.Nil(t, attachment)
}

func TestIssueAttachmentService_List(t *testing.T) {
	bj, err := os.Open("testdata/json/attachment_list.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	issueID := 1234

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   "issues/" + strconv.Itoa(issueID) + "/attachments",
		id:      2,
		name:    "A.png",
		size:    196186,
		created: time.Date(2014, time.September, 11, 6, 26, 5, 0, time.UTC),
	}
	s := &backlog.IssueAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	attachments, err := s.List(backlog.IssueID(issueID))
	assert.NoError(t, err)
	if assert.NotNil(t, attachments) {
		assert.Equal(t, want.id, attachments[0].ID)
		assert.Equal(t, want.name, attachments[0].Name)
		assert.Equal(t, want.size, attachments[0].Size)
		assert.Equal(t, want.size, attachments[0].Size)
		assert.ObjectsAreEqualValues(want.created, attachments[0].Created)
	}
}

func TestIssueAttachmentService_List_param(t *testing.T) {
	cases := map[string]struct {
		target    backlog.IssueIDOrKeyGetter
		wantError bool
	}{
		"valid": {
			target:    backlog.IssueID(1),
			wantError: false,
		},
		"invalid": {
			target:    backlog.IssueID(0),
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			bj, err := os.Open("testdata/json/attachment_list.json")
			if err != nil {
				t.Fatal(err)
			}
			defer bj.Close()

			s := &backlog.IssueAttachmentService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if attachements, err := s.List(tc.target); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, attachements, 2)
				assert.Equal(t, attachements[0].ID, 2)
				assert.Equal(t, attachements[1].ID, 5)
			}
		})
	}
}

func TestIssueAttachmentService_List_clientError(t *testing.T) {
	s := &backlog.IssueAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	attachments, err := s.List(backlog.IssueID(1234))
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestIssueAttachmentService_List_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.IssueAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	attachments, err := s.List(backlog.IssueID(1234))
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestIssueAttachmentService_Remove(t *testing.T) {
	bj, err := os.Open("testdata/json/attachment.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	issueID := 1234
	attachmentID := 8

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   "issues/" + strconv.Itoa(issueID) + "/attachments/" + strconv.Itoa(attachmentID),
		id:      attachmentID,
		name:    "IMG0088.png",
		size:    5563,
		created: time.Date(2014, time.October, 28, 9, 24, 43, 0, time.UTC),
	}

	s := &backlog.IssueAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	attachment, err := s.Remove(backlog.IssueID(issueID), attachmentID)
	assert.NoError(t, err)
	if assert.NotNil(t, attachment) {
		assert.Equal(t, want.id, attachment.ID)
		assert.Equal(t, want.name, attachment.Name)
		assert.Equal(t, want.size, attachment.Size)
		assert.Equal(t, want.size, attachment.Size)
		assert.ObjectsAreEqualValues(want.created, attachment.Created)
	}
}

func TestIssueAttachmentService_Remove_param(t *testing.T) {
	bj, err := os.Open("testdata/json/attachment.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	cases := map[string]struct {
		target       backlog.IssueIDOrKeyGetter
		attachmentID int
		wantError    bool
	}{
		"valid": {
			target:       backlog.IssueKey("test"),
			attachmentID: 8,
			wantError:    false,
		},
		"invalid-target": {
			target:       backlog.IssueKey(""),
			attachmentID: 8,
			wantError:    true,
		},
		"invalid-attachmentID": {
			target:       backlog.IssueKey("test"),
			attachmentID: 0,
			wantError:    true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := &backlog.IssueAttachmentService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Delete: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if attachements, err := s.Remove(tc.target, tc.attachmentID); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, attachements)
			}
		})
	}
}

func TestIssueAttachmentService_Remove_clientError(t *testing.T) {
	s := &backlog.IssueAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	attachment, err := s.Remove(backlog.IssueID(1234), 8)
	assert.Error(t, err)
	assert.Nil(t, attachment)
}

func TestIssueAttachmentService_Remove_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.IssueAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	attachment, err := s.Remove(backlog.IssueID(1234), 8)
	assert.Error(t, err)
	assert.Nil(t, attachment)
}

func TestPullRequestAttachmentService_List(t *testing.T) {
	bj, err := os.Open("testdata/json/attachment_list.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	projectKey := "TEST"
	repoName := "test"
	prNumber := 1234

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   "projects/" + projectKey + "/git/repositories/" + repoName + "/pullRequests/" + strconv.Itoa(prNumber) + "/attachments",
		id:      2,
		name:    "A.png",
		size:    196186,
		created: time.Date(2014, time.September, 11, 6, 26, 5, 0, time.UTC),
	}
	s := &backlog.PullRequestAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	attachments, err := s.List(backlog.ProjectKey(projectKey), backlog.RepositoryName(repoName), prNumber)
	assert.NoError(t, err)
	if assert.NotNil(t, attachments) {
		assert.Equal(t, want.id, attachments[0].ID)
		assert.Equal(t, want.name, attachments[0].Name)
		assert.Equal(t, want.size, attachments[0].Size)
		assert.Equal(t, want.size, attachments[0].Size)
		assert.ObjectsAreEqualValues(want.created, attachments[0].Created)
	}
}

func TestPullRequestAttachmentService_List_param(t *testing.T) {
	bj, err := os.Open("testdata/json/attachment_list.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	cases := map[string]struct {
		targetProject    backlog.ProjectIDOrKeyGetter
		targetRepository backlog.RepositoryIDOrKeyGetter
		prNumber         int
		wantError        bool
	}{
		"valid": {
			targetProject:    backlog.ProjectID(1),
			targetRepository: backlog.RepositoryID(1),
			prNumber:         1,
			wantError:        false,
		},
		"invalid-targetProject": {
			targetProject:    backlog.ProjectID(0),
			targetRepository: backlog.RepositoryID(1),
			prNumber:         1,
			wantError:        true,
		},
		"invalid-targetRepository": {
			targetProject:    backlog.ProjectID(1),
			targetRepository: backlog.RepositoryID(0),
			prNumber:         1,
			wantError:        true,
		},
		"invalid-prNumber": {
			targetProject:    backlog.ProjectID(1),
			targetRepository: backlog.RepositoryID(1),
			prNumber:         0,
			wantError:        true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := &backlog.PullRequestAttachmentService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if attachements, err := s.List(tc.targetProject, tc.targetRepository, tc.prNumber); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, attachements, 2)
				assert.Equal(t, attachements[0].ID, 2)
				assert.Equal(t, attachements[1].ID, 5)
			}
		})
	}
}

func TestPullRequestAttachmentService_List_clientError(t *testing.T) {
	s := &backlog.PullRequestAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	attachments, err := s.List(backlog.ProjectID(1234), backlog.RepositoryName("test"), 10)
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestPullRequestAttachmentService_List_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.PullRequestAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	attachments, err := s.List(backlog.ProjectID(1234), backlog.RepositoryName("test"), 10)
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestPullRequestAttachmentService_Remove(t *testing.T) {
	bj, err := os.Open("testdata/json/attachment.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	projectKey := "TEST"
	repoName := "test"
	prNumber := 1234
	attachmentID := 8

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   "projects/" + projectKey + "/git/repositories/" + repoName + "/pullRequests/" + strconv.Itoa(prNumber) + "/attachments" + strconv.Itoa(attachmentID),
		id:      8,
		name:    "IMG0088.png",
		size:    5563,
		created: time.Date(2014, time.October, 28, 9, 24, 43, 0, time.UTC),
	}
	s := &backlog.PullRequestAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	attachments, err := s.Remove(backlog.ProjectKey(projectKey), backlog.RepositoryName(repoName), prNumber, attachmentID)
	assert.NoError(t, err)
	if assert.NotNil(t, attachments) {
		assert.Equal(t, want.id, attachments.ID)
		assert.Equal(t, want.name, attachments.Name)
		assert.Equal(t, want.size, attachments.Size)
		assert.Equal(t, want.size, attachments.Size)
		assert.ObjectsAreEqualValues(want.created, attachments.Created)
	}
}

func TestPullRequestAttachmentService_Remove_param(t *testing.T) {
	bj, err := os.Open("testdata/json/attachment.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	cases := map[string]struct {
		targetProject    backlog.ProjectIDOrKeyGetter
		targetRepository backlog.RepositoryIDOrKeyGetter
		prNumber         int
		attachmentID     int
		wantError        bool
	}{
		"valid": {
			targetProject:    backlog.ProjectID(1),
			targetRepository: backlog.RepositoryName("test"),
			prNumber:         1,
			attachmentID:     8,
			wantError:        false,
		},
		"invalid-targetProject": {
			targetProject:    backlog.ProjectID(0),
			targetRepository: backlog.RepositoryName("test"),
			prNumber:         1,
			attachmentID:     8,
			wantError:        true,
		},
		"invalid-targetRepository": {
			targetProject:    backlog.ProjectID(1),
			targetRepository: backlog.RepositoryName(""),
			prNumber:         1,
			attachmentID:     8,
			wantError:        true,
		},
		"invalid-prNumber": {
			targetProject:    backlog.ProjectID(1),
			targetRepository: backlog.RepositoryName("test"),
			prNumber:         0,
			attachmentID:     8,
			wantError:        true,
		},
		"invalid-attachmentID": {
			targetProject:    backlog.ProjectID(1),
			targetRepository: backlog.RepositoryName("test"),
			prNumber:         1,
			attachmentID:     0,
			wantError:        true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := &backlog.PullRequestAttachmentService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Delete: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if attachment, err := s.Remove(tc.targetProject, tc.targetRepository, tc.prNumber, tc.attachmentID); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 8, attachment.ID)
			}
		})
	}
}

func TestPullRequestAttachmentService_Remove_clientError(t *testing.T) {
	s := &backlog.PullRequestAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	attachment, err := s.Remove(backlog.ProjectID(1234), backlog.RepositoryName("test"), 10, 8)
	assert.Error(t, err)
	assert.Nil(t, attachment)
}

func TestPullRequestAttachmentService_Remove_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.PullRequestAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	attachment, err := s.Remove(backlog.ProjectID(1234), backlog.RepositoryName("test"), 10, 8)
	assert.Error(t, err)
	assert.Nil(t, attachment)
}
