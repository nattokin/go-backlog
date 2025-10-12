package backlog

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test verifies that OptionService.applyOptions correctly handles
// normal and error flows for both FormOption and QueryOption.
// It covers:
//   - Successful option application
//   - Check() validation errors
//   - set() execution errors
//   - Type consistency (FormOption / QueryOption)
func TestOptionService_applyOptions(t *testing.T) {
	queryOption := newQueryOptionService()
	formOption := newFormOptionService()

	// Dummy error types for type-based validation
	var (
		errCheckFailed = errors.New("check failed")
		errSetFailed   = errors.New("set failed")
	)

	cases := map[string]struct {
		isQuery     bool
		opts        []RequestOption
		expectErr   bool
		expectErrIs error
		wantValues  map[string]string
	}{
		// --- Successful option application ----------------------------------------

		"form-applies-valid-options": {
			isQuery: false,
			opts: []RequestOption{
				formOption.WithName("test"),
				formOption.WithMailAddress("mail@test.com"),
			},
			wantValues: map[string]string{
				"name":        "test",
				"mailAddress": "mail@test.com",
			},
		},

		"query-applies-valid-options": {
			isQuery: true,
			opts: []RequestOption{
				queryOption.WithCount(10),
				queryOption.WithAll(true),
			},
			wantValues: map[string]string{
				"count": "10",
				"all":   "true",
			},
		},

		// --- Validation errors (Check fails) -------------------------------------

		"form-check-fails": {
			isQuery: false,
			opts: []RequestOption{
				&FormOption{
					t: formKey,
					checkFunc: func() error {
						return errCheckFailed
					},
					setFunc: func(form *FormParams) error {
						form.Set("x", "should-not-be-set")
						return nil
					},
				},
			},
			expectErr:   true,
			expectErrIs: errCheckFailed,
		},

		"query-check-fails": {
			isQuery: true,
			opts: []RequestOption{
				&QueryOption{
					t: queryKey,
					checkFunc: func() error {
						return errCheckFailed
					},
					setFunc: func(query *QueryParams) error {
						query.Set("y", "should-not-be-set")
						return nil
					},
				},
			},
			expectErr:   true,
			expectErrIs: errCheckFailed,
		},

		// --- Runtime set() errors -------------------------------------------------

		"form-set-fails": {
			isQuery: false,
			opts: []RequestOption{
				&FormOption{
					t:         formKey,
					checkFunc: func() error { return nil },
					setFunc: func(form *FormParams) error {
						return errSetFailed
					},
				},
			},
			expectErr:   true,
			expectErrIs: errSetFailed,
		},

		"query-set-fails": {
			isQuery: true,
			opts: []RequestOption{
				&QueryOption{
					t:         queryKey,
					checkFunc: func() error { return nil },
					setFunc: func(query *QueryParams) error {
						return errSetFailed
					},
				},
			},
			expectErr:   true,
			expectErrIs: errSetFailed,
		},

		// --- Mixed cases: multiple options, one fails -----------------------------

		"form-mixed-check-error": {
			isQuery: false,
			opts: []RequestOption{
				formOption.WithName("ok"),
				&FormOption{
					t:         formKey,
					checkFunc: func() error { return errCheckFailed },
					setFunc:   func(form *FormParams) error { return nil },
				},
			},
			expectErr:   true,
			expectErrIs: errCheckFailed,
		},

		"query-mixed-set-error": {
			isQuery: true,
			opts: []RequestOption{
				queryOption.WithCount(10),
				&QueryOption{
					t:         queryKey,
					checkFunc: func() error { return nil },
					setFunc: func(query *QueryParams) error {
						return errSetFailed
					},
				},
			},
			expectErr:   true,
			expectErrIs: errSetFailed,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tc.isQuery {
				query := NewQueryParams()
				err := newQueryOptionService().applyOptions(query, toQueryOptions(t, tc.opts)...)

				if tc.expectErr {
					require.Error(t, err)
					require.ErrorIs(t, err, tc.expectErrIs)
					return
				}

				require.NoError(t, err)
				for k, v := range tc.wantValues {
					assert.Equal(t, v, query.Get(k))
				}

			} else {
				form := NewFormParams()
				err := newFormOptionService().applyOptions(form, toFormOptions(t, tc.opts)...)

				if tc.expectErr {
					require.Error(t, err)
					require.ErrorIs(t, err, tc.expectErrIs)
					return
				}

				require.NoError(t, err)
				for k, v := range tc.wantValues {
					assert.Equal(t, v, form.Get(k))
				}
			}
		})
	}
}

