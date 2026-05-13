package issue

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/nattokin/go-backlog/internal/core"
)

// WithCustomField returns a RequestOption that sets a custom field value for
// non-list types (Text, Sentence, Number, Date).
//
// The parameter name is dynamically generated as "customField_{id}".
// Supported value types: string, float64, time.Time.
// time.Time values are formatted as "yyyy-MM-dd".
func WithCustomField[T string | float64 | time.Time](id int, value T) core.RequestOption {
	return &core.APIParamOption{
		Type: core.ParamCustomField,
		CheckFunc: func() error {
			if err := validateCustomFieldID(id); err != nil {
				return err
			}

			name := core.ParamCustomField.Value()
			switch v := any(value).(type) {
			case string:
				if v == "" {
					return core.NewValidationError(fmt.Sprintf("%s value must not be empty", name))
				}
			case time.Time:
				if v.IsZero() {
					return core.NewValidationError(fmt.Sprintf("%s date must not be zero value", name))
				}
			}

			return nil
		},
		SetFunc: func(vals url.Values) error {
			key := fmt.Sprintf("customField_%d", id)
			var serialized string
			switch v := any(value).(type) {
			case string:
				serialized = v
			case float64:
				serialized = strconv.FormatFloat(v, 'f', -1, 64)
			case time.Time:
				serialized = v.Format("2006-01-02")
			}

			vals.Set(key, serialized)
			return nil
		},
	}
}

// WithCustomFieldItem returns a RequestOption that adds a predefined item selection
// for list-type custom fields (Single list, Multiple list, Checkbox, Radio).
//
// The parameter name is dynamically generated as "customField_{id}".
func WithCustomFieldItems(id int, itemIDs []int) core.RequestOption {
	return &core.APIParamOption{
		Type:      core.ParamCustomField,
		CheckFunc: checkCustomFieldFunc(id),
		SetFunc: func(vals url.Values) error {
			key := fmt.Sprintf("customField_%d", id)
			for _, itemID := range itemIDs {
				vals.Add(key, strconv.Itoa(itemID))
			}
			return nil
		},
	}
}

// WithCustomFieldOther returns a RequestOption that sets the free-text "Other" value
// for list-type custom fields where allowInput is enabled.
//
// The parameter name is dynamically generated as "customField_{id}_otherValue".
func WithCustomFieldOther(id int, value string) core.RequestOption {
	return &core.APIParamOption{
		Type:      core.ParamCustomField,
		CheckFunc: checkCustomFieldFunc(id),
		SetFunc: func(vals url.Values) error {
			key := fmt.Sprintf("customField_%d_otherValue", id)
			vals.Set(key, value)
			return nil
		},
	}
}

func checkCustomFieldFunc(id int) func() error {
	return func() error {
		return validateCustomFieldID(id)
	}
}

func validateCustomFieldID(id int) error {
	if id < 1 {
		msg := fmt.Sprintf("invalid %s: must not be less than 1", core.ParamCustomField.Value())
		return core.NewValidationError(msg)
	}
	return nil
}
