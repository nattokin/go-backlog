package validate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/validate"
)

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

func TestValidateIDOrKey(t *testing.T) {
	cases := map[string]struct {
		value   string
		wantErr bool
	}{
		"valid-key":          {value: "PRJ-1"},
		"valid-id":           {value: "1"},
		"invalid-empty":      {value: "", wantErr: true},
		"invalid-whitespace": {value: " \t\n ", wantErr: true},
		"invalid-zero":       {value: "0", wantErr: true},
		"invalid-traversal":  {value: "../projects/PRJ", wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidateIDOrKey("issueIDOrKey", tc.value)
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

func TestValidateIssueExpand(t *testing.T) {
	cases := map[string]struct {
		expand  string
		wantErr bool
	}{
		"valid-childIssueSummary": {expand: "childIssueSummary"},
		"invalid-empty":           {expand: "", wantErr: true},
		"invalid-other":           {expand: "invalid", wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidateIssueExpand("expand", tc.expand)
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

func TestValidatePathSegment(t *testing.T) {
	cases := map[string]struct {
		value      string
		wantErrMsg string
	}{
		"valid-key":            {value: "PRJ-1"},
		"valid-dots-in-name":   {value: "app.v2"},
		"valid-reserved-chars": {value: "a?b#c%d e"},
		"invalid-dot": {
			value:      ".",
			wantErrMsg: "invalid projectIDOrKey: must not be '.' or '..'",
		},
		"invalid-dot-dot": {
			value:      "..",
			wantErrMsg: "invalid projectIDOrKey: must not be '.' or '..'",
		},
		"invalid-slash": {
			value:      "PRJ/users",
			wantErrMsg: `invalid projectIDOrKey: must not contain '/' or '\'`,
		},
		"invalid-backslash": {
			value:      `PRJ\users`,
			wantErrMsg: `invalid projectIDOrKey: must not contain '/' or '\'`,
		},
		"invalid-control-character": {
			value:      "PRJ\x00",
			wantErrMsg: "invalid projectIDOrKey: must not contain control characters",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidatePathSegment("projectIDOrKey", tc.value)
			if tc.wantErrMsg != "" {
				if assert.NotNil(t, ve) {
					assert.Equal(t, "projectIDOrKey", ve.Target())
					assert.Equal(t, tc.wantErrMsg, ve.Message())
				}
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

func TestValidateDirPath(t *testing.T) {
	cases := map[string]struct {
		dirPath    string
		wantErrMsg string
	}{
		"valid-root":             {dirPath: "/"},
		"valid-nested":           {dirPath: "/PressRelease/20091130/"},
		"valid-without-slashes":  {dirPath: "design"},
		"valid-repeated-slashes": {dirPath: "//design//2024/"},
		"valid-dots-in-name":     {dirPath: "/v1.2/"},
		"invalid-empty": {
			dirPath:    "",
			wantErrMsg: "invalid dirPath: must not be empty",
		},
		"invalid-parent-directory": {
			dirPath:    "/design/../../users/1/",
			wantErrMsg: "invalid dirPath: must not be '.' or '..'",
		},
		"invalid-backslash": {
			dirPath:    `/..\design/`,
			wantErrMsg: `invalid dirPath: must not contain '/' or '\'`,
		},
		"invalid-control-character": {
			dirPath:    "/design\x00/",
			wantErrMsg: "invalid dirPath: must not contain control characters",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidateDirPath(tc.dirPath)
			if tc.wantErrMsg != "" {
				if assert.NotNil(t, ve) {
					assert.Equal(t, "dirPath", ve.Target())
					assert.Equal(t, tc.wantErrMsg, ve.Message())
				}
				return
			}
			assert.Nil(t, ve)
		})
	}
}

func TestValidateIssueIDOrKey(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string
		wantErr      bool
	}{
		"valid-key":                         {issueIDOrKey: "PRJ-1"},
		"valid-id":                          {issueIDOrKey: "1"},
		"error-validation-empty":            {issueIDOrKey: "", wantErr: true},
		"error-validation-whitespace-space": {issueIDOrKey: " ", wantErr: true},
		"error-validation-whitespace-tab":   {issueIDOrKey: "\t", wantErr: true},
		"error-validation-whitespace-mixed": {issueIDOrKey: " \t\n ", wantErr: true},
		"error-validation-zero":             {issueIDOrKey: "0", wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidateIssueIDOrKey(tc.issueIDOrKey)
			if tc.wantErr {
				assert.NotNil(t, ve)
				return
			}
			assert.Nil(t, ve)
		})
	}
}

func TestValidateProjectIDOrKey(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		wantErr        bool
	}{
		"valid-key":                         {projectIDOrKey: "PRJ"},
		"valid-id":                          {projectIDOrKey: "1"},
		"error-validation-empty":            {projectIDOrKey: "", wantErr: true},
		"error-validation-whitespace-space": {projectIDOrKey: " ", wantErr: true},
		"error-validation-whitespace-tab":   {projectIDOrKey: "\t", wantErr: true},
		"error-validation-whitespace-mixed": {projectIDOrKey: " \t\n ", wantErr: true},
		"error-validation-zero":             {projectIDOrKey: "0", wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidateProjectIDOrKey(tc.projectIDOrKey)
			if tc.wantErr {
				assert.NotNil(t, ve)
				return
			}
			assert.Nil(t, ve)
		})
	}
}

func TestValidateRepositoryIDOrName(t *testing.T) {
	cases := map[string]struct {
		repositoryIDOrName string
		wantErr            bool
	}{
		"valid-name":                        {repositoryIDOrName: "my-repo"},
		"valid-id":                          {repositoryIDOrName: "1"},
		"error-validation-empty":            {repositoryIDOrName: "", wantErr: true},
		"error-validation-whitespace-space": {repositoryIDOrName: " ", wantErr: true},
		"error-validation-whitespace-tab":   {repositoryIDOrName: "\t", wantErr: true},
		"error-validation-whitespace-mixed": {repositoryIDOrName: " \t\n ", wantErr: true},
		"error-validation-zero":             {repositoryIDOrName: "0", wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ve := validate.ValidateRepositoryIDOrName(tc.repositoryIDOrName)
			if tc.wantErr {
				assert.NotNil(t, ve)
				return
			}
			assert.Nil(t, ve)
		})
	}
}
