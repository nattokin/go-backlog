package core_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
)

func TestOptionService_date(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		option    core.RequestOption
		key       string
		wantValue string
		wantErr   bool
	}{
		"WithCreatedSince-valid": {
			option:    o.WithCreatedSince("2024-03-15"),
			key:       core.ParamCreatedSince.Value(),
			wantValue: "2024-03-15",
		},
		"WithCreatedSince-invalid": {
			option:  o.WithCreatedSince("2024/03/15"),
			wantErr: true,
		},
		"WithCreatedUntil-valid": {
			option:    o.WithCreatedUntil("2024-03-15"),
			key:       core.ParamCreatedUntil.Value(),
			wantValue: "2024-03-15",
		},
		"WithCreatedUntil-invalid": {
			option:  o.WithCreatedUntil(""),
			wantErr: true,
		},
		"WithDueDate-valid": {
			option:    o.WithDueDate("2024-03-15"),
			key:       core.ParamDueDate.Value(),
			wantValue: "2024-03-15",
		},
		"WithDueDate-invalid": {
			option:  o.WithDueDate("20240315"),
			wantErr: true,
		},
		"WithDueDateSince-valid": {
			option:    o.WithDueDateSince("2024-03-15"),
			key:       core.ParamDueDateSince.Value(),
			wantValue: "2024-03-15",
		},
		"WithDueDateSince-invalid": {
			option:  o.WithDueDateSince(""),
			wantErr: true,
		},
		"WithDueDateUntil-valid": {
			option:    o.WithDueDateUntil("2024-03-15"),
			key:       core.ParamDueDateUntil.Value(),
			wantValue: "2024-03-15",
		},
		"WithDueDateUntil-invalid": {
			option:  o.WithDueDateUntil("2024-3-15"),
			wantErr: true,
		},
		"WithReleaseDueDate-valid": {
			option:    o.WithReleaseDueDate("2024-03-15"),
			key:       core.ParamReleaseDueDate.Value(),
			wantValue: "2024-03-15",
		},
		"WithReleaseDueDate-invalid": {
			option:  o.WithReleaseDueDate("2024/03/15"),
			wantErr: true,
		},
		"WithStartDate-valid": {
			option:    o.WithStartDate("2024-03-15"),
			key:       core.ParamStartDate.Value(),
			wantValue: "2024-03-15",
		},
		"WithStartDate-invalid": {
			option:  o.WithStartDate(""),
			wantErr: true,
		},
		"WithStartDateSince-valid": {
			option:    o.WithStartDateSince("2024-03-15"),
			key:       core.ParamStartDateSince.Value(),
			wantValue: "2024-03-15",
		},
		"WithStartDateSince-invalid": {
			option:  o.WithStartDateSince("20240315"),
			wantErr: true,
		},
		"WithStartDateUntil-valid": {
			option:    o.WithStartDateUntil("2024-03-15"),
			key:       core.ParamStartDateUntil.Value(),
			wantValue: "2024-03-15",
		},
		"WithStartDateUntil-invalid": {
			option:  o.WithStartDateUntil("2024-3-5"),
			wantErr: true,
		},
		"WithUpdatedSince-valid": {
			option:    o.WithUpdatedSince("2024-03-15"),
			key:       core.ParamUpdatedSince.Value(),
			wantValue: "2024-03-15",
		},
		"WithUpdatedSince-invalid": {
			option:  o.WithUpdatedSince(""),
			wantErr: true,
		},
		"WithUpdatedUntil-valid": {
			option:    o.WithUpdatedUntil("2024-03-15"),
			key:       core.ParamUpdatedUntil.Value(),
			wantValue: "2024-03-15",
		},
		"WithUpdatedUntil-invalid": {
			option:  o.WithUpdatedUntil("2024/03/15"),
			wantErr: true,
		},
		// WithInitialDate and friends also use dateFormatStringOption
		"WithInitialDate-valid": {
			option:    o.WithInitialDate("2024-01-15"),
			key:       core.ParamInitialDate.Value(),
			wantValue: "2024-01-15",
		},
		"WithInitialDate-invalid-format": {
			option:  o.WithInitialDate("2024/01/15"),
			wantErr: true,
		},
		"WithInitialDate-empty": {
			option:  o.WithInitialDate(""),
			wantErr: true,
		},
		"WithInitialDateMax-valid": {
			option:    o.WithInitialDateMax("2024-12-31"),
			key:       core.ParamMax.Value(),
			wantValue: "2024-12-31",
		},
		"WithInitialDateMax-invalid-format": {
			option:  o.WithInitialDateMax("20241231"),
			wantErr: true,
		},
		"WithInitialDateMax-empty": {
			option:  o.WithInitialDateMax(""),
			wantErr: true,
		},
		"WithInitialDateMin-valid": {
			option:    o.WithInitialDateMin("2024-01-01"),
			key:       core.ParamMin.Value(),
			wantValue: "2024-01-01",
		},
		"WithInitialDateMin-invalid-format": {
			option:  o.WithInitialDateMin("Jan 1 2024"),
			wantErr: true,
		},
		"WithInitialDateMin-empty": {
			option:  o.WithInitialDateMin(""),
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			form := url.Values{}
			err := tc.option.Check()
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			_ = tc.option.Set(form)
			assert.Equal(t, tc.wantValue, form.Get(tc.key))
		})
	}
}
