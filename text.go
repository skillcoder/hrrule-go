package hrrule

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	weekLen        = 7
	workLen        = 5
	initialTextLen = 64
	spaceByte      = ' '
	comaSpace      = ", "
)

var implementedFreq = map[Frequency]struct{}{
	DAILY:   {},
	WEEKLY:  {},
	MONTHLY: {},
	YEARLY:  {},
}

type ListMode int

const (
	IsINT ListMode = iota
	IsNTH
	IsMONTH
)

var (
	langAnd = &i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
		ID:          "And",
		Description: "Used for final delimiter in list",
		Other:       "and",
	}}

	langOr = &i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
		ID:          "Or",
		Description: "Used for final delimiter in list",
		Other:       "or",
	}}

	langOnThe = &i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
		ID:          "OnThe",
		Description: "Used before list bymonthday, byweekday, byyearday and someWeeks",
		Other:       "on the",
	}}
)

type byweekday struct {
	allWeeks  []Weekday
	someWeeks []Weekday
	// isWeekdays work days
	isWeekdays bool
	isEveryDay bool
}

type text struct {
	rule          *ROption
	loc           *i18n.Localizer
	dateFormatter DateFormatter
	bymonthday    []int
	byweekday     *byweekday
	text          strings.Builder
}

func newText(rule *ROption, localizer *i18n.Localizer, dateFormatter DateFormatter) *text {
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

	var weekdays *byweekday
	if len(rule.Byweekday) != 0 {
		weekdays = &byweekday{
			isWeekdays: true,
		}
		everyDay := make(map[Weekday]struct{}, weekLen)
		workDay := make(map[Weekday]struct{}, workLen)
		cntEvery := 0
		cntWork := 0
		for _, w := range rule.Byweekday {
			if _, ok := everyDay[w]; !ok {
				cntEvery++
			}

			if _, ok := workDay[w]; !ok {
				cntWork++
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

		if cntEvery == weekLen {
			weekdays.isEveryDay = true
		}

		if cntWork != workLen {
			weekdays.isWeekdays = false
		}

		sort.Slice(weekdays.allWeeks, func(i, j int) bool { return weekdays.allWeeks[i].Day() < weekdays.allWeeks[j].Day() })
		sort.Slice(weekdays.someWeeks, func(i, j int) bool { return weekdays.someWeeks[i].Day() < weekdays.someWeeks[j].Day() })
	}

	return &text{
		rule:          rule,
		loc:           localizer,
		dateFormatter: dateFormatter,
		bymonthday:    bymonthday,
		byweekday:     weekdays,
	}
}

func (t *text) String() string {
	t.text.Reset()
	t.text.Grow(initialTextLen)

	localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Every",
			Other: "every",
		}})
	t.text.WriteString(localizedString)

	switch t.rule.Freq {
	case DAILY:
		t.daily()
	case WEEKLY:
		t.weekly()
	case MONTHLY:
		t.monthly()
	case YEARLY:
		t.yearly()
	}

	if !t.rule.Until.IsZero() {
		localizedString, _ = t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Until",
				Other: "until",
			}})
		t.add(localizedString)

		until := t.rule.Until
		t.add(t.dateFormatter.Format(until.Year(), until.Month(), until.Day()))
	} else if t.rule.Count != 0 {
		localizedString, _ = t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "TimeCount",
				One:   "for {{.Count}} time",
				Two:   "for {{.Count}} times",
				Few:   "for {{.Count}} times",
				Many:  "for {{.Count}} times",
				Other: "for {{.Count}} times",
			},
			TemplateData: map[string]interface{}{
				"Count": t.rule.Count,
			},
			PluralCount: t.rule.Count,
		})
		t.add(localizedString)
	}

	if !t.isFullyConvertible() {
		localizedString, _ = t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Approximately",
				Other: "(~ approximate)",
			}})
		t.add(localizedString)
	}

	return t.text.String()
}

