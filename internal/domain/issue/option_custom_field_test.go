package issue_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/domain/issue"
)

func TestWithCustomFieldItem(t *testing.T) {
	opt := issue.WithCustomFieldItem(5, 101)
	assert.Equal(t, "customField_5[]", opt.Key())
	require.NoError(t, opt.Check())

	v := url.Values{}
	require.NoError(t, opt.Set(v))
	assert.Equal(t, []string{"101"}, v["customField_5[]"])
}

func TestWithCustomFieldItem_MultipleItems(t *testing.T) {
	v := url.Values{}
	opt1 := issue.WithCustomFieldItem(5, 101)
	opt2 := issue.WithCustomFieldItem(5, 202)
	require.NoError(t, opt1.Set(v))
	require.NoError(t, opt2.Set(v))
	assert.Equal(t, []string{"101", "202"}, v["customField_5[]"])
}

func TestWithCustomFieldOther(t *testing.T) {
	opt := issue.WithCustomFieldOther(5, "other text")
	assert.Equal(t, "customField_5_otherValue", opt.Key())
	require.NoError(t, opt.Check())

	v := url.Values{}
	require.NoError(t, opt.Set(v))
	assert.Equal(t, "other text", v.Get("customField_5_otherValue"))
}
