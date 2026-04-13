package backlog

import (
	"net/url"
	"testing"

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
