package core_test

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
)

func TestOptionService_bool(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		option    core.RequestOption
		key       string
		wantValue bool
	}{
		"WithAll-true": {
			option:    o.WithAll(true),
			key:       core.ParamAll.Value(),
			wantValue: true,
		},
		"WithAll-false": {
			option:    o.WithAll(false),
			key:       core.ParamAll.Value(),
			wantValue: false,
		},
		"WithArchived-true": {
			option:    o.WithArchived(true),
			key:       core.ParamArchived.Value(),
			wantValue: true,
		},
		"WithArchived-false": {
			option:    o.WithArchived(false),
			key:       core.ParamArchived.Value(),
			wantValue: false,
		},
		"WithChartEnabled-true": {
			option:    o.WithChartEnabled(true),
			key:       core.ParamChartEnabled.Value(),
			wantValue: true,
		},
		"WithChartEnabled-false": {
			option:    o.WithChartEnabled(false),
			key:       core.ParamChartEnabled.Value(),
			wantValue: false,
		},
		"WithMailNotify-true": {
			option:    o.WithMailNotify(true),
			key:       core.ParamMailNotify.Value(),
			wantValue: true,
		},
		"WithMailNotify-false": {
			option:    o.WithMailNotify(false),
			key:       core.ParamMailNotify.Value(),
			wantValue: false,
		},
		"WithProjectLeaderCanEditProjectLeader-true": {
			option:    o.WithProjectLeaderCanEditProjectLeader(true),
			key:       core.ParamProjectLeaderCanEditProjectLeader.Value(),
			wantValue: true,
		},
		"WithProjectLeaderCanEditProjectLeader-false": {
			option:    o.WithProjectLeaderCanEditProjectLeader(false),
			key:       core.ParamProjectLeaderCanEditProjectLeader.Value(),
			wantValue: false,
		},
		"WithSendMail-true": {
			option:    o.WithSendMail(true),
			key:       core.ParamSendMail.Value(),
			wantValue: true,
		},
		"WithSendMail-false": {
			option:    o.WithSendMail(false),
			key:       core.ParamSendMail.Value(),
			wantValue: false,
		},
		"WithSubtaskingEnabled-true": {
			option:    o.WithSubtaskingEnabled(true),
			key:       core.ParamSubtaskingEnabled.Value(),
			wantValue: true,
		},
		"WithSubtaskingEnabled-false": {
			option:    o.WithSubtaskingEnabled(false),
			key:       core.ParamSubtaskingEnabled.Value(),
			wantValue: false,
		},
		"WithAttachment-true": {
			option:    o.WithAttachment(true),
			key:       core.ParamAttachment.Value(),
			wantValue: true,
		},
		"WithAttachment-false": {
			option:    o.WithAttachment(false),
			key:       core.ParamAttachment.Value(),
			wantValue: false,
		},
		"WithSharedFile-true": {
			option:    o.WithSharedFile(true),
			key:       core.ParamSharedFile.Value(),
			wantValue: true,
		},
		"WithSharedFile-false": {
			option:    o.WithSharedFile(false),
			key:       core.ParamSharedFile.Value(),
			wantValue: false,
		},
		"WithHasDueDate-true": {
			option:    o.WithHasDueDate(true),
			key:       core.ParamHasDueDate.Value(),
			wantValue: true,
		},
		"WithHasDueDate-false": {
			option:    o.WithHasDueDate(false),
			key:       core.ParamHasDueDate.Value(),
			wantValue: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			form := url.Values{}
			err := tc.option.Check()
			require.NoError(t, err)
			_ = tc.option.Set(form)
			assert.Equal(t, strconv.FormatBool(tc.wantValue), form.Get(tc.key))
		})
	}

}
