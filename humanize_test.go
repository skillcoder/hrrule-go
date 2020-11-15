package hrrule

import (
	"testing"
)

func TestHRRule_Humanize(t *testing.T) {
	tests := []struct {
		name       string
		lang       string
		inputRRule string
		want       string
		wantErr    bool
	}{
		{
			name:       "every 5 days",
			lang:       "en-US",
			inputRRule: "FREQ=DAILY;INTERVAL=5",
			want:       "every 5 days",
			wantErr:    false,
		},
		{
			name:       "every 2 days on Thursday, Friday and on the 2nd last Tuesday",
			lang:       "en-US",
			inputRRule: "FREQ=DAILY;INTERVAL=2;WKST=MO;BYDAY=-2TU,TH,FR",
			want:       "every 2 days on Thursday, Friday and on the 2nd last Tuesday",
			wantErr:    false,
		},
	}

	bundle, err := NewI18NBundle()
	if err != nil {
		t.Fatalf("create bundle: %v", err)
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imp := &HRRule{
				bundle: bundle,
			}

			rOption, err := StrToROption(tt.inputRRule)
			if err != nil {
				t.Fatalf("str to option: %v", err)
			}

			got, err := imp.Humanize(rOption, tt.lang)
			if (err != nil) != tt.wantErr {
				t.Errorf("Humanize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Humanize()\n\thave: %v\n\twant: %v", got, tt.want)
			}
		})
	}
}
