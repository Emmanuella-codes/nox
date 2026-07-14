package db

import "testing"

func TestNormalizeDatabaseURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "adds sslmode disable for localhost",
			in:   "postgresql://admin:secret@localhost:5432/nox_db",
			want: "postgresql://admin:secret@localhost:5432/nox_db?sslmode=disable",
		},
		{
			name: "adds sslmode disable for loopback ip",
			in:   "postgres://admin:secret@127.0.0.1:5432/nox_db",
			want: "postgres://admin:secret@127.0.0.1:5432/nox_db?sslmode=disable",
		},
		{
			name: "preserves explicit sslmode",
			in:   "postgresql://admin:secret@localhost:5432/nox_db?sslmode=require",
			want: "postgresql://admin:secret@localhost:5432/nox_db?sslmode=require",
		},
		{
			name: "preserves remote hosts",
			in:   "postgresql://admin:secret@db.example.com:5432/nox_db",
			want: "postgresql://admin:secret@db.example.com:5432/nox_db",
		},
		{
			name: "preserves invalid urls",
			in:   "not a url",
			want: "not a url",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := normalizeDatabaseURL(tt.in)
			if got != tt.want {
				t.Fatalf("normalizeDatabaseURL(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
