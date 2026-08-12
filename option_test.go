package backlog_test

import (
	"context"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestActivityOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	o := c.User.Activity.Option

	// --- Integer options ------------------------------------------------------------
	t.Run("integer-options", func(t *testing.T) {
		cases := map[string]struct {
			option    backlog.RequestOption
			key       string
			wantValue int
		}{
			"with-query-min-id": {
				option:    o.WithMinID(5),
				key:       option.ParamMinID.Value(),
				wantValue: 5,
			},
			"with-query-max-id": {
				option:    o.WithMaxID(10),
				key:       option.ParamMaxID.Value(),
				wantValue: 10,
			},
			"with-query-count": {
				option:    o.WithCount(25),
				key:       option.ParamCount.Value(),
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
			option    backlog.RequestOption
			key       string
			wantValue string
		}{
			"with-query-order-asc": {
				option:    o.WithOrder(backlog.OrderAsc),
				key:       option.ParamOrder.Value(),
				wantValue: string(backlog.OrderAsc),
			},
			"with-query-order-desc": {
				option:    o.WithOrder(backlog.OrderDesc),
				key:       option.ParamOrder.Value(),
				wantValue: string(backlog.OrderDesc),
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
			option    backlog.RequestOption
			key       string
			wantValue []int
		}{
			"with-query-activity-type-ids": {
				option:    o.WithActivityTypeIDs([]int{1, 2, 3}),
				key:       option.ParamActivityTypeIDs.Value(),
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

func TestRequestOption_Check(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	o := c.User.Activity.Option

	t.Run("valid-returns-nil", func(t *testing.T) {
		t.Parallel()
		opt := o.WithCount(20)
		assert.Nil(t, opt.Check())
	})

	t.Run("invalid-returns-ValidationError", func(t *testing.T) {
		t.Parallel()
		// count=0 is invalid; Check() should propagate the ValidationError
		opt := o.WithCount(0)
		ve := opt.Check()
		require.NotNil(t, ve)
		assert.NotEmpty(t, ve.Target())
		assert.NotEmpty(t, ve.Message())
	})
}

// failingCheckOption is a custom RequestOption whose Check always returns a ValidationError.
// It is used to verify that toCoreOption correctly wraps option.Check() in CheckFunc,
// propagating the error through the internal pipeline back to the caller as a *ValidationError.
type failingCheckOption struct {
	key string
}

func (o *failingCheckOption) Key() string { return o.key }
func (o *failingCheckOption) Check() *backlog.ValidationError {
	return backlog.NewValidationError(o.key, "always fails")
}
func (o *failingCheckOption) Set(url.Values) error { return nil }

func Test_toCoreOption_CheckFunc(t *testing.T) {
	t.Parallel()

	c, err := backlog.NewClient(
		"https://example.backlog.com", "token",
		backlog.WithDoer(&mock.Doer{DoFunc: mock.NewNotFoundDoFunc()}),
	)
	require.NoError(t, err)

	opt := &failingCheckOption{key: option.ParamCount.Value()}
	_, err = c.Issue.Comment.List(context.Background(), "PRJ-1", opt)

	require.Error(t, err)
	var ve *backlog.ValidationError
	assert.ErrorAs(t, err, &ve)
	assert.Equal(t, option.ParamCount.Value(), ve.Target())
	assert.Equal(t, "always fails", ve.Message())
}
