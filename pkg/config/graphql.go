package config

import (
	"fmt"
	"strings"
)

// GraphQL configures the optional gqlgen HTTP endpoint (Fiber + net/http adaptor).
type GraphQL struct {
	Enabled        bool   `mapstructure:"enabled"`
	Path           string `mapstructure:"path"`
	Playground     bool   `mapstructure:"playground"`
	PlaygroundPath string `mapstructure:"playground_path"`
}

func (c *Config) applyGraphQLDefaults() {
	g := &c.GraphQL
	if strings.TrimSpace(g.Path) == "" {
		g.Path = "/graphql"
	}
	if strings.TrimSpace(g.PlaygroundPath) == "" {
		g.PlaygroundPath = "/playground"
	}
}

func (c *Config) validateGraphQL() error {
	if !c.GraphQL.Enabled {
		return nil
	}
	if !strings.HasPrefix(strings.TrimSpace(c.GraphQL.Path), "/") {
		return fmt.Errorf("graphql.path must start with / (got %q)", c.GraphQL.Path)
	}
	if c.GraphQL.Playground && !strings.HasPrefix(strings.TrimSpace(c.GraphQL.PlaygroundPath), "/") {
		return fmt.Errorf("graphql.playground_path must start with / (got %q)", c.GraphQL.PlaygroundPath)
	}
	return nil
}
