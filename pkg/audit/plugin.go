package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const instanceOldJSON = "zatrano_audit_old_json"

// RegisterGORM installs create/update/delete callbacks for types registered with RegisterSubject.
// w should be the same Writer used for HTTP audit (typically from NewWriter).
// Safe to call once at process startup; registering twice replaces callbacks with the same names.
func RegisterGORM(db *gorm.DB, w Writer, log *zap.Logger) {
	if db == nil || w == nil {
		return
	}
	if log == nil {
		log = zap.NewNop()
	}
	p := &gormPlugin{w: w, log: log}
	db.Callback().Create().After("gorm:create").Register("zatrano:audit:after_create", p.afterCreate)
	db.Callback().Update().Before("gorm:update").Register("zatrano:audit:before_update", p.beforeUpdate)
	db.Callback().Update().After("gorm:update").Register("zatrano:audit:after_update", p.afterUpdate)
	db.Callback().Delete().Before("gorm:delete").Register("zatrano:audit:before_delete", p.beforeDelete)
	db.Callback().Delete().After("gorm:delete").Register("zatrano:audit:after_delete", p.afterDelete)
}

type gormPlugin struct {
	w   Writer
	log *zap.Logger
}

func skipAuditTable(name string) bool {
	switch name {
	case "zatrano_activity_logs", "zatrano_http_audit_logs":
		return true
	default:
		return false
	}
}

func (p *gormPlugin) shouldSkip(db *gorm.DB) bool {
	if db.Statement == nil || db.Statement.Schema == nil {
		return true
	}
	if skipAuditTable(db.Statement.Table) {
		return true
	}
	if _, ok := subjectForModelType(db.Statement.Schema.ModelType); !ok {
		return true
	}
	if skipFromContext(db.Statement.Context) {
		return true
	}
	return false
}

func (p *gormPlugin) afterCreate(db *gorm.DB) {
	if p.shouldSkip(db) {
		return
	}
	st, _ := subjectForModelType(db.Statement.Schema.ModelType)
	id := primaryIDString(db.Statement)
	if id == "" {
		return
	}
	rv := reflect.Indirect(db.Statement.ReflectValue)
	nb, err := json.Marshal(rv.Interface())
	if err != nil {
		return
	}
	patch, err := DiffJSONPatch([]byte("{}"), nb)
	if err != nil {
		return
	}
	row := &ActivityLog{
		SubjectType: st,
		SubjectID:   id,
		Action:      "create",
		Changes:     patch,
	}
	p.fillMeta(db.Statement.Context, row)
	if err := p.w.WriteActivity(db.Statement.Context, row); err != nil && p.log != nil {
		p.log.Warn("audit write activity", zap.Error(err))
	}
}

func (p *gormPlugin) beforeUpdate(db *gorm.DB) {
	if p.shouldSkip(db) {
		return
	}
	stmt := db.Statement
	if stmt.Schema.PrioritizedPrimaryField == nil {
		return
	}
	id := primaryIDString(stmt)
	if id == "" {
		return
	}
	pk := stmt.Schema.PrioritizedPrimaryField.DBName
	sub := db.Session(&gorm.Session{NewDB: true, SkipHooks: true}).WithContext(stmt.Context)
	var old map[string]any
	if err := sub.Table(stmt.Schema.Table).Where(fmt.Sprintf("%s = ?", pk), id).Take(&old).Error; err != nil {
		return
	}
	b, err := json.Marshal(old)
	if err != nil {
		return
	}
	db.InstanceSet(instanceOldJSON, b)
}

func (p *gormPlugin) afterUpdate(db *gorm.DB) {
	if p.shouldSkip(db) {
		return
	}
	v, ok := db.InstanceGet(instanceOldJSON)
	if !ok {
		return
	}
	oldB, ok := v.([]byte)
	if !ok || len(oldB) == 0 {
		return
	}
	rv := reflect.Indirect(db.Statement.ReflectValue)
	nb, err := json.Marshal(rv.Interface())
	if err != nil {
		return
	}
	patch, err := DiffJSONPatch(oldB, nb)
	if err != nil {
		return
	}
	st, _ := subjectForModelType(db.Statement.Schema.ModelType)
	id := primaryIDString(db.Statement)
	row := &ActivityLog{
		SubjectType: st,
		SubjectID:   id,
		Action:      "update",
		Changes:     patch,
	}
	p.fillMeta(db.Statement.Context, row)
	if err := p.w.WriteActivity(db.Statement.Context, row); err != nil && p.log != nil {
		p.log.Warn("audit write activity", zap.Error(err))
	}
}

func (p *gormPlugin) beforeDelete(db *gorm.DB) {
	if p.shouldSkip(db) {
		return
	}
	stmt := db.Statement
	if stmt.Schema.PrioritizedPrimaryField == nil {
		return
	}
	id := primaryIDString(stmt)
	if id == "" {
		return
	}
	pk := stmt.Schema.PrioritizedPrimaryField.DBName
	sub := db.Session(&gorm.Session{NewDB: true, SkipHooks: true}).WithContext(stmt.Context)
	var old map[string]any
	if err := sub.Table(stmt.Schema.Table).Where(fmt.Sprintf("%s = ?", pk), id).Take(&old).Error; err != nil {
		return
	}
	b, err := json.Marshal(old)
	if err != nil {
		return
	}
	db.InstanceSet(instanceOldJSON, b)
}

func (p *gormPlugin) afterDelete(db *gorm.DB) {
	if p.shouldSkip(db) {
		return
	}
	v, ok := db.InstanceGet(instanceOldJSON)
	if !ok {
		return
	}
	oldB, ok := v.([]byte)
	if !ok || len(oldB) == 0 {
		return
	}
	patch, err := DiffJSONPatch(oldB, []byte("{}"))
	if err != nil {
		return
	}
	st, _ := subjectForModelType(db.Statement.Schema.ModelType)
	id := primaryIDString(db.Statement)
	row := &ActivityLog{
		SubjectType: st,
		SubjectID:   id,
		Action:      "delete",
		Changes:     patch,
	}
	p.fillMeta(db.Statement.Context, row)
	if err := p.w.WriteActivity(db.Statement.Context, row); err != nil && p.log != nil {
		p.log.Warn("audit write activity", zap.Error(err))
	}
}

func (p *gormPlugin) fillMeta(ctx context.Context, row *ActivityLog) {
	if u := UserFromContext(ctx); u != "" {
		row.UserID = &u
	}
	rid, ip := requestFromContext(ctx)
	if rid != "" {
		row.RequestID = &rid
	}
	if ip != "" {
		row.IP = &ip
	}
}

func primaryIDString(stmt *gorm.Statement) string {
	if stmt == nil || stmt.Schema == nil || stmt.Schema.PrioritizedPrimaryField == nil {
		return ""
	}
	rv := stmt.ReflectValue
	if !rv.IsValid() {
		return ""
	}
	rv = reflect.Indirect(rv)
	if rv.Kind() != reflect.Struct {
		return ""
	}
	pf := stmt.Schema.PrioritizedPrimaryField
	val, zero := pf.ValueOf(stmt.Context, rv)
	if zero {
		return ""
	}
	return fmt.Sprint(val)
}
