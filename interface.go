package hrrule

import (
	"time"
)

type Humanizer interface {
	Humanize(rule *ROption, lang string) (string, error)
}

type DateFormatter interface {
	Format(year int, month time.Month, day int) string
	MonthName(month time.Month) string
	Nth(i int) string
	WeekdayName(wDay Weekday) string
}
