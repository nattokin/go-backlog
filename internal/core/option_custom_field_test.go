package core_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithCustomField_String(t *testing.T) {
	opt := core.WithCustomField(1, "hello")
	assert.Equal(t, "customField_1", opt.Key())
	require.NoError(t, opt.Check())

	v := url.Values{}
	require.NoError(t, opt.Set(v))
	assert.Equal(t, "hello", v.Get("customField_1"))
}

func TestWithCustomField_Int(t *testing.T) {
	opt := core.WithCustomField(42, 99)
	assert.Equal(t, "customField_42", opt.Key())

	v := url.Values{}
	require.NoError(t, opt.Set(v))
	assert.Equal(t, "99", v.Get("customField_42"))
}

func TestWithCustomField_Float64(t *testing.T) {
	opt := core.WithCustomField(3, 1.5)
	assert.Equal(t, "customField_3", opt.Key())

	v := url.Values{}
	require.NoError(t, opt.Set(v))
	assert.Equal(t, "1.5", v.Get("customField_3"))
}

func TestWithCustomField_Time(t *testing.T) {
	date := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	opt := core.WithCustomField(7, date)
	assert.Equal(t, "customField_7", opt.Key())

	v := url.Values{}
	require.NoError(t, opt.Set(v))
	assert.Equal(t, "2024-03-15", v.Get("customField_7"))
}

func TestWithCustomFieldItem(t *testing.T) {
	opt := core.WithCustomFieldItem(5, 101)
	assert.Equal(t, "customField_5[]", opt.Key())
	require.NoError(t, opt.Check())

	v := url.Values{}
	require.NoError(t, opt.Set(v))
	assert.Equal(t, []string{"101"}, v["customField_5[]"])
}

func TestWithCustomFieldItem_MultipleItems(t *testing.T) {
	v := url.Values{}
	opt1 := core.WithCustomFieldItem(5, 101)
	opt2 := core.WithCustomFieldItem(5, 202)
	require.NoError(t, opt1.Set(v))
	require.NoError(t, opt2.Set(v))
	assert.Equal(t, []string{"101", "202"}, v["customField_5[]"])
}

func TestWithCustomFieldOther(t *testing.T) {
	opt := core.WithCustomFieldOther(5, "other text")
	assert.Equal(t, "customField_5_otherValue", opt.Key())
	require.NoError(t, opt.Check())

	v := url.Values{}
	require.NoError(t, opt.Set(v))
	assert.Equal(t, "other text", v.Get("customField_5_otherValue"))
}

func TestSplitCustomFieldOptions(t *testing.T) {
	option := &core.OptionService{}
	normal1 := option.WithSummary("test")
	normal2 := option.WithDescription("desc")
	cf1 := core.WithCustomField(1, "val")
	cf2 := core.WithCustomFieldItem(2, 10)
	cf3 := core.WithCustomFieldOther(3, "other")

	regular, custom := core.SplitCustomFieldOptions([]core.RequestOption{normal1, cf1, normal2, cf2, cf3})

	assert.Len(t, regular, 2)
	assert.Equal(t, normal1.Key(), regular[0].Key())
	assert.Equal(t, normal2.Key(), regular[1].Key())

	assert.Len(t, custom, 3)
	assert.Equal(t, "customField_1", custom[0].Key())
	assert.Equal(t, "customField_2[]", custom[1].Key())
	assert.Equal(t, "customField_3_otherValue", custom[2].Key())
}

func TestApplyCustomFieldOptions(t *testing.T) {
	v := url.Values{}
	cf1 := core.WithCustomFieldItem(10, 100)
	cf2 := core.WithCustomFieldOther(10, "custom text")

	_, customs := core.SplitCustomFieldOptions([]core.RequestOption{cf1, cf2})
	require.NoError(t, core.ApplyCustomFieldOptions(v, customs))

	assert.Equal(t, []string{"100"}, v["customField_10[]"])
	assert.Equal(t, "custom text", v.Get("customField_10_otherValue"))
}
