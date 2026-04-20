package notifications

import (
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDatabaseChannel_MarkAsRead(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := db.AutoMigrate(&NotificationRecord{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	channel := NewDatabaseChannel(db)
	ctx := context.Background()

	// Create a notification
	notif := NewNotification("Test Subject", "Test Body", "test@example.com")
	if err := channel.Send(ctx, notif); err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	var records []NotificationRecord
	db.Find(&records)
	if len(records) == 0 {
		t.Fatalf("expected notification to be saved")
	}

	notifID := records[0].ID
	if err := channel.MarkAsRead(ctx, notifID); err != nil {
		t.Fatalf("MarkAsRead failed: %v", err)
	}

	var updated NotificationRecord
	db.First(&updated, notifID)
	if updated.ReadAt == nil {
		t.Fatalf("expected ReadAt to be set")
	}
}
