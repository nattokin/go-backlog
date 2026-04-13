package backlog

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/project"
	"github.com/nattokin/go-backlog/internal/user"
	"github.com/nattokin/go-backlog/internal/wiki"
)

func TestActivityOptionService(t *testing.T) {
	o := activity.NewActivityOptionService(&core.OptionService{})

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
				option:    o.WithOrder(OrderAsc),
				key:       core.ParamOrder.Value(),
				wantValue: string(OrderAsc),
			},
			"with-query-order-desc": {
				option:    o.WithOrder(OrderDesc),
				key:       core.ParamOrder.Value(),
				wantValue: string(OrderDesc),
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
	s := project.NewProjectOptionService(&core.OptionService{})

	// --- Form boolean options -------------------------------------------------------
	t.Run("form-boolean-options", func(t *testing.T) {
		cases := map[string]struct {
			option    core.RequestOption
			key       string
			wantValue bool
		}{
			"with-form-archived": {
				option:    s.WithArchived(true),
				key:       core.ParamArchived.Value(),
				wantValue: true,
			},
			"with-form-chart-enabled": {
				option:    s.WithChartEnabled(true),
				key:       core.ParamChartEnabled.Value(),
				wantValue: true,
			},
			"with-form-subtasking-enabled": {
				option:    s.WithSubtaskingEnabled(false),
				key:       core.ParamSubtaskingEnabled.Value(),
				wantValue: false,
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

	// --- Query boolean options ------------------------------------------------------
	t.Run("query-boolean-options", func(t *testing.T) {
		cases := map[string]struct {
			option    core.RequestOption
			key       string
			wantValue bool
		}{
			"with-query-archived": {
				option:    s.WithArchived(true),
				key:       core.ParamArchived.Value(),
				wantValue: true,
			},
			"with-query-all": {
				option:    s.WithAll(true),
				key:       core.ParamAll.Value(),
				wantValue: true,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				query := url.Values{}
				err := tc.option.Set(query)
				require.NoError(t, err)
				assert.Equal(t, strconv.FormatBool(tc.wantValue), query.Get(tc.key))
			})
		}
	})

	// --- Form string options --------------------------------------------------------
	t.Run("form-string-options", func(t *testing.T) {
		cases := map[string]struct {
			option    core.RequestOption
			key       string
			wantValue string
		}{
			"with-form-name": {
				option:    s.WithName("demo-project"),
				key:       core.ParamName.Value(),
				wantValue: "demo-project",
			},
			"with-form-key": {
				option:    s.WithKey("DEMO"),
				key:       core.ParamKey.Value(),
				wantValue: "DEMO",
			},
			"with-form-text-formatting-rule": {
				option:    s.WithTextFormattingRule("markdown"),
				key:       core.ParamTextFormattingRule.Value(),
				wantValue: "markdown",
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

func TestUserOptionService(t *testing.T) {
	o := user.NewUserOptionService(&core.OptionService{})

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
	s := wiki.NewWikiOptionService(&core.OptionService{})

	// --- Query options ------------------------------------------------------------
	t.Run("query-options", func(t *testing.T) {
		cases := map[string]struct {
			option    core.RequestOption
			key       string
			wantValue string
		}{
			"with-query-keyword": {
				option:    s.WithKeyword("backlog"),
				key:       core.ParamKeyword.Value(),
				wantValue: "backlog",
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

	// --- Form string options ------------------------------------------------------
	t.Run("form-string-options", func(t *testing.T) {
		cases := map[string]struct {
			option    core.RequestOption
			key       string
			wantValue string
		}{
			"with-form-content": {
				option:    s.WithContent("Wiki page content"),
				key:       core.ParamContent.Value(),
				wantValue: "Wiki page content",
			},
			"with-form-name": {
				option:    s.WithName("How to Use Backlog"),
				key:       core.ParamName.Value(),
				wantValue: "How to Use Backlog",
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

	// --- Form boolean options -----------------------------------------------------
	t.Run("form-boolean-options", func(t *testing.T) {
		cases := map[string]struct {
			option    core.RequestOption
			key       string
			wantValue bool
		}{
			"with-form-mail-notify": {
				option:    s.WithMailNotify(true),
				key:       core.ParamMailNotify.Value(),
				wantValue: true,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.Set(form)
				require.NoError(t, err)
				assert.Equal(t, "true", form.Get(tc.key))
			})
		}
	})
}
