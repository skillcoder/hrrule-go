package hrrule

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func NewI18NBundle() (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	err := filepath.Walk("l10n", getWalkFunc(bundle))
	if err != nil {
		return nil, err
	}

	return bundle, nil
}

func getWalkFunc(b *i18n.Bundle) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if filepath.Ext(path) == ".toml" {
				_, err := b.LoadMessageFile(path)
				return err
			}
		}
		return nil
	}
}
