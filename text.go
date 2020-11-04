package hrrule

import (
	"sort"
)

const weekLen = 7

type byweekday struct {
	allWeeks  []Weekday
	someWeeks []Weekday
	// isWeekdays work days
	isWeekdays bool
	isEveryDay bool
}

type text struct {
	rule ROption
	// TODO: i18n
	// TODO: dateFormatter
	lang       string
	bymonthday []int
	byweekday  byweekday
	text       []string
}

func newText(rule ROption, lang string) text {

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

	return text{
		rule:       rule,
		lang:       lang,
		bymonthday: bymonthday,
		byweekday:  weekdays,
	}
}

func (t text) String() string {
	// FIXME
	return ""
}
