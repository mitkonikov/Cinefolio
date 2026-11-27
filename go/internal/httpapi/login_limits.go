package httpapi

import (
	"net"
	"net/http"
	"net/netip"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	loginBurst       = 8
	loginClientLimit = 4096
	loginClientIdle  = time.Minute
)

type loginClient struct {
	limiter *rate.Limiter
	seen    time.Time
}

type loginLimits struct {
	mu        sync.Mutex
	clients   map[netip.Addr]loginClient
	nextSweep time.Time
}

func (l *loginLimits) allow(ip netip.Addr, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.clients == nil {
		l.clients = make(map[netip.Addr]loginClient)
	}
	if !now.Before(l.nextSweep) {
		for key, client := range l.clients {
			if now.Sub(client.seen) >= loginClientIdle {
				delete(l.clients, key)
			}
		}
		l.nextSweep = now.Add(loginClientIdle)
	}
	client, exists := l.clients[ip]
	if !exists {
		if len(l.clients) >= loginClientLimit {
			return false
		}
		client.limiter = rate.NewLimiter(rate.Every(2*time.Second), loginBurst)
	}
	client.seen = now
	l.clients[ip] = client
	return client.limiter.AllowN(now, 1)
}

func loginIP(r *http.Request, trustLoopbackProxy bool) netip.Addr {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	peer, _ := netip.ParseAddr(host)
	ip := peer.Unmap()
	// Nginx overwrites the header with one address. Untrusted peers and malformed
	// headers fall back to the socket address.
	if trustLoopbackProxy && ip.IsLoopback() {
		if values := r.Header.Values("X-Forwarded-For"); len(values) == 1 {
			if forwarded, err := netip.ParseAddr(values[0]); err == nil && forwarded.Zone() == "" {
				ip = forwarded.Unmap()
			}
		}
	}
	if ip.Is6() {
		ip = netip.PrefixFrom(ip, 64).Masked().Addr()
	}
	return ip
}
