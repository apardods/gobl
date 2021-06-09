package tax

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
)

// Code represents a string used to uniquely identify the data we're looking
// at. We use "code" instead of "id", to reenforce the fact that codes should
// be more easily set and used by humans than IDs or UUIDs.
type Code string

// Region defines the holding structure for a regions categories and subsequent
// Rates and Values.
type Region struct {
	Code Code        `json:"code" jsonschema:"title=Code"`
	Name i18n.String `json:"name" jsonschema:"title=Name"`

	Categories []Category `json:"categories"`
}

// Category
type Category struct {
	Code Code        `json:"code" jsonschema:"title=Code"`
	Name i18n.String `json:"name" jsonschema:"title=Name"`
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`

	// Retained when true implies that the tax amount will be retained
	// by the buyer on behalf of the supplier, and thus subtracted from
	// the invoice taxable base total.
	Retained bool `json:"retained,omitempty"`

	// Rates array
	Defs []Def
}

// Def defines a tax combination of category and rate.
type Def struct {
	// Code identifies this rate within the system
	Code Code `json:"code" jsonschema:"title=Code"`

	Name i18n.String `json:"name" jsonschema:"title=Name"`
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`

	// Values contains a list of Value objects that contain the
	// current and historical percentage values for the rate.
	// Order is important, newer values should come before
	// older values.
	Values []Value `json:"values" jsonschema:"title=Values"`
}

// Value contains a percentage rate or fixed amount for a given date range.
// Fiscal policy changes mean that rates are not static so we need to
// be able to apply the correct rate for a given period.
type Value struct {
	Since    org.Date       `json:"since,omitempty"`
	Percent  num.Percentage `json:"percent"`
	Disabled bool           `json:"disabled,omitempty"`
}

// Validate enures the basic region definition is valid.
func (r Region) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
	)
	return err
}

// Validate checks that our tax definition is valid. This is only really
// meant to be used when testing new regional tax definitions.
func (d Def) Validate() error {
	err := validation.ValidateStruct(&d,
		validation.Field(&d.Code, validation.Required),
		validation.Field(&d.Name, validation.Required),
		validation.Field(&d.Values, validation.Required, validation.By(checkDefValuesOrder)),
	)
	return err
}

// Validate ensures the tax rate contains all the required fields.
func (v Value) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Percent, validation.Required),
	)
}

func checkDefValuesOrder(list interface{}) error {
	values, ok := list.([]Value)
	if !ok {
		return errors.New("must be a tax rate value array")
	}
	var date org.Date
	// loop through and check order of Since value
	for i := range values {
		v := &values[i]
		if date.IsValid() {
			if v.Since.IsValid() && !v.Since.Before(date.Date) {
				return errors.New("invalid date order")
			}
		}
		date = v.Since
	}
	return nil
}

// Category provides the requested category by its code.
func (r Region) Category(code Code) (Category, bool) {
	for _, c := range r.Categories {
		if c.Code == code {
			return c, true
		}
	}
	return Category{}, false
}

// Def provides the rate definition with a matching code for
// the category.
func (c Category) Def(code Code) (Def, bool) {
	for _, d := range c.Defs {
		if d.Code == code {
			return d, true
		}
	}
	return Def{}, false
}

// On determines the tax rate value for the provided date.
func (d Def) On(date org.Date) (Value, bool) {
	for _, v := range d.Values {
		if !v.Since.IsValid() || v.Since.Before(date.Date) {
			return v, true
		}
	}
	return Value{}, false
}
