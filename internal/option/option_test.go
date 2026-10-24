package option_test

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/validation"
)

// TestAPIParamOptionType_Value pins every Param* constant to the literal
// parameter name documented in the Backlog API. Other tests usually compare
// against option.ParamXxx.Value(), which cannot catch a constant whose value
// is wrong; this table is the single place where the literal values are checked.
func TestAPIParamOptionType_Value(t *testing.T) {
	want := map[option.APIParamOptionType]string{
		option.ParamActivityTypeIDs:                   "activityTypeId[]",
		option.ParamActualHours:                       "actualHours",
		option.ParamAll:                               "all",
		option.ParamAllEvent:                          "allEvent",
		option.ParamAllowAddItem:                      "allowAddItem",
		option.ParamAllowInput:                        "allowInput",
		option.ParamApplicableIssueTypeIDs:            "applicableIssueTypes[]",
		option.ParamArchived:                          "archived",
		option.ParamAssigneeID:                        "assigneeId",
		option.ParamAssigneeIDs:                       "assigneeId[]",
		option.ParamAttachment:                        "attachment",
		option.ParamAttachmentIDs:                     "attachmentId[]",
		option.ParamBase:                              "base",
		option.ParamBranch:                            "branch",
		option.ParamCategoryIDs:                       "categoryId[]",
		option.ParamChartEnabled:                      "chartEnabled",
		option.ParamColor:                             "color",
		option.ParamComment:                           "comment",
		option.ParamCommentID:                         "commentId",
		option.ParamContent:                           "content",
		option.ParamCount:                             "count",
		option.ParamCreatedSince:                      "createdSince",
		option.ParamCreatedUntil:                      "createdUntil",
		option.ParamCreatedUserIDs:                    "createdUserId[]",
		option.ParamCustomField:                       "customField",
		option.ParamDescription:                       "description",
		option.ParamDueDate:                           "dueDate",
		option.ParamDueDateSince:                      "dueDateSince",
		option.ParamDueDateUntil:                      "dueDateUntil",
		option.ParamEstimatedHours:                    "estimatedHours",
		option.ParamExcludeGroupMembers:               "excludeGroupMembers",
		option.ParamHasDueDate:                        "hasDueDate",
		option.ParamHookURL:                           "hookUrl",
		option.ParamIDs:                               "id[]",
		option.ParamInitialDate:                       "initialDate",
		option.ParamInitialShift:                      "initialShift",
		option.ParamInitialValue:                      "initialValue",
		option.ParamInitialValueType:                  "initialValueType",
		option.ParamIssueID:                           "issueId",
		option.ParamIssueIDs:                          "issueId[]",
		option.ParamIssueTypeID:                       "issueTypeId",
		option.ParamIssueTypeIDs:                      "issueTypeId[]",
		option.ParamItems:                             "items[]",
		option.ParamKey:                               "key",
		option.ParamKeyword:                           "keyword",
		option.ParamMailAddress:                       "mailAddress",
		option.ParamMailNotify:                        "mailNotify",
		option.ParamMax:                               "max",
		option.ParamMaxID:                             "maxId",
		option.ParamMilestoneIDs:                      "milestoneId[]",
		option.ParamMin:                               "min",
		option.ParamMinID:                             "minId",
		option.ParamName:                              "name",
		option.ParamNotifiedUserIDs:                   "notifiedUserId[]",
		option.ParamOffset:                            "offset",
		option.ParamOrder:                             "order",
		option.ParamParentChild:                       "parentChild",
		option.ParamParentIssueID:                     "parentIssueId",
		option.ParamParentIssueIDs:                    "parentIssueId[]",
		option.ParamPassword:                          "password",
		option.ParamPriorityID:                        "priorityId",
		option.ParamPriorityIDs:                       "priorityId[]",
		option.ParamProjectIDs:                        "projectId[]",
		option.ParamProjectLeaderCanEditProjectLeader: "projectLeaderCanEditProjectLeader",
		option.ParamPullRequestCommentID:              "pullRequestCommentId",
		option.ParamPullRequestID:                     "pullRequestId",
		option.ParamReleaseDueDate:                    "releaseDueDate",
		option.ParamRequired:                          "required",
		option.ParamResolutionID:                      "resolutionId",
		option.ParamResolutionIDs:                     "resolutionId[]",
		option.ParamRoleType:                          "roleType",
		option.ParamSharedFile:                        "sharedFile",
		option.ParamSort:                              "sort",
		option.ParamStartDate:                         "startDate",
		option.ParamStartDateSince:                    "startDateSince",
		option.ParamStartDateUntil:                    "startDateUntil",
		option.ParamStatusID:                          "statusId",
		option.ParamStatusIDs:                         "statusId[]",
		option.ParamSubtaskingEnabled:                 "subtaskingEnabled",
		option.ParamSummary:                           "summary",
		option.ParamTemplateDescription:               "templateDescription",
		option.ParamTemplateSummary:                   "templateSummary",
		option.ParamTextFormattingRule:                "textFormattingRule",
		option.ParamTypeID:                            "typeId",
		option.ParamUnit:                              "unit",
		option.ParamUpdatedSince:                      "updatedSince",
		option.ParamUpdatedUntil:                      "updatedUntil",
		option.ParamUserID:                            "userId",
		option.ParamVersionIDs:                        "versionId[]",
		option.ParamWebhookActivityTypeIDs:            "activityTypeIds[]",
		option.ParamWikiID:                            "wikiId",
	}

	for typ, value := range want {
		t.Run(value, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, value, typ.Value())
		})
	}

	t.Run("covers-all-constants", func(t *testing.T) {
		t.Parallel()

		got := declaredParamConstants(t)
		assert.Len(t, want, len(got), "every Param* constant in option.go must be listed in this table")
	})
}

