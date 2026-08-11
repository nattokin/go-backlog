package core_test

import (
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
)

func TestAPIParamOption_Key(t *testing.T) {
	cases := map[string]struct {
		option  *core.APIParamOption
		wantKey string
	}{
		"Type-only": {
			option: &core.APIParamOption{
				Type: core.ParamKey,
			},
			wantKey: core.ParamKey.Value(),
		},
		"KeyFunc-overrides-Type": {
			option: &core.APIParamOption{
				Type:    core.ParamKey,
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
		option      *core.APIParamOption
		expectPanic bool
	}{
		"SetFunc-nil": {
			option: &core.APIParamOption{
				Type:      core.ParamKey,
				CheckFunc: func() *core.ValidationError { return nil },
				SetFunc:   nil,
			},
			expectPanic: true,
		},
		"CheckFunc-nil": {
			option: &core.APIParamOption{
				Type:      core.ParamKey,
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
			core.ApplyOptions(v, []core.APIParamOptionType{core.ParamKey}, tc.option)
		})
	}
}

func TestApplyOptions(t *testing.T) {
	validTypes := []core.APIParamOptionType{core.ParamKey, core.ParamName}

	cases := map[string]struct {
		opts        []*core.APIParamOption
		wantErr     bool
		wantErrType any
	}{
		"nilOption": {
			opts:        []*core.APIParamOption{nil},
			wantErr:     true,
			wantErrType: &core.InvalidOptionError{},
		},
		"nilOption-second": {
			opts: []*core.APIParamOption{
				{
					Type:    core.ParamKey,
					SetFunc: func(_ url.Values) error { return nil },
				},
				nil,
			},
			wantErr:     true,
			wantErrType: &core.InvalidOptionError{},
		},
		"invalidKey": {
			opts: []*core.APIParamOption{
				{
					Type:    core.ParamOffset,
					SetFunc: func(_ url.Values) error { return nil },
				},
			},
			wantErr:     true,
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"checkError-single": {
			opts: []*core.APIParamOption{
				{
					Type:      core.ParamKey,
					CheckFunc: func() *core.ValidationError { return core.NewValidationError("key", "check failed") },
					SetFunc:   func(_ url.Values) error { return nil },
				},
			},
			wantErr:     true,
			wantErrType: core.ValidationErrors(nil),
		},
		"checkError-multiple": {
			opts: []*core.APIParamOption{
				{
					Type:      core.ParamKey,
					CheckFunc: func() *core.ValidationError { return core.NewValidationError("key", "key is empty") },
					SetFunc:   func(_ url.Values) error { return nil },
				},
				{
					Type:      core.ParamName,
					CheckFunc: func() *core.ValidationError { return core.NewValidationError("name", "name is empty") },
					SetFunc:   func(_ url.Values) error { return nil },
				},
			},
			wantErr:     true,
			wantErrType: core.ValidationErrors(nil),
		},
		"setError": {
			opts: []*core.APIParamOption{
				{
					Type:    core.ParamKey,
					SetFunc: func(_ url.Values) error { return errors.New("set failed") },
				},
			},
			wantErr: true,
		},
		"success": {
			opts: []*core.APIParamOption{
				{
					Type:    core.ParamKey,
					SetFunc: func(v url.Values) error { v.Set(core.ParamKey.Value(), "val"); return nil },
				},
			},
			wantErr: false,
		},
		"noOptions": {
			opts:    []*core.APIParamOption{},
			wantErr: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := url.Values{}
			err := core.ApplyOptions(v, validTypes, tc.opts...)

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
	validTypes := []core.APIParamOptionType{core.ParamKey, core.ParamName}

	opt1 := &core.APIParamOption{
		Type:      core.ParamKey,
		CheckFunc: func() *core.ValidationError { return core.NewValidationError("key", "key is empty") },
		SetFunc:   func(_ url.Values) error { return nil },
	}
	opt2 := &core.APIParamOption{
		Type:      core.ParamName,
		CheckFunc: func() *core.ValidationError { return core.NewValidationError("name", "name is empty") },
		SetFunc:   func(_ url.Values) error { return nil },
	}

	v := url.Values{}
	err := core.ApplyOptions(v, validTypes, opt1, opt2)
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
	failFastErr := core.NewInvalidOptionError("nil option is not allowed")

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

			err := core.MergeValidationErrors(tc.ves, tc.optErr)

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
