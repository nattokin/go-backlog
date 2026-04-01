package backlog

import (
	"errors"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryOptionService_applyOptions(t *testing.T) {
	cases := map[string]struct {
		opts      []*QueryOption
		wantErr   bool
		wantValue string
	}{
		"success": {
			opts: []*QueryOption{
				{
					checkFunc: func() error { return nil },
					setFunc: func(f url.Values) error {
						f.Set("k", "v")
						return nil
					},
				},
			},
			wantValue: "v",
		},

		"check_error": {
			opts: []*QueryOption{
				{
					checkFunc: func() error { return errors.New("check error") },
					setFunc: func(f url.Values) error {
						f.Set("k", "v")
						return nil
					},
				},
			},
			wantErr: true,
		},

		"set_error": {
			opts: []*QueryOption{
				{
					checkFunc: func() error { return nil },
					setFunc: func(f url.Values) error {
						return errors.New("set error")
					},
				},
			},
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := &QueryOptionService{}
			query := url.Values{}

			err := s.applyOptions(query, tc.opts...)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantValue, query.Get("k"))
		})
	}
}

func TestFormOptionService_applyOptions(t *testing.T) {
	cases := map[string]struct {
		opts      []*FormOption
		wantErr   bool
		wantValue string
	}{
		"success": {
			opts: []*FormOption{
				{
					checkFunc: func() error { return nil },
					setFunc: func(f url.Values) error {
						f.Set("k", "v")
						return nil
					},
				},
			},
			wantValue: "v",
		},

		"check_error": {
			opts: []*FormOption{
				{
					checkFunc: func() error { return errors.New("check error") },
					setFunc: func(f url.Values) error {
						f.Set("k", "v")
						return nil
					},
				},
			},
			wantErr: true,
		},

		"set_error": {
			opts: []*FormOption{
				{
					checkFunc: func() error { return nil },
					setFunc: func(f url.Values) error {
						return errors.New("set error")
					},
				},
			},
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := &FormOptionService{}
			form := url.Values{}

			err := s.applyOptions(form, tc.opts...)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantValue, form.Get("k"))
		})
	}
}

func TestActivityOptionService(t *testing.T) {
	o := newActivityOptionService()

	// --- Integer options ------------------------------------------------------------
	t.Run("integer-options", func(t *testing.T) {
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
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				query := url.Values{}
				err := tc.option.set(query)
				require.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.wantValue), query.Get(tc.key))
			})
		}
	})

	// --- Enum options ---------------------------------------------------------------
	t.Run("enum-options", func(t *testing.T) {
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
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				query := url.Values{}
				err := tc.option.set(query)
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, query.Get(tc.key))
			})
		}
	})

	// --- Special options -------------------------------------------------------------
	t.Run("special-options", func(t *testing.T) {
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
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				query := url.Values{}
				err := tc.option.set(query)
				require.NoError(t, err)

				expected := make([]string, len(tc.wantValue))
				for i, v := range tc.wantValue {
					expected[i] = strconv.Itoa(v)
				}

				// Compare joined values (manual extraction)
				values := (query)[tc.key]
				assert.Equal(t, expected, values)
			})
		}
	})
}

func TestProjectOptionService(t *testing.T) {
	s := newProjectOptionService()

	// --- Form boolean options -------------------------------------------------------
	t.Run("form-boolean-options", func(t *testing.T) {
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
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.set(form)
				require.NoError(t, err)
				assert.Equal(t, strconv.FormatBool(tc.wantValue), form.Get(tc.key))
			})
		}
	})

	// --- Query boolean options ------------------------------------------------------
	t.Run("query-boolean-options", func(t *testing.T) {
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
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				query := url.Values{}
				err := tc.option.set(query)
				require.NoError(t, err)
				assert.Equal(t, strconv.FormatBool(tc.wantValue), query.Get(tc.key))
			})
		}
	})

	// --- Form string options --------------------------------------------------------
	t.Run("form-string-options", func(t *testing.T) {
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
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.set(form)
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, form.Get(tc.key))
			})
		}
	})
}

func TestUserOptionService(t *testing.T) {
	o := newUserOptionService()

	// --- Boolean options ------------------------------------------------------------
	t.Run("boolean-options", func(t *testing.T) {
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
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.set(form)
				require.NoError(t, err)
				assert.Equal(t, strconv.FormatBool(tc.wantValue), form.Get(tc.key))
			})
		}
	})

	// --- Integer options ------------------------------------------------------------
	t.Run("integer-options", func(t *testing.T) {
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
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.set(form)
				require.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.wantValue), form.Get(tc.key))
			})
		}
	})

	// --- String options -------------------------------------------------------------
	t.Run("string-options", func(t *testing.T) {
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
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.set(form)
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, form.Get(tc.key))
			})
		}
	})
}

func TestWikiOptionService(t *testing.T) {
	s := &WikiOptionService{
		registry: &optionRegistry{
			form:  newFormOptionService(),
			query: newQueryOptionService(),
		},
	}

	// --- Query options ------------------------------------------------------------
	t.Run("query-options", func(t *testing.T) {
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
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				query := url.Values{}
				err := tc.option.set(query)
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, query.Get(tc.key))
			})
		}
	})

	// --- Form string options ------------------------------------------------------
	t.Run("form-string-options", func(t *testing.T) {
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
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.set(form)
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, form.Get(tc.key))
			})
		}
	})

	// --- Form boolean options -----------------------------------------------------
	t.Run("form-boolean-options", func(t *testing.T) {
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
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.set(form)
				require.NoError(t, err)
				assert.Equal(t, "true", form.Get(tc.key))
			})
		}
	})
}
