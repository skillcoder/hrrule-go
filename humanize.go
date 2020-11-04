package hrrule

type HRRule struct {
}

func New() Humanizer {
	return &HRRule{}
}

func (imp *HRRule) Humanize(rule ROption, lang string) (string, error) {
	txt := newText(rule, lang)

	return txt.String(), nil
}
