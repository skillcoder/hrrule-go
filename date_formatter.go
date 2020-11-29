package hrrule

import (
	"strconv"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	commaByte     = ','
	monthLayout   = "January"
	weekDayLayout = "Monday"
	UNKNOWN       = "UNKNOWN"
)

var longMonthNames = []string{
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

var longDayNames = []string{
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
	"Sunday",
}

type formatterImpl struct {
	loc *i18n.Localizer
}

func NewDateFormatterSimple(loc *i18n.Localizer) DateFormatter {
	return &formatterImpl{
		loc: loc,
	}
}

func (df *formatterImpl) Format(year int, month time.Month, day int) string {
	var date strings.Builder
	date.Grow(32)

	date.WriteString(strconv.Itoa(day))
	date.WriteByte(spaceByte)
	date.WriteString(monthName(month, df.loc))
	date.WriteByte(spaceByte)
	date.WriteString(strconv.Itoa(year))
	// TODO: weekday from year, month, day
	// timing.WriteString(weekdayName(weekday, lang))

	return date.String()
}

func (df *formatterImpl) MonthName(month time.Month) string {
	return monthName(month, df.loc)
}

// Nth return int with suffix
// TODO: implement me
func (df *formatterImpl) Nth(i int) string {
	var last string

	if i < 0 {
		last = df.loc.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
			ID:    "Last",
			Other: "last",
		}})
	}

	if i == -1 {
		return last
	}

	nPos := abs(i)
	var nth strings.Builder
	nth.WriteString(strconv.Itoa(nPos))
	switch nPos {
	case 1, 21, 31:
		nth.WriteString(df.loc.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
			ID:    "Suffix for 1, 21, 31",
			Other: "st",
		}}))
	case 2, 22:
		nth.WriteString(df.loc.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
			ID:    "Suffix for 2, 22",
			Other: "nd",
		}}))
	case 3, 23:
		nth.WriteString(df.loc.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
			ID:    "Suffix for 3, 23",
			Other: "rd",
		}}))
	default:
		nth.WriteString(df.loc.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
			ID:    "Suffix for other numbers",
			Other: "th",
		}}))
	}

	if i < 0 {
		nth.WriteByte(spaceByte)
		nth.WriteString(last)
		return nth.String()
	}

	return nth.String()
}

// TODO: implement me
func (df *formatterImpl) WeekdayName(wDay Weekday) string {
	weekday := getGoWeekday(wDay.weekday)
	var sb strings.Builder
	sb.Grow(16)
	if wDay.n != 0 {
		sb.WriteString(df.Nth(wDay.n))
		sb.WriteByte(spaceByte)
	}

	sb.WriteString(time.Weekday(weekday).String())

	return sb.String()
}

func monthName(month time.Month, _ *i18n.Localizer) string {
	return month.String()
}

func getGoWeekday(num int) int {
	if num == 6 {
		return 0
	}

	return num + 1
}

func abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}