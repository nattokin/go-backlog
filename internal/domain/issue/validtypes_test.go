package issue

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

	wantFilter := []option.APIParamOptionType{
		option.ParamSort,
		option.ParamOrder,
		option.ParamExpand,
	}
	wantList := []option.APIParamOptionType{
		option.ParamOffset,
		option.ParamCount,
	}
	wantUpdate := []option.APIParamOptionType{
		option.ParamStatusID,
		option.ParamResolutionID,
		option.ParamComment,
	}

	assert.Equal(t, wantFilter, filterValidTypes[len(countValidTypes):])
	assert.Equal(t, wantList, listValidTypes[len(filterValidTypes):])
	assert.Equal(t, wantUpdate, updateValidTypes[len(createValidTypes):])

	// Appending to a derived slice must not rewrite the tail of its source.
	tail := append(filterValidTypes, option.ParamCount)
	assert.Equal(t, wantList, listValidTypes[len(filterValidTypes):])
	assert.Equal(t, wantFilter, filterValidTypes[len(countValidTypes):])
	assert.Equal(t, option.ParamCount, tail[len(tail)-1])
}
