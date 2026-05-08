package activity_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/space"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestBaseActivityService_GetList(t *testing.T) {
	o := &core.OptionService{}
	type want struct {
		activityTypeID []string
		minID          string
		maxID          string
		count          string
		order          string
	}
	cases := map[string]struct {
		opts      []core.RequestOption
		wantError bool
		want      want
	}{
		"success-no-option": {
			opts:      []core.RequestOption{},
			wantError: false,
			want: want{
				activityTypeID: nil,
				minID:          "",
				maxID:          "",
				count:          "",
				order:          "",
			},
		},
		"success-withActivityTypeIDs": {
			opts: []core.RequestOption{
				o.WithActivityTypeIDs([]int{1}),
			},
			wantError: false,
			want: want{
				activityTypeID: []string{"1"},
				minID:          "",
				maxID:          "",
				count:          "",
				order:          "",
			},
		},
		"success-withMinID": {
			opts: []core.RequestOption{
				o.WithMinID(1),
			},
			wantError: false,
			want: want{
				activityTypeID: nil,
				minID:          "1",
				maxID:          "",
				count:          "",
				order:          "",
			},
		},
		"success-withMaxID": {
			opts: []core.RequestOption{
				o.WithMaxID(1),
			},
			wantError: false,
			want: want{
				activityTypeID: nil,
				minID:          "",
				maxID:          "1",
				count:          "",
				order:          "",
			},
		},
		"success-withCount": {
			opts: []core.RequestOption{
				o.WithCount(1),
			},
			wantError: false,
			want: want{
				activityTypeID: nil,
				minID:          "",
				maxID:          "",
				count:          "1",
				order:          "",
			},
		},
		"success-withOrder": {
			opts: []core.RequestOption{
				o.WithOrder(model.OrderAsc),
			},
			wantError: false,
			want: want{
				activityTypeID: nil,
				minID:          "",
				maxID:          "",
				count:          "",
				order:          "asc",
			},
		},
		"success-multiple-options": {
			opts: []core.RequestOption{
				o.WithActivityTypeIDs([]int{1, 2}),
				o.WithMinID(1),
				o.WithMaxID(26),
				o.WithCount(20),
				o.WithOrder(model.OrderAsc),
			},
			wantError: false,
			want: want{
				activityTypeID: []string{"1", "2"},
				minID:          "1",
				maxID:          "26",
				count:          "20",
				order:          "asc",
			},
		},
		"error-option-invalid-value": {
			opts: []core.RequestOption{
				o.WithCount(0),
			},
			wantError: true,
			want:      want{},
		},
		"error-option-invalid-type": {
			opts:      []core.RequestOption{mock.NewInvalidTypeOption()},
			wantError: true,
			want:      want{},
		},
		"error-option-set-failed": {
			opts: []core.RequestOption{
				mock.NewFailingSetOption(core.ParamCount),
			},
			wantError: true,
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, tc.want.activityTypeID, (query)["activityTypeId[]"])
				assert.Equal(t, tc.want.minID, query.Get("minId"))
				assert.Equal(t, tc.want.maxID, query.Get("maxId"))
				assert.Equal(t, tc.want.count, query.Get("count"))
				assert.Equal(t, tc.want.order, query.Get("order"))
				return mock.NewJSONResponse(fixture.Activity.ListJSON), nil
			}
			s := space.NewActivityService(method)

			if resp, err := s.List(context.Background(), tc.opts...); tc.wantError {
				require.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}
