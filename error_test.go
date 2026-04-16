package backlog_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/core"
)

// convertError is unexported, so tests drive it indirectly via wiki.All which
// calls convertError on every error path. For the root-package wrapper types
// themselves we construct them by passing a core error through a real service
// call that returns convertError's output, using errors.As to extract the
// typed value.
//
// Helper: run a Wiki.All call with a doer that returns a given HTTP status,
// and return the resulting error.
func callWikiAllWithStatus(t *testing.T, statusCode int) error {
	t.Helper()
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(&mockDoer{do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"not found","code":6,"moreInfo":""}]}`)),
			}, nil
		}}),
	)
	require.NoError(t, err)
	_, err = c.Wiki.All(context.Background(), "PROJECT")
	return err
}

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

// callWikiAllWithInvalidOption drives convertError via an invalid option key,
// which the service layer converts to *backlog.InvalidOptionKeyError.
func callWikiAllWithInvalidOption(t *testing.T) error {
	t.Helper()
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	// WithCount is not a valid option for Wiki.All — triggers InvalidOptionKeyError.
	_, err = c.Wiki.All(context.Background(), "PROJECT", c.Wiki.Option.WithContent("x"))
	return err
}

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
//  ValidationError
// ──────────────────────────────────────────────────────────────

func TestValidationError_Error(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	// wikiID=0 triggers a ValidationError in the internal layer.
	_, err = c.Wiki.One(context.Background(), 0)
	require.Error(t, err)

	var target *backlog.ValidationError
	require.True(t, errors.As(err, &target))
	assert.NotEmpty(t, target.Error())
}

// ──────────────────────────────────────────────────────────────
//  InternalClientError
// ──────────────────────────────────────────────────────────────

func TestInternalClientError_Error(t *testing.T) {
	// Empty baseURL triggers InternalClientError from NewClient.
	_, err := backlog.NewClient("", "token")
	require.Error(t, err)

	var target *backlog.InternalClientError
	require.True(t, errors.As(err, &target))
	assert.NotEmpty(t, target.Error())
}
