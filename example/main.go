package main

import (
	"fmt"
	"log"

	"github.com/skillcoder/hrrule-go"
)

func main() {
	bundle, err := hrrule.NewI18NBundle("../l10n")
	if err != nil {
		log.Fatalf("filed to init i18n bundle: %v", err)
	}

	hRule, err := hrrule.New(bundle)
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
