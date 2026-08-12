package option_test

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/option"
)

func TestOptionService_bool(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		option    *option.APIParamOption
		key       string
		wantValue bool
	}{
		"WithAll-false": {
			option:    o.WithAll(false),
			key:       option.ParamAll.Value(),
			wantValue: false,
		},
		"WithAll-true": {
			option:    o.WithAll(true),
			key:       option.ParamAll.Value(),
			wantValue: true,
		},
		"WithAllEvent-false": {
			option:    o.WithAllEvent(false),
			key:       option.ParamAllEvent.Value(),
			wantValue: false,
		},
		"WithAllEvent-true": {
			option:    o.WithAllEvent(true),
			key:       option.ParamAllEvent.Value(),
			wantValue: true,
		},
		"WithAllowAddItem-false": {
			option:    o.WithAllowAddItem(false),
			key:       option.ParamAllowAddItem.Value(),
			wantValue: false,
		},
		"WithAllowAddItem-true": {
			option:    o.WithAllowAddItem(true),
			key:       option.ParamAllowAddItem.Value(),
			wantValue: true,
		},
		"WithAllowInput-false": {
			option:    o.WithAllowInput(false),
			key:       option.ParamAllowInput.Value(),
			wantValue: false,
		},
		"WithAllowInput-true": {
			option:    o.WithAllowInput(true),
			key:       option.ParamAllowInput.Value(),
			wantValue: true,
		},
		"WithArchived-false": {
			option:    o.WithArchived(false),
			key:       option.ParamArchived.Value(),
			wantValue: false,
		},
		"WithArchived-true": {
			option:    o.WithArchived(true),
			key:       option.ParamArchived.Value(),
			wantValue: true,
		},
		"WithAttachment-false": {
			option:    o.WithAttachment(false),
			key:       option.ParamAttachment.Value(),
			wantValue: false,
		},
		"WithAttachment-true": {
			option:    o.WithAttachment(true),
			key:       option.ParamAttachment.Value(),
			wantValue: true,
		},
		"WithChartEnabled-false": {
			option:    o.WithChartEnabled(false),
			key:       option.ParamChartEnabled.Value(),
			wantValue: false,
		},
		"WithChartEnabled-true": {
			option:    o.WithChartEnabled(true),
			key:       option.ParamChartEnabled.Value(),
			wantValue: true,
		},
		"WithExcludeGroupMembers-false": {
			option:    o.WithExcludeGroupMembers(false),
			key:       option.ParamExcludeGroupMembers.Value(),
			wantValue: false,
		},
		"WithExcludeGroupMembers-true": {
			option:    o.WithExcludeGroupMembers(true),
			key:       option.ParamExcludeGroupMembers.Value(),
			wantValue: true,
		},
		"WithHasDueDate-false": {
			option:    o.WithHasDueDate(false),
			key:       option.ParamHasDueDate.Value(),
			wantValue: false,
		},
		"WithHasDueDate-true": {
			option:    o.WithHasDueDate(true),
			key:       option.ParamHasDueDate.Value(),
			wantValue: true,
		},
		"WithMailNotify-false": {
			option:    o.WithMailNotify(false),
			key:       option.ParamMailNotify.Value(),
			wantValue: false,
		},
		"WithMailNotify-true": {
			option:    o.WithMailNotify(true),
			key:       option.ParamMailNotify.Value(),
			wantValue: true,
		},
		"WithProjectLeaderCanEditProjectLeader-false": {
			option:    o.WithProjectLeaderCanEditProjectLeader(false),
			key:       option.ParamProjectLeaderCanEditProjectLeader.Value(),
			wantValue: false,
		},
		"WithProjectLeaderCanEditProjectLeader-true": {
			option:    o.WithProjectLeaderCanEditProjectLeader(true),
			key:       option.ParamProjectLeaderCanEditProjectLeader.Value(),
			wantValue: true,
		},
		"WithRequired-false": {
			option:    o.WithRequired(false),
			key:       option.ParamRequired.Value(),
			wantValue: false,
		},
		"WithRequired-true": {
			option:    o.WithRequired(true),
			key:       option.ParamRequired.Value(),
			wantValue: true,
		},
		"WithSendMail-false": {
			option:    o.WithSendMail(false),
			key:       option.ParamSendMail.Value(),
			wantValue: false,
		},
		"WithSendMail-true": {
			option:    o.WithSendMail(true),
			key:       option.ParamSendMail.Value(),
			wantValue: true,
		},
		"WithSharedFile-false": {
			option:    o.WithSharedFile(false),
			key:       option.ParamSharedFile.Value(),
			wantValue: false,
		},
		"WithSharedFile-true": {
			option:    o.WithSharedFile(true),
			key:       option.ParamSharedFile.Value(),
			wantValue: true,
		},
		"WithSubtaskingEnabled-false": {
			option:    o.WithSubtaskingEnabled(false),
			key:       option.ParamSubtaskingEnabled.Value(),
			wantValue: false,
		},
		"WithSubtaskingEnabled-true": {
			option:    o.WithSubtaskingEnabled(true),
			key:       option.ParamSubtaskingEnabled.Value(),
			wantValue: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			form := url.Values{}
			require.Nil(t, tc.option.Check())
			_ = tc.option.Set(form)
			assert.Equal(t, strconv.FormatBool(tc.wantValue), form.Get(tc.key))
		})
	}
}
