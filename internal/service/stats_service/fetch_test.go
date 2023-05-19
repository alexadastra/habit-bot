package stats_service

import (
	"testing"
	"time"
)

func Test_findWeekRange(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name     string
		args     args
		wantFrom time.Time
		wantTo   time.Time
	}{
		{
			name: "wednesday",
			args: args{
				t: time.Unix(1684331504, 0), // Wed May 17 2023 13:51:44 GMT+0000
			},
			wantFrom: time.Unix(1684098000, 0), // Sun May 14 2023 21:00:00 GMT+0000
			wantTo:   time.Unix(1684702800, 0), // Sun May 21 2023 21:00:00 GMT+0000
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := findWeekRange(tt.args.t)
			if got.Unix() != tt.wantFrom.Unix() {
				t.Errorf("findWeekRange() got = %v, want %v", got, tt.wantFrom)
			}
			if got1.Unix() != tt.wantTo.Unix() {
				t.Errorf("findWeekRange() got1 = %v, want %v", got1, tt.wantTo)
			}
		})
	}
}
