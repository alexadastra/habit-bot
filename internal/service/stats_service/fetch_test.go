package stats_service

import (
	"reflect"
	"testing"
	"time"
)

func Test_findWeekRange(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name  string
		args  args
		wantFrom  time.Time
		wantTo time.Time
	}{
		{
			name: "wednesday",
			args: args{
				t: time.Unix(1684331504, 0), 	// Wed May 17 2023 13:51:44 GMT+0000
			},
			wantFrom: time.Unix(1684011600, 0), // Sat May 13 2023 21:00:00 GMT+0000
			wantTo: time.Unix(1684702800, 0), 	// Sun May 21 2023 21:00:00 GMT+0000
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := findWeekRange(tt.args.t)
			if !reflect.DeepEqual(got, tt.wantFrom) {
				t.Errorf("findWeekRange() got = %v, want %v", got, tt.wantFrom)
			}
			if !reflect.DeepEqual(got1, tt.wantTo) {
				t.Errorf("findWeekRange() got1 = %v, want %v", got1, tt.wantTo)
			}
		})
	}
}
