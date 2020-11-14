package hrrule

import (
	"strconv"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Frequency denotes the period on which the rule is evaluated.
type Frequency int

// Constants
const (
	YEARLY Frequency = iota
	MONTHLY
	WEEKLY
	DAILY
	HOURLY
	MINUTELY
	SECONDLY
)

// Weekday specifying the nth weekday.
// Field N could be positive or negative (like MO(+2) or MO(-3).
// Not specifying N (0) is the same as specifying +1.
type Weekday struct {
	weekday int
	n       int
}

// Nth return the nth weekday
// __call__ - Cannot call the object directly,
// do it through e.g. TH.nth(-1) instead,
func (wday *Weekday) Nth(n int) Weekday {
	return Weekday{wday.weekday, n}
}

// N returns index of the week, e.g. for 3MO, N() will return 3
func (wday *Weekday) N() int {
	return wday.n
}

// Day returns index of the day in a week (0 for MO, 6 for SU)
func (wday *Weekday) Day() int {
	return wday.weekday
}

const plusByte = '+'

var dayNames = []string{
	"MO",
	"TU",
	"WE",
	"TH",
	"FR",
	"SA",
	"SU",
}

// String convert struct to human readable string
// String returns the English rrule name of the Weekday with number
func (wday *Weekday) String(loc *i18n.Localizer) string {
	if MO.weekday <= wday.weekday && wday.weekday <= SU.weekday {
		var sb strings.Builder
		sb.Grow(32)
		if wday.n != 0 {
			if wday.n > 0 {
				sb.WriteByte(plusByte)
			}
			sb.WriteString(strconv.Itoa(wday.n))
			sb.WriteString(dayNames[wday.weekday])
			return sb.String()
		}

		sb.WriteString(dayNames[wday.weekday])
		return sb.String()
	}

	return "%!Weekday(" + strconv.Itoa(wday.weekday) + ", " +strconv.Itoa(wday.n)+ ")"
}

// Weekdays
var (
	MO = Weekday{weekday: 0}
	TU = Weekday{weekday: 1}
	WE = Weekday{weekday: 2}
	TH = Weekday{weekday: 3}
	FR = Weekday{weekday: 4}
	SA = Weekday{weekday: 5}
	SU = Weekday{weekday: 6}
)

type ROption struct {
	Freq       Frequency
	Dtstart    time.Time
	Interval   int
	Wkst       Weekday
	Count      int
	Until      time.Time
	Bysetpos   []int
	Bymonth    []int
	Bymonthday []int
	Byyearday  []int
	Byweekno   []int
	Byweekday  []Weekday
	Byhour     []int
	Byminute   []int
	Bysecond   []int
	Byeaster   []int
}
