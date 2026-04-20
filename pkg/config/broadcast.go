package config

import "strings"

// Broadcast configures WebSocket hub, Pusher-compatible framing, and optional SSE.
type Broadcast struct {
	// Enabled turns on /broadcast routes and Hub allocation at bootstrap.
	Enabled bool `mapstructure:"enabled"`
	// PathPrefix is the URL prefix for broadcast routes (default /broadcast).
	PathPrefix string `mapstructure:"path_prefix"`
	// JWTQueryParam is the query key for Bearer token on WebSocket/SSE URLs (browser-friendly).
	JWTQueryParam string `mapstructure:"jwt_query_param"`
	// SSEEnabled registers GET {prefix}/sse/:channel when broadcast is enabled.
	SSEEnabled bool `mapstructure:"sse_enabled"`
	// AllowOrigins lists acceptable Origin values for WebSocket upgrade (empty = allow any).
	AllowOrigins []string `mapstructure:"allow_origins"`
}

func (c *Config) applyBroadcastDefaults() {
	b := &c.Broadcast
	if strings.TrimSpace(b.PathPrefix) == "" {
		b.PathPrefix = "/broadcast"
	}
	if strings.TrimSpace(b.JWTQueryParam) == "" {
		b.JWTQueryParam = "access_token"
	}
}
