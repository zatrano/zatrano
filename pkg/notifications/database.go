package notifications

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// DatabaseChannel stores notifications in the database.
type DatabaseChannel struct {
	db *gorm.DB
}

// NewDatabaseChannel creates a database notification channel.
func NewDatabaseChannel(db *gorm.DB) *DatabaseChannel {
	return &DatabaseChannel{db: db}
}

// Name implements Channel.
func (c *DatabaseChannel) Name() string {
	return "database"
}

// Send implements Channel.
func (c *DatabaseChannel) Send(ctx context.Context, notif Notification) error {
	dataJSON, _ := json.Marshal(notif.Data())
	record := NotificationRecord{
		Type:      "notification",
		Subject:   notif.Subject(),
		Body:      notif.Body(),
		Data:      string(dataJSON),
		CreatedAt: time.Now(),
	}
	return c.db.WithContext(ctx).Create(&record).Error
}

// MarkAsRead marks a notification as read.
func (c *DatabaseChannel) MarkAsRead(ctx context.Context, notificationID uint) error {
	now := time.Now()
	return c.db.WithContext(ctx).Model(&NotificationRecord{}).
		Where("id = ?", notificationID).
		Update("read_at", now).Error
}

// GetUnread retrieves unread notifications for a user.
func (c *DatabaseChannel) GetUnread(ctx context.Context, userID uint) ([]NotificationRecord, error) {
	var records []NotificationRecord
	err := c.db.WithContext(ctx).
		Where("user_id = ? AND read_at IS NULL", userID).
		Order("created_at DESC").
		Find(&records).Error
	return records, err
}

// GetNotifications retrieves all notifications for a user (paginated).
func (c *DatabaseChannel) GetNotifications(ctx context.Context, userID uint, limit int, offset int) ([]NotificationRecord, error) {
	var records []NotificationRecord
	err := c.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&records).Error
	return records, err
}

// DeleteNotification deletes a notification.
func (c *DatabaseChannel) DeleteNotification(ctx context.Context, notificationID uint) error {
	return c.db.WithContext(ctx).Delete(&NotificationRecord{}, notificationID).Error
}

// ClearOldNotifications deletes notifications older than days.
func (c *DatabaseChannel) ClearOldNotifications(ctx context.Context, days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	return c.db.WithContext(ctx).
		Where("created_at < ?", cutoff).
		Delete(&NotificationRecord{}).Error
}
