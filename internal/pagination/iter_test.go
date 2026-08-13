package pagination_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/pagination"
)

func TestAll(t *testing.T) {
	t.Run("single-short-page", func(t *testing.T) {
		t.Parallel()

		calls := 0
		fetch := func(ctx context.Context, offset int) ([]*int, error) {
			calls++
			assert.Equal(t, 0, offset)
			one, two := 1, 2
			return []*int{&one, &two}, nil
		}

		var got []int
		for v, err := range pagination.All(context.Background(), 10, fetch) {
			require.NoError(t, err)
			got = append(got, *v)
		}

		assert.Equal(t, []int{1, 2}, got)
		assert.Equal(t, 1, calls, "a page shorter than perPage must stop iteration after the first fetch")
	})

	t.Run("multiple-pages", func(t *testing.T) {
		t.Parallel()

		const perPage = 2
		pages := [][]int{{1, 2}, {3, 4}, {5}}
		var offsetsSeen []int

		fetch := func(ctx context.Context, offset int) ([]*int, error) {
			offsetsSeen = append(offsetsSeen, offset)
			page := pages[0]
			pages = pages[1:]
			out := make([]*int, len(page))
			for i, v := range page {
				v := v
				out[i] = &v
			}
			return out, nil
		}

		var got []int
		for v, err := range pagination.All(context.Background(), perPage, fetch) {
			require.NoError(t, err)
			got = append(got, *v)
		}

		assert.Equal(t, []int{1, 2, 3, 4, 5}, got)
		assert.Equal(t, []int{0, 2, 4}, offsetsSeen, "offset must advance by the number of items returned on each page")
	})

	t.Run("fetch-error", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("fetch failed")
		calls := 0
		fetch := func(ctx context.Context, offset int) ([]*int, error) {
			calls++
			return nil, wantErr
		}

		var gotErr error
		var got []int
		for v, err := range pagination.All(context.Background(), 10, fetch) {
			if err != nil {
				gotErr = err
				continue
			}
			got = append(got, *v)
		}

		assert.Equal(t, wantErr, gotErr)
		assert.Empty(t, got)
		assert.Equal(t, 1, calls, "iteration must stop after a fetch error")
	})

	t.Run("early-stop", func(t *testing.T) {
		t.Parallel()

		const perPage = 2
		pages := [][]int{{1, 2}, {3, 4}, {5}}
		calls := 0

		fetch := func(ctx context.Context, offset int) ([]*int, error) {
			calls++
			page := pages[0]
			pages = pages[1:]
			out := make([]*int, len(page))
			for i, v := range page {
				v := v
				out[i] = &v
			}
			return out, nil
		}

		var got []int
		for v, err := range pagination.All(context.Background(), perPage, fetch) {
			require.NoError(t, err)
			got = append(got, *v)
			if len(got) == 1 {
				break
			}
		}

		assert.Equal(t, []int{1}, got)
		assert.Equal(t, 1, calls, "breaking out of the range loop must stop further fetches")
	})
}
