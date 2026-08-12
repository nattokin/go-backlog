package option_test

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/option"
)

func TestOptionService_float(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		option    *option.APIParamOption
		key       string
		wantValue float64
	}{
		"WithMin-positive": {
			option:    o.WithMin(1.5),
			key:       option.ParamMin.Value(),
			wantValue: 1.5,
		},
		"WithMin-zero": {
			option:    o.WithMin(0),
			key:       option.ParamMin.Value(),
			wantValue: 0,
		},
		"WithMin-negative": {
			option:    o.WithMin(-10),
			key:       option.ParamMin.Value(),
			wantValue: -10,
		},
		"WithMax-positive": {
			option:    o.WithMax(100.5),
			key:       option.ParamMax.Value(),
			wantValue: 100.5,
		},
		"WithMax-zero": {
			option:    o.WithMax(0),
			key:       option.ParamMax.Value(),
			wantValue: 0,
		},
		"WithMax-negative": {
			option:    o.WithMax(-1),
			key:       option.ParamMax.Value(),
			wantValue: -1,
		},
		"WithInitialValue-positive": {
			option:    o.WithInitialValue(3.14),
			key:       option.ParamInitialValue.Value(),
			wantValue: 3.14,
		},
		"WithInitialValue-zero": {
			option:    o.WithInitialValue(0),
			key:       option.ParamInitialValue.Value(),
			wantValue: 0,
		},
		"WithInitialValue-negative": {
			option:    o.WithInitialValue(-5.5),
			key:       option.ParamInitialValue.Value(),
			wantValue: -5.5,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			form := url.Values{}
			require.Nil(t, tc.option.Check())
			_ = tc.option.Set(form)
			assert.Equal(t, strconv.FormatFloat(tc.wantValue, 'f', -1, 64), form.Get(tc.key))
		})
	}
}
