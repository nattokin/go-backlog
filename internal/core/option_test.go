package core_test

import (
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
)

func TestAPIParamOption(t *testing.T) {
	cases := map[string]struct {
		option      core.RequestOption
		expectPanic bool
	}{
		"SetFunc-nil": {
			option: &core.APIParamOption{
				Type:      core.ParamKey,
				CheckFunc: func() error { return nil },
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
	validTypes := []core.APIParamOptionType{core.ParamKey}

	cases := map[string]struct {
		opts        []core.RequestOption
		wantErr     bool
		wantErrType any
	}{
		"nilOption": {
			opts:        []core.RequestOption{nil},
			wantErr:     true,
			wantErrType: &core.ValidationError{},
		},
		"nilOption-second": {
			opts: []core.RequestOption{
				&core.APIParamOption{
					Type:    core.ParamKey,
					SetFunc: func(_ url.Values) error { return nil },
				},
				nil,
			},
			wantErr:     true,
			wantErrType: &core.ValidationError{},
		},
		"invalidKey": {
			opts: []core.RequestOption{
				&core.APIParamOption{
					Type:    core.ParamName,
					SetFunc: func(_ url.Values) error { return nil },
				},
			},
			wantErr:     true,
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"checkError": {
			opts: []core.RequestOption{
				&core.APIParamOption{
					Type:      core.ParamKey,
					CheckFunc: func() error { return errors.New("check failed") },
					SetFunc:   func(_ url.Values) error { return nil },
				},
			},
			wantErr: true,
		},
		"success": {
			opts: []core.RequestOption{
				&core.APIParamOption{
					Type:    core.ParamKey,
					SetFunc: func(v url.Values) error { v.Set(core.ParamKey.Value(), "val"); return nil },
				},
			},
			wantErr: false,
		},
		"noOptions": {
			opts:    []core.RequestOption{},
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
