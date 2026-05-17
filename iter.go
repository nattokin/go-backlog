package backlog

import (
	"context"
	"iter"

	"github.com/nattokin/go-backlog/internal/core"
)

// allSeq returns an iter.Seq2 that drives offset-based pagination.
// It delegates to core.AllSeq and converts each element via the provided convert function.
func allSeq[M, T any](
	ctx context.Context,
	perPage int,
	fetch func(ctx context.Context, offset int) ([]*M, error),
	convert func(*M) *T,
) iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		for v, err := range core.AllSeq(ctx, perPage, fetch) {
			if !yield(convert(v), err) {
				return
			}
		}
	}
}
