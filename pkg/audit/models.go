package audit

import (
	"encoding/json"
	"time"
)

// ActivityLog maps to zatrano_activity_logs (written by GORM hooks).
type ActivityLog struct {
	ID          uint64          `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt   time.Time       `gorm:"column:created_at" json:"created_at"`
	UserID      *string         `gorm:"column:user_id;size:128" json:"user_id,omitempty"`
	SubjectType string          `gorm:"column:subject_type;size:128;not null" json:"subject_type"`
	SubjectID   string          `gorm:"column:subject_id;size:64;not null" json:"subject_id"`
	Action      string          `gorm:"column:action;size:32;not null" json:"action"`
	Changes     json.RawMessage `gorm:"column:changes;type:jsonb" json:"changes,omitempty"`
	RequestID   *string         `gorm:"column:request_id;size:128" json:"request_id,omitempty"`
	IP          *string         `gorm:"column:ip;size:64" json:"ip,omitempty"`
	Metadata    json.RawMessage `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
}

func (ActivityLog) TableName() string { return "zatrano_activity_logs" }

// HTTPAuditLog maps to zatrano_http_audit_logs.
type HTTPAuditLog struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UserID     *string   `gorm:"column:user_id;size:128" json:"user_id,omitempty"`
	Method     string    `gorm:"column:method;size:16;not null" json:"method"`
	Path       string    `gorm:"column:path;size:2048;not null" json:"path"`
	URLQuery   *string   `gorm:"column:url_query;size:2048" json:"url_query,omitempty"`
	Status     int       `gorm:"column:status;not null" json:"status"`
	DurationMs int       `gorm:"column:duration_ms;not null" json:"duration_ms"`
	RequestID  *string   `gorm:"column:request_id;size:128" json:"request_id,omitempty"`
	IP         *string   `gorm:"column:ip;size:64" json:"ip,omitempty"`
}

func (HTTPAuditLog) TableName() string { return "zatrano_http_audit_logs" }
