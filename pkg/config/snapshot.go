package config

import "strings"

// SanitizedSnapshot returns a nested map safe to print (secrets masked). Stable keys for scripting.
func SanitizedSnapshot(c *Config) map[string]any {
	out := map[string]any{
		"env":               c.Env,
		"app_name":          c.AppName,
		"http_addr":         c.HTTPAddr,
		"http_read_timeout": c.HTTPReadTimeout.String(),
		"database_url":      MaskConnectionURL(c.DatabaseURL),
		"database_required": c.DatabaseRequired,
		"redis_url":         MaskConnectionURL(c.RedisURL),
		"redis_required":    c.RedisRequired,
		"log_level":         c.LogLevel,
		"log_development":   c.LogDevelopment,
		"static_path":       c.StaticPath,
		"static_url_prefix": c.StaticURLPrefix,
		"migrations_dir":    c.MigrationsDir,
		"seeds_dir":         c.SeedsDir,
		"openapi_path":      c.OpenAPIPath,
		"http": map[string]any{
			"cors_enabled":              c.HTTP.CORSEnabled,
			"cors_allow_origins":        c.HTTP.CORSAllowOrigins,
			"cors_allow_credentials":    c.HTTP.CORSAllowCredentials,
			"rate_limit_enabled":        c.HTTP.RateLimitEnabled,
			"rate_limit_max":            c.HTTP.RateLimitMax,
			"rate_limit_window":         c.HTTP.RateLimitWindow.String(),
			"rate_limit_redis":          c.HTTP.RateLimitRedis,
			"rate_limit_by_ip":          c.HTTP.RateLimitByIP,
			"rate_limit_by_jwt_subject": c.HTTP.RateLimitByJWTSub,
			"request_timeout":           c.HTTP.RequestTimeout.String(),
			"body_limit":                c.HTTP.BodyLimit,
		},
		"security": map[string]any{
			"session_enabled":     c.Security.SessionEnabled,
			"csrf_enabled":        c.Security.CSRFEnabled,
			"csrf_skip_prefixes":  c.Security.CSRFSkipPrefixes,
			"trusted_origins":     c.Security.TrustedOrigins,
			"jwt_secret":          MaskSecret(c.Security.JWTSecret),
			"jwt_issuer":          c.Security.JWTIssuer,
			"jwt_expiry":          c.Security.JWTExpiry.String(),
			"cookie_secure":       c.Security.CookieSecure,
			"demo_token_endpoint": c.Security.DemoTokenEndpoint,
			"api_keys_enabled":    c.Security.APIKeysEnabled,
			"api_key_header":      c.Security.APIKeyHeader,
		},
		"oauth": oauthSnapshot(&c.OAuth),
		"i18n": map[string]any{
			"enabled":           c.I18n.Enabled,
			"default_locale":    c.I18n.DefaultLocale,
			"supported_locales": c.I18n.SupportedLocales,
			"locales_dir":       c.I18n.LocalesDir,
			"cookie_name":       c.I18n.CookieName,
			"query_key":         c.I18n.QueryKey,
		},
		"view": map[string]any{
			"root":           c.View.Root,
			"extension":      c.View.Extension,
			"components_dir": c.View.ComponentsDir,
			"layouts_dir":    c.View.LayoutsDir,
			"dev_mode":       c.View.DevMode,
			"asset": map[string]any{
				"public_dir":    c.View.Asset.PublicDir,
				"public_url":    c.View.Asset.PublicURL,
				"vite_manifest": c.View.Asset.ViteManifest,
				"vite_dev_url":  c.View.Asset.ViteDevURL,
			},
		},
		"broadcast": map[string]any{
			"enabled":           c.Broadcast.Enabled,
			"path_prefix":       c.Broadcast.PathPrefix,
			"jwt_query_param":   c.Broadcast.JWTQueryParam,
			"sse_enabled":       c.Broadcast.SSEEnabled,
			"allow_origins":     c.Broadcast.AllowOrigins,
		},
		"tenant": map[string]any{
			"enabled":            c.Tenant.Enabled,
			"mode":               c.Tenant.Mode,
			"header_name":        c.Tenant.HeaderName,
			"subdomain_suffix":   c.Tenant.SubdomainSuffix,
			"required":           c.Tenant.Required,
			"isolation":          c.Tenant.Isolation,
			"row_column":         c.Tenant.RowColumn,
			"schema_prefix":      c.Tenant.SchemaPrefix,
		},
		"audit": map[string]any{
			"enabled":        c.Audit.Enabled,
			"model_enabled":  c.Audit.ModelEnabled,
			"http_enabled":   c.Audit.HttpEnabled,
			"http_driver":    c.Audit.HttpDriver,
			"http_file_path": c.Audit.HttpFilePath,
		},
		"search": map[string]any{
			"enabled":                 c.Search.Enabled,
			"driver":                  strings.TrimSpace(c.Search.Driver),
			"default_index_prefix":    c.Search.DefaultIndexPrefix,
			"meilisearch_url":         strings.TrimSpace(c.Search.MeilisearchURL),
			"typesense_url":           strings.TrimSpace(c.Search.TypesenseURL),
			"postgres_fts_language":   c.Search.PostgresFTSLanguage,
			"meilisearch_api_key":     MaskSecret(c.Search.MeilisearchAPIKey),
			"typesense_api_key":       MaskSecret(c.Search.TypesenseAPIKey),
		},
		"features": map[string]any{
			"enabled":            c.Features.Enabled,
			"source":             c.Features.Source,
			"definitions_count": len(c.Features.Definitions),
		},
		"graphql": map[string]any{
			"enabled":          c.GraphQL.Enabled,
			"path":             c.GraphQL.Path,
			"playground":       c.GraphQL.Playground,
			"playground_path":  c.GraphQL.PlaygroundPath,
		},
	}
	return out
}

func oauthSnapshot(o *OAuth) map[string]any {
	prov := map[string]any{
		"google": map[string]any{
			"client_id":     strings.TrimSpace(o.Providers.Google.ClientID),
			"client_secret": MaskSecret(o.Providers.Google.ClientSecret),
			"scopes":        o.Providers.Google.Scopes,
		},
		"github": map[string]any{
			"client_id":     strings.TrimSpace(o.Providers.Github.ClientID),
			"client_secret": MaskSecret(o.Providers.Github.ClientSecret),
			"scopes":        o.Providers.Github.Scopes,
		},
	}
	return map[string]any{
		"enabled":   o.Enabled,
		"base_url":  strings.TrimSpace(o.BaseURL),
		"providers": prov,
	}
}
