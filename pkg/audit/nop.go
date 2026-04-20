package audit

import "context"

// NopWriter discards audit rows (used when audit is disabled).
type NopWriter struct{}

func (NopWriter) WriteActivity(context.Context, *ActivityLog) error { return nil }
func (NopWriter) WriteHTTP(context.Context, *HTTPAuditLog) error    { return nil }
