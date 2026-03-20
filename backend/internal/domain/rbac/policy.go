package rbac

import (
	"time"
)

// Policy - правило доступа
type Policy struct {
	ID        uint   `gorm:"primarykey"`
	PType     string `gorm:"column:ptype;size:100;index"` // p или g
	V0        string `gorm:"column:v0;size:255;index"`    // субъект (роль или пользователь)
	V1        string `gorm:"column:v1;size:255;index"`    // объект или роль
	V2        string `gorm:"column:v2;size:255"`          // действие или домен
	V3        string `gorm:"column:v3;size:255"`
	V4        string `gorm:"column:v4;size:255"`
	V5        string `gorm:"column:v5;size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName задаёт имя таблицы для GORM
func (Policy) TableName() string {
	return "casbin_rule"
}
