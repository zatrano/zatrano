package database

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"zatrano/app/auth" // Auth helper'larımıza erişim (veya context key)
)

// AuditingPlugin, GORM için bir eklentidir ve denetim alanlarını otomatik olarak doldurur.
type AuditingPlugin struct {
	Context *fiber.Ctx
}

func (p *AuditingPlugin) Name() string {
	return "auditingPlugin"
}

func (p *AuditingPlugin) Initialize(db *gorm.DB) error {
	// Create işlemi öncesi çalışacak Hook
	db.Callback().Create().Before("gorm:before_create").Register("auditing:set_created_by", p.setCreatedBy)
	
	// Update işlemi öncesi çalışacak Hook
	db.Callback().Update().Before("gorm:before_update").Register("auditing:set_updated_by", p.setUpdatedBy)

	// Delete işlemi öncesi çalışacak Hook
	db.Callback().Delete().Before("gorm:before_delete").Register("auditing:set_deleted_by", p.setDeletedBy)
	
	return nil
}

func (p *AuditingPlugin) getAuthUserID() *uint {
	if p.Context == nil {
		return nil
	}
	// `auth.ID` helper'ı, context'te kullanıcı yoksa 0 döndürür.
	// Biz ise NULL olabilen bir pointer istiyoruz.
	if user := auth.User(p.Context); user != nil {
		id := user.ID
		return &id
	}
	return nil
}

func (p *AuditingPlugin) setCreatedBy(db *gorm.DB) {
	if userID := p.getAuthUserID(); userID != nil {
		// `SetColumn` metodu, modelin `CreatedBy` alanını, eğer varsa, ayarlar.
		db.Statement.SetColumn("CreatedBy", userID)
		db.Statement.SetColumn("UpdatedBy", userID)
	}
}

func (p *AuditingPlugin) setUpdatedBy(db *gorm.DB) {
	if userID := p.getAuthUserID(); userID != nil {
		db.Statement.SetColumn("UpdatedBy", userID)
	}
}

func (p *AuditingPlugin) setDeletedBy(db *gorm.DB) {
	if userID := p.getAuthUserID(); userID != nil {
		// Soft delete için, `DeletedBy` alanını ayarla.
		// Bu, GORM'un `deleted_at` alanını doldurmadan hemen önce çalışır.
		// `Clauses` ile sadece güncellenecek alanları belirtmek daha güvenli olabilir.
		db.Statement.SetColumn("DeletedBy", userID)
	}
}