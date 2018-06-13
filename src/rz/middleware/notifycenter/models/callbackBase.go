package models

import "time"

type CallbackBasePo struct {
	Finished      bool                     `gorm:"column:finished"`
	FinishedTime  time.Time                `gorm:"column:finishedTime"`
	States        string                   `gorm:"column:states"`
	ErrorMessages string                   `gorm:"column:errorMessages"`
}