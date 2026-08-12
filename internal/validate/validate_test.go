package validate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/validate"
)

func TestValidateIssueIDOrKey(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string
		wantErr      bool
	}{
		"valid-key":              {issueIDOrKey: "PRJ-1"},
		"valid-id":               {issueIDOrKey: "1"},
		"error-validation-empty": {issueIDOrKey: "", wantErr: true},
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
		"valid-key":              {projectIDOrKey: "PRJ"},
		"valid-id":               {projectIDOrKey: "1"},
		"error-validation-empty": {projectIDOrKey: "", wantErr: true},
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
		"valid-name":             {repositoryIDOrName: "my-repo"},
		"valid-id":               {repositoryIDOrName: "1"},
		"error-validation-empty": {repositoryIDOrName: "", wantErr: true},
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
