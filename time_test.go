package backlog_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
)

func TestTimestamp_embedsTimeTime(t *testing.T) {
	t.Parallel()

	base := time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC)
	ts := backlog.Timestamp{Time: base}

	assert.Equal(t, base.Year(), ts.Year())
	assert.Equal(t, base.Month(), ts.Month())
	assert.Equal(t, base.Day(), ts.Day())
	assert.True(t, base.Equal(ts.Time))
}

func TestTimestamp_zero(t *testing.T) {
	t.Parallel()

	var ts backlog.Timestamp
	assert.True(t, ts.IsZero())
}

func TestNewDate(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		input   string
		wantErr bool
	}{
		"valid": {
			input:   "2024-03-31",
			wantErr: false,
		},
		"invalid_format": {
			input:   "2024/03/31",
			wantErr: true,
		},
		"invalid_month": {
			input:   "2024-13-01",
			wantErr: true,
		},
		"invalid_day": {
			input:   "2024-01-32",
			wantErr: true,
		},
		"empty": {
			input:   "",
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			d, err := backlog.NewDate(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				var target *backlog.InvalidDateStringError
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, tc.input, target.Value)
				assert.True(t, d.IsZero())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.input, d.String())
				assert.False(t, d.IsZero())
			}
		})
	}
}

func TestDate_String(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		date backlog.Date
		want string
	}{
		"with_value": {
			date: func() backlog.Date {
				d, _ := backlog.NewDate("2024-03-31")
				return d
			}(),
			want: "2024-03-31",
		},
		"zero": {
			date: backlog.Date{},
			want: "",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, tc.date.String())
		})
	}
}

func TestDate_IsZero(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		date backlog.Date
		want bool
	}{
		"with_value": {
			date: func() backlog.Date {
				d, _ := backlog.NewDate("2024-03-31")
				return d
			}(),
			want: false,
		},
		"zero": {
			date: backlog.Date{},
			want: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, tc.date.IsZero())
		})
	}
}

func TestInvalidDateStringError_Error(t *testing.T) {
	t.Parallel()

	err := &backlog.InvalidDateStringError{Value: "2024/03/31"}
	assert.Equal(t, `backlog: invalid date string "2024/03/31": expected "YYYY-MM-DD" format`, err.Error())
}