func (t *text) daily() {
	if t.rule.Interval != 1 {
		t.add(strconv.Itoa(t.rule.Interval))
	}

	if t.byweekday != nil && t.byweekday.isWeekdays {
		localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "IntervalCountWeekday",
				Description: "Used after interval count",
				One:         "weekday",
				Two:         "weekdays",
				Few:         "weekdays",
				Many:        "weekdays",
				Other:       "weekdays",
			},
			PluralCount: t.rule.Interval,
		})
		t.add(localizedString)
	} else {
		localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "IntervalCountDay",
				Description: "Used after interval count",
				One:         "day",
				Two:         "days",
				Few:         "days",
				Many:        "days",
				Other:       "days",
			},
			PluralCount: t.rule.Interval,
		})
		t.add(localizedString)
	}

	if len(t.rule.Bymonth) != 0 {
		localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "In",
				Description: "Used before by month count",
				Other:       "in",
			}})
		t.add(localizedString)

		t.addByMonth()
	}

	if len(t.bymonthday) != 0 {
		t.addByMonthday()
	} else if t.byweekday != nil {
		t.addByWeekday()
	} else if len(t.rule.Byhour) != 0 {
		t.addByHour()
	}
}

func (t *text) weekly() {
	if t.rule.Interval != 1 {
		t.add(strconv.Itoa(t.rule.Interval))
		localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "IntervalCountWeeks",
				Description: "Used after interval count",
				One:         "week",
				Two:         "weeks",
				Few:         "weeks",
				Many:        "weeks",
				Other:       "weeks",
			},
			PluralCount: t.rule.Interval,
		})
		t.add(localizedString)
	}

	if t.byweekday != nil && t.byweekday.isWeekdays {
		if t.rule.Interval == 1 {
			localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:          "IntervalCountWeekdays",
					Description: "Used before weekdays",
					One:         "weekday",
					Two:         "weekdays",
					Few:         "weekdays",
					Many:        "weekdays",
					Other:       "weekdays",
				},
				PluralCount: t.rule.Interval,
			})
			t.add(localizedString)
		} else {
			localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:          "OnWeekdays",
					Description: "Used before weekdays list",
					Other:       "on weekdays",
				}})
			t.add(localizedString)
		}
	} else if t.byweekday != nil && t.byweekday.isEveryDay {
		localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "IntervalCountDays",
				Description: "Used after interval count if every day",
				One:         "day",
				Two:         "days",
				Few:         "days",
				Many:        "days",
				Other:       "days",
			},
			PluralCount: t.rule.Interval,
		})
		t.add(localizedString)
	} else {
		if t.rule.Interval == 1 {
			localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:          "Week",
					Description: "Used after interval count",
					Other:       "week",
				}})
			t.add(localizedString)
		}

		if len(t.rule.Bymonth) != 0 {
			localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:          "In",
					Description: "Used before by month count",
					Other:       "in",
				}})
			t.add(localizedString)

			t.addByMonth()
		}

		if len(t.bymonthday) != 0 {
			t.addByMonthday()
		} else if t.byweekday != nil {
			t.addByWeekday()
		}
	}
}

func (t *text) monthly() {
	if len(t.rule.Bymonth) != 0 {
		if t.rule.Interval != 1 {
			localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:          "IntervalCountMonths",
					Description: "Used for months by month",
					One:         "{{.Interval}} month",
					Two:         "{{.Interval}} months in",
					Few:         "{{.Interval}} months in",
					Many:        "{{.Interval}} months in",
					Other:       "{{.Interval}} months in",
				},
				TemplateData: map[string]interface{}{
					"Interval": t.rule.Interval,
				},
				PluralCount: t.rule.Interval,
			})
			t.add(localizedString)
		}

		t.addByMonth()
	} else {
		if t.rule.Interval != 1 {
			t.add(strconv.Itoa(t.rule.Interval))
		}

		localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "IntervalMonths",
				Description: "Used after interval monthly",
				One:         "month",
				Two:         "months",
				Few:         "months",
				Many:        "months",
				Other:       "months",
			},
			PluralCount: t.rule.Interval,
		})
		t.add(localizedString)
	}

	if len(t.bymonthday) != 0 {
		t.addByMonthday()
	} else if t.byweekday != nil && t.byweekday.isWeekdays {
		localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "OnWeekdays",
				Description: "Used if on every work weekdays on week",
				Other:       "on weekdays",
			}})
		t.add(localizedString)
	} else if t.byweekday != nil {
		t.addByWeekday()
	}
}

