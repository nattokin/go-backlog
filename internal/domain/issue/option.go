package issue

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/validate"
)

// WithCustomField returns a *core.APIParamOption that sets a custom field value for
// non-list types (Text, Sentence, Number, Date).
func WithCustomField[T string | float64 | time.Time](id int, value T) *core.APIParamOption {
	return &core.APIParamOption{
		Type: core.ParamCustomField,
		CheckFunc: func() *core.ValidationError {
			if ve := validate.ValidateCustomFieldID(id); ve != nil {
				return ve
			}

			name := core.ParamCustomField.Value()
			switch v := any(value).(type) {
			case string:
				if v == "" {
					return core.NewValidationError(name, fmt.Sprintf("%s value must not be empty", name))
				}
			case time.Time:
				if v.IsZero() {
					return core.NewValidationError(name, fmt.Sprintf("%s date must not be zero value", name))
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

// WithCustomFieldItems returns a *core.APIParamOption that sets predefined item selections
// for list-type custom fields.
func WithCustomFieldItems(id int, itemIDs []int) *core.APIParamOption {
	return &core.APIParamOption{
		Type: core.ParamCustomField,
		CheckFunc: func() *core.ValidationError {
			if ve := validate.ValidateCustomFieldID(id); ve != nil {
				return ve
			}
			return validateItemIDs(itemIDs)
		},
		SetFunc: func(vals url.Values) error {
			key := fmt.Sprintf("customField_%d", id)
			for _, itemID := range itemIDs {
				vals.Add(key, strconv.Itoa(itemID))
			}
			return nil
		},
	}
}

// WithCustomFieldOther returns a *core.APIParamOption that sets the free-text "Other"
// value for list-type custom fields where allowInput is enabled.
func WithCustomFieldOther(id int, value string) *core.APIParamOption {
	return &core.APIParamOption{
		Type: core.ParamCustomField,
		CheckFunc: func() *core.ValidationError {
			return validate.ValidateCustomFieldID(id)
		},
		SetFunc: func(vals url.Values) error {
			key := fmt.Sprintf("customField_%d_otherValue", id)
			vals.Set(key, value)
			return nil
		},
	}
}

func validateItemIDs(ids []int) *core.ValidationError {
	for _, id := range ids {
		if id < 1 {
			return core.NewValidationError("customField_itemID", fmt.Sprintf("customField itemID must not be less than 1, got %d", id))
		}
	}
	return nil
}
