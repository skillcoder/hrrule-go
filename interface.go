package hrrule

type Humanizer interface {
	Humanize(rule ROption, lang string) (string, error)
}
