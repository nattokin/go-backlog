package backlog

import (
	"errors"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
// ──────────────────────────────────────────────────────────────
//  TestApiOption (formerly TestFormOption / TestQueryOption)
// ──────────────────────────────────────────────────────────────
//

func TestApiOption(t *testing.T) {
	cases := map[string]struct {
		option         RequestOption
		expectCheckErr bool
		expectSetErr   bool
		wantValue      string
		wantCheckErr   error
		wantSetErr     error
	}{
		"Success": {
			option: &apiOption{
				t:         formKey,
				checkFunc: func() error { return nil },
				setFunc: func(form url.Values) error {
					form.Set(formKey.Value(), "success")
					return nil
				},
			},
			wantValue: "success",
		},
		"Check-error": {
			option:         newFormOptionWithCheckError(formKey),
			expectCheckErr: true,
			wantCheckErr:   errors.New("check error"),
		},
		"set-error": {
			option:       newFormOptionWithSetError(formKey),
			expectSetErr: true,
			wantSetErr:   errors.New("set error"),
		},
		"checkFunc-nil": {
			option: &apiOption{
				t:         formKey,
				checkFunc: nil,
				setFunc: func(form url.Values) error {
					form.Set(formKey.Value(), "checkFunc nil")
					return nil
				},
			},
			wantValue: "checkFunc nil",
		},
		"set-nil": {
			option: &apiOption{
				t:         formKey,
				checkFunc: func() error { return nil },
				setFunc:   nil,
			},
			expectSetErr: true,
			wantSetErr:   newValidationError("option has no setter"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			form := url.Values{}
			err := tc.option.Check()
			if tc.expectCheckErr {
				require.Error(t, err)
				assert.IsType(t, tc.wantCheckErr, err)
				return
			}
			require.NoError(t, err)

			if err := tc.option.Set(form); tc.expectSetErr {
				require.Error(t, err)
				assert.IsType(t, tc.wantSetErr, err)
			} else {
				require.NoError(t, err)
				ao := tc.option.(*apiOption)
				assert.Equal(t, tc.wantValue, form.Get(ao.t.(formType).Value()))
			}
		})
	}
}

// ──────────────────────────────────────────────────────────────
//	TestOptionService (form options)
// ──────────────────────────────────────────────────────────────

func TestOptionService_FormOptions(t *testing.T) {
	o := newOptionService()

	// --- Boolean options ------------------------------------------------------------
	t.Run("boolean-options", func(t *testing.T) {
		cases := map[string]struct {
			option    RequestOption
			key       string
			wantValue bool
		}{
			"WithArchived-true": {
				option:    o.WithArchived(true),
				key:       formArchived.Value(),
				wantValue: true,
			},
			"WithArchived-false": {
				option:    o.WithArchived(false),
				key:       formArchived.Value(),
				wantValue: false,
			},
			"WithChartEnabled-true": {
				option:    o.WithChartEnabled(true),
				key:       formChartEnabled.Value(),
				wantValue: true,
			},
			"WithChartEnabled-false": {
				option:    o.WithChartEnabled(false),
				key:       formChartEnabled.Value(),
				wantValue: false,
			},
			"WithMailNotify-true": {
				option:    o.WithMailNotify(true),
				key:       formMailNotify.Value(),
				wantValue: true,
			},
			"WithMailNotify-false": {
				option:    o.WithMailNotify(false),
				key:       formMailNotify.Value(),
				wantValue: false,
			},
			"WithProjectLeaderCanEditProjectLeader-true": {
				option:    o.WithProjectLeaderCanEditProjectLeader(true),
				key:       formProjectLeaderCanEditProjectLeader.Value(),
				wantValue: true,
			},
			"WithProjectLeaderCanEditProjectLeader-false": {
				option:    o.WithProjectLeaderCanEditProjectLeader(false),
				key:       formProjectLeaderCanEditProjectLeader.Value(),
				wantValue: false,
			},
			"WithSendMail-true": {
				option:    o.WithSendMail(true),
				key:       formSendMail.Value(),
				wantValue: true,
			},
			"WithSendMail-false": {
				option:    o.WithSendMail(false),
				key:       formSendMail.Value(),
				wantValue: false,
			},
			"WithSubtaskingEnabled-true": {
				option:    o.WithSubtaskingEnabled(true),
				key:       formSubtaskingEnabled.Value(),
				wantValue: true,
			},
			"WithSubtaskingEnabled-false": {
				option:    o.WithSubtaskingEnabled(false),
				key:       formSubtaskingEnabled.Value(),
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
	})

	// --- Integer options ------------------------------------------------------------
	t.Run("integer-options", func(t *testing.T) {
		cases := map[string]struct {
			option    RequestOption
			key       string
			wantValue int
			wantErr   bool
		}{
			"WithUserID-valid-1": {
				option:    o.WithUserID(1),
				key:       formUserID.Value(),
				wantValue: 1,
			},
			"WithUserID-valid-2": {
				option:    o.WithUserID(2),
				key:       formUserID.Value(),
				wantValue: 2,
			},
			"WithUserID-invalid-0": {
				option:  o.WithUserID(0),
				key:     formUserID.Value(),
				wantErr: true,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.Check()
				if tc.wantErr {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				_ = tc.option.Set(form)
				assert.Equal(t, strconv.Itoa(tc.wantValue), form.Get(tc.key))
			})
		}
	})

	// --- String options ------------------------------------------------------------
	t.Run("string-options", func(t *testing.T) {
		cases := map[string]struct {
			option    RequestOption
			key       string
			wantValue string
			wantErr   bool
		}{
			"WithContent-valid": {
				option:    o.WithContent("Hello"),
				key:       formContent.Value(),
				wantValue: "Hello",
			},
			"WithContent-empty": {
				option:  o.WithContent(""),
				key:     formContent.Value(),
				wantErr: true,
			},
			"WithKey-valid": {
				option:    o.WithKey("ABC"),
				key:       formKey.Value(),
				wantValue: "ABC",
			},
			"WithKey-empty": {
				option:  o.WithKey(""),
				key:     formKey.Value(),
				wantErr: true,
			},
			"WithMailAddress-valid": {
				option:    o.WithMailAddress("test@example.com"),
				key:       formMailAddress.Value(),
				wantValue: "test@example.com",
			},
			"WithMailAddress-empty": {
				option:  o.WithMailAddress(""),
				key:     formMailAddress.Value(),
				wantErr: true,
			},
			"WithName-valid": {
				option:    o.WithName("testname"),
				key:       formName.Value(),
				wantValue: "testname",
			},
			"WithName-empty": {
				option:  o.WithName(""),
				key:     formName.Value(),
				wantErr: true,
			},
			"WithPassword-valid-8chars": {
				option:    o.WithPassword("abcdefgh"),
				key:       formPassword.Value(),
				wantValue: "abcdefgh",
			},
			"WithPassword-valid-9chars": {
				option:    o.WithPassword("abcdefghi"),
				key:       formPassword.Value(),
				wantValue: "abcdefghi",
			},
			"WithPassword-valid-7chars": {
				option:  o.WithPassword("abcdefg"),
				key:     formPassword.Value(),
				wantErr: true,
			},
			"WithPassword-invalid-empty": {
				option:  o.WithPassword(""),
				key:     formPassword.Value(),
				wantErr: true,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.Check()
				if tc.wantErr {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				_ = tc.option.Set(form)
				assert.Equal(t, tc.wantValue, form.Get(tc.key))
			})
		}
	})

	// --- Enum or special options ------------------------------------------------------------
	t.Run("enum-or-special-options", func(t *testing.T) {
		cases := map[string]struct {
			option    RequestOption
			key       string
			wantValue string
			wantErr   bool
		}{
			"WithRoleType-valid-1": {
				option:    o.WithRoleType(1),
				key:       formRoleType.Value(),
				wantValue: "1",
			},
			"WithRoleType-valid-6": {
				option:    o.WithRoleType(6),
				key:       formRoleType.Value(),
				wantValue: "6",
			},
			"WithRoleType-invalid-0": {
				option:  o.WithRoleType(0),
				key:     formRoleType.Value(),
				wantErr: true,
			},
			"WithRoleType-invalid-7": {
				option:  o.WithRoleType(7),
				key:     formRoleType.Value(),
				wantErr: true,
			},
			"WithTextFormattingRule-valid-backlog": {
				option:    o.WithTextFormattingRule(FormatBacklog),
				key:       formTextFormattingRule.Value(),
				wantValue: string(FormatBacklog),
			},
			"WithTextFormattingRule-valid-markdown": {
				option:    o.WithTextFormattingRule(FormatMarkdown),
				key:       formTextFormattingRule.Value(),
				wantValue: string(FormatMarkdown),
			},
			"WithTextFormattingRule-invalid": {
				option:  o.WithTextFormattingRule("invalid"),
				key:     formTextFormattingRule.Value(),
				wantErr: true,
			},
			"WithTextFormattingRule-invalid-empty": {
				option:  o.WithTextFormattingRule(""),
				key:     formTextFormattingRule.Value(),
				wantErr: true,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.Check()
				if tc.wantErr {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				_ = tc.option.Set(form)
				assert.Equal(t, tc.wantValue, form.Get(tc.key))
			})
		}
	})
}