// declaredParamConstants returns the names of all APIParamOptionType constants
// declared in option.go, so that adding a constant without pinning it fails.
func declaredParamConstants(t *testing.T) []string {
	t.Helper()

	f, err := parser.ParseFile(token.NewFileSet(), "option.go", nil, 0)
	require.NoError(t, err)

	var names []string
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			ident, ok := vs.Type.(*ast.Ident)
			if !ok || ident.Name != "APIParamOptionType" {
				continue
			}
			for _, n := range vs.Names {
				names = append(names, n.Name)
			}
		}
	}
	return names
}

func TestAPIParamOption_Key(t *testing.T) {
	cases := map[string]struct {
		option  *option.APIParamOption
		wantKey string
	}{
		"Type-only": {
			option: &option.APIParamOption{
				Type: option.ParamKey,
			},
			wantKey: option.ParamKey.Value(),
		},
		"KeyFunc-overrides-Type": {
			option: &option.APIParamOption{
				Type:    option.ParamKey,
				KeyFunc: func() string { return "customKey" },
			},
			wantKey: "customKey",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}

func TestAPIParamOption(t *testing.T) {
	cases := map[string]struct {
		option      *option.APIParamOption
		expectPanic bool
	}{
		"SetFunc-nil": {
			option: &option.APIParamOption{
				Type:      option.ParamKey,
				CheckFunc: func() *validation.Error { return nil },
				SetFunc:   nil,
			},
			expectPanic: true,
		},
		"CheckFunc-nil": {
			option: &option.APIParamOption{
				Type:      option.ParamKey,
				CheckFunc: nil,
				SetFunc:   func(_ url.Values) error { return nil },
			},
			expectPanic: false,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				r := recover()

				if tc.expectPanic && r == nil {
					t.Fatalf("expected panic")
				}

				if !tc.expectPanic && r != nil {
					t.Fatalf("unexpected panic: %v", r)
				}
			}()

			v := url.Values{}
			option.ApplyOptions(v, []option.APIParamOptionType{option.ParamKey}, tc.option)
		})
	}
}

