package backlog

import (
	"errors"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionService_applyQueryOptions(t *testing.T) {
	cases := map[string]struct {
		opts      []RequestOption
		wantErr   bool
		wantValue string
	}{
		"success": {
			opts: []RequestOption{
				&apiOption{
					t:         queryKeyword,
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
			opts: []RequestOption{
				&apiOption{
					t: queryKeyword,
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
			opts: []RequestOption{
				&apiOption{
					t:         queryKeyword,
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

			s := newOptionService()
			query := url.Values{}

			err := s.applyQueryOptions(query, []queryType{queryKeyword}, tc.opts...)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantValue, query.Get("k"))
		})
	}
}

func TestOptionService_applyFormOptions(t *testing.T) {
	cases := map[string]struct {
		opts      []RequestOption
		wantErr   bool
		wantValue string
	}{
		"success": {
			opts: []RequestOption{
				&apiOption{
					t:         formName,
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
			opts: []RequestOption{
				&apiOption{
					t: formName,
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
			opts: []RequestOption{
				&apiOption{
					t:         formName,
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

			s := newOptionService()
			form := url.Values{}

			err := s.applyFormOptions(form, []formType{formName}, tc.opts...)

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
			option    RequestOption
			key       string
			wantValue int
		}{
			"with-query-min-id": {
				option:    o.WithMinID(5),
				key:       queryMinID.Value(),
				wantValue: 5,
			},
			"with-query-max-id": {
				option:    o.WithMaxID(10),
				key:       queryMaxID.Value(),
				wantValue: 10,
			},
			"with-query-count": {
				option:    o.WithCount(25),
				key:       queryCount.Value(),
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
			option    RequestOption
			key       string
			wantValue string
		}{
			"with-query-order-asc": {
				option:    o.WithOrder(OrderAsc),
				key:       queryOrder.Value(),
				wantValue: string(OrderAsc),
			},
			"with-query-order-desc": {
				option:    o.WithOrder(OrderDesc),
				key:       queryOrder.Value(),
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
			option    RequestOption
			key       string
			wantValue []int
		}{
			"with-query-activity-type-ids": {
				option:    o.WithActivityTypeIDs([]int{1, 2, 3}),
				key:       queryActivityTypeIDs.Value(),
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
	s := newProjectOptionService()

	// --- Form boolean options -------------------------------------------------------
	t.Run("form-boolean-options", func(t *testing.T) {
		cases := map[string]struct {
			option    RequestOption
			key       string
			wantValue bool
		}{
			"with-form-archived": {
				option:    s.WithArchived(true),
				key:       formArchived.Value(),
				wantValue: true,
			},
			"with-form-chart-enabled": {
				option:    s.WithChartEnabled(true),
				key:       formChartEnabled.Value(),
				wantValue: true,
			},
			"with-form-subtasking-enabled": {
				option:    s.WithSubtaskingEnabled(false),
				key:       formSubtaskingEnabled.Value(),
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
			option    RequestOption
			key       string
			wantValue bool
		}{
			"with-query-archived": {
				option:    s.WithArchived(true),
				key:       queryArchived.Value(),
				wantValue: true,
			},
			"with-query-all": {
				option:    s.WithAll(true),
				key:       queryAll.Value(),
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
			option    RequestOption
			key       string
			wantValue string
		}{
			"with-form-name": {
				option:    s.WithName("demo-project"),
				key:       formName.Value(),
				wantValue: "demo-project",
			},
			"with-form-key": {
				option:    s.WithKey("DEMO"),
				key:       formKey.Value(),
				wantValue: "DEMO",
			},
			"with-form-text-formatting-rule": {
				option:    s.WithTextFormattingRule("markdown"),
				key:       formTextFormattingRule.Value(),
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
	o := newUserOptionService()

	// --- Boolean options ------------------------------------------------------------
	t.Run("boolean-options", func(t *testing.T) {
		cases := map[string]struct {
			option    RequestOption
			key       string
			wantValue bool
		}{
			"with-form-send-mail": {
				option:    o.WithSendMail(true),
				key:       formSendMail.Value(),
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
			option    RequestOption
			key       string
			wantValue int
		}{
			"with-form-user-id": {
				option:    o.WithUserID(1),
				key:       formUserID.Value(),
				wantValue: 1,
			},
			"with-form-role-type": {
				option:    o.WithRoleType(2),
				key:       formRoleType.Value(),
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
			option    RequestOption
			key       string
			wantValue string
		}{
			"with-form-name": {
				option:    o.WithName("example-user"),
				key:       formName.Value(),
				wantValue: "example-user",
			},
			"with-form-mail-address": {
				option:    o.WithMailAddress("user@example.com"),
				key:       formMailAddress.Value(),
				wantValue: "user@example.com",
			},
			"with-form-password": {
				option:    o.WithPassword("securepass"),
				key:       formPassword.Value(),
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
	s := newWikiOptionService()

	// --- Query options ------------------------------------------------------------
	t.Run("query-options", func(t *testing.T) {
		cases := map[string]struct {
			option    RequestOption
			key       string
			wantValue string
		}{
			"with-query-keyword": {
				option:    s.WithKeyword("backlog"),
				key:       queryKeyword.Value(),
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
			option    RequestOption
			key       string
			wantValue string
		}{
			"with-form-content": {
				option:    s.WithContent("Wiki page content"),
				key:       formContent.Value(),
				wantValue: "Wiki page content",
			},
			"with-form-name": {
				option:    s.WithName("How to Use Backlog"),
				key:       formName.Value(),
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
			option    RequestOption
			key       string
			wantValue bool
		}{
			"with-form-mail-notify": {
				option:    s.WithMailNotify(true),
				key:       formMailNotify.Value(),
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
