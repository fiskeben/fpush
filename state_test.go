package main

import (
	"testing"
	"time"
)

func TestGetLastPushTime(t *testing.T) {
	nowfunc = func() time.Time {
		return time.Date(2019, 7, 13, 22, 12, 17, 0, time.UTC)
	}

	tests := []struct {
		name    string
		args    string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "gets time of last push",
			args:    "testdata/statefile",
			want:    time.Date(2019, 7, 13, 23, 58, 10, 0, time.UTC), ///time.Parse(time.RFC3339, "2019-07-13T23:58:10+02:00"),
			wantErr: false,
		},
		{
			name:    "gets default time when file does not exist",
			args:    "testdata/missing-statefile",
			want:    time.Date(2019, 7, 13, 21, 12, 17, 0, time.UTC),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		res, err := readStateFile(tt.args)
		if err != nil && !tt.wantErr {
			t.Errorf("%s got unexpected error: %v", tt.name, err)
			continue
		}
		if err == nil && tt.wantErr {
			t.Errorf("%s expected an error", tt.name)
			continue
		}
		if res != tt.want {
			t.Errorf("%s got %v expected %v", tt.name, res, tt.want)
		}
	}
}
