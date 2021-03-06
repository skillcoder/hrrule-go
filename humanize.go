package hrrule

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type HRRule struct {
	bundle *i18n.Bundle
}

func New(bundle *i18n.Bundle) (Humanizer, error) {
	if bundle == nil {
		var err error
		bundle, err = NewI18NBundle("./l10n")
		if err != nil {
			return nil, err
		}
	}

	return &HRRule{
		bundle: bundle,
	}, nil
}

func (imp *HRRule) Humanize(rule *ROption, lang string) (string, error) {
	localizer := i18n.NewLocalizer(imp.bundle, lang)
	dateFormatter := NewDateFormatterSimple(localizer)

	txt := newText(rule, localizer, dateFormatter)

	return txt.String(), nil
}
