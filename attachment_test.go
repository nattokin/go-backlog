package backlog_test

import (
	"errors"
	"net/http"
	"os"
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

	want := struct {
		spath string
		id    int
		name  string
		size  int
	}{
		spath: "space/attachment",
		id:    1,
		name:  "test.txt",
		size:  8857,
	}
	cm := &backlog.ExportClientMethod{
		Uploade: func(spath, fPath, fName string) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewAttachmentService(cm)
	attachment, err := s.Uploade("fpath", "fname")
	assert.Nil(t, err)
	if assert.NotNil(t, attachment) {
		assert.Equal(t, want.id, attachment.ID)
		assert.Equal(t, want.name, attachment.Name)
		assert.Equal(t, want.size, attachment.Size)
	}
}

func TestAttachmentService_Uploade_clientError(t *testing.T) {
	cm := &backlog.ExportClientMethod{
		Uploade: func(spath, fPath, fName string) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewAttachmentService(cm)
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

	cm := &backlog.ExportClientMethod{
		Uploade: func(spath, fPath, fName string) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewAttachmentService(cm)
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

	want := struct {
		spath   string
		id      int
		name    string
		size    int
		created time.Time
	}{
		spath:   "space/attachment",
		id:      2,
		name:    "Duke.png",
		size:    196186,
		created: time.Date(2014, time.September, 11, 6, 26, 5, 0, time.UTC),
	}
	cm := &backlog.ExportClientMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewWikiAttachmentService(cm)
	attachments, err := s.Attach(1234, []int{2})
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
	cm := &backlog.ExportClientMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewWikiAttachmentService(cm)
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

	cm := &backlog.ExportClientMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewWikiAttachmentService(cm)
	attachments, err := s.Attach(1234, []int{2})
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestWikiAttachmentService_List(t *testing.T) {
	t.Log("")
}

func TestWikiAttachmentService_List_clientError(t *testing.T) {
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewWikiAttachmentService(cm)
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

	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewWikiAttachmentService(cm)
	attachments, err := s.List(1234)
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestWikiAttachmentService_Remove(t *testing.T) {
	t.Log("")
}

func TestWikiAttachmentService_Remove_clientError(t *testing.T) {
	cm := &backlog.ExportClientMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewWikiAttachmentService(cm)
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

	cm := &backlog.ExportClientMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewWikiAttachmentService(cm)
	attachment, err := s.Remove(1234, 8)
	assert.Error(t, err)
	assert.Nil(t, attachment)
}

func TestIssueAttachmentService_List(t *testing.T) {
	t.Log("")
}

func TestIssueAttachmentService_List_clientError(t *testing.T) {
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewIssueAttachmentService(cm)
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

	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewIssueAttachmentService(cm)
	attachments, err := s.List("1234")
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestIssueAttachmentService_Remove(t *testing.T) {
	t.Log("")
}

func TestIssueAttachmentService_Remove_clientError(t *testing.T) {
	cm := &backlog.ExportClientMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewIssueAttachmentService(cm)
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

	cm := &backlog.ExportClientMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewIssueAttachmentService(cm)
	attachment, err := s.Remove("1234", 8)
	assert.Error(t, err)
	assert.Nil(t, attachment)
}

func TestPullRequestAttachmentService_List(t *testing.T) {
	t.Log("")
}

func TestPullRequestAttachmentService_List_clientError(t *testing.T) {
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewPullRequestAttachmentService(cm)
	attachments, err := s.List("1234", "TEST", 10)
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestPullRequestAttachmentService_List_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewPullRequestAttachmentService(cm)
	attachments, err := s.List("1234", "TEST", 10)
	assert.Error(t, err)
	assert.Nil(t, attachments)
}

func TestPullRequestAttachmentService_Remove(t *testing.T) {
	t.Log("")
}

func TestPullRequestAttachmentService_Remove_clientError(t *testing.T) {
	cm := &backlog.ExportClientMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewPullRequestAttachmentService(cm)
	attachment, err := s.Remove("1234", "TEST", 10, 8)
	assert.Error(t, err)
	assert.Nil(t, attachment)
}

func TestPullRequestAttachmentService_Remove_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	cm := &backlog.ExportClientMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewPullRequestAttachmentService(cm)
	attachment, err := s.Remove("1234", "TEST", 10, 8)
	assert.Error(t, err)
	assert.Nil(t, attachment)
}