func TestApplyOptions(t *testing.T) {
	validTypes := []option.APIParamOptionType{option.ParamKey, option.ParamName}

	cases := map[string]struct {
		opts        []*option.APIParamOption
		wantErr     bool
		wantErrType any
	}{
		"nilOption": {
			opts:        []*option.APIParamOption{nil},
			wantErr:     true,
			wantErrType: &option.InvalidOptionError{},
		},
		"nilOption-second": {
			opts: []*option.APIParamOption{
				{
					Type:    option.ParamKey,
					SetFunc: func(_ url.Values) error { return nil },
				},
				nil,
			},
			wantErr:     true,
			wantErrType: &option.InvalidOptionError{},
		},
		"invalidKey": {
			opts: []*option.APIParamOption{
				{
					Type:    option.ParamOffset,
					SetFunc: func(_ url.Values) error { return nil },
				},
			},
			wantErr:     true,
			wantErrType: &option.InvalidOptionKeyError{},
		},
		"checkError-single": {
			opts: []*option.APIParamOption{
				{
					Type:      option.ParamKey,
					CheckFunc: func() *validation.Error { return validation.NewError("key", "check failed") },
					SetFunc:   func(_ url.Values) error { return nil },
				},
			},
			wantErr:     true,
			wantErrType: validation.Errors(nil),
		},
		"checkError-multiple": {
			opts: []*option.APIParamOption{
				{
					Type:      option.ParamKey,
					CheckFunc: func() *validation.Error { return validation.NewError("key", "key is empty") },
					SetFunc:   func(_ url.Values) error { return nil },
				},
				{
					Type:      option.ParamName,
					CheckFunc: func() *validation.Error { return validation.NewError("name", "name is empty") },
					SetFunc:   func(_ url.Values) error { return nil },
				},
			},
			wantErr:     true,
			wantErrType: validation.Errors(nil),
		},
		"setError": {
			opts: []*option.APIParamOption{
				{
					Type:    option.ParamKey,
					SetFunc: func(_ url.Values) error { return errors.New("set failed") },
				},
			},
			wantErr: true,
		},
		"success": {
			opts: []*option.APIParamOption{
				{
					Type:    option.ParamKey,
					SetFunc: func(v url.Values) error { v.Set(option.ParamKey.Value(), "val"); return nil },
				},
			},
			wantErr: false,
		},
		"noOptions": {
			opts:    []*option.APIParamOption{},
			wantErr: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := url.Values{}
			err := option.ApplyOptions(v, validTypes, tc.opts...)

			if tc.wantErr {
				require.Error(t, err)
				if tc.wantErrType != nil {
					assert.True(t, errors.As(err, &tc.wantErrType))
				}
				return
			}
			require.NoError(t, err)
		})
	}
}

// TestApplyOptions_multipleValidationErrors verifies that when multiple options
// fail Check(), all errors are collected and returned as ValidationErrors.
func TestApplyOptions_multipleValidationErrors(t *testing.T) {
	validTypes := []option.APIParamOptionType{option.ParamKey, option.ParamName}

	opt1 := &option.APIParamOption{
		Type:      option.ParamKey,
		CheckFunc: func() *validation.Error { return validation.NewError("key", "key is empty") },
		SetFunc:   func(_ url.Values) error { return nil },
	}
	opt2 := &option.APIParamOption{
		Type:      option.ParamName,
		CheckFunc: func() *validation.Error { return validation.NewError("name", "name is empty") },
		SetFunc:   func(_ url.Values) error { return nil },
	}

	v := url.Values{}
	err := option.ApplyOptions(v, validTypes, opt1, opt2)
	require.Error(t, err)

	var ves validation.Errors
	require.True(t, errors.As(err, &ves))
	require.Len(t, ves, 2)
	assert.Equal(t, "key", ves[0].Target())
	assert.Equal(t, "name", ves[1].Target())
}

