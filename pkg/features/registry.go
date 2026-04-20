package features

import (
	"context"
	"errors"
	"strings"

	"github.com/zatrano/zatrano/pkg/config"
	"gorm.io/gorm"
)

// Registry resolves feature flags from config and/or database.
type Registry struct {
	cfg     *config.Features
	source  string
	db      *gorm.DB
	static  map[string]config.FeatureDefinition
	enabled bool
}

// NewRegistry builds a registry. db may be nil when source is config only.
func NewRegistry(cfg *config.Config, db *gorm.DB) *Registry {
	if cfg == nil {
		return &Registry{}
	}
	f := &cfg.Features
	static := make(map[string]config.FeatureDefinition)
	for _, d := range f.Definitions {
		k := normalizeKey(d.Key)
		if k == "" {
			continue
		}
		static[k] = d
	}
	return &Registry{
		cfg:     f,
		source:  strings.ToLower(strings.TrimSpace(f.Source)),
		db:      db,
		static:  static,
		enabled: f.Enabled,
	}
}

// Enabled reports whether the features subsystem is turned on in config.
func (r *Registry) Enabled() bool { return r != nil && r.enabled }

func (r *Registry) definitionFromConfig(key string) (config.FeatureDefinition, bool) {
	if r == nil {
		return config.FeatureDefinition{}, false
	}
	d, ok := r.static[normalizeKey(key)]
	return d, ok
}

func (r *Registry) loadDB(ctx context.Context, key string) (Definition, bool) {
	if r == nil || r.db == nil {
		return Definition{}, false
	}
	var row DBFlag
	tx := r.db.WithContext(ctx).Where("key = ?", key).Limit(1).First(&row)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return Definition{}, false
		}
		return Definition{}, false
	}
	if tx.RowsAffected == 0 {
		return Definition{}, false
	}
	return Definition{
		Enabled:        row.Enabled,
		RolloutPercent: row.RolloutPercent,
		AllowedRoles:   row.AllowedRoles,
	}, true
}

func (r *Registry) resolveDefinition(ctx context.Context, key string) (Definition, bool) {
	k := normalizeKey(key)
	if k == "" {
		return Definition{}, false
	}
	if !r.Enabled() {
		return Definition{}, false
	}
	switch r.source {
	case "config":
		d, ok := r.definitionFromConfig(k)
		if !ok {
			return Definition{}, false
		}
		return Definition{
			Enabled:        d.Enabled,
			RolloutPercent: d.RolloutPercent,
			AllowedRoles:   d.AllowedRoles,
		}, true
	case "db":
		return r.loadDB(ctx, k)
	case "both":
		if def, ok := r.loadDB(ctx, k); ok {
			return def, true
		}
		d, ok := r.definitionFromConfig(k)
		if !ok {
			return Definition{}, false
		}
		return Definition{
			Enabled:        d.Enabled,
			RolloutPercent: d.RolloutPercent,
			AllowedRoles:   d.AllowedRoles,
		}, true
	default:
		return Definition{}, false
	}
}
