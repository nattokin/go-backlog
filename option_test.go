package backlog

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ──────────────────────────────────────────────────────────────
//  OptionService.applyOptions internal tests
// ──────────────────────────────────────────────────────────────
//
// This test verifies that OptionService.applyOptions correctly handles
// normal and error flows for both FormOption and QueryOption.
// It covers:
//   - Successful option application
//   - Check() validation errors
//   - set() execution errors
//   - Type consistency (FormOption / QueryOption)
//

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

func TestActivityOptionService(t *testing.T) {
	// TODO
}

func TestProjectOptionService(t *testing.T) {
	// TODO
}

func TestUserOptionService(t *testing.T) {
	o := newUserOptionService()

	t.Run("Boolean options", func(t *testing.T) {
		cases := map[string]struct {
			option    *FormOption
			key       string
			wantValue bool
		}{
			"WithFormSendMail": {
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

	t.Run("Integer options", func(t *testing.T) {
		cases := map[string]struct {
			option    *FormOption
			key       string
			wantValue int
		}{
			"WithFormUserID": {
				option:    o.WithFormUserID(1),
				key:       formUserID.Value(),
				wantValue: 1,
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
}

func TestWikiOptionService(t *testing.T) {
	// TODO
}

// --- Option Service Helpers ---

// newQueryOptionService returns a test instance of QueryOptionService.
func newQueryOptionService() *QueryOptionService {
	return &QueryOptionService{}
}

// newFormOptionService returns a test instance of FormOptionService.
func newFormOptionService() *FormOptionService {
	return &FormOptionService{}
}

// newActivityOptionService returns a test instance of ActivityOptionService.
func newActivityOptionService() *ActivityOptionService {
	return &ActivityOptionService{
		support: &optionSupport{
			query: newQueryOptionService(),
			form:  newFormOptionService(),
		},
	}
}

// newProjectOptionService returns a test instance of ProjectOptionService.
func newProjectOptionService() *ProjectOptionService {
	return &ProjectOptionService{
		support: &optionSupport{
			query: newQueryOptionService(),
			form:  newFormOptionService(),
		},
	}
}

// newUserOptionService returns a test instance of UserOptionService.
func newUserOptionService() *UserOptionService {
	return &UserOptionService{
		support: &optionSupport{
			query: newQueryOptionService(),
			form:  newFormOptionService(),
		},
	}
}

// newWikiOptionService returns a test instance of WikiOptionService.
func newWikiOptionService() *WikiOptionService {
	return &WikiOptionService{
		support: &optionSupport{
			query: newQueryOptionService(),
			form:  newFormOptionService(),
		},
	}
}
