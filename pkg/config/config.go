package config

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config is the runtime configuration for ZATRANO (HTTP, data stores, logging, security).
// Values are loaded from optional .env, config/{env}.yaml, and environment variables
// (nested YAML keys map to env like SECURITY_JWT_SECRET).
type Config struct {
	Env     string `mapstructure:"env"`
	AppName string `mapstructure:"app_name"`

	HTTPAddr        string        `mapstructure:"http_addr"`
	HTTPReadTimeout time.Duration `mapstructure:"http_read_timeout"`

	DatabaseURL      string `mapstructure:"database_url"`
	// DatabaseDriver selects the SQL backend: postgres | mysql | sqlserver | sqlite (default: postgres).
	// Override with DATABASE_DRIVER or config database_driver.
	DatabaseDriver   string `mapstructure:"database_driver"`
	DatabaseRequired bool   `mapstructure:"database_required"`

	RedisURL      string `mapstructure:"redis_url"`
	RedisRequired bool   `mapstructure:"redis_required"`

	LogLevel       string `mapstructure:"log_level"`
	LogDevelopment bool   `mapstructure:"log_development"`

	// StaticPath is the local directory for public assets (optional).
	StaticPath string `mapstructure:"static_path"`
	// StaticURLPrefix is the URL prefix for static files (e.g. /static).
	StaticURLPrefix string `mapstructure:"static_url_prefix"`

	MigrationsDir string `mapstructure:"migrations_dir"`
	// MigrationsSource: embed (default) uses driver-specific SQL from pkg/migrations; file uses migrations_dir on disk.
	MigrationsSource string `mapstructure:"migrations_source"`
	SeedsDir         string `mapstructure:"seeds_dir"`
	OpenAPIPath   string `mapstructure:"openapi_path"`

	Security  Security  `mapstructure:"security"`
	OAuth     OAuth     `mapstructure:"oauth"`
	API       API       `mapstructure:"api"`
	HTTP      HTTP      `mapstructure:"http"`
	I18n      I18n      `mapstructure:"i18n"`
	Mail      Mail      `mapstructure:"mail"`
	View      View      `mapstructure:"view"`
	Broadcast Broadcast `mapstructure:"broadcast"`
	Tenant    Tenant    `mapstructure:"tenant"`
	Audit     Audit     `mapstructure:"audit"`
	Search    Search    `mapstructure:"search"`
	Features  Features  `mapstructure:"features"`
	GraphQL   GraphQL   `mapstructure:"graphql"`
}

type API struct {
	Version string `mapstructure:"version"`
	Prefix  string `mapstructure:"prefix"`
}

// LoadOptions controls where configuration is read from.
type LoadOptions struct {
	// Env is the profile name (e.g. dev, prod). Defaults to ENV or "dev".
	Env string
	// ConfigDir is the directory containing {env}.yaml (default "config").
	ConfigDir string
	// DotEnv, if true, loads .env from the working directory when present.
	DotEnv bool
}

