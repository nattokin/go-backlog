package backlog_test

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	backlog "github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestAttachmentService_Uploade(t *testing.T) {
	bj, err := os.Open("testdata/json/upload_attachment.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	fpath := "fpath"
	fname := "fname"

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
		name:  "test.txt",
		size:  8857,
	}
	m := &backlog.ExportMethod{
		Uploade: func(spath, fpath, fname string) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.fpath, fpath)
			assert.Equal(t, want.fname, fname)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewAttachmentService(m)
	attachment, err := s.Uploade(fpath, fname)
	assert.Nil(t, err)
	if assert.NotNil(t, attachment) {
		assert.Equal(t, want.id, attachment.ID)
		assert.Equal(t, want.name, attachment.Name)
		assert.Equal(t, want.size, attachment.Size)
	}
}

func TestAttachmentService_Uploade_clientError(t *testing.T) {
	m := &backlog.ExportMethod{
		Uploade: func(spath, fpath, fname string) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewAttachmentService(m)
	attachement, err := s.Uploade("fpath", "fname")
	assert.Error(t, err)
	assert.Nil(t, attachement)
}

func TestAttachmentService_Uploade_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	m := &backlog.ExportMethod{
		Uploade: func(spath, fpath, fname string) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewAttachmentService(m)
	attachement, err := s.Uploade("fpath", "fname")
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
		name:    "Duke.png",
		size:    196186,
		created: time.Date(2014, time.September, 11, 6, 26, 5, 0, time.UTC),
	}
	m := &backlog.ExportMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			v := *params.ExportURLValues()
			assert.Equal(t, []string{"2"}, v["attachmentId[]"])
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewWikiAttachmentService(m)
	attachments, err := s.Attach(wikiID, []int{2})
	assert.Nil(t, err)
	if assert.NotNil(t, attachments) {
		assert.Equal(t, want.id, attachments[0].ID)
		assert.Equal(t, want.name, attachments[0].Name)
		assert.Equal(t, want.size, attachments[0].Size)
		assert.Equal(t, want.size, attachments[0].Size)
		assert.ObjectsAreEqualValues(want.created, attachments[0].Created)
	}
}

func TestWikiAttachmentService_Attach_clientError(t *testing.T) {
	m := &backlog.ExportMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewWikiAttachmentService(m)
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

	m := &backlog.ExportMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewWikiAttachmentService(m)
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
		name:    "Duke.png",
		size:    196186,
		created: time.Date(2014, time.September, 11, 6, 26, 5, 0, time.UTC),
	}
	m := &backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewWikiAttachmentService(m)
	attachments, err := s.List(wikiID)
	assert.Nil(t, err)
	if assert.NotNil(t, attachments) {
		assert.Equal(t, want.id, attachments[0].ID)
		assert.Equal(t, want.name, attachments[0].Name)
		assert.Equal(t, want.size, attachments[0].Size)
		assert.Equal(t, want.size, attachments[0].Size)
		assert.ObjectsAreEqualValues(want.created, attachments[0].Created)
	}
}

func TestWikiAttachmentService_List_clientError(t *testing.T) {
	m := &backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewWikiAttachmentService(m)
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

	m := &backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewWikiAttachmentService(m)
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
	m := &backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewWikiAttachmentService(m)
	attachments, err := s.Remove(wikiID, attachmentID)
	assert.Nil(t, err)
	if assert.NotNil(t, attachments) {
		assert.Equal(t, want.id, attachments.ID)
		assert.Equal(t, want.name, attachments.Name)
		assert.Equal(t, want.size, attachments.Size)
		assert.Equal(t, want.size, attachments.Size)
		assert.ObjectsAreEqualValues(want.created, attachments.Created)
	}
}

func TestWikiAttachmentService_Remove_clientError(t *testing.T) {
	m := &backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewWikiAttachmentService(m)
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

	m := &backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewWikiAttachmentService(m)
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

	issueIDOrKey := "1234"

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   "issues/" + issueIDOrKey + "/attachments",
		id:      2,
		name:    "Duke.png",
		size:    196186,
		created: time.Date(2014, time.September, 11, 6, 26, 5, 0, time.UTC),
	}
	m := &backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewIssueAttachmentService(m)
	attachments, err := s.List(issueIDOrKey)
	assert.Nil(t, err)
	if assert.NotNil(t, attachments) {
		assert.Equal(t, want.id, attachments[0].ID)
		assert.Equal(t, want.name, attachments[0].Name)
		assert.Equal(t, want.size, attachments[0].Size)
		assert.Equal(t, want.size, attachments[0].Size)
		assert.ObjectsAreEqualValues(want.created, attachments[0].Created)
	}
}