func (t *text) yearly() {
	if len(t.rule.Bymonth) != 0 {
		if t.rule.Interval != 1 {
			localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:          "IntervalYearlyByMonth",
					Description: "Used for years by month",
					One:         "{{.Interval}} year",
					Two:         "{{.Interval}} years",
					Few:         "{{.Interval}} years",
					Many:        "{{.Interval}} years",
					Other:       "{{.Interval}} years",
				},
				TemplateData: map[string]interface{}{
					"Interval": t.rule.Interval,
				},
				PluralCount: t.rule.Interval,
			})
			t.add(localizedString)
		}

		t.addByMonth()
	} else {
		if t.rule.Interval != 1 {
			t.add(strconv.Itoa(t.rule.Interval))
		}

		localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "IntervalYears",
				Description: "Used after interval yearly",
				One:         "year",
				Two:         "years",
				Few:         "years",
				Many:        "years",
				Other:       "years",
			},
			PluralCount: t.rule.Interval,
		})
		t.add(localizedString)
	}

	if len(t.bymonthday) != 0 {
		t.addByMonthday()
	} else if t.byweekday != nil {
		t.addByWeekday()
	}

	if len(t.rule.Byyearday) != 0 {
		localizedString, _ := t.loc.Localize(langOnThe)
		t.add(localizedString)

		localizedString, _ = t.loc.Localize(langAnd)
		t.add(t.list(t.rule.Byyearday, IsNTH, localizedString))

		localizedString, _ = t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "ByYearDay",
				Description: "Used after by year day",
				Other:       "day",
			}})
		t.add(localizedString)
	}

	if len(t.rule.Byweekno) != 0 {
		localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "InWeekNo",
				Description: "Used before by week no",
				Other:       "in",
			}})
		t.add(localizedString)

		localizedString, _ = t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "CountWeekNo",
				Description: "Used before WeekNo list",
				One:         "week",
				Two:         "weeks",
				Few:         "weeks",
				Many:        "weeks",
				Other:       "weeks",
			},
			PluralCount: len(t.rule.Byweekno),
		})
		t.add(localizedString)

		localizedString, _ = t.loc.Localize(langOnThe)
		t.add(t.list(t.rule.Byweekno, IsINT, localizedString))
	}
}

func (t *text) addByMonth() {
	localizedString, _ := t.loc.Localize(langAnd)
	t.add(t.list(t.rule.Bymonth, IsMONTH, localizedString))
}

func (t *text) addByMonthday() {
	if t.byweekday != nil && len(t.byweekday.allWeeks) != 0 {
		localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "OnAllWeeks",
				Description: "Used before all weeks list",
				Other:       "on",
			}})
		t.add(localizedString)

		localizedString, _ = t.loc.Localize(langOr)
		t.add(t.listWeekday(t.byweekday.allWeeks, localizedString))

		localizedString, _ = t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "TheAllWeeks",
				Description: "Used between all weeks list and bymonthday",
				Other:       "the",
			}})
		t.add(localizedString)

		localizedString, _ = t.loc.Localize(langOr)
		t.add(t.list(t.bymonthday, IsNTH, localizedString))
	} else {
		localizedString, _ := t.loc.Localize(langOnThe)
		t.add(localizedString)

		localizedString, _ = t.loc.Localize(langAnd)
		t.add(t.list(t.bymonthday, IsNTH, localizedString))
	}
}

