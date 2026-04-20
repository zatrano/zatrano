package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/zatrano/zatrano/pkg/config"
)

// Writer persists activity and HTTP audit records.
type Writer interface {
	WriteActivity(ctx context.Context, row *ActivityLog) error
	WriteHTTP(ctx context.Context, row *HTTPAuditLog) error
}

// auditWriter is the concrete implementation selected by config.
type auditWriter struct {
	cfg *config.Config
	db  *gorm.DB
	log *zap.Logger
	mu  sync.Mutex // file append
}

// NewWriter builds a Writer from config.
func NewWriter(cfg *config.Config, db *gorm.DB, log *zap.Logger) (Writer, error) {
	if cfg == nil || !cfg.Audit.Enabled {
		return NopWriter{}, nil
	}
	if log == nil {
		log = zap.NewNop()
	}
	if cfg.Audit.ModelEnabled && db == nil {
		return nil, fmt.Errorf("audit.model_enabled requires database connection")
	}
	if cfg.Audit.HttpEnabled && strings.EqualFold(cfg.Audit.HttpDriver, "db") && db == nil {
		return nil, fmt.Errorf("audit.http_enabled with driver db requires database connection")
	}
	if cfg.Audit.HttpEnabled && strings.EqualFold(cfg.Audit.HttpDriver, "file") && strings.TrimSpace(cfg.Audit.HttpFilePath) == "" {
		return nil, fmt.Errorf("audit.http_file_path required for file driver")
	}
	return &auditWriter{cfg: cfg, db: db, log: log}, nil
}

func (w *auditWriter) WriteActivity(ctx context.Context, row *ActivityLog) error {
	if w == nil || !w.cfg.Audit.ModelEnabled || w.db == nil {
		return nil
	}
	return w.db.Session(&gorm.Session{SkipHooks: true}).WithContext(ctx).Create(row).Error
}

func (w *auditWriter) WriteHTTP(ctx context.Context, row *HTTPAuditLog) error {
	if w == nil || !w.cfg.Audit.HttpEnabled {
		return nil
	}
	switch strings.ToLower(strings.TrimSpace(w.cfg.Audit.HttpDriver)) {
	case "file":
		return w.appendHTTPJSONL(row)
	case "db":
		if w.db == nil {
			return fmt.Errorf("audit http db writer: nil db")
		}
		return w.db.Session(&gorm.Session{SkipHooks: true}).WithContext(ctx).Create(row).Error
	default:
		return fmt.Errorf("unknown audit http driver %q", w.cfg.Audit.HttpDriver)
	}
}

func (w *auditWriter) appendHTTPJSONL(row *HTTPAuditLog) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	f, err := os.OpenFile(w.cfg.Audit.HttpFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	b, err := json.Marshal(row)
	if err != nil {
		return err
	}
	_, err = f.Write(append(b, '\n'))
	return err
}
