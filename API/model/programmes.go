package model

import "time"

type Programme struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UUID        string     `gorm:"column:uuid;type:char(36);not null;uniqueIndex:programmes_uuid_unique" json:"uuid"`
	Name        string     `gorm:"column:name;type:mediumtext;not null" json:"name"`
	Description *string    `gorm:"column:description;type:longtext" json:"description"`
	Active      bool       `gorm:"column:active;not null;default:true" json:"active"`
	Settings    *string    `gorm:"column:settings;type:longtext" json:"settings"`
	SortOrder   int        `gorm:"column:sort_order;not null;default:1" json:"sort_order"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt   *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Programme) TableName() string {
	return "programmes"
}
