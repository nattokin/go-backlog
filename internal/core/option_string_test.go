package core_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

func TestOptionService_string(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		option    core.RequestOption
		key       string
		wantValue string
		wantErr   bool
	}{
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

		"WithOrder-asc": {
			option:    o.WithOrder(model.OrderAsc),
			key:       core.ParamOrder.Value(),
			wantValue: "asc",
		},
		"WithOrder-desc": {
			option:    o.WithOrder(model.OrderDesc),
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
			option:    o.WithTextFormattingRule(model.FormatBacklog),
			key:       core.ParamTextFormattingRule.Value(),
			wantValue: string(model.FormatBacklog),
		},
		"WithTextFormattingRule-valid-markdown": {
			option:    o.WithTextFormattingRule(model.FormatMarkdown),
			key:       core.ParamTextFormattingRule.Value(),
			wantValue: string(model.FormatMarkdown),
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
			sort    model.IssueSort
			wantErr bool
		}{
			"actualHours":    {sort: model.IssueSortActualHours},
			"assignee":       {sort: model.IssueSortAssignee},
			"attachment":     {sort: model.IssueSortAttachment},
			"category":       {sort: model.IssueSortCategory},
			"childIssue":     {sort: model.IssueSortChildIssue},
			"created":        {sort: model.IssueSortCreated},
			"createdUser":    {sort: model.IssueSortCreatedUser},
			"dueDate":        {sort: model.IssueSortDueDate},
			"estimatedHours": {sort: model.IssueSortEstimatedHours},
			"issueType":      {sort: model.IssueSortIssueType},
			"milestone":      {sort: model.IssueSortMilestone},
			"priority":       {sort: model.IssueSortPriority},
			"sharedFile":     {sort: model.IssueSortSharedFile},
			"startDate":      {sort: model.IssueSortStartDate},
			"status":         {sort: model.IssueSortStatus},
			"summary":        {sort: model.IssueSortSummary},
			"updated":        {sort: model.IssueSortUpdated},
			"updatedUser":    {sort: model.IssueSortUpdatedUser},
			"version":        {sort: model.IssueSortVersion},

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
				assert.Equal(t, string(tc.sort), q.Get(core.ParamSort.Value()))
			})
		}
	})

}