// TestMergeValidationErrors covers every branch of MergeValidationErrors
// directly, independent of any ApplyOptions call site. Some call sites
// (e.g. functions that build their options internally rather than
// accepting them from the caller) can never actually produce a fail-fast
// optErr in practice, so the fail-fast branch is only reachable through
// this direct test.
func TestMergeValidationErrors(t *testing.T) {
	fixedVe := validation.NewError("fixed", "fixed arg is invalid")
	optVe := validation.NewError("opt", "opt is invalid")
	failFastErr := option.NewInvalidOptionError("nil option is not allowed")

	cases := map[string]struct {
		ves    validation.Errors
		optErr error

		wantNil                bool
		wantValidationErrCount int
		wantFailFastErr        error
	}{
		"no ves and no optErr": {
			ves:     nil,
			optErr:  nil,
			wantNil: true,
		},
		"ves only": {
			ves:                    validation.Errors{fixedVe},
			optErr:                 nil,
			wantValidationErrCount: 1,
		},
		"optErr ValidationErrors only": {
			ves:                    nil,
			optErr:                 validation.Errors{optVe},
			wantValidationErrCount: 1,
		},
		"ves and optErr ValidationErrors are merged": {
			ves:                    validation.Errors{fixedVe},
			optErr:                 validation.Errors{optVe},
			wantValidationErrCount: 2,
		},
		"fail-fast optErr discards ves": {
			ves:             validation.Errors{fixedVe},
			optErr:          failFastErr,
			wantFailFastErr: failFastErr,
		},
		"fail-fast optErr with empty ves": {
			ves:             nil,
			optErr:          failFastErr,
			wantFailFastErr: failFastErr,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := option.MergeValidationErrors(tc.ves, tc.optErr)

			if tc.wantNil {
				assert.NoError(t, err)
				return
			}

			if tc.wantFailFastErr != nil {
				assert.Equal(t, tc.wantFailFastErr, err)
				var ves validation.Errors
				assert.False(t, errors.As(err, &ves))
				return
			}

			require.Error(t, err)
			var ves validation.Errors
			require.True(t, errors.As(err, &ves))
			assert.Len(t, ves, tc.wantValidationErrCount)
		})
	}
}

func TestInvalidOptionKeyError_Error_form(t *testing.T) {
	e := &option.InvalidOptionKeyError{
		Invalid: option.ParamKey.Value(),
		ValidList: []string{
			option.ParamName.Value(),
			option.ParamKey.Value(),
			option.ParamChartEnabled.Value(),
		},
	}
	assert.EqualError(t, e, "invalid option key:key, allowed option keys:name,key,chartEnabled")
}

func TestInvalidOptionKeyError_Error_query(t *testing.T) {
	e := &option.InvalidOptionKeyError{
		Invalid: option.ParamActivityTypeIDs.Value(),
		ValidList: []string{
			option.ParamAll.Value(),
			option.ParamArchived.Value(),
			option.ParamOrder.Value(),
		},
	}
	assert.EqualError(t, e, "invalid option key:activityTypeId[], allowed option keys:all,archived,order")
}

func TestInvalidOptionKeyError_errorsAs_query(t *testing.T) {
	err := option.NewInvalidOptionKeyError(option.ParamActivityTypeIDs.Value(), []option.APIParamOptionType{option.ParamAll, option.ParamArchived})
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *option.InvalidOptionKeyError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, option.ParamActivityTypeIDs.Value(), target.Invalid)
}

func TestInvalidOptionKeyError_errorsAs_form(t *testing.T) {
	err := option.NewInvalidOptionKeyError(option.ParamKey.Value(), []option.APIParamOptionType{option.ParamName, option.ParamChartEnabled})
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *option.InvalidOptionKeyError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, option.ParamKey.Value(), target.Invalid)
}

func TestInvalidOptionError_errorsAs(t *testing.T) {
	err := option.NewInvalidOptionError("nil option is not allowed")
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *option.InvalidOptionError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, "nil option is not allowed", target.Error())
}
