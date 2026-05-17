package core

import (
	"context"
	"iter"
)

// AllSeq returns an iter.Seq2 that drives offset-based pagination over any list endpoint.
// fetch must accept (ctx, offset) and return a page of results.
// Iteration stops when the returned page is shorter than perPage, signalling the last page.
func AllSeq[T any](
	ctx context.Context,
	perPage int,
	fetch func(ctx context.Context, offset int) ([]*T, error),
) iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		offset := 0
		for {
			items, err := fetch(ctx, offset)
			if err != nil {
				yield(nil, err)
				return
			}
			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}
			if len(items) < perPage {
				return
			}
			offset += len(items)
		}
	}
}
