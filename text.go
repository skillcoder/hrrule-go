package hrrule

import (
	"sort"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const weekLen = 7
const initialTextLen = 64
const spaceByte = ' '

type byweekday struct {
	allWeeks  []Weekday
	someWeeks []Weekday
	// isWeekdays work days
	isWeekdays bool
	isEveryDay bool
}

type text struct {
	rule ROption
	loc  *i18n.Localizer
	// TODO: dateFormatter
	bymonthday []int
	byweekday  byweekday
	text       strings.Builder
}

func newText(rule ROption, localizer *i18n.Localizer) *text {
	bymonthday := make([]int, 0, len(rule.Bymonthday))
	if len(rule.Bymonthday) != 0 {
		// 1, 2, 3, ... , -5, -4, -3, ...
		neg := make([]int, 0, len(rule.Bymonthday))
		pos := make([]int, 0, len(rule.Bymonthday))
		for _, x := range rule.Bymonthday {
			if x >= 0 {
				pos = append(pos, x)
				continue
			}
			neg = append(neg, x)
		}

		sort.Ints(pos)
		sort.Ints(neg)

		bymonthday = append(bymonthday, pos...)
		bymonthday = append(bymonthday, neg...)
	}

	var weekdays byweekday
	if len(rule.Byweekday) != 0 {
		weekdays.isWeekdays = true
		everyDay := make(map[Weekday]struct{}, weekLen)
		cnt := 0
		for _, w := range rule.Byweekday {
			if _, ok := everyDay[w]; !ok {
				cnt++
			}

			if w == SU || w == SA {
				weekdays.isWeekdays = false
			}

			if w.N() == 0 {
				weekdays.allWeeks = append(weekdays.allWeeks, w)
				continue
			}

			weekdays.someWeeks = append(weekdays.someWeeks, w)
		}

		if cnt == weekLen {
			weekdays.isEveryDay = true
		}

		sort.Slice(weekdays.allWeeks, func(i, j int) bool { return weekdays.allWeeks[i].Day() < weekdays.allWeeks[j].Day() })
		sort.Slice(weekdays.someWeeks, func(i, j int) bool { return weekdays.someWeeks[i].Day() < weekdays.someWeeks[j].Day() })
	}

	return &text{
		rule:       rule,
		loc:        localizer,
		bymonthday: bymonthday,
		byweekday:  weekdays,
	}
}

func (t *text) String() string {
	t.text.Reset()
	t.text.Grow(initialTextLen)

	t.text.WriteString(t.loc.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
		ID: "every",
		Other: "every",
	}}))

	switch t.rule.Freq {
	case DAILY:

	}

	return t.text.String()
}

func (t *text) add(s string) {
	t.text.WriteByte(spaceByte)
	t.text.WriteString(s)
}

// TODO: change to i18n
func (t text) getText(s string) string {
	return s
}
