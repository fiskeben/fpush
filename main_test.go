package main

import (
	"reflect"
	"testing"
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
