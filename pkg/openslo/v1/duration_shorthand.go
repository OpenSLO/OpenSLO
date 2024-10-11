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
	n, err := fmt.Sscanf(string(text), "%d%s", &d.value, &d.unit)
	if err != nil {
		return err
	}
	if n != 2 {
		return fmt.Errorf("invalid duration shorthand: %s, expected [0-9]+[mhdwMQY]", text)
	}
	return nil
}

// UnmarshalText implements [encoding.TextMarshaler].
func (d DurationShorthand) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

// String implements [fmt.Stringer].
func (d DurationShorthand) String() string {
	return fmt.Sprintf("%d%s", d.value, d.unit)
}

// Duration returns the [time.Duration] representation of [DurationShorthand].
func (d DurationShorthand) Duration() time.Duration {
	switch d.unit {
	case "m":
		return time.Duration(d.value) * time.Minute
	case "h":
		return time.Duration(d.value) * time.Hour
	case "d":
		return time.Duration(d.value) * 24 * time.Hour
	case "w":
		return time.Duration(d.value) * 7 * 24 * time.Hour
	case "M":
		return time.Duration(d.value) * 30 * 24 * time.Hour
	case "Q":
		return time.Duration(d.value) * 90 * 24 * time.Hour
	case "Y":
		return time.Duration(d.value) * 365 * 24 * time.Hour
	default:
		panic("invalid unit")
	}
}

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

var validDiuartionUnits = []DurationShorthandUnit{"m", "h", "d", "w", "M", "Q", "Y"}

// Validate checks if [DurationShorthand] is correct.
func (d DurationShorthand) Validate() error {
	return durationShortHandValidation.Validate(d)
}

var durationShortHandValidation = govy.New(
	govy.For(func(d DurationShorthand) DurationShorthandUnit { return d.unit }).
		Required().
		Rules(rules.OneOf(validDiuartionUnits...)),
	govy.For(func(d DurationShorthand) int { return d.value }).
		Rules(rules.GTE(0)),
)
