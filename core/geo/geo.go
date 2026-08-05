package geo

import (
	"net"
	"strings"
	"sync"
)

// Location is a coarse geo result.
type Location struct {
	IP          string  `json:"ip"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	City        string  `json:"city,omitempty"`
	Lat         float64 `json:"lat,omitempty"`
	Lon         float64 `json:"lon,omitempty"`
	Source      string  `json:"source"`
}

// Resolver looks up IP geolocation (stub / private-range aware).
type Resolver struct {
	mu    sync.RWMutex
	exact map[string]Location
}

// New creates a geo resolver with a few demo mappings.
func New() *Resolver {
	r := &Resolver{exact: make(map[string]Location)}
	r.Put("8.8.8.8", Location{Country: "United States", CountryCode: "US", City: "Mountain View", Lat: 37.386, Lon: -122.084, Source: "stub"})
	r.Put("1.1.1.1", Location{Country: "Australia", CountryCode: "AU", City: "Sydney", Lat: -33.8688, Lon: 151.2093, Source: "stub"})
	return r
}

// Put registers an exact IP mapping.
func (r *Resolver) Put(ip string, loc Location) {
	r.mu.Lock()
	defer r.mu.Unlock()
	loc.IP = ip
	if loc.Source == "" {
		loc.Source = "manual"
	}
	r.exact[ip] = loc
}

// Lookup resolves an IP address.
func (r *Resolver) Lookup(ip string) Location {
	ip = strings.TrimSpace(ip)
	if host, _, err := net.SplitHostPort(ip); err == nil {
		ip = host
	}
	r.mu.RLock()
	if loc, ok := r.exact[ip]; ok {
		r.mu.RUnlock()
		return loc
	}
	r.mu.RUnlock()

	parsed := net.ParseIP(ip)
	if parsed == nil {
		return Location{IP: ip, Country: "Unknown", CountryCode: "XX", Source: "invalid"}
	}
	if parsed.IsLoopback() {
		return Location{IP: ip, Country: "Local", CountryCode: "LO", City: "Loopback", Source: "private"}
	}
	if parsed.IsPrivate() {
		return Location{IP: ip, Country: "Private Network", CountryCode: "PR", Source: "private"}
	}
	if parsed.IsLinkLocalUnicast() || parsed.IsLinkLocalMulticast() {
		return Location{IP: ip, Country: "Link Local", CountryCode: "LL", Source: "private"}
	}
	return Location{IP: ip, Country: "Unknown", CountryCode: "XX", Source: "stub"}
}
