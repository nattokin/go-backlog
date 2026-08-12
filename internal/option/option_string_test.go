package option_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/option"
)

func TestOptionService_string(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		option    *option.APIParamOption
		key       string
		wantValue string
		wantErr   bool
	}{
		"WithBase-empty": {
			option:  o.WithBase(""),
			key:     option.ParamBase.Value(),
			wantErr: true,
		},
		"WithBase-whitespace": {
			option:  o.WithBase("   "),
			key:     option.ParamBase.Value(),
			wantErr: true,
		},
		"WithBase-valid": {
			option:    o.WithBase("main"),
			key:       option.ParamBase.Value(),
			wantValue: "main",
		},
		"WithBranch-empty": {
			option:  o.WithBranch(""),
			key:     option.ParamBranch.Value(),
			wantErr: true,
		},
		"WithBranch-whitespace": {
			option:  o.WithBranch("   "),
			key:     option.ParamBranch.Value(),
			wantErr: true,
		},
		"WithBranch-valid": {
			option:    o.WithBranch("feature/foo"),
			key:       option.ParamBranch.Value(),
			wantValue: "feature/foo",
		},
		"WithColor-empty": {
			option:  o.WithColor(""),
			key:     option.ParamColor.Value(),
			wantErr: true,
		},
		"WithColor-whitespace": {
			option:  o.WithColor("   "),
			key:     option.ParamColor.Value(),
			wantErr: true,
		},
		"WithColor-valid": {
			option:    o.WithColor("#e30000"),
			key:       option.ParamColor.Value(),
			wantValue: "#e30000",
		},
		"WithComment-empty": {
			option:    o.WithComment(""),
			key:       option.ParamComment.Value(),
			wantValue: "",
		},
		"WithComment-valid": {
			option:    o.WithComment("looks good"),
			key:       option.ParamComment.Value(),
			wantValue: "looks good",
		},
		"WithContent-empty": {
			option:  o.WithContent(""),
			key:     option.ParamContent.Value(),
			wantErr: true,
		},
		"WithContent-whitespace": {
			option:  o.WithContent("   "),
			key:     option.ParamContent.Value(),
			wantErr: true,
		},
		"WithContent-valid": {
			option:    o.WithContent("Hello"),
			key:       option.ParamContent.Value(),
			wantValue: "Hello",
		},
		"WithDescription-empty": {
			option:    o.WithDescription(""),
			key:       option.ParamDescription.Value(),
			wantValue: "",
		},
		"WithDescription-non-empty": {
			option:    o.WithDescription("desc"),
			key:       option.ParamDescription.Value(),
			wantValue: "desc",
		},
		"WithHookURL-empty": {
			option:  o.WithHookURL(""),
			key:     option.ParamHookURL.Value(),
			wantErr: true,
		},
		"WithHookURL-whitespace": {
			option:  o.WithHookURL("   "),
			key:     option.ParamHookURL.Value(),
			wantErr: true,
		},
		"WithHookURL-valid": {
			option:    o.WithHookURL("https://example.com/webhook"),
			key:       option.ParamHookURL.Value(),
			wantValue: "https://example.com/webhook",
		},
		"WithKey-empty": {
			option:  o.WithKey(""),
			key:     option.ParamKey.Value(),
			wantErr: true,
		},
		"WithKey-whitespace": {
			option:  o.WithKey("   "),
			key:     option.ParamKey.Value(),
			wantErr: true,
		},
		"WithKey-valid": {
			option:    o.WithKey("ABC"),
			key:       option.ParamKey.Value(),
			wantValue: "ABC",
		},
		"WithKeyword-empty": {
			option:    o.WithKeyword(""),
			key:       option.ParamKeyword.Value(),
			wantValue: "",
		},
		"WithKeyword-non-empty": {
			option:    o.WithKeyword("backlog"),
			key:       option.ParamKeyword.Value(),
			wantValue: "backlog",
		},
		"WithMailAddress-empty": {
			option:  o.WithMailAddress(""),
			key:     option.ParamMailAddress.Value(),
			wantErr: true,
		},
		"WithMailAddress-valid": {
			option:    o.WithMailAddress("test@example.com"),
			key:       option.ParamMailAddress.Value(),
			wantValue: "test@example.com",
		},
		"WithMailAddress-valid-plus": {
			option:    o.WithMailAddress("user+tag@example.co.jp"),
			key:       option.ParamMailAddress.Value(),
			wantValue: "user+tag@example.co.jp",
		},
		"WithMailAddress-invalid-no-at": {
			option:  o.WithMailAddress("notanemail"),
			key:     option.ParamMailAddress.Value(),
			wantErr: true,
		},
		"WithMailAddress-invalid-display-name": {
			option:  o.WithMailAddress("John Doe <john@example.com>"),
			key:     option.ParamMailAddress.Value(),
			wantErr: true,
		},
		"WithMailAddress-invalid-angle-bracket": {
			option:  o.WithMailAddress("<john@example.com>"),
			key:     option.ParamMailAddress.Value(),
			wantErr: true,
		},
		"WithName-empty": {
			option:  o.WithName(""),
			key:     option.ParamName.Value(),
			wantErr: true,
		},
		"WithName-whitespace": {
			option:  o.WithName("   "),
			key:     option.ParamName.Value(),
			wantErr: true,
		},
		"WithName-valid": {
			option:    o.WithName("testname"),
			key:       option.ParamName.Value(),
			wantValue: "testname",
		},
		"WithSummary-empty": {
			option:  o.WithSummary(""),
			key:     option.ParamSummary.Value(),
			wantErr: true,
		},
		"WithSummary-whitespace": {
			option:  o.WithSummary("   "),
			key:     option.ParamSummary.Value(),
			wantErr: true,
		},
		"WithSummary-valid": {
			option:    o.WithSummary("summary"),
			key:       option.ParamSummary.Value(),
			wantValue: "summary",
		},
		"WithTemplateDescription-empty": {
			option:    o.WithTemplateDescription(""),
			key:       option.ParamTemplateDescription.Value(),
			wantValue: "",
		},
		"WithTemplateDescription-non-empty": {
			option:    o.WithTemplateDescription("default description"),
			key:       option.ParamTemplateDescription.Value(),
			wantValue: "default description",
		},
		"WithTemplateSummary-empty": {
			option:    o.WithTemplateSummary(""),
			key:       option.ParamTemplateSummary.Value(),
			wantValue: "",
		},
		"WithTemplateSummary-non-empty": {
			option:    o.WithTemplateSummary("default summary"),
			key:       option.ParamTemplateSummary.Value(),
			wantValue: "default summary",
		},
		"WithUnit-empty": {
			option:    o.WithUnit(""),
			key:       option.ParamUnit.Value(),
			wantValue: "",
		},
		"WithUnit-valid": {
			option:    o.WithUnit("kg"),
			key:       option.ParamUnit.Value(),
			wantValue: "kg",
		},
		"WithOrder-asc": {
			option:    o.WithOrder("asc"),
			key:       option.ParamOrder.Value(),
			wantValue: "asc",
		},
		"WithOrder-desc": {
			option:    o.WithOrder("desc"),
			key:       option.ParamOrder.Value(),
			wantValue: "desc",
		},
		"WithOrder-empty": {
			option:  o.WithOrder(""),
			key:     option.ParamOrder.Value(),
			wantErr: true,
		},
		"WithOrder-invalid": {
			option:  o.WithOrder("invalid"),
			key:     option.ParamOrder.Value(),
			wantErr: true,
		},
		"WithPassword-invalid-empty": {
			option:  o.WithPassword(""),
			key:     option.ParamPassword.Value(),
			wantErr: true,
		},
		"WithPassword-valid-7chars": {
			option:  o.WithPassword("abcdefg"),
			key:     option.ParamPassword.Value(),
			wantErr: true,
		},
		"WithPassword-valid-8chars": {
			option:    o.WithPassword("abcdefgh"),
			key:       option.ParamPassword.Value(),
			wantValue: "abcdefgh",
		},
		"WithPassword-valid-9chars": {
			option:    o.WithPassword("abcdefghi"),
			key:       option.ParamPassword.Value(),
			wantValue: "abcdefghi",
		},
		"WithTextFormattingRule-invalid": {
			option:  o.WithTextFormattingRule("invalid"),
			key:     option.ParamTextFormattingRule.Value(),
			wantErr: true,
		},
		"WithTextFormattingRule-invalid-empty": {
			option:  o.WithTextFormattingRule(""),
			key:     option.ParamTextFormattingRule.Value(),
			wantErr: true,
		},
		"WithTextFormattingRule-valid-backlog": {
			option:    o.WithTextFormattingRule("backlog"),
			key:       option.ParamTextFormattingRule.Value(),
			wantValue: "backlog",
		},
		"WithTextFormattingRule-valid-markdown": {
			option:    o.WithTextFormattingRule("markdown"),
			key:       option.ParamTextFormattingRule.Value(),
			wantValue: "markdown",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			form := url.Values{}
			ve := tc.option.Check()
			if tc.wantErr {
				assert.NotNil(t, ve)
				return
			}
			require.Nil(t, ve)
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
				ve := opt.Check()
				if tc.wantErr {
					assert.NotNil(t, ve)
					return
				}
				require.Nil(t, ve)
				_ = opt.Set(q)
				assert.Equal(t, tc.sort, q.Get(option.ParamSort.Value()))
			})
		}
	})
}