// Load reads configuration and returns a validated Config.
func Load(opts LoadOptions) (*Config, error) {
	if opts.ConfigDir == "" {
		opts.ConfigDir = "config"
	}
	envName := strings.TrimSpace(opts.Env)
	if envName == "" {
		envName = strings.TrimSpace(os.Getenv("ENV"))
	}
	if envName == "" {
		envName = "dev"
	}

	if opts.DotEnv {
		_ = godotenv.Load()
	}

	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetConfigName(envName)
	v.SetConfigType("yaml")
	v.AddConfigPath(opts.ConfigDir)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			v.Set("env", envName)
		} else {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	// Defaults when file is missing or keys omitted
	v.SetDefault("env", envName)
	v.SetDefault("app_name", "ZATRANO")
	v.SetDefault("http_addr", ":8080")
	v.SetDefault("http_read_timeout", 30*time.Second)
	v.SetDefault("database_required", false)
	v.SetDefault("database_driver", "")
	v.SetDefault("redis_required", false)
	v.SetDefault("log_level", "info")
	v.SetDefault("log_development", envName == "dev")
	v.SetDefault("static_path", "public")
	v.SetDefault("static_url_prefix", "/static")
	v.SetDefault("migrations_dir", "migrations")
	v.SetDefault("migrations_source", "embed")
	v.SetDefault("seeds_dir", "db/seeds")
	v.SetDefault("openapi_path", "api/openapi.yaml")
	v.SetDefault("security.session_enabled", true)
	v.SetDefault("security.csrf_enabled", true)
	v.SetDefault("security.jwt_issuer", "zatrano")
	v.SetDefault("security.jwt_expiry", 60*time.Minute)
	v.SetDefault("security.demo_token_endpoint", false)
	v.SetDefault("oauth.enabled", false)
	v.SetDefault("http.cors_enabled", false)
	v.SetDefault("http.rate_limit_enabled", false)
	v.SetDefault("http.rate_limit_max", 100)
	v.SetDefault("http.rate_limit_window", time.Minute)
	v.SetDefault("http.rate_limit_redis", false)
	v.SetDefault("http.request_timeout", time.Duration(0))
	v.SetDefault("http.body_limit", 0)
	v.SetDefault("http.shutdown_timeout", 15*time.Second)
	v.SetDefault("http.graceful_restart", false)
	v.SetDefault("http.graceful_restart_pid_file", "")
	v.SetDefault("i18n.enabled", false)
	v.SetDefault("i18n.default_locale", "en")
	v.SetDefault("i18n.locales_dir", "locales")
	v.SetDefault("i18n.cookie_name", "zatrano_lang")
	v.SetDefault("i18n.query_key", "lang")
	v.SetDefault("mail.driver", "log")
	v.SetDefault("mail.from_name", "ZATRANO")
	v.SetDefault("mail.from_email", "noreply@example.com")
	v.SetDefault("mail.templates_dir", "views/mails")
	v.SetDefault("mail.smtp.host", "localhost")
	v.SetDefault("mail.smtp.port", 587)
	v.SetDefault("mail.smtp.encryption", "tls")

	v.SetDefault("view.root", "views")
	v.SetDefault("view.extension", ".html")
	v.SetDefault("view.components_dir", "components")
	v.SetDefault("view.layouts_dir", "layouts")
	v.SetDefault("view.dev_mode", envName == "dev")
	v.SetDefault("view.asset.public_dir", "public")
	v.SetDefault("view.asset.public_url", "/public")
	v.SetDefault("view.asset.vite_manifest", "")
	v.SetDefault("view.asset.vite_dev_url", "")
	v.SetDefault("api.version", "v1")
	v.SetDefault("broadcast.enabled", false)
	v.SetDefault("broadcast.path_prefix", "/broadcast")
	v.SetDefault("broadcast.jwt_query_param", "access_token")
	v.SetDefault("broadcast.sse_enabled", true)
	v.SetDefault("tenant.enabled", false)
	v.SetDefault("tenant.mode", "header")
	v.SetDefault("tenant.header_name", "X-Tenant-ID")
	v.SetDefault("tenant.isolation", "row")
	v.SetDefault("tenant.row_column", "tenant_id")
	v.SetDefault("tenant.schema_prefix", "tenant_")
	v.SetDefault("audit.enabled", false)
	v.SetDefault("audit.model_enabled", false)
	v.SetDefault("audit.http_enabled", false)
	v.SetDefault("audit.http_driver", "db")
	v.SetDefault("search.enabled", false)
	v.SetDefault("search.driver", "")
	v.SetDefault("search.default_index_prefix", "zatrano_")
	v.SetDefault("search.postgres_fts_language", "simple")
	v.SetDefault("features.enabled", false)
	v.SetDefault("features.source", "config")
	v.SetDefault("graphql.enabled", false)
	v.SetDefault("graphql.path", "/graphql")
	v.SetDefault("graphql.playground", false)
	v.SetDefault("graphql.playground_path", "/playground")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	cfg.applyDerivedDefaults()

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) validate() error {
	switch strings.ToLower(strings.TrimSpace(c.LogLevel)) {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("invalid log_level %q (use debug, info, warn, error)", c.LogLevel)
	}
	if c.DatabaseRequired && strings.TrimSpace(c.DatabaseURL) == "" {
		return fmt.Errorf("database_required is true but database_url is empty (set DATABASE_URL or config/database_url)")
	}
	if err := c.validateDatabaseDriver(); err != nil {
		return err
	}
	if err := c.validateMigrationsSource(); err != nil {
		return err
	}
	if c.RedisRequired && strings.TrimSpace(c.RedisURL) == "" {
		return fmt.Errorf("redis_required is true but redis_url is empty (set REDIS_URL or config/redis_url)")
	}
	if c.Security.DemoTokenEndpoint && strings.EqualFold(strings.TrimSpace(c.Env), "prod") {
		return fmt.Errorf("security.demo_token_endpoint cannot be true when env is prod")
	}
	if c.OAuth.Enabled {
		if strings.TrimSpace(c.RedisURL) == "" {
			return fmt.Errorf("oauth.enabled requires redis_url (state storage)")
		}
		if strings.TrimSpace(c.OAuth.BaseURL) == "" {
			return fmt.Errorf("oauth.base_url is required when oauth.enabled is true")
		}
		if !oauthProviderConfigured(c.OAuth.Providers.Google) && !oauthProviderConfigured(c.OAuth.Providers.Github) {
			return fmt.Errorf("oauth.enabled requires at least one provider with client_id (google or github)")
		}
	}
	if c.Security.APIKeysEnabled && strings.TrimSpace(c.DatabaseURL) == "" {
		return fmt.Errorf("security.api_keys_enabled requires database_url")
	}
	if err := c.validateHTTP(); err != nil {
		return err
	}
	if err := c.validateI18n(); err != nil {
		return err
	}
	if err := c.validateTenant(); err != nil {
		return err
	}
	if err := c.validateAudit(); err != nil {
		return err
	}
	if err := c.validateSearch(); err != nil {
		return err
	}
	if err := c.validateFeatures(); err != nil {
		return err
	}
	if err := c.validateGraphQL(); err != nil {
		return err
	}
	return nil
}

