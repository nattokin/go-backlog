package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomFieldBound_UnmarshalJSON(t *testing.T) {
	number := 1.5
	date := time.Date(2024, 1, 14, 15, 0, 0, 0, time.UTC)

	cases := map[string]struct {
		input   string
		want    *CustomFieldBound
		wantErr bool
	}{
		"number": {
			input: `{"min": 1.5}`,
			want:  &CustomFieldBound{Number: &number},
		},
		"date": {
			input: `{"min": "2024-01-14T15:00:00Z"}`,
			want:  &CustomFieldBound{Date: &date},
		},
		"null": {
			input: `{"min": null}`,
			want:  nil,
		},
		"absent": {
			input: `{}`,
			want:  nil,
		},
		"error-invalid-date": {
			input:   `{"min": "2024-01-15"}`,
			wantErr: true,
		},
		"error-invalid-type": {
			input:   `{"min": true}`,
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var got CustomField
			err := json.Unmarshal([]byte(tc.input), &got)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got.Min)
		})
	}
}