func (t *text) addByWeekday() {
	if t.byweekday != nil && len(t.byweekday.allWeeks) != 0 && !t.byweekday.isWeekdays {
		localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "OnWeekday",
				Description: "Used before by weekday list",
				Other:       "on",
			}})
		t.add(localizedString)

		t.add(t.listWeekday(t.byweekday.allWeeks, ""))
	}
	if t.byweekday != nil && len(t.byweekday.someWeeks) != 0 {
		if len(t.byweekday.allWeeks) != 0 {
			localizedString, _ := t.loc.Localize(langAnd)
			t.add(localizedString)
		}

		localizedString, _ := t.loc.Localize(langOnThe)
		t.add(localizedString)

		localizedString, _ = t.loc.Localize(langAnd)
		t.add(t.listWeekday(t.byweekday.someWeeks, localizedString))
	}
}

func (t *text) addByHour() {
	localizedString, _ := t.loc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "At",
			Description: "Used before hour list",
			Other:       "at",
		}})
	t.add(localizedString)

	localizedString, _ = t.loc.Localize(langAnd)
	t.add(t.list(t.rule.Byhour, IsINT, localizedString))
}

func (t *text) isFullyConvertible() bool {
	if _, ok := implementedFreq[t.rule.Freq]; !ok {
		return false
	}

	if !t.rule.Until.IsZero() && t.rule.Count != 0 {
		return false
	}

	// TODO: do not know how check exist t.rule.Wkst
	// TODO: do not know how check exist t.rule.Freq
	if !t.rule.Dtstart.IsZero() {
		return true
	}

	// TODO: implement checks for important values for each Freq
	//  if (!contains(t.implementedFreq[t.rule.Freq], key)) { return false }

	return true
}

func (t *text) list(arr []int, mode ListMode, finalDelimiter string) string {
	delimJoin := func(array []string, finDelimiter string) string {
		var sb strings.Builder
		sb.Grow(4*len(arr) + len(finDelimiter) + 2)
		for i := range arr {
			if i != 0 {
				if i == len(array)-1 {
					sb.WriteByte(spaceByte)
					sb.WriteString(finalDelimiter)
					sb.WriteByte(spaceByte)
				} else {
					sb.WriteString(comaSpace)
				}
			}

			sb.WriteString(array[i])
		}

		return sb.String()
	}

	if finalDelimiter != "" {
		return delimJoin(t.intToStringSlice(arr, mode), finalDelimiter)
	}

	return strings.Join(t.intToStringSlice(arr, mode), comaSpace)
}

func (t *text) listWeekday(arr []Weekday, finalDelimiter string) string {
	delimJoin := func(array []string, finDelimiter string) string {
		var sb strings.Builder
		sb.Grow(4*len(arr) + len(finDelimiter) + 2)
		for i := range arr {
			if i != 0 {
				if i == len(array)-1 {
					sb.WriteByte(spaceByte)
					sb.WriteString(finalDelimiter)
					sb.WriteByte(spaceByte)
				} else {
					sb.WriteString(comaSpace)
				}
			}

			sb.WriteString(array[i])
		}

		return sb.String()
	}

	if finalDelimiter != "" {
		return delimJoin(t.weekdayToStringSlice(arr), finalDelimiter)
	}

	return strings.Join(t.weekdayToStringSlice(arr), comaSpace)
}

func (t *text) add(s string) {
	t.text.WriteByte(spaceByte)
	t.text.WriteString(s)
}

// TODO: change to i18n
func (t text) getText(s string) string {
	return s
}

func (t text) weekdayToStringSlice(arr []Weekday) []string {
	str := make([]string, 0, len(arr))
	for i := range arr {
		str = append(str, t.dateFormatter.WeekdayName(arr[i]))
	}
	return str
}

func (t text) intToStringSlice(arr []int, mode ListMode) []string {
	str := make([]string, 0, len(arr))
	for _, i := range arr {
		str = append(str, t.intToString(i, mode))
	}
	return str
}

func (t text) intToString(i int, mode ListMode) string {
	switch mode {
	case IsINT:
		return strconv.Itoa(i)
	case IsNTH:
		return t.dateFormatter.Nth(i)
	case IsMONTH:
		return t.dateFormatter.MonthName(time.Month(i))
	}

	return UNKNOWN
}
