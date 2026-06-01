package backlog_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

// ──────────────────────────────────────────────────────────────
//  APIResponseError
// ──────────────────────────────────────────────────────────────

func TestAPIResponseError_Error(t *testing.T) {
	err := callWikiAllWithStatus(t, 404)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Status Code:404")
	assert.Contains(t, err.Error(), "not found")
}

func TestAPIResponseError_StatusCode(t *testing.T) {
	err := callWikiAllWithStatus(t, 403)
	require.Error(t, err)

	var target *backlog.APIResponseError
	require.True(t, errors.As(err, &target))
	assert.Equal(t, 403, target.StatusCode())
}

func TestAPIResponseError_Errors(t *testing.T) {
	err := callWikiAllWithStatus(t, 404)
	require.Error(t, err)

	var target *backlog.APIResponseError
	require.True(t, errors.As(err, &target))

	errs := target.Errors()
	require.Len(t, errs, 1)
	assert.Equal(t, "not found", errs[0].Message)
	assert.Equal(t, 6, errs[0].Code)
}

// ──────────────────────────────────────────────────────────────
//  InvalidOptionKeyError
// ──────────────────────────────────────────────────────────────

func TestInvalidOptionKeyError_Error(t *testing.T) {
	err := callWikiAllWithInvalidOption(t)
	require.Error(t, err)

	var target *backlog.InvalidOptionKeyError
	require.True(t, errors.As(err, &target))
	assert.Contains(t, target.Error(), "invalid option key")
}

func TestInvalidOptionKeyError_InvalidKey(t *testing.T) {
	err := callWikiAllWithInvalidOption(t)
	require.Error(t, err)

	var target *backlog.InvalidOptionKeyError
	require.True(t, errors.As(err, &target))
	assert.Equal(t, core.ParamContent.Value(), target.InvalidKey())
}

func TestInvalidOptionKeyError_AllowKeys(t *testing.T) {
	err := callWikiAllWithInvalidOption(t)
	require.Error(t, err)

	var target *backlog.InvalidOptionKeyError
	require.True(t, errors.As(err, &target))
	assert.NotEmpty(t, target.AllowKeys())
	assert.Contains(t, target.AllowKeys(), core.ParamKeyword.Value())
}

// ──────────────────────────────────────────────────────────────
//  InvalidOptionError
// ──────────────────────────────────────────────────────────────

func TestInvalidOptionError_Error(t *testing.T) {
	err := callIssueCommentListWithNilOption(t)
	require.Error(t, err)

	var target *backlog.InvalidOptionError
	require.True(t, errors.As(err, &target))
	assert.NotEmpty(t, target.Error())
}

// ──────────────────────────────────────────────────────────────
//  ValidationError
// ──────────────────────────────────────────────────────────────

func TestValidationError_Error(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	// wikiID=0 is invalid and triggers a ValidationError in the internal layer.
	_, err = c.Wiki.One(context.Background(), 0)
	require.Error(t, err)

	var target *backlog.ValidationError
	require.True(t, errors.As(err, &target))
	assert.NotEmpty(t, target.Error())
}

func TestValidationErrors_Error(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	// issueIDOrKey="" and WithCount(0) each produce a ValidationError,
	// so convertError receives a ValidationErrors (2 elements) and joins them.
	_, err = c.Issue.Comment.List(context.Background(), "", c.Issue.Comment.Option.WithCount(0))
	require.Error(t, err)

	// errors.Join result: each element is reachable via errors.As
	var ve *backlog.ValidationError
	assert.True(t, errors.As(err, &ve))
	assert.NotEmpty(t, err.Error())
}

// ──────────────────────────────────────────────────────────────
//  InternalClientError
// ──────────────────────────────────────────────────────────────

func TestInternalClientError_Error(t *testing.T) {
	// An empty baseURL triggers InternalClientError from NewClient.
	_, err := backlog.NewClient("", "token")
	require.Error(t, err)

	var target *backlog.InternalClientError
	require.True(t, errors.As(err, &target))
	assert.NotEmpty(t, target.Error())
}

// ──────────────────────────────────────────────────────────────
//  InvalidDateStringError
// ──────────────────────────────────────────────────────────────

func TestInvalidDateStringError_Error(t *testing.T) {
	t.Parallel()

	_, err := backlog.NewDate("2024/03/31")
	require.Error(t, err)
	assert.Equal(t, `backlog: invalid date string "2024/03/31": expected "YYYY-MM-DD" format`, err.Error())
}

// ──────────────────────────────────────────────────────────────
//  convertError (indirect via service methods)
// ──────────────────────────────────────────────────────────────

func Test_convertError_default_passthroughsUnknownError(t *testing.T) {
	sentinel := errors.New("network error")
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(&mock.Doer{DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, sentinel
		}}),
	)
	require.NoError(t, err)
	_, err = c.Wiki.List(context.Background(), "PROJECT")
	assert.True(t, errors.Is(err, sentinel))
}

func Test_convertError_InvalidOptionError(t *testing.T) {
	err := callIssueCommentListWithNilOption(t)
	require.Error(t, err)

	var target *backlog.InvalidOptionError
	assert.True(t, errors.As(err, &target))
}

// ──────────────────────────────────────────────────────────────
//  Test helpers
// ──────────────────────────────────────────────────────────────

// convertError is unexported, so tests drive it indirectly via service methods
// which call convertError on every error path. errors.As is used to extract
// the typed wrapper value for assertion.

// callWikiAllWithStatus runs Wiki.All with a doer that returns the given HTTP
// status code and a single-element errors array, then returns the error.
func callWikiAllWithStatus(t *testing.T, statusCode int) error {
	t.Helper()
	body := `{"errors":[{"message":"not found","code":6,"moreInfo":""}]}`
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(&mock.Doer{DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(strings.NewReader(body)),
			}, nil
		}}),
	)
	require.NoError(t, err)
	_, err = c.Wiki.List(context.Background(), "PROJECT")
	return err
}

// callWikiAllWithInvalidOption drives convertError via an invalid option key.
// WithContent is not valid for Wiki.All, triggering InvalidOptionKeyError.
func callWikiAllWithInvalidOption(t *testing.T) error {
	t.Helper()
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	_, err = c.Wiki.List(context.Background(), "PROJECT", c.Wiki.Option.WithContent("x"))
	return err
}

// callIssueCommentListWithNilOption passes a nil RequestOption to trigger InvalidOptionError.
func callIssueCommentListWithNilOption(t *testing.T) error {
	t.Helper()
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	_, err = c.Issue.Comment.List(context.Background(), "PRJ-1", nil)
	return err
}
