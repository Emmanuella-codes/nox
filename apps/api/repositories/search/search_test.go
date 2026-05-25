package search

import "testing"

func TestNormalizeOptionsDefaultsAndCapsLimit(t *testing.T) {
	tests := []struct {
		name string
		in   Options
		want Options
	}{
		{name: "defaults limit", in: Options{}, want: Options{Limit: 10}},
		{name: "caps limit", in: Options{Limit: 100, Offset: 2}, want: Options{Limit: 30, Offset: 2}},
		{name: "resets negative offset", in: Options{Limit: 5, Offset: -4}, want: Options{Limit: 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeOptions(tt.in)
			if got != tt.want {
				t.Fatalf("expected %+v, got %+v", tt.want, got)
			}
		})
	}
}
