package search

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"gorm.io/gorm"
)

// Importer loads all rows for a model and pushes them to the search driver.
type Importer func(ctx context.Context, db *gorm.DB, drv Driver) error

var (
	importersMu sync.RWMutex
	importers   = map[string]Importer{}
)

// RegisterImporter registers a bulk import handler for a logical model name (case-insensitive).
func RegisterImporter(modelName string, fn Importer) {
	importersMu.Lock()
	defer importersMu.Unlock()
	key := strings.ToLower(strings.TrimSpace(modelName))
	if key == "" {
		panic("search: RegisterImporter: empty model name")
	}
	importers[key] = fn
}

// ImporterFor returns the registered importer or nil.
func ImporterFor(modelName string) Importer {
	importersMu.RLock()
	defer importersMu.RUnlock()
	return importers[strings.ToLower(strings.TrimSpace(modelName))]
}

// RegisteredModels returns sorted registered keys for CLI help.
func RegisteredModels() []string {
	importersMu.RLock()
	defer importersMu.RUnlock()
	out := make([]string, 0, len(importers))
	for k := range importers {
		out = append(out, k)
	}
	return out
}

// RunImport executes the named importer.
func RunImport(ctx context.Context, db *gorm.DB, drv Driver, modelName string) error {
	fn := ImporterFor(modelName)
	if fn == nil {
		keys := RegisteredModels()
		sort.Strings(keys)
		if len(keys) == 0 {
			return fmt.Errorf("no search importer registered for %q (no importers registered yet; call search.RegisterImporter in init())", modelName)
		}
		return fmt.Errorf("no search importer registered for %q; registered models: %v", modelName, keys)
	}
	return fn(ctx, db, drv)
}
