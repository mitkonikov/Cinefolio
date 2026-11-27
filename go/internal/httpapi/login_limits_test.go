package httpapi

import (
	"net/http/httptest"
	"net/netip"
	"testing"
)

func TestLoginIPTrustBoundary(t *testing.T) {
	for _, tt := range []struct {
		peer    string
		headers []string
		trust   bool
		want    string
	}{
		{"198.51.100.10:1234", []string{"198.51.100.20"}, true, "198.51.100.10"},
		{"127.0.0.1:1234", []string{"198.51.100.20"}, false, "127.0.0.1"},
		{"127.0.0.1:1234", []string{"198.51.100.20"}, true, "198.51.100.20"},
		{"127.0.0.1:1234", []string{"198.51.100.20, 198.51.100.21"}, true, "127.0.0.1"},
		{"127.0.0.1:1234", []string{"198.51.100.20", "198.51.100.21"}, true, "127.0.0.1"},
		{"[::1]:1234", []string{"2001:db8:a:b::1"}, true, "2001:db8:a:b::"},
	} {
		r := httptest.NewRequest("POST", "/api/admin/login", nil)
		r.RemoteAddr = tt.peer
		for _, v := range tt.headers {
			r.Header.Add("X-Forwarded-For", v)
		}
		if got := loginIP(r, tt.trust); got != netip.MustParseAddr(tt.want) {
			t.Errorf("peer %s headers %v trust %t: got %s, want %s", tt.peer, tt.headers, tt.trust, got, tt.want)
		}
	}
}
