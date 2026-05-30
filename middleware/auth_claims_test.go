package middleware

import (
	"encoding/json"
	"testing"
)

func TestToIntClaim(t *testing.T) {
	cases := []struct {
		value interface{}
		ok    bool
		want  int
	}{
		{value: 1, ok: true, want: 1},
		{value: int64(2), ok: true, want: 2},
		{value: float64(3), ok: true, want: 3},
		{value: json.Number("4"), ok: true, want: 4},
		{value: "5", ok: true, want: 5},
		{value: "bad", ok: false, want: 0},
		{value: nil, ok: false, want: 0},
	}

	for _, tc := range cases {
		got, ok := toIntClaim(tc.value)
		if ok != tc.ok || got != tc.want {
			t.Fatalf("toIntClaim(%v) = (%d, %v), want (%d, %v)", tc.value, got, ok, tc.want, tc.ok)
		}
	}
}
