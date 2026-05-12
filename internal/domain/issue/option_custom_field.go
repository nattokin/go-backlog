package issue

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
)

// WithCustomFieldItem returns a RequestOption that adds a predefined item selection
// for list-type custom fields (Single list, Multiple list, Checkbox, Radio).
//
// The parameter name is dynamically generated as "customField_{id}[]".
// Can be called multiple times with the same id to select multiple items.
func WithCustomFieldItem(id int, itemID int) core.RequestOption {
	key := fmt.Sprintf("customField_%d[]", id)
	return &core.APIParamOption{
		Type: core.ParamCustomField,
		SetFunc: func(vals url.Values) error {
			vals.Add(key, strconv.Itoa(itemID))
			return nil
		},
	}
}

// WithCustomFieldOther returns a RequestOption that sets the free-text "Other" value
// for list-type custom fields where allowInput is enabled.
//
// The parameter name is dynamically generated as "customField_{id}_otherValue".
func WithCustomFieldOther(id int, value string) core.RequestOption {
	key := fmt.Sprintf("customField_%d_otherValue", id)
	return &core.APIParamOption{
		Type: core.ParamCustomField,
		SetFunc: func(vals url.Values) error {
			vals.Set(key, value)
			return nil
		},
	}
}