func oauthProviderConfigured(p OAuthProvider) bool {
	return strings.TrimSpace(p.ClientID) != "" && strings.TrimSpace(p.ClientSecret) != ""
}

func (c *Config) applyDerivedDefaults() {
	c.DatabaseDriver = strings.ToLower(strings.TrimSpace(c.DatabaseDriver))
	if runtime.GOOS == "windows" {
		c.HTTP.GracefulRestart = false
	}
	if strings.TrimSpace(c.RedisURL) == "" {
		c.Security.SessionEnabled = false
		c.Security.CSRFEnabled = false
	}
	if len(c.Security.CSRFSkipPrefixes) == 0 {
		c.Security.CSRFSkipPrefixes = []string{"/api/"}
	}
	if c.Security.JWTExpiry <= 0 {
		c.Security.JWTExpiry = 60 * time.Minute
	}
	if strings.TrimSpace(c.Security.JWTIssuer) == "" {
		c.Security.JWTIssuer = "zatrano"
	}
	if strings.TrimSpace(c.API.Version) == "" {
		c.API.Version = "v1"
	}
	if strings.TrimSpace(c.API.Prefix) == "" {
		c.API.Prefix = "/api/" + strings.TrimPrefix(strings.TrimSpace(c.API.Version), "/")
	}
	if strings.TrimSpace(c.Security.APIKeyHeader) == "" {
		c.Security.APIKeyHeader = "X-API-Key"
	}
	if strings.TrimSpace(c.MigrationsDir) == "" {
		c.MigrationsDir = "migrations"
	}
	if strings.TrimSpace(c.MigrationsSource) == "" {
		c.MigrationsSource = "embed"
	}
	c.MigrationsSource = strings.ToLower(strings.TrimSpace(c.MigrationsSource))
	if strings.TrimSpace(c.SeedsDir) == "" {
		c.SeedsDir = "db/seeds"
	}
	if strings.TrimSpace(c.OpenAPIPath) == "" {
		c.OpenAPIPath = "api/openapi.yaml"
	}
	c.Security.CSRFSkipPrefixes = appendUniquePrefix(c.Security.CSRFSkipPrefixes, "/auth/oauth/")
	c.applyHTTPDefaults()
	c.applyI18nDefaults()
	c.applyBroadcastDefaults()
	c.applyTenantDefaults()
	c.applyAuditDefaults()
	c.applySearchDefaults()
	c.applyFeaturesDefaults()
	c.applyGraphQLDefaults()
}

func appendUniquePrefix(s []string, v string) []string {
	for _, x := range s {
		if x == v {
			return s
		}
	}
	return append(s, v)
}

// NormalizedDatabaseDriver returns the active SQL driver (default postgres).
func (c *Config) NormalizedDatabaseDriver() string {
	if strings.TrimSpace(c.DatabaseDriver) == "" {
		return "postgres"
	}
	return c.DatabaseDriver
}

func (c *Config) validateDatabaseDriver() error {
	switch c.NormalizedDatabaseDriver() {
	case "postgres", "mysql", "sqlserver", "sqlite":
		return nil
	default:
		return fmt.Errorf("invalid database_driver %q (use postgres, mysql, sqlserver, sqlite, or omit for postgres)", c.DatabaseDriver)
	}
}

func (c *Config) validateMigrationsSource() error {
	switch c.MigrationsSource {
	case "embed", "file":
		return nil
	default:
		return fmt.Errorf("invalid migrations_source %q (use embed or file)", c.MigrationsSource)
	}
}
