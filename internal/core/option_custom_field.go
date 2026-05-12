package core

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// CustomFieldOption is a marker interface implemented by options that set
// custom field parameters (customField_{id}, customField_{id}[], customField_{id}_otherValue).
// These options use dynamically generated keys and are applied separately from
// the static ValidateOption check in ApplyOptions.
type CustomFieldOption interface {
	RequestOption
	isCustomFieldOption()
}

// customFieldOption is the concrete implementation of CustomFieldOption.
type customFieldOption struct {
	key     string
	setFunc func(url.Values) error
}

func (o *customFieldOption) Key() string            { return o.key }
func (o *customFieldOption) Check() error           { return nil }
func (o *customFieldOption) Set(v url.Values) error { return o.setFunc(v) }
func (o *customFieldOption) isCustomFieldOption()   {}

// SplitCustomFieldOptions separates custom field options from regular options.
// Custom field options are returned separately so callers can apply them without
// the static key validation performed by ApplyOptions.
func SplitCustomFieldOptions(opts []RequestOption) (regular []RequestOption, custom []CustomFieldOption) {
	for _, opt := range opts {
		if cf, ok := opt.(CustomFieldOption); ok {
			custom = append(custom, cf)
		} else {
			regular = append(regular, opt)
		}
	}
	return regular, custom
}

// ApplyCustomFieldOptions applies custom field options directly to url.Values
// without any key validation.
func ApplyCustomFieldOptions(v url.Values, opts []CustomFieldOption) error {
	for _, opt := range opts {
		if err := opt.Check(); err != nil {
			return err
		}
		if err := opt.Set(v); err != nil {
			return err
		}
	}
	return nil
}

// WithCustomField returns a RequestOption that sets a custom field value for
// non-list types (Text, Sentence, Number, Date).
//
// The parameter name is dynamically generated as "customField_{id}".
// Supported value types: string, int, float64, time.Time.
// time.Time values are formatted as "yyyy-MM-dd".
func WithCustomField[T string | int | float64 | time.Time](id int, value T) RequestOption {
	key := fmt.Sprintf("customField_%d", id)
	var serialized string
	switch v := any(value).(type) {
	case string:
		serialized = v
	case int:
		serialized = strconv.Itoa(v)
	case float64:
		serialized = strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		serialized = v.Format("2006-01-02")
	}
	return &customFieldOption{
		key: key,
		setFunc: func(vals url.Values) error {
			vals.Set(key, serialized)
			return nil
		},
	}
}

// WithCustomFieldItem returns a RequestOption that adds a predefined item selection
// for list-type custom fields (Single list, Multiple list, Checkbox, Radio).
//
// The parameter name is dynamically generated as "customField_{id}[]".
// Can be called multiple times with the same id to select multiple items.
func WithCustomFieldItem(id int, itemID int) RequestOption {
	key := fmt.Sprintf("customField_%d[]", id)
	return &customFieldOption{
		key: key,
		setFunc: func(vals url.Values) error {
			vals.Add(key, strconv.Itoa(itemID))
			return nil
		},
	}
}

// WithCustomFieldOther returns a RequestOption that sets the free-text "Other" value
// for list-type custom fields where allowInput is enabled.
//
// The parameter name is dynamically generated as "customField_{id}_otherValue".
func WithCustomFieldOther(id int, value string) RequestOption {
	key := fmt.Sprintf("customField_%d_otherValue", id)
	return &customFieldOption{
		key: key,
		setFunc: func(vals url.Values) error {
			vals.Set(key, value)
			return nil
		},
	}
}
