package search

import "testing"

func TestNormalizeOptionsDefaultsAndCapsLimit(t *testing.T) {
	tests := []struct {
		name string
		in   Options
		want Options
	}{
		{name: "defaults limit", in: Options{}, want: Options{Limit: 10, Scope: "all"}},
		{name: "caps limit", in: Options{Limit: 100, Offset: 2}, want: Options{Limit: 30, Offset: 2, Scope: "all"}},
		{name: "resets negative offset", in: Options{Limit: 5, Offset: -4}, want: Options{Limit: 5, Scope: "all"}},
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