func TestIssueAttachmentService_List_clientError(t *testing.T) {
	m := &backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewIssueAttachmentService(m)
	attachments, err := s.List("1234")
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestIssueAttachmentService_List_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	m := &backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewIssueAttachmentService(m)
	attachments, err := s.List("1234")
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestIssueAttachmentService_Remove(t *testing.T) {
	bj, err := os.Open("testdata/json/attachment.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	issueIDOrKey := "1234"
	attachmentID := 8

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   "issues/" + issueIDOrKey + "/attachments/" + strconv.Itoa(attachmentID),
		id:      8,
		name:    "IMG0088.png",
		size:    5563,
		created: time.Date(2014, time.October, 28, 9, 24, 43, 0, time.UTC),
	}
	m := &backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewIssueAttachmentService(m)
	attachments, err := s.Remove(issueIDOrKey, attachmentID)
	assert.Nil(t, err)
	if assert.NotNil(t, attachments) {
		assert.Equal(t, want.id, attachments.ID)
		assert.Equal(t, want.name, attachments.Name)
		assert.Equal(t, want.size, attachments.Size)
		assert.Equal(t, want.size, attachments.Size)
		assert.ObjectsAreEqualValues(want.created, attachments.Created)
	}
}

func TestIssueAttachmentService_Remove_clientError(t *testing.T) {
	m := &backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewIssueAttachmentService(m)
	attachment, err := s.Remove("1234", 8)
	assert.Error(t, err)
	assert.Nil(t, attachment)
}

func TestIssueAttachmentService_Remove_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	m := &backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewIssueAttachmentService(m)
	attachment, err := s.Remove("1234", 8)
	assert.Error(t, err)
	assert.Nil(t, attachment)
}

func TestPullRequestAttachmentService_List(t *testing.T) {
	bj, err := os.Open("testdata/json/attachment_list.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	projectIDOrKey := "1234"
	repoIDOrName := "test"
	prNumber := 1234

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   "projects/" + projectIDOrKey + "/git/repositories/" + repoIDOrName + "/pullRequests/" + strconv.Itoa(prNumber) + "/attachments",
		id:      2,
		name:    "Duke.png",
		size:    196186,
		created: time.Date(2014, time.September, 11, 6, 26, 5, 0, time.UTC),
	}
	m := &backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewPullRequestAttachmentService(m)
	attachments, err := s.List(projectIDOrKey, repoIDOrName, prNumber)
	assert.Nil(t, err)
	if assert.NotNil(t, attachments) {
		assert.Equal(t, want.id, attachments[0].ID)
		assert.Equal(t, want.name, attachments[0].Name)
		assert.Equal(t, want.size, attachments[0].Size)
		assert.Equal(t, want.size, attachments[0].Size)
		assert.ObjectsAreEqualValues(want.created, attachments[0].Created)
	}
}

func TestPullRequestAttachmentService_List_clientError(t *testing.T) {
	m := &backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewPullRequestAttachmentService(m)
	attachments, err := s.List("1234", "test", 10)
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestPullRequestAttachmentService_List_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	m := &backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewPullRequestAttachmentService(m)
	attachments, err := s.List("1234", "test", 10)
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestPullRequestAttachmentService_Remove(t *testing.T) {
	bj, err := os.Open("testdata/json/attachment.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	projectIDOrKey := "1234"
	repoIDOrName := "test"
	prNumber := 1234
	attachmentID := 8

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   "projects/" + projectIDOrKey + "/git/repositories/" + repoIDOrName + "/pullRequests/" + strconv.Itoa(prNumber) + "/attachments" + strconv.Itoa(attachmentID),
		id:      8,
		name:    "IMG0088.png",
		size:    5563,
		created: time.Date(2014, time.October, 28, 9, 24, 43, 0, time.UTC),
	}
	m := &backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewPullRequestAttachmentService(m)
	attachments, err := s.Remove(projectIDOrKey, repoIDOrName, prNumber, attachmentID)
	assert.Nil(t, err)
	if assert.NotNil(t, attachments) {
		assert.Equal(t, want.id, attachments.ID)
		assert.Equal(t, want.name, attachments.Name)
		assert.Equal(t, want.size, attachments.Size)
		assert.Equal(t, want.size, attachments.Size)
		assert.ObjectsAreEqualValues(want.created, attachments.Created)
	}
}

func TestPullRequestAttachmentService_Remove_clientError(t *testing.T) {
	m := &backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewPullRequestAttachmentService(m)
	attachment, err := s.Remove("1234", "test", 10, 8)
	assert.Error(t, err)
	assert.Nil(t, attachment)
}

func TestPullRequestAttachmentService_Remove_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	m := &backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewPullRequestAttachmentService(m)
	attachment, err := s.Remove("1234", "test", 10, 8)
	assert.Error(t, err)
	assert.Nil(t, attachment)
}
