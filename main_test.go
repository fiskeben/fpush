package main

import (
	"reflect"
	"testing"
	"time"
)

func TestLimitFiles(t *testing.T) {
	tests := []struct {
		name string
		arg  []string
		want []string
	}{
		{
			name: "returns slice smaller than 5",
			arg:  []string{"a", "b", "c", "d"},
			want: []string{"a", "b", "c", "d"},
		},
		{
			name: "5 item slice is limited",
			arg:  []string{"a", "b", "c", "d", "e"},
			want: []string{"a", "b", "c", "e"},
		},
		{
			name: "7 item slice is limited",
			arg:  []string{"a", "b", "c", "d", "e", "f", "g"},
			want: []string{"a", "c", "e", "g"},
		},
		{
			name: "8 item slice is limited",
			arg:  []string{"a", "b", "c", "d", "e", "f", "g", "h"},
			want: []string{"a", "c", "f", "h"},
		},
	}

	for _, tt := range tests {
		if res := limitFiles(tt.arg); !reflect.DeepEqual(res, tt.want) {
			t.Errorf("%s failed, %v != %v", tt.name, res, tt.want)
		}
	}
}

type testFileInfo struct {
	modTime time.Time
}

func (t testFileInfo) ModTime() time.Time {
	return t.modTime
}

func TestCheckFile(t *testing.T) {
	now := time.Now()

	type args struct {
		f        fileInfo
		now      time.Time
		lastPush time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1 minute old last push 60 minutes ago",
			args: args{
				f:        testFileInfo{now.Add(-1 * time.Minute)},
				now:      now,
				lastPush: now.Add(-60 * time.Minute),
			},
			want: true,
		},
		{
			name: "1 minute old last push 59 minutes ago",
			args: args{
				f:        testFileInfo{now.Add(-1 * time.Minute)},
				now:      now,
				lastPush: now.Add(-59 * time.Minute),
			},
			want: false,
		},
		{
			name: "5 minutes old last push 120 minutes ago",
			args: args{
				f:        testFileInfo{now.Add(-5 * time.Minute)},
				now:      now,
				lastPush: now.Add(-120 * time.Minute),
			},
			want: true,
		},
		{
			name: "6 minutes old last push 120 minutes ago",
			args: args{
				f:        testFileInfo{now.Add(-6 * time.Minute)},
				now:      now,
				lastPush: now.Add(-120 * time.Minute),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		if checkFile(tt.args.f, tt.args.now, tt.args.lastPush) != tt.want {
			t.Errorf("%s failed", tt.name)
		}
	}
}
