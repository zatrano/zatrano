package features

import "time"

// DBFlag maps zatrano_feature_flags (dynamic overrides / experiments).
type DBFlag struct {
	ID             uint      `gorm:"primaryKey"`
	Key            string    `gorm:"size:190;uniqueIndex;not null"`
	Enabled        bool      `gorm:"not null;default:false"`
	RolloutPercent int       `gorm:"not null;default:0"`
	AllowedRoles   []string  `gorm:"type:jsonb;serializer:json"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (DBFlag) TableName() string { return "zatrano_feature_flags" }
