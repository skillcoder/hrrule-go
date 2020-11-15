package hrrule

import (
	"fmt"

	"github.com/teambition/rrule-go"
)

// StrToROption parse rrule string to ROption
func StrToROption(rRuleStr string) (*ROption, error) {
	rOption, err := rrule.StrToROption(rRuleStr)
	if err != nil {
		return nil, fmt.Errorf("str to rrule: %w", err)
	}
	return FromROption(rOption), nil
}

// FromRRule covert teambition/rrule-go to ROption
func FromROption(o *rrule.ROption) *ROption {
	return &ROption{
		Freq:     Frequency(o.Freq),
		Dtstart:  o.Dtstart,
		Interval: o.Interval,
		Wkst: Weekday{
			weekday: o.Wkst.Day(),
			n:       o.Wkst.N(),
		},
		Count:      o.Count,
		Until:      o.Until,
		Bysetpos:   o.Bysetpos,
		Bymonth:    o.Bymonth,
		Bymonthday: o.Bymonthday,
		Byyearday:  o.Byyearday,
		Byweekno:   o.Byweekno,
		Byweekday:  FromWeekdaySlice(o.Byweekday),
		Byhour:     o.Byhour,
		Byminute:   o.Byminute,
		Bysecond:   o.Bysecond,
		Byeaster:   o.Byeaster,
	}
}

// FromWeekdaySlice covert from teambition/rrule-go slice Weekday slice
func FromWeekdaySlice(weekdays []rrule.Weekday) []Weekday {
	wd := make([]Weekday, 0, len(weekdays))
	for i := range weekdays {
		wd = append(wd, Weekday{
			weekday: weekdays[i].Day(),
			n:       weekdays[i].N(),
		})
	}

	return wd
}
