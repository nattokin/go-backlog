package backlog_test

import (
	"errors"
	"io"
	"net/http"
	"os"
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
		Upload: func(spath, fileName string, r io.Reader) (*http.Response, error) {
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

	f, err := os.Open("testdata/testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	attachment, err := s.Upload(fpath, f)
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
		Upload: func(spath, fileName string, r io.Reader) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})

	f, err := os.Open("testdata/testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	attachement, err := s.Upload("fpath", f)
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
		Upload: func(spath, fileName string, r io.Reader) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})

	f, err := os.Open("testdata/testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	attachement, err := s.Upload("fpath", f)
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
	spath := "wikis/1234/attachments"

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   spath,
		id:      2,
		name:    "A.png",
		size:    196186,
		created: time.Date(2014, time.September, 11, 6, 26, 5, 0, time.UTC),
	}
	s := &backlog.WikiAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			v := *form.ExportURLValues()
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

func TestWikiAttachmentService_Attach_param(t *testing.T) {
	cases := map[string]struct {
		wikiID        int
		attachmentIDs []int
		wantError     bool
	}{
		"valid": {
			wikiID:        1,
			attachmentIDs: []int{1, 2},
			wantError:     false,
		},
		"wikiID_invalid": {
			wikiID:        0,
			attachmentIDs: []int{1, 2},
			wantError:     true,
		},
		"attachmentIDs_invalid": {
			wikiID:        1,
			attachmentIDs: []int{0, 1, 2},
			wantError:     true,
		},
		"attachmentIDs_empty": {
			wikiID:        1,
			attachmentIDs: []int{},
			wantError:     true,
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
				Post: func(spath string, form *backlog.FormParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if attachements, err := s.Attach(tc.wikiID, tc.attachmentIDs); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, attachements, 2)
			}
		})
	}
}

func TestWikiAttachmentService_Attach_clientError(t *testing.T) {
	s := &backlog.WikiAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
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
		Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
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
	spath := "wikis/1234/attachments"

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   spath,
		id:      2,
		name:    "A.png",
		size:    196186,
		created: time.Date(2014, time.September, 11, 6, 26, 5, 0, time.UTC),
	}
	s := &backlog.WikiAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
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
				Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
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
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
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
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
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
	spath := "wikis/1234/attachments/8"

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   spath,
		id:      8,
		name:    "IMG0088.png",
		size:    5563,
		created: time.Date(2014, time.October, 28, 9, 24, 43, 0, time.UTC),
	}
	s := &backlog.WikiAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
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
				Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
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
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
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
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
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
	spath := "issues/1234/attachments"

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   spath,
		id:      2,
		name:    "A.png",
		size:    196186,
		created: time.Date(2014, time.September, 11, 6, 26, 5, 0, time.UTC),
	}
	s := &backlog.IssueAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
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
		issueID   backlog.IssueID
		wantError bool
	}{
		"valid": {
			issueID:   1,
			wantError: false,
		},
		"invalid": {
			issueID:   0,
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
				Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if attachements, err := s.List(tc.issueID); tc.wantError {
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
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
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
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
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
	spath := "issues/1234/attachments/8"

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   spath,
		id:      attachmentID,
		name:    "IMG0088.png",
		size:    5563,
		created: time.Date(2014, time.October, 28, 9, 24, 43, 0, time.UTC),
	}

	s := &backlog.IssueAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
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
		issueKey     backlog.IssueKey
		attachmentID int
		wantError    bool
	}{
		"valid": {
			issueKey:     backlog.IssueKey("test"),
			attachmentID: 8,
			wantError:    false,
		},
		"invalid-issueKey": {
			issueKey:     backlog.IssueKey(""),
			attachmentID: 8,
			wantError:    true,
		},
		"invalid-attachmentID": {
			issueKey:     backlog.IssueKey("test"),
			attachmentID: 0,
			wantError:    true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := &backlog.IssueAttachmentService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if attachements, err := s.Remove(tc.issueKey, tc.attachmentID); tc.wantError {
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
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
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
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
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
	spath := "projects/TEST/git/repositories/test/pullRequests/1234/attachments"

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   spath,
		id:      2,
		name:    "A.png",
		size:    196186,
		created: time.Date(2014, time.September, 11, 6, 26, 5, 0, time.UTC),
	}
	s := &backlog.PullRequestAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
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
		project    backlog.ProjectID
		repository backlog.RepositoryID
		prNumber   int
		wantError  bool
	}{
		"valid": {
			project:    1,
			repository: 1,
			prNumber:   1,
			wantError:  false,
		},
		"invalid-project": {
			project:    0,
			repository: 1,
			prNumber:   1,
			wantError:  true,
		},
		"invalid-repository": {
			project:    1,
			repository: 0,
			prNumber:   1,
			wantError:  true,
		},
		"invalid-prNumber": {
			project:    1,
			repository: 1,
			prNumber:   0,
			wantError:  true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := &backlog.PullRequestAttachmentService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if attachements, err := s.List(tc.project, tc.repository, tc.prNumber); tc.wantError {
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
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
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
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
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
	spath := "projects/TEST/git/repositories/test/pullRequests/1234/attachments/8"

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   spath,
		id:      8,
		name:    "IMG0088.png",
		size:    5563,
		created: time.Date(2014, time.October, 28, 9, 24, 43, 0, time.UTC),
	}
	s := &backlog.PullRequestAttachmentService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
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
		project      backlog.ProjectID
		repository   backlog.RepositoryName
		prNumber     int
		attachmentID int
		wantError    bool
	}{
		"valid": {
			project:      1,
			repository:   "test",
			prNumber:     1,
			attachmentID: 8,
			wantError:    false,
		},
		"invalid-project": {
			project:      0,
			repository:   "test",
			prNumber:     1,
			attachmentID: 8,
			wantError:    true,
		},
		"invalid-repository": {
			project:      1,
			repository:   "",
			prNumber:     1,
			attachmentID: 8,
			wantError:    true,
		},
		"invalid-prNumber": {
			project:      1,
			repository:   "test",
			prNumber:     0,
			attachmentID: 8,
			wantError:    true,
		},
		"invalid-attachmentID": {
			project:      1,
			repository:   "test",
			prNumber:     1,
			attachmentID: 0,
			wantError:    true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := &backlog.PullRequestAttachmentService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if attachment, err := s.Remove(tc.project, tc.repository, tc.prNumber, tc.attachmentID); tc.wantError {
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
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
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
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
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
