package core_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
)

func TestOptionService_string(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		option    core.RequestOption
		key       string
		wantValue string
		wantErr   bool
	}{
		"WithBase-empty": {
			option:  o.WithBase(""),
			key:     core.ParamBase.Value(),
			wantErr: true,
		},
		"WithBase-valid": {
			option:    o.WithBase("main"),
			key:       core.ParamBase.Value(),
			wantValue: "main",
		},
		"WithBranch-empty": {
			option:  o.WithBranch(""),
			key:     core.ParamBranch.Value(),
			wantErr: true,
		},
		"WithBranch-valid": {
			option:    o.WithBranch("feature/foo"),
			key:       core.ParamBranch.Value(),
			wantValue: "feature/foo",
		},
		"WithComment-empty": {
			option:    o.WithComment(""),
			key:       core.ParamComment.Value(),
			wantValue: "",
		},
		"WithComment-valid": {
			option:    o.WithComment("looks good"),
			key:       core.ParamComment.Value(),
			wantValue: "looks good",
		},
		"WithContent-empty": {
			option:  o.WithContent(""),
			key:     core.ParamContent.Value(),
			wantErr: true,
		},
		"WithContent-valid": {
			option:    o.WithContent("Hello"),
			key:       core.ParamContent.Value(),
			wantValue: "Hello",
		},
		"WithDescription-empty": {
			option:    o.WithDescription(""),
			key:       core.ParamDescription.Value(),
			wantValue: "",
		},
		"WithDescription-non-empty": {
			option:    o.WithDescription("desc"),
			key:       core.ParamDescription.Value(),
			wantValue: "desc",
		},
		"WithHookURL-empty": {
			option:  o.WithHookURL(""),
			key:     core.ParamHookURL.Value(),
			wantErr: true,
		},
		"WithHookURL-valid": {
			option:    o.WithHookURL("https://example.com/webhook"),
			key:       core.ParamHookURL.Value(),
			wantValue: "https://example.com/webhook",
		},
		"WithKey-empty": {
			option:  o.WithKey(""),
			key:     core.ParamKey.Value(),
			wantErr: true,
		},
		"WithKey-valid": {
			option:    o.WithKey("ABC"),
			key:       core.ParamKey.Value(),
			wantValue: "ABC",
		},
		"WithKeyword-empty": {
			option:    o.WithKeyword(""),
			key:       core.ParamKeyword.Value(),
			wantValue: "",
		},
		"WithKeyword-non-empty": {
			option:    o.WithKeyword("backlog"),
			key:       core.ParamKeyword.Value(),
			wantValue: "backlog",
		},
		"WithMailAddress-empty": {
			option:  o.WithMailAddress(""),
			key:     core.ParamMailAddress.Value(),
			wantErr: true,
		},
		"WithMailAddress-valid": {
			option:    o.WithMailAddress("test@example.com"),
			key:       core.ParamMailAddress.Value(),
			wantValue: "test@example.com",
		},
		"WithMailAddress-valid-plus": {
			option:    o.WithMailAddress("user+tag@example.co.jp"),
			key:       core.ParamMailAddress.Value(),
			wantValue: "user+tag@example.co.jp",
		},
		"WithMailAddress-invalid-no-at": {
			option:  o.WithMailAddress("notanemail"),
			key:     core.ParamMailAddress.Value(),
			wantErr: true,
		},
		"WithMailAddress-invalid-display-name": {
			option:  o.WithMailAddress("John Doe <john@example.com>"),
			key:     core.ParamMailAddress.Value(),
			wantErr: true,
		},
		"WithMailAddress-invalid-angle-bracket": {
			option:  o.WithMailAddress("<john@example.com>"),
			key:     core.ParamMailAddress.Value(),
			wantErr: true,
		},
		"WithName-empty": {
			option:  o.WithName(""),
			key:     core.ParamName.Value(),
			wantErr: true,
		},
		"WithName-valid": {
			option:    o.WithName("testname"),
			key:       core.ParamName.Value(),
			wantValue: "testname",
		},
		"WithSummary-empty": {
			option:  o.WithSummary(""),
			key:     core.ParamSummary.Value(),
			wantErr: true,
		},
		"WithSummary-valid": {
			option:    o.WithSummary("summary"),
			key:       core.ParamSummary.Value(),
			wantValue: "summary",
		},
		"WithTemplateDescription-empty": {
			option:    o.WithTemplateDescription(""),
			key:       core.ParamTemplateDescription.Value(),
			wantValue: "",
		},
		"WithTemplateDescription-non-empty": {
			option:    o.WithTemplateDescription("default description"),
			key:       core.ParamTemplateDescription.Value(),
			wantValue: "default description",
		},
		"WithTemplateSummary-empty": {
			option:    o.WithTemplateSummary(""),
			key:       core.ParamTemplateSummary.Value(),
			wantValue: "",
		},
		"WithTemplateSummary-non-empty": {
			option:    o.WithTemplateSummary("default summary"),
			key:       core.ParamTemplateSummary.Value(),
			wantValue: "default summary",
		},
		"WithUnit-empty": {
			option:    o.WithUnit(""),
			key:       core.ParamUnit.Value(),
			wantValue: "",
		},
		"WithUnit-valid": {
			option:    o.WithUnit("kg"),
			key:       core.ParamUnit.Value(),
			wantValue: "kg",
		},
		"WithOrder-asc": {
			option:    o.WithOrder("asc"),
			key:       core.ParamOrder.Value(),
			wantValue: "asc",
		},
		"WithOrder-desc": {
			option:    o.WithOrder("desc"),
			key:       core.ParamOrder.Value(),
			wantValue: "desc",
		},
		"WithOrder-empty": {
			option:  o.WithOrder(""),
			key:     core.ParamOrder.Value(),
			wantErr: true,
		},
		"WithOrder-invalid": {
			option:  o.WithOrder("invalid"),
			key:     core.ParamOrder.Value(),
			wantErr: true,
		},
		"WithPassword-invalid-empty": {
			option:  o.WithPassword(""),
			key:     core.ParamPassword.Value(),
			wantErr: true,
		},
		"WithPassword-valid-7chars": {
			option:  o.WithPassword("abcdefg"),
			key:     core.ParamPassword.Value(),
			wantErr: true,
		},
		"WithPassword-valid-8chars": {
			option:    o.WithPassword("abcdefgh"),
			key:       core.ParamPassword.Value(),
			wantValue: "abcdefgh",
		},
		"WithPassword-valid-9chars": {
			option:    o.WithPassword("abcdefghi"),
			key:       core.ParamPassword.Value(),
			wantValue: "abcdefghi",
		},
		"WithTextFormattingRule-invalid": {
			option:  o.WithTextFormattingRule("invalid"),
			key:     core.ParamTextFormattingRule.Value(),
			wantErr: true,
		},
		"WithTextFormattingRule-invalid-empty": {
			option:  o.WithTextFormattingRule(""),
			key:     core.ParamTextFormattingRule.Value(),
			wantErr: true,
		},
		"WithTextFormattingRule-valid-backlog": {
			option:    o.WithTextFormattingRule("backlog"),
			key:       core.ParamTextFormattingRule.Value(),
			wantValue: "backlog",
		},
		"WithTextFormattingRule-valid-markdown": {
			option:    o.WithTextFormattingRule("markdown"),
			key:       core.ParamTextFormattingRule.Value(),
			wantValue: "markdown",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			form := url.Values{}
			err := tc.option.Check()
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			_ = tc.option.Set(form)
			assert.Equal(t, tc.wantValue, form.Get(tc.key))
		})
	}

	// --- IssueSort option ---------------------------------------------------------
	t.Run("WithIssueSort", func(t *testing.T) {
		cases := map[string]struct {
			sort    string
			wantErr bool
		}{
			"actualHours":    {sort: "actualHours"},
			"assignee":       {sort: "assignee"},
			"attachment":     {sort: "attachment"},
			"category":       {sort: "category"},
			"childIssue":     {sort: "childIssue"},
			"created":        {sort: "created"},
			"createdUser":    {sort: "createdUser"},
			"dueDate":        {sort: "dueDate"},
			"estimatedHours": {sort: "estimatedHours"},
			"issueType":      {sort: "issueType"},
			"milestone":      {sort: "milestone"},
			"priority":       {sort: "priority"},
			"sharedFile":     {sort: "sharedFile"},
			"startDate":      {sort: "startDate"},
			"status":         {sort: "status"},
			"summary":        {sort: "summary"},
			"updated":        {sort: "updated"},
			"updatedUser":    {sort: "updatedUser"},
			"version":        {sort: "version"},

			"empty":   {sort: "", wantErr: true},
			"invalid": {sort: "invalid", wantErr: true},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				opt := o.WithIssueSort(tc.sort)
				q := url.Values{}
				err := opt.Check()
				if tc.wantErr {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				_ = opt.Set(q)
				assert.Equal(t, tc.sort, q.Get(core.ParamSort.Value()))
			})
		}
	})
}
