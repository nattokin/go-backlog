package backlog_test

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

func TestActivityOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	o := c.User.Activity.Option

	// --- Integer options ------------------------------------------------------------
	t.Run("integer-options", func(t *testing.T) {
		cases := map[string]struct {
			option    core.RequestOption
			key       string
			wantValue int
		}{
			"with-query-min-id": {
				option:    o.WithMinID(5),
				key:       core.ParamMinID.Value(),
				wantValue: 5,
			},
			"with-query-max-id": {
				option:    o.WithMaxID(10),
				key:       core.ParamMaxID.Value(),
				wantValue: 10,
			},
			"with-query-count": {
				option:    o.WithCount(25),
				key:       core.ParamCount.Value(),
				wantValue: 25,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				query := url.Values{}
				err := tc.option.Set(query)
				require.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.wantValue), query.Get(tc.key))
			})
		}
	})

	// --- Enum options ---------------------------------------------------------------
	t.Run("enum-options", func(t *testing.T) {
		cases := map[string]struct {
			option    core.RequestOption
			key       string
			wantValue string
		}{
			"with-query-order-asc": {
				option:    o.WithOrder(model.OrderAsc),
				key:       core.ParamOrder.Value(),
				wantValue: string(model.OrderAsc),
			},
			"with-query-order-desc": {
				option:    o.WithOrder(model.OrderDesc),
				key:       core.ParamOrder.Value(),
				wantValue: string(model.OrderDesc),
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				query := url.Values{}
				err := tc.option.Set(query)
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, query.Get(tc.key))
			})
		}
	})

	// --- Special options -------------------------------------------------------------
	t.Run("special-options", func(t *testing.T) {
		cases := map[string]struct {
			option    core.RequestOption
			key       string
			wantValue []int
		}{
			"with-query-activity-type-ids": {
				option:    o.WithActivityTypeIDs([]int{1, 2, 3}),
				key:       core.ParamActivityTypeIDs.Value(),
				wantValue: []int{1, 2, 3},
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				query := url.Values{}
				err := tc.option.Set(query)
				require.NoError(t, err)

				expected := make([]string, len(tc.wantValue))
				for i, v := range tc.wantValue {
					expected[i] = strconv.Itoa(v)
				}

				values := (query)[tc.key]
				assert.Equal(t, expected, values)
			})
		}
	})
}

func TestProjectOptionService(t *testing.T) {
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

func TestUserOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	o := c.User.Option

	// --- Boolean options ------------------------------------------------------------
	t.Run("boolean-options", func(t *testing.T) {
		cases := map[string]struct {
			option    core.RequestOption
			key       string
			wantValue bool
		}{
			"with-form-send-mail": {
				option:    o.WithSendMail(true),
				key:       core.ParamSendMail.Value(),
				wantValue: true,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.Set(form)
				require.NoError(t, err)
				assert.Equal(t, strconv.FormatBool(tc.wantValue), form.Get(tc.key))
			})
		}
	})

	// --- Integer options ------------------------------------------------------------
	t.Run("integer-options", func(t *testing.T) {
		cases := map[string]struct {
			option    core.RequestOption
			key       string
			wantValue int
		}{
			"with-form-user-id": {
				option:    o.WithUserID(1),
				key:       core.ParamUserID.Value(),
				wantValue: 1,
			},
			"with-form-role-type": {
				option:    o.WithRoleType(2),
				key:       core.ParamRoleType.Value(),
				wantValue: 2,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.Set(form)
				require.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.wantValue), form.Get(tc.key))
			})
		}
	})

	// --- String options -------------------------------------------------------------
	t.Run("string-options", func(t *testing.T) {
		cases := map[string]struct {
			option    core.RequestOption
			key       string
			wantValue string
		}{
			"with-form-name": {
				option:    o.WithName("example-user"),
				key:       core.ParamName.Value(),
				wantValue: "example-user",
			},
			"with-form-mail-address": {
				option:    o.WithMailAddress("user@example.com"),
				key:       core.ParamMailAddress.Value(),
				wantValue: "user@example.com",
			},
			"with-form-password": {
				option:    o.WithPassword("securepass"),
				key:       core.ParamPassword.Value(),
				wantValue: "securepass",
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.Set(form)
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, form.Get(tc.key))
			})
		}
	})
}

func TestWikiOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	s := c.Wiki.Option

	cases := map[string]struct {
		option  core.RequestOption
		wantKey string
	}{
		"WithKeyword": {
			option:  s.WithKeyword("backlog"),
			wantKey: core.ParamKeyword.Value(),
		},
		"WithContent": {
			option:  s.WithContent("Wiki page content"),
			wantKey: core.ParamContent.Value(),
		},
		"WithName": {
			option:  s.WithName("How to Use Backlog"),
			wantKey: core.ParamName.Value(),
		},
		"WithMailNotify": {
			option:  s.WithMailNotify(true),
			wantKey: core.ParamMailNotify.Value(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}
