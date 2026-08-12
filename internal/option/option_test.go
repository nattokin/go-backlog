package option_test

import (
	"errors"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/option"
)

func TestAPIParamOption_Key(t *testing.T) {
	cases := map[string]struct {
		option  *option.APIParamOption
		wantKey string
	}{
		"Type-only": {
			option: &option.APIParamOption{
				Type: option.ParamKey,
			},
			wantKey: option.ParamKey.Value(),
		},
		"KeyFunc-overrides-Type": {
			option: &option.APIParamOption{
				Type:    option.ParamKey,
				KeyFunc: func() string { return "customKey" },
			},
			wantKey: "customKey",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}

func TestAPIParamOption(t *testing.T) {
	cases := map[string]struct {
		option      *option.APIParamOption
		expectPanic bool
	}{
		"SetFunc-nil": {
			option: &option.APIParamOption{
				Type:      option.ParamKey,
				CheckFunc: func() *core.ValidationError { return nil },
				SetFunc:   nil,
			},
			expectPanic: true,
		},
		"CheckFunc-nil": {
			option: &option.APIParamOption{
				Type:      option.ParamKey,
				CheckFunc: nil,
				SetFunc:   func(_ url.Values) error { return nil },
			},
			expectPanic: false,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				r := recover()

				if tc.expectPanic && r == nil {
					t.Fatalf("expected panic")
				}

				if !tc.expectPanic && r != nil {
					t.Fatalf("unexpected panic: %v", r)
				}
			}()

			v := url.Values{}
			option.ApplyOptions(v, []option.APIParamOptionType{option.ParamKey}, tc.option)
		})
	}
}

func TestApplyOptions(t *testing.T) {
	validTypes := []option.APIParamOptionType{option.ParamKey, option.ParamName}

	cases := map[string]struct {
		opts        []*option.APIParamOption
		wantErr     bool
		wantErrType any
	}{
		"nilOption": {
			opts:        []*option.APIParamOption{nil},
			wantErr:     true,
			wantErrType: &option.InvalidOptionError{},
		},
		"nilOption-second": {
			opts: []*option.APIParamOption{
				{
					Type:    option.ParamKey,
					SetFunc: func(_ url.Values) error { return nil },
				},
				nil,
			},
			wantErr:     true,
			wantErrType: &option.InvalidOptionError{},
		},
		"invalidKey": {
			opts: []*option.APIParamOption{
				{
					Type:    option.ParamOffset,
					SetFunc: func(_ url.Values) error { return nil },
				},
			},
			wantErr:     true,
			wantErrType: &option.InvalidOptionKeyError{},
		},
		"checkError-single": {
			opts: []*option.APIParamOption{
				{
					Type:      option.ParamKey,
					CheckFunc: func() *core.ValidationError { return core.NewValidationError("key", "check failed") },
					SetFunc:   func(_ url.Values) error { return nil },
				},
			},
			wantErr:     true,
			wantErrType: core.ValidationErrors(nil),
		},
		"checkError-multiple": {
			opts: []*option.APIParamOption{
				{
					Type:      option.ParamKey,
					CheckFunc: func() *core.ValidationError { return core.NewValidationError("key", "key is empty") },
					SetFunc:   func(_ url.Values) error { return nil },
				},
				{
					Type:      option.ParamName,
					CheckFunc: func() *core.ValidationError { return core.NewValidationError("name", "name is empty") },
					SetFunc:   func(_ url.Values) error { return nil },
				},
			},
			wantErr:     true,
			wantErrType: core.ValidationErrors(nil),
		},
		"setError": {
			opts: []*option.APIParamOption{
				{
					Type:    option.ParamKey,
					SetFunc: func(_ url.Values) error { return errors.New("set failed") },
				},
			},
			wantErr: true,
		},
		"success": {
			opts: []*option.APIParamOption{
				{
					Type:    option.ParamKey,
					SetFunc: func(v url.Values) error { v.Set(option.ParamKey.Value(), "val"); return nil },
				},
			},
			wantErr: false,
		},
		"noOptions": {
			opts:    []*option.APIParamOption{},
			wantErr: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := url.Values{}
			err := option.ApplyOptions(v, validTypes, tc.opts...)

			if tc.wantErr {
				require.Error(t, err)
				if tc.wantErrType != nil {
					assert.True(t, errors.As(err, &tc.wantErrType))
				}
				return
			}
			require.NoError(t, err)
		})
	}
}

