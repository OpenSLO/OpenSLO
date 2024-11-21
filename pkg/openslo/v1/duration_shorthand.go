package v1

import (
	"fmt"
	"time"

	"github.com/nobl9/govy/pkg/govy"
	"github.com/nobl9/govy/pkg/rules"
)

// ParseDurationShorthand parses a string representation of [DurationShorthand].
func ParseDurationShorthand(s string) (DurationShorthand, error) {
	d := new(DurationShorthand)
	err := d.UnmarshalText([]byte(s))
	return *d, err
}

// NewDurationShorthand creates a new [DurationShorthand] instance.
func NewDurationShorthand(value int, unit DurationShorthandUnit) DurationShorthand {
	return DurationShorthand{
		unit:  unit,
		value: value,
	}
}

// DurationShorthand is a shorthand representation of time duration.
// It consists of a value and unit, e.g. '1m' (1 minute), '10d' (10 days).
type DurationShorthand struct {
	unit  DurationShorthandUnit
	value int
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (d *DurationShorthand) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}
	if n, err := fmt.Sscanf(string(text), "%d%s", &d.value, &d.unit); err != nil || n != 2 {
		return fmt.Errorf("invalid duration shorthand: %s, expected [0-9]+[mhdwMQY]", text)
	}
	return nil
}

// MarshalText implements [encoding.TextMarshaler].
func (d DurationShorthand) MarshalText() ([]byte, error) {
	if d.value == 0 {
		return []byte{}, nil
	}
	return []byte(d.String()), nil
}

// String implements [fmt.Stringer].
func (d DurationShorthand) String() string {
	if d.value == 0 {
		return ""
	}
	return fmt.Sprintf("%d%s", d.value, d.unit)
}

// Duration returns the [time.Duration] representation of [DurationShorthand].
func (d DurationShorthand) Duration() time.Duration {
	switch d.unit {
	case DurationShorthandUnitMinute:
		return time.Duration(d.value) * time.Minute
	case DurationShorthandUnitHour:
		return time.Duration(d.value) * time.Hour
	case DurationShorthandUnitDay:
		return time.Duration(d.value) * 24 * time.Hour
	case DurationShorthandUnitWeek:
		return time.Duration(d.value) * 7 * 24 * time.Hour
	case DurationShorthandUnitMonth:
		return time.Duration(d.value) * 30 * 24 * time.Hour
	case DurationShorthandUnitQuarter:
		return time.Duration(d.value) * 90 * 24 * time.Hour
	case DurationShorthandUnitYear:
		return time.Duration(d.value) * 365 * 24 * time.Hour
	default:
		panic("invalid unit")
	}
}

// DurationShorthandUnit is a unit of [DurationShorthand].
type DurationShorthandUnit string

const (
	DurationShorthandUnitMinute  DurationShorthandUnit = "m"
	DurationShorthandUnitHour    DurationShorthandUnit = "h"
	DurationShorthandUnitDay     DurationShorthandUnit = "d"
	DurationShorthandUnitWeek    DurationShorthandUnit = "w"
	DurationShorthandUnitMonth   DurationShorthandUnit = "M"
	DurationShorthandUnitQuarter DurationShorthandUnit = "Q"
	DurationShorthandUnitYear    DurationShorthandUnit = "Y"
)

var validDurationUnits = []DurationShorthandUnit{
	DurationShorthandUnitMinute,
	DurationShorthandUnitHour,
	DurationShorthandUnitDay,
	DurationShorthandUnitWeek,
	DurationShorthandUnitMonth,
	DurationShorthandUnitQuarter,
	DurationShorthandUnitYear,
}

// Validate checks if [DurationShorthand] is correct.
func (d DurationShorthand) Validate() error {
	return durationShortHandValidation.Validate(d)
}

var durationShortHandValidation = govy.New(
	govy.For(func(d DurationShorthand) DurationShorthandUnit { return d.unit }).
		WithName("unit").
		Required().
		Rules(rules.OneOf(validDurationUnits...)),
	govy.For(func(d DurationShorthand) int { return d.value }).
		WithName("value").
		Rules(rules.GTE(0)),
)
