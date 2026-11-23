package pullrequest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/option"
)

// TestValidTypesOwnBackingArray guards the derived validTypes slices against
// append sharing a backing array: deriving a slice from another one with spare
// capacity overwrites elements of the source slice.
func TestValidTypesOwnBackingArray(t *testing.T) {
	t.Parallel()

	wantList := []option.APIParamOptionType{
		option.ParamOffset,
		option.ParamCount,
	}

	assert.Equal(t, wantList, listValidTypes[len(filterValidTypes):])

	// Appending to the derived slice must not rewrite the tail of its source.
	tail := append(listValidTypes, option.ParamExpand)
	assert.Equal(t, wantList, listValidTypes[len(filterValidTypes):])
	assert.Equal(t, option.ParamExpand, tail[len(tail)-1])
}