// TestApplyOptions_multipleValidationErrors verifies that when multiple options
// fail Check(), all errors are collected and returned as ValidationErrors.
func TestApplyOptions_multipleValidationErrors(t *testing.T) {
	validTypes := []option.APIParamOptionType{option.ParamKey, option.ParamName}

	opt1 := &option.APIParamOption{
		Type:      option.ParamKey,
		CheckFunc: func() *core.ValidationError { return core.NewValidationError("key", "key is empty") },
		SetFunc:   func(_ url.Values) error { return nil },
	}
	opt2 := &option.APIParamOption{
		Type:      option.ParamName,
		CheckFunc: func() *core.ValidationError { return core.NewValidationError("name", "name is empty") },
		SetFunc:   func(_ url.Values) error { return nil },
	}

	v := url.Values{}
	err := option.ApplyOptions(v, validTypes, opt1, opt2)
	require.Error(t, err)

	var ves core.ValidationErrors
	require.True(t, errors.As(err, &ves))
	require.Len(t, ves, 2)
	assert.Equal(t, "key", ves[0].Target())
	assert.Equal(t, "name", ves[1].Target())
}

// TestMergeValidationErrors covers every branch of MergeValidationErrors
// directly, independent of any ApplyOptions call site. Some call sites
// (e.g. functions that build their options internally rather than
// accepting them from the caller) can never actually produce a fail-fast
// optErr in practice, so the fail-fast branch is only reachable through
// this direct test.
func TestMergeValidationErrors(t *testing.T) {
	fixedVe := core.NewValidationError("fixed", "fixed arg is invalid")
	optVe := core.NewValidationError("opt", "opt is invalid")
	failFastErr := option.NewInvalidOptionError("nil option is not allowed")

	cases := map[string]struct {
		ves    core.ValidationErrors
		optErr error

		wantNil                bool
		wantValidationErrCount int
		wantFailFastErr        error
	}{
		"no ves and no optErr": {
			ves:     nil,
			optErr:  nil,
			wantNil: true,
		},
		"ves only": {
			ves:                    core.ValidationErrors{fixedVe},
			optErr:                 nil,
			wantValidationErrCount: 1,
		},
		"optErr ValidationErrors only": {
			ves:                    nil,
			optErr:                 core.ValidationErrors{optVe},
			wantValidationErrCount: 1,
		},
		"ves and optErr ValidationErrors are merged": {
			ves:                    core.ValidationErrors{fixedVe},
			optErr:                 core.ValidationErrors{optVe},
			wantValidationErrCount: 2,
		},
		"fail-fast optErr discards ves": {
			ves:             core.ValidationErrors{fixedVe},
			optErr:          failFastErr,
			wantFailFastErr: failFastErr,
		},
		"fail-fast optErr with empty ves": {
			ves:             nil,
			optErr:          failFastErr,
			wantFailFastErr: failFastErr,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := option.MergeValidationErrors(tc.ves, tc.optErr)

			if tc.wantNil {
				assert.NoError(t, err)
				return
			}

			if tc.wantFailFastErr != nil {
				assert.Equal(t, tc.wantFailFastErr, err)
				var ves core.ValidationErrors
				assert.False(t, errors.As(err, &ves))
				return
			}

			require.Error(t, err)
			var ves core.ValidationErrors
			require.True(t, errors.As(err, &ves))
			assert.Len(t, ves, tc.wantValidationErrCount)
		})
	}
}

func TestInvalidOptionKeyError_Error_form(t *testing.T) {
	e := &option.InvalidOptionKeyError{
		Invalid: option.ParamKey.Value(),
		ValidList: []string{
			option.ParamName.Value(),
			option.ParamKey.Value(),
			option.ParamChartEnabled.Value(),
		},
	}
	assert.EqualError(t, e, "invalid option key:key, allowed option keys:name,key,chartEnabled")
}

func TestInvalidOptionKeyError_Error_query(t *testing.T) {
	e := &option.InvalidOptionKeyError{
		Invalid: option.ParamActivityTypeIDs.Value(),
		ValidList: []string{
			option.ParamAll.Value(),
			option.ParamArchived.Value(),
			option.ParamOrder.Value(),
		},
	}
	assert.EqualError(t, e, "invalid option key:activityTypeId[], allowed option keys:all,archived,order")
}

func TestInvalidOptionKeyError_errorsAs_query(t *testing.T) {
	err := option.NewInvalidOptionKeyError(option.ParamActivityTypeIDs.Value(), []option.APIParamOptionType{option.ParamAll, option.ParamArchived})
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *option.InvalidOptionKeyError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, option.ParamActivityTypeIDs.Value(), target.Invalid)
}

func TestInvalidOptionKeyError_errorsAs_form(t *testing.T) {
	err := option.NewInvalidOptionKeyError(option.ParamKey.Value(), []option.APIParamOptionType{option.ParamName, option.ParamChartEnabled})
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *option.InvalidOptionKeyError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, option.ParamKey.Value(), target.Invalid)
}

func TestInvalidOptionError_errorsAs(t *testing.T) {
	err := option.NewInvalidOptionError("nil option is not allowed")
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *option.InvalidOptionError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, "nil option is not allowed", target.Error())
}
