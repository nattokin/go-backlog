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
