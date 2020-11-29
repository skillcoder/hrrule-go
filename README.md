# hrrule-go 🧐👀⚛️📜

[![made with Go](https://img.shields.io/badge/made%20with-Go-brightgreen)](http://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-brightgreen)](https://github.com/skillcoder/hrrule-go/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/skillcoder/hrrule-go)](https://github.com/skillcoder/hrrule-go/issues)

Library for make **h**uman **r**eadable **r**ecurrence **rule**s from iCalendar RRULE (RFC5545) in Golang  
It supports serialization of recurrence rules to natural language, with internationalisation!

It is a partial port of the rrule module from [rrule.js](https://github.com/jakubroztocil/rrule) library.

## Using
```
package main

import (
	"fmt"
	"log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/skillcoder/hrrule-go"
)

func main() {
	hRule, err := hrrule.New(i18n.NewBundle(language.AmericanEnglish))
	if err != nil {
		log.Fatalf("filed to init rrule humanizer: %v", err)
	}

	rOption, err := hrrule.StrToROption("FREQ=MONTHLY;INTERVAL=1;BYDAY=-1FR;UNTIL=20241229T155400Z")
	if err != nil {
		log.Fatalf("rrule string to option: %v", err)
	}

	nlString, err := hRule.Humanize(rOption, "en-US")
	if err != nil {
		log.Fatalf("humanize rrule to string: %v", err)
	}
	fmt.Println(nlString)
}
```

## Translation to new language
See docs in https://github.com/nicksnyder/go-i18n
1. `touch l10n/translate.ru.toml`
2. `goi18n merge l10n/active.en-US.toml translate.ru.toml`
3. After `translate.ru.toml` has been translated, move it to `l10n/active.ru-RU.toml`.

🚧 It is necessary to agreement on the declension, cases and kind in the languages in which they exist 🚧

## TODO
 * Day of the week translation support
 * Months translation support
