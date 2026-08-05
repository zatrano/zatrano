package config

import "github.com/zatrano/framework/core/env"

// Auth returns authentication guard and provider defaults.
func Auth() map[string]any {
	return map[string]any{
		"defaults": map[string]any{
			"guard":    env.Get("AUTH_GUARD", "web"),
			"provider": env.Get("AUTH_PROVIDER", "users"),
		},
		"guards": map[string]any{
			"web": map[string]any{
				"driver":   "session",
				"provider": "users",
			},
			"api": map[string]any{
				"driver":   "session",
				"provider": "users",
			},
		},
		"providers": map[string]any{
			"users": map[string]any{
				"driver": "database",
				"table":  "users",
			},
		},
	}
}
