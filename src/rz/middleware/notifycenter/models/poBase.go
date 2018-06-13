package models

import "time"

type PoBase struct {
	Comment       string    `gorm:"column:comment"`
	OperationUser string    `gorm:"column:operationUser"`
	CreatedTime   time.Time `gorm:"column:createdTime"`
	UpdatedTime   time.Time `gorm:"column:updatedTime"`
	Deleted       bool      `gorm:"column:deleted"`
	Version       int       `gorm:"column:version"`
}