// This test verifies that each WithQueryXxx method in ActivityOptionService
// correctly builds and applies the expected query parameters.
// Since these methods only delegate to QueryOptionService,
// one success case per method is sufficient.
func TestActivityOptionService(t *testing.T) {
	o := newActivityOptionService()

	// --- Integer options ------------------------------------------------------------
	t.Run("Integer options", func(t *testing.T) {
		cases := map[string]struct {
			option    *QueryOption
			key       string
			wantValue int
		}{
			"with-query-min-id": {
				option:    o.WithQueryMinID(5),
				key:       queryMinID.Value(),
				wantValue: 5,
			},
			"with-query-max-id": {
				option:    o.WithQueryMaxID(10),
				key:       queryMaxID.Value(),
				wantValue: 10,
			},
			"with-query-count": {
				option:    o.WithQueryCount(25),
				key:       queryCount.Value(),
				wantValue: 25,
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				query := NewQueryParams()
				err := tc.option.set(query)
				require.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.wantValue), query.Get(tc.key))
			})
		}
	})

	// --- Enum options ---------------------------------------------------------------
	t.Run("Enum options", func(t *testing.T) {
		cases := map[string]struct {
			option    *QueryOption
			key       string
			wantValue string
		}{
			"with-query-order-asc": {
				option:    o.WithQueryOrder(OrderAsc),
				key:       queryOrder.Value(),
				wantValue: string(OrderAsc),
			},
			"with-query-order-desc": {
				option:    o.WithQueryOrder(OrderDesc),
				key:       queryOrder.Value(),
				wantValue: string(OrderDesc),
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				query := NewQueryParams()
				err := tc.option.set(query)
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, query.Get(tc.key))
			})
		}
	})

	// --- Special options -------------------------------------------------------------
	t.Run("Special options", func(t *testing.T) {
		cases := map[string]struct {
			option    *QueryOption
			key       string
			wantValue []int
		}{
			"with-query-activity-type-ids": {
				option:    o.WithQueryActivityTypeIDs([]int{1, 2, 3}),
				key:       queryActivityTypeIDs.Value(),
				wantValue: []int{1, 2, 3},
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				query := NewQueryParams()
				err := tc.option.set(query)
				require.NoError(t, err)

				expected := make([]string, len(tc.wantValue))
				for i, v := range tc.wantValue {
					expected[i] = strconv.Itoa(v)
				}

				// Compare joined values (manual extraction)
				values := (*query.Values)[tc.key]
				assert.Equal(t, expected, values)
			})
		}
	})
}

// This test verifies that each WithXxx method in ProjectOptionService
// correctly builds and applies the expected form and query parameters.
// Each option is tested for one success case.
//

func TestProjectOptionService(t *testing.T) {
	s := newProjectOptionService()

	// --- Form boolean options -------------------------------------------------------
	t.Run("Form boolean options", func(t *testing.T) {
		cases := map[string]struct {
			option    *FormOption
			key       string
			wantValue bool
		}{
			"with-form-archived": {
				option:    s.WithFormArchived(true),
				key:       formArchived.Value(),
				wantValue: true,
			},
			"with-form-chart-enabled": {
				option:    s.WithFormChartEnabled(true),
				key:       formChartEnabled.Value(),
				wantValue: true,
			},
			"with-form-subtasking-enabled": {
				option:    s.WithFormSubtaskingEnabled(false),
				key:       formSubtaskingEnabled.Value(),
				wantValue: false,
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := NewFormParams()
				err := tc.option.set(form)
				require.NoError(t, err)
				assert.Equal(t, strconv.FormatBool(tc.wantValue), form.Get(tc.key))
			})
		}
	})

	// --- Query boolean options ------------------------------------------------------
	t.Run("Query boolean options", func(t *testing.T) {
		cases := map[string]struct {
			option    *QueryOption
			key       string
			wantValue bool
		}{
			"with-query-archived": {
				option:    s.WithQueryArchived(true),
				key:       queryArchived.Value(),
				wantValue: true,
			},
			"with-query-all": {
				option:    s.WithQueryAll(true),
				key:       queryAll.Value(),
				wantValue: true,
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				query := NewQueryParams()
				err := tc.option.set(query)
				require.NoError(t, err)
				assert.Equal(t, strconv.FormatBool(tc.wantValue), query.Get(tc.key))
			})
		}
	})

	// --- Form string options --------------------------------------------------------
	t.Run("Form string options", func(t *testing.T) {
		cases := map[string]struct {
			option    *FormOption
			key       string
			wantValue string
		}{
			"with-form-name": {
				option:    s.WithFormName("demo-project"),
				key:       formName.Value(),
				wantValue: "demo-project",
			},
			"with-form-key": {
				option:    s.WithFormKey("DEMO"),
				key:       formKey.Value(),
				wantValue: "DEMO",
			},
			"with-form-text-formatting-rule": {
				option:    s.WithFormTextFormattingRule("markdown"),
				key:       formTextFormattingRule.Value(),
				wantValue: "markdown",
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := NewFormParams()
				err := tc.option.set(form)
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, form.Get(tc.key))
			})
		}
	})
}

