package core_test

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
)

func TestOptionService_float(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		option    core.RequestOption
		key       string
		wantValue float64
	}{
		"WithMin-positive": {
			option:    o.WithMin(1.5),
			key:       core.ParamMin.Value(),
			wantValue: 1.5,
		},
		"WithMin-zero": {
			option:    o.WithMin(0),
			key:       core.ParamMin.Value(),
			wantValue: 0,
		},
		"WithMin-negative": {
			option:    o.WithMin(-10),
			key:       core.ParamMin.Value(),
			wantValue: -10,
		},
		"WithMax-positive": {
			option:    o.WithMax(100.5),
			key:       core.ParamMax.Value(),
			wantValue: 100.5,
		},
		"WithMax-zero": {
			option:    o.WithMax(0),
			key:       core.ParamMax.Value(),
			wantValue: 0,
		},
		"WithMax-negative": {
			option:    o.WithMax(-1),
			key:       core.ParamMax.Value(),
			wantValue: -1,
		},
		"WithInitialValue-positive": {
			option:    o.WithInitialValue(3.14),
			key:       core.ParamInitialValue.Value(),
			wantValue: 3.14,
		},
		"WithInitialValue-zero": {
			option:    o.WithInitialValue(0),
			key:       core.ParamInitialValue.Value(),
			wantValue: 0,
		},
		"WithInitialValue-negative": {
			option:    o.WithInitialValue(-5.5),
			key:       core.ParamInitialValue.Value(),
			wantValue: -5.5,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			form := url.Values{}
			err := tc.option.Check()
			require.NoError(t, err)
			_ = tc.option.Set(form)
			assert.Equal(t, strconv.FormatFloat(tc.wantValue, 'f', -1, 64), form.Get(tc.key))
		})
	}
}
