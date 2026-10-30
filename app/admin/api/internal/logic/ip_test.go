package logic

import "testing"

func TestIPBlocked(t *testing.T) {
	for _, tc := range []struct {
		ip, rules string
		want      bool
	}{
		{"192.168.1.20", "127.0.0.1;192.168.1.*", true},
		{"192.168.2.20", "192.168.1.*", false},
		{"10.0.0.5", "10.0.0.9-10.0.0.5", true},
		{"10.0.1.1", "10.0.0.5-10.0.0.9", false},
		{"::1", "::1", true},
	} {
		if got := ipBlocked(tc.ip, tc.rules); got != tc.want {
			t.Errorf("ipBlocked(%q,%q) = %v", tc.ip, tc.rules, got)
		}
	}
}
