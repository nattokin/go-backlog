package backlog_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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
