package option_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/option"
)

func TestOptionService_date(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		option    *option.APIParamOption
		key       string
		wantValue string
		wantErr   bool
	}{
		"WithCreatedSince-valid": {
			option:    o.WithCreatedSince("2024-03-15"),
			key:       option.ParamCreatedSince.Value(),
			wantValue: "2024-03-15",
		},
		"WithCreatedSince-invalid-empty": {
			option:  o.WithCreatedSince(""),
			wantErr: true,
		},
		"WithCreatedSince-invalid-format": {
			option:  o.WithCreatedSince("2024/03/15"),
			wantErr: true,
		},
		"WithCreatedUntil-valid": {
			option:    o.WithCreatedUntil("2024-03-15"),
			key:       option.ParamCreatedUntil.Value(),
			wantValue: "2024-03-15",
		},
		"WithCreatedUntil-invalid-empty": {
			option:  o.WithCreatedUntil(""),
			wantErr: true,
		},
		"WithCreatedUntil-invalid-format": {
			option:  o.WithCreatedUntil("20240315"),
			wantErr: true,
		},
		"WithDueDate-valid": {
			option:    o.WithDueDate("2024-03-15"),
			key:       option.ParamDueDate.Value(),
			wantValue: "2024-03-15",
		},
		"WithDueDate-valid-empty": {
			option:    o.WithDueDate(""),
			key:       option.ParamDueDate.Value(),
			wantValue: "",
		},
		"WithDueDate-invalid-format": {
			option:  o.WithDueDate("20240315"),
			wantErr: true,
		},
		"WithDueDateSince-valid": {
			option:    o.WithDueDateSince("2024-03-15"),
			key:       option.ParamDueDateSince.Value(),
			wantValue: "2024-03-15",
		},
		"WithDueDateSince-invalid-empty": {
			option:  o.WithDueDateSince(""),
			wantErr: true,
		},
		"WithDueDateSince-invalid-format": {
			option:  o.WithDueDateSince("2024/03/15"),
			wantErr: true,
		},
		"WithDueDateUntil-valid": {
			option:    o.WithDueDateUntil("2024-03-15"),
			key:       option.ParamDueDateUntil.Value(),
			wantValue: "2024-03-15",
		},
		"WithDueDateUntil-invalid-empty": {
			option:  o.WithDueDateUntil(""),
			wantErr: true,
		},
		"WithDueDateUntil-invalid-format": {
			option:  o.WithDueDateUntil("2024-3-15"),
			wantErr: true,
		},
		"WithReleaseDueDate-valid": {
			option:    o.WithReleaseDueDate("2024-03-15"),
			key:       option.ParamReleaseDueDate.Value(),
			wantValue: "2024-03-15",
		},
		"WithReleaseDueDate-valid-empty": {
			option:    o.WithReleaseDueDate(""),
			key:       option.ParamReleaseDueDate.Value(),
			wantValue: "",
		},
		"WithReleaseDueDate-invalid-format": {
			option:  o.WithReleaseDueDate("2024/03/15"),
			wantErr: true,
		},
		"WithStartDate-valid": {
			option:    o.WithStartDate("2024-03-15"),
			key:       option.ParamStartDate.Value(),
			wantValue: "2024-03-15",
		},
		"WithStartDate-valid-empty": {
			option:    o.WithStartDate(""),
			key:       option.ParamStartDate.Value(),
			wantValue: "",
		},
		"WithStartDate-invalid-format": {
			option:  o.WithStartDate("20240315"),
			wantErr: true,
		},
		"WithStartDateSince-valid": {
			option:    o.WithStartDateSince("2024-03-15"),
			key:       option.ParamStartDateSince.Value(),
			wantValue: "2024-03-15",
		},
		"WithStartDateSince-invalid-empty": {
			option:  o.WithStartDateSince(""),
			wantErr: true,
		},
		"WithStartDateSince-invalid-format": {
			option:  o.WithStartDateSince("20240315"),
			wantErr: true,
		},
		"WithStartDateUntil-valid": {
			option:    o.WithStartDateUntil("2024-03-15"),
			key:       option.ParamStartDateUntil.Value(),
			wantValue: "2024-03-15",
		},
		"WithStartDateUntil-invalid-empty": {
			option:  o.WithStartDateUntil(""),
			wantErr: true,
		},
		"WithStartDateUntil-invalid-format": {
			option:  o.WithStartDateUntil("2024-3-5"),
			wantErr: true,
		},
		"WithUpdatedSince-valid": {
			option:    o.WithUpdatedSince("2024-03-15"),
			key:       option.ParamUpdatedSince.Value(),
			wantValue: "2024-03-15",
		},
		"WithUpdatedSince-invalid-empty": {
			option:  o.WithUpdatedSince(""),
			wantErr: true,
		},
		"WithUpdatedSince-invalid-format": {
			option:  o.WithUpdatedSince("2024/03/15"),
			wantErr: true,
		},
		"WithUpdatedUntil-valid": {
			option:    o.WithUpdatedUntil("2024-03-15"),
			key:       option.ParamUpdatedUntil.Value(),
			wantValue: "2024-03-15",
		},
		"WithUpdatedUntil-invalid-empty": {
			option:  o.WithUpdatedUntil(""),
			wantErr: true,
		},
		"WithUpdatedUntil-invalid-format": {
			option:  o.WithUpdatedUntil("2024/03/15"),
			wantErr: true,
		},
		"WithInitialDate-valid": {
			option:    o.WithInitialDate("2024-01-15"),
			key:       option.ParamInitialDate.Value(),
			wantValue: "2024-01-15",
		},
		"WithInitialDate-invalid-empty": {
			option:  o.WithInitialDate(""),
			wantErr: true,
		},
		"WithInitialDate-invalid-format": {
			option:  o.WithInitialDate("2024/01/15"),
			wantErr: true,
		},
		"WithMaxDate-valid": {
			option:    o.WithMaxDate("2024-12-31"),
			key:       option.ParamMax.Value(),
			wantValue: "2024-12-31",
		},
		"WithMaxDate-invalid-empty": {
			option:  o.WithMaxDate(""),
			wantErr: true,
		},
		"WithMaxDate-invalid-format": {
			option:  o.WithMaxDate("20241231"),
			wantErr: true,
		},
		"WithMinDate-valid": {
			option:    o.WithMinDate("2024-01-01"),
			key:       option.ParamMin.Value(),
			wantValue: "2024-01-01",
		},
		"WithMinDate-invalid-empty": {
			option:  o.WithMinDate(""),
			wantErr: true,
		},
		"WithMinDate-invalid-format": {
			option:  o.WithMinDate("Jan 1 2024"),
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			form := url.Values{}
			ve := tc.option.Check()
			if tc.wantErr {
				assert.NotNil(t, ve)
				return
			}
			require.Nil(t, ve)
			_ = tc.option.Set(form)
			require.True(t, form.Has(tc.key))
			assert.Equal(t, tc.wantValue, form.Get(tc.key))
		})
	}
}