// This test verifies that each WithXxx method in UserOptionService
// correctly sets its associated FormOption key and value.
// Since these methods only wrap internal core functions,
// one success case per option is sufficient.
func TestUserOptionService(t *testing.T) {
	o := newUserOptionService()

	// --- Boolean options ------------------------------------------------------------
	t.Run("Boolean options", func(t *testing.T) {
		cases := map[string]struct {
			option    *FormOption
			key       string
			wantValue bool
		}{
			"with-form-send-mail": {
				option:    o.WithFormSendMail(true),
				key:       formSendMail.Value(),
				wantValue: true,
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := NewFormParams()
				err := tc.option.set(form)
				require.NoError(t, err)
				assert.Equal(t, strconv.FormatBool(tc.wantValue), form.Get(tc.key))
			})
		}
	})

	// --- Integer options ------------------------------------------------------------
	t.Run("Integer options", func(t *testing.T) {
		cases := map[string]struct {
			option    *FormOption
			key       string
			wantValue int
		}{
			"with-form-user-id": {
				option:    o.WithFormUserID(1),
				key:       formUserID.Value(),
				wantValue: 1,
			},
			"with-form-role-type": {
				option:    o.WithFormRoleType(2),
				key:       formRoleType.Value(),
				wantValue: 2,
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := NewFormParams()
				err := tc.option.set(form)
				require.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.wantValue), form.Get(tc.key))
			})
		}
	})

	// --- String options -------------------------------------------------------------
	t.Run("String options", func(t *testing.T) {
		cases := map[string]struct {
			option    *FormOption
			key       string
			wantValue string
		}{
			"with-form-name": {
				option:    o.WithFormName("example-user"),
				key:       formName.Value(),
				wantValue: "example-user",
			},
			"with-form-mail-address": {
				option:    o.WithFormMailAddress("user@example.com"),
				key:       formMailAddress.Value(),
				wantValue: "user@example.com",
			},
			"with-form-password": {
				option:    o.WithFormPassword("securepass"),
				key:       formPassword.Value(),
				wantValue: "securepass",
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := NewFormParams()
				err := tc.option.set(form)
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, form.Get(tc.key))
			})
		}
	})
}

// This test verifies that each WithXxx method in WikiOptionService
// correctly builds and applies the expected form and query parameters.
// Since these methods delegate to FormOptionService or QueryOptionService,
// one success case per method is sufficient.
func TestWikiOptionService(t *testing.T) {
	s := &WikiOptionService{
		support: &optionSupport{
			form:  newFormOptionService(),
			query: newQueryOptionService(),
		},
	}

	// --- Query options ------------------------------------------------------------
	t.Run("Query options", func(t *testing.T) {
		cases := map[string]struct {
			option    *QueryOption
			key       string
			wantValue string
		}{
			"with-query-keyword": {
				option:    s.WithQueryKeyword("backlog"),
				key:       queryKeyword.Value(),
				wantValue: "backlog",
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				query := NewQueryParams()
				err := tc.option.set(query)
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, query.Get(tc.key))
			})
		}
	})

	// --- Form string options ------------------------------------------------------
	t.Run("Form string options", func(t *testing.T) {
		cases := map[string]struct {
			option    *FormOption
			key       string
			wantValue string
		}{
			"with-form-content": {
				option:    s.WithFormContent("Wiki page content"),
				key:       formContent.Value(),
				wantValue: "Wiki page content",
			},
			"with-form-name": {
				option:    s.WithFormName("How to Use Backlog"),
				key:       formName.Value(),
				wantValue: "How to Use Backlog",
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := NewFormParams()
				err := tc.option.set(form)
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, form.Get(tc.key))
			})
		}
	})

	// --- Form boolean options -----------------------------------------------------
	t.Run("Form boolean options", func(t *testing.T) {
		cases := map[string]struct {
			option    *FormOption
			key       string
			wantValue bool
		}{
			"with-form-mail-notify": {
				option:    s.WithFormMailNotify(true),
				key:       formMailNotify.Value(),
				wantValue: true,
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := NewFormParams()
				err := tc.option.set(form)
				require.NoError(t, err)
				assert.Equal(t, "true", form.Get(tc.key))
			})
		}
	})
}
