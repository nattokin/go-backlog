package core_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
)

func TestOptionService_time(t *testing.T) {
	o := &core.OptionService{}

	date := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	want := "2024-03-15"

	cases := map[string]struct {
		option core.RequestOption
		key    string
	}{
		"WithCreatedSince":   {option: o.WithCreatedSince(date), key: core.ParamCreatedSince.Value()},
		"WithCreatedUntil":   {option: o.WithCreatedUntil(date), key: core.ParamCreatedUntil.Value()},
		"WithUpdatedSince":   {option: o.WithUpdatedSince(date), key: core.ParamUpdatedSince.Value()},
		"WithUpdatedUntil":   {option: o.WithUpdatedUntil(date), key: core.ParamUpdatedUntil.Value()},
		"WithStartDateSince": {option: o.WithStartDateSince(date), key: core.ParamStartDateSince.Value()},
		"WithStartDateUntil": {option: o.WithStartDateUntil(date), key: core.ParamStartDateUntil.Value()},
		"WithDueDateSince":   {option: o.WithDueDateSince(date), key: core.ParamDueDateSince.Value()},
		"WithDueDateUntil":   {option: o.WithDueDateUntil(date), key: core.ParamDueDateUntil.Value()},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			q := url.Values{}
			err := tc.option.Check()
			require.NoError(t, err)
			_ = tc.option.Set(q)
			assert.Equal(t, want, q.Get(tc.key))
		})
	}
}
