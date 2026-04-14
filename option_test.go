package backlog_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

func TestProjectOptionService_Keys(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	s := c.Project.Option

	cases := map[string]struct {
		option  core.RequestOption
		wantKey string
	}{
		"WithAll": {
			option:  s.WithAll(true),
			wantKey: core.ParamAll.Value(),
		},
		"WithArchived": {
			option:  s.WithArchived(true),
			wantKey: core.ParamArchived.Value(),
		},
		"WithChartEnabled": {
			option:  s.WithChartEnabled(true),
			wantKey: core.ParamChartEnabled.Value(),
		},
		"WithKey": {
			option:  s.WithKey("TEST"),
			wantKey: core.ParamKey.Value(),
		},
		"WithName": {
			option:  s.WithName("test"),
			wantKey: core.ParamName.Value(),
		},
		"WithProjectLeaderCanEditProjectLeader": {
			option:  s.WithProjectLeaderCanEditProjectLeader(true),
			wantKey: core.ParamProjectLeaderCanEditProjectLeader.Value(),
		},
		"WithSubtaskingEnabled": {
			option:  s.WithSubtaskingEnabled(true),
			wantKey: core.ParamSubtaskingEnabled.Value(),
		},
		"WithTextFormattingRule": {
			option:  s.WithTextFormattingRule(model.FormatBacklog),
			wantKey: core.ParamTextFormattingRule.Value(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}
