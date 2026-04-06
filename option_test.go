package backlog

import (
	"net/url"
	"testing"
)

func Test_apiParamOption(t *testing.T) {
	cases := map[string]struct {
		option      RequestOption
		expectPanic bool
	}{
		"setFunc-nil": {
			option: &apiParamOption{
				t:         paramKey,
				checkFunc: func() error { return nil },
				setFunc:   nil,
			},
			expectPanic: true,
		},
		"checkFunc-nil": {
			option: &apiParamOption{
				t:         paramKey,
				checkFunc: nil,
				setFunc:   func(_ url.Values) error { return nil },
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
			applyOptions(v, []apiParamOptionType{paramKey}, tc.option)
		})
	}

}
