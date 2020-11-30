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
		{
			name:       "every week in April on Tuesday for 30 times",
			lang:       "en-US",
			inputRRule: "FREQ=WEEKLY;COUNT=30;INTERVAL=1;WKST=MO;BYDAY=TU;BYMONTH=4",
			want:       "every week in April on Tuesday for 30 times",
			wantErr:    false,
		},
		{
			name:       "every April, June and September on Tuesday the last for 30 times",
			lang:       "en-US",
			inputRRule: "FREQ=MONTHLY;COUNT=30;INTERVAL=1;WKST=MO;BYDAY=TU;BYMONTH=4,6,9;BYMONTHDAY=-1",
			want:       "every April, June and September on Tuesday the last for 30 times",
			wantErr:    false,
		},
		{
			name:       "every 2 years March, April and June on Monday, Tuesday, Wednesday, Thursday or Friday the 2nd last for 30 times",
			lang:       "en-US",
			inputRRule: "FREQ=YEARLY;COUNT=30;INTERVAL=2;WKST=MO;BYDAY=MO,TU,WE,TH,FR;BYMONTH=3,4,6;BYMONTHDAY=-2",
			want:       "every 2 years March, April and June on Monday, Tuesday, Wednesday, Thursday or Friday the 2nd last for 30 times",
			wantErr:    false,
		},
		{
			name:       "every week on Tuesday or Wednesday the 2nd last until 29 December 2024",
			lang:       "en-US",
			inputRRule: "FREQ=WEEKLY;UNTIL=20241229T155400Z;INTERVAL=1;WKST=MO;BYDAY=TU,WE;BYMONTHDAY=-2;BYWEEKNO=13,20",
			want:       "every week on Tuesday or Wednesday the 2nd last until 29 December 2024",
			wantErr:    false,
		},
		{
			name:       "every year on the 256th day for 10 times",
			lang:       "en-US",
			inputRRule: "FREQ=YEARLY;COUNT=10;INTERVAL=1;WKST=MO;BYYEARDAY=256",
			want:       "every year on the 256th day for 10 times",
			wantErr:    false,
		},
		{
			name:       "every month on the last Friday until 29 December 2024",
			lang:       "en-US",
			inputRRule: "FREQ=MONTHLY;INTERVAL=1;BYDAY=-1FR;UNTIL=20241229T155400Z",
			want:       "every month on the last Friday until 29 December 2024",
			wantErr:    false,
		},
		{
			name:       "every month on Thursday until 3 December 2021",
			lang:       "en-US",
			inputRRule: "FREQ=MONTHLY;UNTIL=20211203T160000Z;INTERVAL=1;WKST=MO;BYDAY=TH;BYSETPOS=1",
			want:       "every month on Thursday until 3 December 2021",
			wantErr:    false,
		},
	}

	bundle, err := NewI18NBundle("./l10n")
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
