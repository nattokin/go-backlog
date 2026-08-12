package validate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/validate"
)

func TestValidatePositiveInt(t *testing.T) {
	cases := map[string]struct {
		value   int
		wantErr bool
	}{
		"valid-1":     {value: 1},
		"valid-2":     {value: 2},
		"invalid-0":   {value: 0, wantErr: true},
		"invalid-neg": {value: -1, wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidatePositiveInt("userId", tc.value)
			if tc.wantErr {
				assert.NotNil(t, ve)
				assert.Equal(t, "userId", ve.Target())
				return
			}
			assert.Nil(t, ve)
		})
	}
}

func TestValidatePositiveInts(t *testing.T) {
	cases := map[string]struct {
		values  []int
		wantErr bool
	}{
		"valid-empty":     {values: []int{}},
		"valid-all-1":     {values: []int{1, 2, 3}},
		"invalid-has-0":   {values: []int{1, 0, 3}, wantErr: true},
		"invalid-has-neg": {values: []int{1, -1}, wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidatePositiveInts("userId", tc.values)
			if tc.wantErr {
				assert.NotNil(t, ve)
				return
			}
			assert.Nil(t, ve)
		})
	}
}

func TestValidateIntRange(t *testing.T) {
	cases := map[string]struct {
		value   int
		wantErr bool
	}{
		"valid-min":     {value: 1},
		"valid-mid":     {value: 50},
		"valid-max":     {value: 100},
		"invalid-below": {value: 0, wantErr: true},
		"invalid-above": {value: 101, wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidateIntRange("count", tc.value, 1, 100)
			if tc.wantErr {
				assert.NotNil(t, ve)
				return
			}
			assert.Nil(t, ve)
		})
	}
}

func TestValidateActivityTypeID(t *testing.T) {
	cases := map[string]struct {
		id      int
		wantErr bool
	}{
		"valid-min":     {id: 1},
		"valid-max":     {id: 26},
		"invalid-below": {id: 0, wantErr: true},
		"invalid-above": {id: 27, wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidateActivityTypeID("activityTypeIds", tc.id, 26)
			if tc.wantErr {
				assert.NotNil(t, ve)
				return
			}
			assert.Nil(t, ve)
		})
	}
}

func TestValidatePositiveFloat64(t *testing.T) {
	cases := map[string]struct {
		value   float64
		wantErr bool
	}{
		"valid-small":  {value: 0.1},
		"valid-large":  {value: 100.5},
		"invalid-zero": {value: 0, wantErr: true},
		"invalid-neg":  {value: -1, wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidatePositiveFloat64("estimatedHours", tc.value)
			if tc.wantErr {
				assert.NotNil(t, ve)
				return
			}
			assert.Nil(t, ve)
		})
	}
}

func TestValidateDateFormat(t *testing.T) {
	cases := map[string]struct {
		date    string
		wantErr bool
	}{
		"valid":              {date: "2024-01-01"},
		"invalid-empty":      {date: "", wantErr: true},
		"invalid-slash":      {date: "2024/01/01", wantErr: true},
		"invalid-no-padding": {date: "2024-1-1", wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidateDateFormat("dueDate", tc.date)
			if tc.wantErr {
				assert.NotNil(t, ve)
				return
			}
			assert.Nil(t, ve)
		})
	}
}

func TestValidateNonEmptyString(t *testing.T) {
	cases := map[string]struct {
		value   string
		wantErr bool
	}{
		"valid":              {value: "main"},
		"invalid-empty":      {value: "", wantErr: true},
		"invalid-whitespace": {value: "   ", wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidateNonEmptyString("branch", tc.value)
			if tc.wantErr {
				assert.NotNil(t, ve)
				return
			}
			assert.Nil(t, ve)
		})
	}
}

func TestValidateEmail(t *testing.T) {
	cases := map[string]struct {
		value   string
		wantErr bool
	}{
		"valid":                 {value: "test@example.com"},
		"valid-plus":            {value: "user+tag@example.co.jp"},
		"invalid-empty":         {value: "", wantErr: true},
		"invalid-no-at":         {value: "notanemail", wantErr: true},
		"invalid-display-name":  {value: "John Doe <john@example.com>", wantErr: true},
		"invalid-angle-bracket": {value: "<john@example.com>", wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidateEmail("mailAddress", tc.value)
			if tc.wantErr {
				assert.NotNil(t, ve)
				return
			}
			assert.Nil(t, ve)
		})
	}
}

func TestValidateOrder(t *testing.T) {
	cases := map[string]struct {
		order   string
		wantErr bool
	}{
		"valid-asc":     {order: "asc"},
		"valid-desc":    {order: "desc"},
		"invalid-empty": {order: "", wantErr: true},
		"invalid-other": {order: "invalid", wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidateOrder("order", tc.order)
			if tc.wantErr {
				assert.NotNil(t, ve)
				return
			}
			assert.Nil(t, ve)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	cases := map[string]struct {
		password string
		wantErr  bool
	}{
		"invalid-empty":  {password: "", wantErr: true},
		"invalid-7chars": {password: "abcdefg", wantErr: true},
		"valid-8chars":   {password: "abcdefgh"},
		"valid-9chars":   {password: "abcdefghi"},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidatePassword("password", tc.password)
			if tc.wantErr {
				assert.NotNil(t, ve)
				return
			}
			assert.Nil(t, ve)
		})
	}
}

func TestValidateIssueSort(t *testing.T) {
	cases := map[string]struct {
		sort    string
		wantErr bool
	}{
		"valid-summary": {sort: "summary"},
		"valid-status":  {sort: "status"},
		"invalid-empty": {sort: "", wantErr: true},
		"invalid-other": {sort: "invalid", wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidateIssueSort("sort", tc.sort)
			if tc.wantErr {
				assert.NotNil(t, ve)
				return
			}
			assert.Nil(t, ve)
		})
	}
}

func TestValidateTextFormattingRule(t *testing.T) {
	cases := map[string]struct {
		format  string
		wantErr bool
	}{
		"valid-backlog":  {format: "backlog"},
		"valid-markdown": {format: "markdown"},
		"invalid-empty":  {format: "", wantErr: true},
		"invalid-other":  {format: "invalid", wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidateTextFormattingRule("textFormattingRule", tc.format)
			if tc.wantErr {
				assert.NotNil(t, ve)
				return
			}
			assert.Nil(t, ve)
		})
	}
}
