package logic

import (
	"net/netip"
	"strings"
)

// RuoYi blacklist entries are separated by semicolons. IPv4 ranges include both endpoints.
func ipBlocked(ip, rules string) bool {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return false
	}
	addr = addr.Unmap()
	for _, rule := range strings.Split(rules, ";") {
		rule = strings.TrimSpace(rule)
		if addr.String() == rule {
			return true
		}
		if start, end, ok := strings.Cut(rule, "-"); ok {
			a, e1 := netip.ParseAddr(start)
			b, e2 := netip.ParseAddr(end)
			if e1 != nil || e2 != nil || !a.Is4() || !b.Is4() || !addr.Is4() {
				continue
			}
			if b.Less(a) {
				a, b = b, a
			}
			if !addr.Less(a) && !b.Less(addr) {
				return true
			}
		}
		if addr.Is4() && strings.Contains(rule, "*") {
			parts := strings.Split(rule, ".")
			if len(parts) != 4 {
				continue
			}
			octets := strings.Split(addr.String(), ".")
			match := true
			for i, part := range parts {
				if part != "*" && part != octets[i] {
					match = false
					break
				}
			}
			if match {
				return true
			}
		}
	}
	return false
}
